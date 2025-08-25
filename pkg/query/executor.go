package query

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"smart-log-analyser/pkg/parser"
)

// Executor executes queries against log entries
type Executor struct {
	logs []*parser.LogEntry
}

// NewExecutor creates a new query executor
func NewExecutor(logs []*parser.LogEntry) *Executor {
	return &Executor{logs: logs}
}

// Execute executes a parsed query and returns results
func (e *Executor) Execute(stmt *SelectStatement) (*QueryResult, error) {
	// Filter logs based on WHERE clause
	filteredLogs := e.logs
	if stmt.Where != nil {
		var err error
		filteredLogs, err = e.filterLogs(filteredLogs, stmt.Where)
		if err != nil {
			return nil, fmt.Errorf("error filtering logs: %w", err)
		}
	}

	// Handle GROUP BY
	if len(stmt.GroupBy) > 0 {
		return e.executeGroupBy(stmt, filteredLogs)
	}

	// Handle simple SELECT
	return e.executeSelect(stmt, filteredLogs)
}

// filterLogs filters logs based on WHERE clause
func (e *Executor) filterLogs(logs []*parser.LogEntry, where Expression) ([]*parser.LogEntry, error) {
	var filtered []*parser.LogEntry

	for _, log := range logs {
		result, err := where.Evaluate(log)
		if err != nil {
			continue // Skip logs that cause evaluation errors
		}

		match, err := toBool(result)
		if err != nil {
			continue // Skip non-boolean results
		}

		if match {
			filtered = append(filtered, log)
		}
	}

	return filtered, nil
}

// executeSelect executes a simple SELECT without GROUP BY
func (e *Executor) executeSelect(stmt *SelectStatement, logs []*parser.LogEntry) (*QueryResult, error) {
	result := &QueryResult{}

	// Determine columns
	if len(stmt.Fields) == 1 && stmt.Fields[0].Expression.String() == "*" {
		// SELECT * - return all fields
		result.Columns = []string{"IP", "Timestamp", "Method", "URL", "Protocol", "Status", "Size", "Referer", "UserAgent"}
	} else {
		// SELECT specific fields
		for _, field := range stmt.Fields {
			if field.Alias != "" {
				result.Columns = append(result.Columns, field.Alias)
			} else {
				result.Columns = append(result.Columns, field.Expression.String())
			}
		}
	}

	// Process each log entry
	for _, log := range logs {
		var row []Value

		if len(stmt.Fields) == 1 && stmt.Fields[0].Expression.String() == "*" {
			// SELECT * - return all field values
			row = []Value{
				{Type: ValueString, StringVal: log.IP},
				{Type: ValueTime, TimeVal: log.Timestamp},
				{Type: ValueString, StringVal: log.Method},
				{Type: ValueString, StringVal: log.URL},
				{Type: ValueString, StringVal: log.Protocol},
				{Type: ValueInt, IntVal: int64(log.Status)},
				{Type: ValueInt, IntVal: log.Size},
				{Type: ValueString, StringVal: log.Referer},
				{Type: ValueString, StringVal: log.UserAgent},
			}
		} else {
			// SELECT specific fields
			for _, field := range stmt.Fields {
				value, err := field.Expression.Evaluate(log)
				if err != nil {
					// Use NULL-like value for errors
					value = Value{Type: ValueString, StringVal: ""}
				}
				row = append(row, value)
			}
		}

		result.Rows = append(result.Rows, row)
	}

	// Apply ORDER BY
	if len(stmt.OrderBy) > 0 {
		err := e.sortRows(result, stmt.OrderBy, logs)
		if err != nil {
			return nil, fmt.Errorf("error sorting results: %w", err)
		}
	}

	// Apply LIMIT
	if stmt.Limit != nil {
		limit := int(*stmt.Limit)
		if limit < len(result.Rows) {
			result.Rows = result.Rows[:limit]
		}
	}

	result.Count = len(result.Rows)
	return result, nil
}

// executeGroupBy executes a SELECT with GROUP BY
func (e *Executor) executeGroupBy(stmt *SelectStatement, logs []*parser.LogEntry) (*QueryResult, error) {
	// Group logs by GROUP BY expressions
	groups, err := e.groupLogs(logs, stmt.GroupBy)
	if err != nil {
		return nil, fmt.Errorf("error grouping logs: %w", err)
	}

	result := &QueryResult{}

	// Build column names
	for _, expr := range stmt.GroupBy {
		result.Columns = append(result.Columns, expr.String())
	}

	for _, field := range stmt.Fields {
		// Skip aggregate functions that are already in GROUP BY
		if !e.isGroupByExpression(field.Expression, stmt.GroupBy) {
			if field.Alias != "" {
				result.Columns = append(result.Columns, field.Alias)
			} else {
				result.Columns = append(result.Columns, field.Expression.String())
			}
		}
	}

	// Process each group
	for _, group := range groups {
		// Build row starting with group key values
		var row []Value
		for _, keyValue := range group.KeyValues {
			row = append(row, keyValue)
		}

		// Evaluate aggregate functions for this group
		for _, field := range stmt.Fields {
			if !e.isGroupByExpression(field.Expression, stmt.GroupBy) {
				value, err := e.evaluateAggregate(field.Expression, group.Logs)
				if err != nil {
					value = Value{Type: ValueString, StringVal: ""}
				}
				row = append(row, value)
			}
		}

		// Apply HAVING filter if present
		if stmt.Having != nil {
			// Create a dummy log entry with aggregated values for HAVING evaluation
			dummyLog := &parser.LogEntry{}
			havingResult, err := stmt.Having.Evaluate(dummyLog)
			if err != nil {
				continue
			}
			match, err := toBool(havingResult)
			if err != nil || !match {
				continue
			}
		}

		result.Rows = append(result.Rows, row)
	}

	// Apply ORDER BY
	if len(stmt.OrderBy) > 0 {
		err := e.sortGroupedRows(result, stmt.OrderBy)
		if err != nil {
			return nil, fmt.Errorf("error sorting results: %w", err)
		}
	}

	// Apply LIMIT
	if stmt.Limit != nil {
		limit := int(*stmt.Limit)
		if limit < len(result.Rows) {
			result.Rows = result.Rows[:limit]
		}
	}

	result.Count = len(result.Rows)
	return result, nil
}

// GroupData represents grouped log data
type GroupData struct {
	KeyValues []Value
	Logs []*parser.LogEntry
}

// groupLogs groups logs by the specified expressions
func (e *Executor) groupLogs(logs []*parser.LogEntry, groupBy []Expression) (map[string]GroupData, error) {
	groups := make(map[string]GroupData)

	for _, log := range logs {
		// Evaluate group key
		var keyValues []Value
		for _, expr := range groupBy {
			value, err := expr.Evaluate(log)
			if err != nil {
				// Use empty string for errors
				value = Value{Type: ValueString, StringVal: ""}
			}
			keyValues = append(keyValues, value)
		}

		// Create string key for grouping
		keyStr := e.groupKeyToString(keyValues)

		// Add to group
		if group, exists := groups[keyStr]; exists {
			group.Logs = append(group.Logs, log)
			groups[keyStr] = group
		} else {
			groups[keyStr] = GroupData{
				KeyValues: keyValues,
				Logs:      []*parser.LogEntry{log},
			}
		}
	}

	return groups, nil
}

// groupKeyToString converts group key values to a string
func (e *Executor) groupKeyToString(values []Value) string {
	var parts []string
	for _, value := range values {
		parts = append(parts, value.String())
	}
	return strings.Join(parts, "|")
}

// isGroupByExpression checks if an expression is in the GROUP BY clause
func (e *Executor) isGroupByExpression(expr Expression, groupBy []Expression) bool {
	exprStr := expr.String()
	for _, groupExpr := range groupBy {
		if groupExpr.String() == exprStr {
			return true
		}
	}
	return false
}

// evaluateAggregate evaluates aggregate functions
func (e *Executor) evaluateAggregate(expr Expression, logs []*parser.LogEntry) (Value, error) {
	// Check if it's a function expression
	if funcExpr, ok := expr.(*FunctionExpression); ok {
		switch strings.ToUpper(funcExpr.Name) {
		case "COUNT":
			return Value{Type: ValueInt, IntVal: int64(len(logs))}, nil

		case "SUM":
			if len(funcExpr.Arguments) != 1 {
				return Value{}, fmt.Errorf("SUM requires exactly 1 argument")
			}
			var sum float64
			for _, log := range logs {
				val, err := funcExpr.Arguments[0].Evaluate(log)
				if err != nil {
					continue
				}
				switch val.Type {
				case ValueInt:
					sum += float64(val.IntVal)
				case ValueFloat:
					sum += val.FloatVal
				}
			}
			return Value{Type: ValueFloat, FloatVal: sum}, nil

		case "AVG":
			if len(funcExpr.Arguments) != 1 {
				return Value{}, fmt.Errorf("AVG requires exactly 1 argument")
			}
			var sum float64
			count := 0
			for _, log := range logs {
				val, err := funcExpr.Arguments[0].Evaluate(log)
				if err != nil {
					continue
				}
				switch val.Type {
				case ValueInt:
					sum += float64(val.IntVal)
					count++
				case ValueFloat:
					sum += val.FloatVal
					count++
				}
			}
			if count == 0 {
				return Value{Type: ValueFloat, FloatVal: 0}, nil
			}
			return Value{Type: ValueFloat, FloatVal: sum / float64(count)}, nil

		case "MIN", "MAX":
			if len(funcExpr.Arguments) != 1 {
				return Value{}, fmt.Errorf("%s requires exactly 1 argument", funcExpr.Name)
			}
			if len(logs) == 0 {
				return Value{Type: ValueInt, IntVal: 0}, nil
			}
			
			firstVal, err := funcExpr.Arguments[0].Evaluate(logs[0])
			if err != nil {
				return Value{Type: ValueInt, IntVal: 0}, nil
			}
			
			result := firstVal
			for _, log := range logs[1:] {
				val, err := funcExpr.Arguments[0].Evaluate(log)
				if err != nil {
					continue
				}
				
				cmp := e.compareValues(val, result)
				if (funcExpr.Name == "MIN" && cmp < 0) || (funcExpr.Name == "MAX" && cmp > 0) {
					result = val
				}
			}
			return result, nil
		}
	}

	// For non-aggregate expressions, return the first value
	if len(logs) > 0 {
		return expr.Evaluate(logs[0])
	}

	return Value{Type: ValueString, StringVal: ""}, nil
}

// sortRows sorts result rows based on ORDER BY clause
func (e *Executor) sortRows(result *QueryResult, orderBy []OrderByClause, logs []*parser.LogEntry) error {
	sort.Slice(result.Rows, func(i, j int) bool {
		for _, clause := range orderBy {
			// Evaluate the expression for both rows
			val1, err1 := clause.Expression.Evaluate(logs[i])
			val2, err2 := clause.Expression.Evaluate(logs[j])

			if err1 != nil || err2 != nil {
				continue
			}

			// Compare values
			cmp := e.compareValues(val1, val2)
			if cmp == 0 {
				continue
			}

			if clause.Descending {
				return cmp > 0
			}
			return cmp < 0
		}
		return false
	})

	return nil
}

// sortGroupedRows sorts grouped results
func (e *Executor) sortGroupedRows(result *QueryResult, orderBy []OrderByClause) error {
	sort.Slice(result.Rows, func(i, j int) bool {
		for _, clause := range orderBy {
			// For grouped results, we need to match column names
			colIndex := e.findColumnIndex(result.Columns, clause.Expression.String())
			if colIndex == -1 {
				continue
			}

			if colIndex >= len(result.Rows[i]) || colIndex >= len(result.Rows[j]) {
				continue
			}

			val1 := result.Rows[i][colIndex]
			val2 := result.Rows[j][colIndex]

			cmp := e.compareValues(val1, val2)
			if cmp == 0 {
				continue
			}

			if clause.Descending {
				return cmp > 0
			}
			return cmp < 0
		}
		return false
	})

	return nil
}

// findColumnIndex finds the index of a column by name
func (e *Executor) findColumnIndex(columns []string, name string) int {
	for i, col := range columns {
		if col == name {
			return i
		}
	}
	return -1
}

// compareValues compares two values and returns -1, 0, or 1
func (e *Executor) compareValues(v1, v2 Value) int {
	// Type coercion if needed
	if v1.Type != v2.Type {
		v1, v2, _ = coerceValues(v1, v2)
	}

	switch v1.Type {
	case ValueString:
		return strings.Compare(v1.StringVal, v2.StringVal)
	case ValueInt:
		if v1.IntVal == v2.IntVal {
			return 0
		} else if v1.IntVal < v2.IntVal {
			return -1
		} else {
			return 1
		}
	case ValueFloat:
		if v1.FloatVal == v2.FloatVal {
			return 0
		} else if v1.FloatVal < v2.FloatVal {
			return -1
		} else {
			return 1
		}
	case ValueTime:
		if v1.TimeVal.Equal(v2.TimeVal) {
			return 0
		} else if v1.TimeVal.Before(v2.TimeVal) {
			return -1
		} else {
			return 1
		}
	}
	return 0
}

// ExecuteQuery is a convenience function to execute a query string
func ExecuteQuery(query string, logs []*parser.LogEntry) (*QueryResult, error) {
	// Parse the query
	stmt, err := ParseQuery(query)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	// Execute the query
	executor := NewExecutor(logs)
	return executor.Execute(stmt)
}

// FormatResult formats a query result for display
func FormatResult(result *QueryResult, format string) (string, error) {
	switch strings.ToLower(format) {
	case "table", "":
		return formatAsTable(result), nil
	case "csv":
		return formatAsCSV(result), nil
	case "json":
		return formatAsJSON(result), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// formatAsTable formats result as a table
func formatAsTable(result *QueryResult) string {
	if len(result.Rows) == 0 {
		return "No results found."
	}

	var output strings.Builder

	// Header
	output.WriteString(strings.Join(result.Columns, " | "))
	output.WriteString("\n")
	output.WriteString(strings.Repeat("-", len(strings.Join(result.Columns, " | "))))
	output.WriteString("\n")

	// Rows
	for _, row := range result.Rows {
		var rowStrs []string
		for _, value := range row {
			rowStrs = append(rowStrs, formatValue(value))
		}
		output.WriteString(strings.Join(rowStrs, " | "))
		output.WriteString("\n")
	}

	output.WriteString(fmt.Sprintf("\nTotal: %d rows\n", result.Count))
	return output.String()
}

// formatAsCSV formats result as CSV
func formatAsCSV(result *QueryResult) string {
	var output strings.Builder

	// Header
	output.WriteString(strings.Join(result.Columns, ","))
	output.WriteString("\n")

	// Rows
	for _, row := range result.Rows {
		var rowStrs []string
		for _, value := range row {
			val := formatValue(value)
			// Escape CSV values
			if strings.Contains(val, ",") || strings.Contains(val, "\"") {
				val = "\"" + strings.ReplaceAll(val, "\"", "\"\"") + "\""
			}
			rowStrs = append(rowStrs, val)
		}
		output.WriteString(strings.Join(rowStrs, ","))
		output.WriteString("\n")
	}

	return output.String()
}

// formatAsJSON formats result as JSON (simplified)
func formatAsJSON(result *QueryResult) string {
	var output strings.Builder
	output.WriteString("{\n")
	output.WriteString(fmt.Sprintf("  \"count\": %d,\n", result.Count))
	output.WriteString("  \"columns\": [")
	for i, col := range result.Columns {
		if i > 0 {
			output.WriteString(", ")
		}
		output.WriteString(fmt.Sprintf("\"%s\"", col))
	}
	output.WriteString("],\n")
	output.WriteString("  \"rows\": [\n")

	for i, row := range result.Rows {
		if i > 0 {
			output.WriteString(",\n")
		}
		output.WriteString("    [")
		for j, value := range row {
			if j > 0 {
				output.WriteString(", ")
			}
			output.WriteString(formatValueAsJSON(value))
		}
		output.WriteString("]")
	}

	output.WriteString("\n  ]\n}")
	return output.String()
}

// formatValue formats a value for display
func formatValue(value Value) string {
	switch value.Type {
	case ValueString:
		return value.StringVal
	case ValueInt:
		return strconv.FormatInt(value.IntVal, 10)
	case ValueFloat:
		return strconv.FormatFloat(value.FloatVal, 'f', 2, 64)
	case ValueBool:
		return strconv.FormatBool(value.BoolVal)
	case ValueTime:
		return value.TimeVal.Format("2006-01-02 15:04:05")
	default:
		return ""
	}
}

// formatValueAsJSON formats a value as JSON
func formatValueAsJSON(value Value) string {
	switch value.Type {
	case ValueString:
		return fmt.Sprintf("\"%s\"", strings.ReplaceAll(value.StringVal, "\"", "\\\""))
	case ValueInt:
		return strconv.FormatInt(value.IntVal, 10)
	case ValueFloat:
		return strconv.FormatFloat(value.FloatVal, 'f', 2, 64)
	case ValueBool:
		return strconv.FormatBool(value.BoolVal)
	case ValueTime:
		return fmt.Sprintf("\"%s\"", value.TimeVal.Format("2006-01-02T15:04:05Z"))
	default:
		return "null"
	}
}