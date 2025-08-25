package query

import (
	"fmt"
	"strings"

	"smart-log-analyser/pkg/parser"
)

// QueryEngine provides a high-level interface for executing queries
type QueryEngine struct {
	logs []*parser.LogEntry
}

// NewQueryEngine creates a new query engine
func NewQueryEngine(logs []*parser.LogEntry) *QueryEngine {
	return &QueryEngine{logs: logs}
}

// Query executes a query string and returns formatted results
func (qe *QueryEngine) Query(queryStr string, format string) (string, error) {
	result, err := qe.ExecuteQuery(queryStr)
	if err != nil {
		return "", err
	}

	return FormatResult(result, format)
}

// ExecuteQuery executes a query string and returns raw results
func (qe *QueryEngine) ExecuteQuery(queryStr string) (*QueryResult, error) {
	return ExecuteQuery(queryStr, qe.logs)
}

// ValidateQuery validates a query without executing it
func (qe *QueryEngine) ValidateQuery(queryStr string) error {
	_, err := ParseQuery(queryStr)
	return err
}

// GetAvailableFields returns the list of available fields for querying
func (qe *QueryEngine) GetAvailableFields() []string {
	return []string{
		"ip", "timestamp", "method", "url", "protocol",
		"status", "size", "referer", "user_agent",
	}
}

// GetAvailableFunctions returns the list of available functions
func (qe *QueryEngine) GetAvailableFunctions() []string {
	return []string{
		// Aggregate functions
		"COUNT", "SUM", "AVG", "MIN", "MAX",
		// Time functions
		"HOUR", "DAY", "WEEKDAY", "DATE", "TIME_DIFF",
		// String functions
		"UPPER", "LOWER", "LENGTH", "SUBSTR",
		// Network functions
		"IP_TO_INT", "IS_PRIVATE_IP", "COUNTRY",
	}
}

// GetAvailableOperators returns the list of available operators
func (qe *QueryEngine) GetAvailableOperators() []string {
	return []string{
		// Comparison operators
		"=", "!=", "<>", "<", "<=", ">", ">=",
		// String operators
		"LIKE", "MATCHES", "CONTAINS", "STARTS_WITH", "ENDS_WITH",
		// Special operators
		"IN", "BETWEEN", "IN_RANGE", "IS_BOT", "IS_ERROR", "IS_SUCCESS",
		// Logical operators
		"AND", "OR", "NOT",
	}
}

// GetSampleQueries returns example queries for documentation
func (qe *QueryEngine) GetSampleQueries() map[string]string {
	return map[string]string{
		"Basic filtering": `
			SELECT * FROM logs WHERE status = 404
		`,
		"Error analysis": `
			SELECT url, COUNT() as error_count 
			FROM logs 
			WHERE IS_ERROR(status) 
			GROUP BY url 
			ORDER BY error_count DESC 
			LIMIT 10
		`,
		"Bot traffic": `
			SELECT ip, COUNT() as requests 
			FROM logs 
			WHERE IS_BOT(user_agent) 
			GROUP BY ip 
			ORDER BY requests DESC
		`,
		"Time-based analysis": `
			SELECT HOUR(timestamp) as hour, COUNT() as requests 
			FROM logs 
			GROUP BY hour 
			ORDER BY hour
		`,
		"Security analysis": `
			SELECT ip, COUNT() as attempts 
			FROM logs 
			WHERE status = 401 AND url LIKE '/admin*' 
			GROUP BY ip 
			HAVING attempts > 5 
			ORDER BY attempts DESC
		`,
		"Large requests": `
			SELECT url, AVG(size) as avg_size, COUNT() as count 
			FROM logs 
			WHERE size > 100000 
			GROUP BY url 
			ORDER BY avg_size DESC
		`,
		"Geographic analysis": `
			SELECT COUNTRY(ip) as country, COUNT() as requests 
			FROM logs 
			WHERE NOT IS_PRIVATE_IP(ip) 
			GROUP BY country 
			ORDER BY requests DESC 
			LIMIT 20
		`,
		"Complex filtering": `
			SELECT method, status, COUNT() as count 
			FROM logs 
			WHERE timestamp BETWEEN '2024-08-20 00:00:00' AND '2024-08-20 23:59:59' 
			AND (status >= 400 OR size > 1000000) 
			GROUP BY method, status 
			ORDER BY count DESC
		`,
	}
}

// BuildQuery provides a fluent interface for building queries
type QueryBuilder struct {
	selectFields []string
	fromTable    string
	whereClause  string
	groupByFields []string
	orderByFields []string
	havingClause string
	limitValue   *int64
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		fromTable: "logs", // Default table
	}
}

// Select specifies the fields to select
func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.selectFields = fields
	return qb
}

// From specifies the table name
func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.fromTable = table
	return qb
}

// Where adds a WHERE condition
func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
	if qb.whereClause != "" {
		qb.whereClause += " AND " + condition
	} else {
		qb.whereClause = condition
	}
	return qb
}

// WhereOr adds a WHERE condition with OR
func (qb *QueryBuilder) WhereOr(condition string) *QueryBuilder {
	if qb.whereClause != "" {
		qb.whereClause += " OR " + condition
	} else {
		qb.whereClause = condition
	}
	return qb
}

// GroupBy specifies GROUP BY fields
func (qb *QueryBuilder) GroupBy(fields ...string) *QueryBuilder {
	qb.groupByFields = fields
	return qb
}

// OrderBy specifies ORDER BY fields
func (qb *QueryBuilder) OrderBy(fields ...string) *QueryBuilder {
	qb.orderByFields = fields
	return qb
}

// Having adds a HAVING condition
func (qb *QueryBuilder) Having(condition string) *QueryBuilder {
	qb.havingClause = condition
	return qb
}

// Limit sets the LIMIT value
func (qb *QueryBuilder) Limit(limit int64) *QueryBuilder {
	qb.limitValue = &limit
	return qb
}

// Build constructs the final query string
func (qb *QueryBuilder) Build() string {
	var query strings.Builder

	// SELECT clause
	query.WriteString("SELECT ")
	if len(qb.selectFields) == 0 {
		query.WriteString("*")
	} else {
		query.WriteString(strings.Join(qb.selectFields, ", "))
	}

	// FROM clause
	query.WriteString(" FROM " + qb.fromTable)

	// WHERE clause
	if qb.whereClause != "" {
		query.WriteString(" WHERE " + qb.whereClause)
	}

	// GROUP BY clause
	if len(qb.groupByFields) > 0 {
		query.WriteString(" GROUP BY " + strings.Join(qb.groupByFields, ", "))
	}

	// HAVING clause
	if qb.havingClause != "" {
		query.WriteString(" HAVING " + qb.havingClause)
	}

	// ORDER BY clause
	if len(qb.orderByFields) > 0 {
		query.WriteString(" ORDER BY " + strings.Join(qb.orderByFields, ", "))
	}

	// LIMIT clause
	if qb.limitValue != nil {
		query.WriteString(fmt.Sprintf(" LIMIT %d", *qb.limitValue))
	}

	return query.String()
}

// PrebuiltQueries provides common query templates
type PrebuiltQueries struct{}

// NewPrebuiltQueries creates a new instance
func NewPrebuiltQueries() *PrebuiltQueries {
	return &PrebuiltQueries{}
}

// ErrorAnalysis returns a query for analyzing errors
func (pq *PrebuiltQueries) ErrorAnalysis() *QueryBuilder {
	return NewQueryBuilder().
		Select("url", "status", "COUNT() as error_count").
		Where("IS_ERROR(status)").
		GroupBy("url", "status").
		OrderBy("error_count DESC").
		Limit(20)
}

// TopIPs returns a query for finding top IP addresses
func (pq *PrebuiltQueries) TopIPs() *QueryBuilder {
	return NewQueryBuilder().
		Select("ip", "COUNT() as requests").
		GroupBy("ip").
		OrderBy("requests DESC").
		Limit(20)
}

// BotTraffic returns a query for analyzing bot traffic
func (pq *PrebuiltQueries) BotTraffic() *QueryBuilder {
	return NewQueryBuilder().
		Select("ip", "user_agent", "COUNT() as requests").
		Where("IS_BOT(user_agent)").
		GroupBy("ip", "user_agent").
		OrderBy("requests DESC").
		Limit(10)
}

// HourlyTraffic returns a query for hourly traffic analysis
func (pq *PrebuiltQueries) HourlyTraffic() *QueryBuilder {
	return NewQueryBuilder().
		Select("HOUR(timestamp) as hour", "COUNT() as requests").
		GroupBy("hour").
		OrderBy("hour")
}

// LargeRequests returns a query for finding large requests
func (pq *PrebuiltQueries) LargeRequests(minSize int64) *QueryBuilder {
	return NewQueryBuilder().
		Select("url", "method", "AVG(size) as avg_size", "COUNT() as count").
		Where(fmt.Sprintf("size > %d", minSize)).
		GroupBy("url", "method").
		OrderBy("avg_size DESC").
		Limit(10)
}

// SecurityThreats returns a query for security analysis
func (pq *PrebuiltQueries) SecurityThreats() *QueryBuilder {
	return NewQueryBuilder().
		Select("ip", "url", "COUNT() as attempts").
		Where("(status = 401 OR status = 403) AND (url LIKE '/admin*' OR url LIKE '/login*')").
		GroupBy("ip", "url").
		Having("attempts > 3").
		OrderBy("attempts DESC")
}

// StatusCodeDistribution returns a query for status code analysis
func (pq *PrebuiltQueries) StatusCodeDistribution() *QueryBuilder {
	return NewQueryBuilder().
		Select("status", "COUNT() as count").
		GroupBy("status").
		OrderBy("count DESC")
}

// GeographicAnalysis returns a query for geographic analysis
func (pq *PrebuiltQueries) GeographicAnalysis() *QueryBuilder {
	return NewQueryBuilder().
		Select("COUNTRY(ip) as country", "COUNT() as requests").
		Where("NOT IS_PRIVATE_IP(ip)").
		GroupBy("country").
		OrderBy("requests DESC").
		Limit(20)
}

// MethodAnalysis returns a query for HTTP method analysis
func (pq *PrebuiltQueries) MethodAnalysis() *QueryBuilder {
	return NewQueryBuilder().
		Select("method", "COUNT() as count", "AVG(size) as avg_size").
		GroupBy("method").
		OrderBy("count DESC")
}

// TimeRangeAnalysis returns a query for analyzing a specific time range
func (pq *PrebuiltQueries) TimeRangeAnalysis(startTime, endTime string) *QueryBuilder {
	condition := fmt.Sprintf("timestamp BETWEEN '%s' AND '%s'", startTime, endTime)
	return NewQueryBuilder().
		Select("HOUR(timestamp) as hour", "COUNT() as requests", "AVG(size) as avg_size").
		Where(condition).
		GroupBy("hour").
		OrderBy("hour")
}

// QueryHelper provides utilities for query construction and validation
type QueryHelper struct{}

// NewQueryHelper creates a new query helper
func NewQueryHelper() *QueryHelper {
	return &QueryHelper{}
}

// EscapeString escapes a string for use in queries
func (qh *QueryHelper) EscapeString(str string) string {
	return strings.ReplaceAll(str, "'", "''")
}

// ValidateFieldName checks if a field name is valid
func (qh *QueryHelper) ValidateFieldName(field string) bool {
	validFields := map[string]bool{
		"ip": true, "timestamp": true, "method": true, "url": true,
		"protocol": true, "status": true, "size": true, "referer": true, "user_agent": true,
	}
	return validFields[strings.ToLower(field)]
}

// ValidateFunctionName checks if a function name is valid
func (qh *QueryHelper) ValidateFunctionName(function string) bool {
	validFunctions := map[string]bool{
		"count": true, "sum": true, "avg": true, "min": true, "max": true,
		"hour": true, "day": true, "weekday": true, "date": true,
		"upper": true, "lower": true, "length": true, "substr": true,
		"is_private_ip": true, "country": true,
	}
	return validFunctions[strings.ToLower(function)]
}

// SuggestCorrection suggests corrections for common query errors
func (qh *QueryHelper) SuggestCorrection(err error) string {
	errorMsg := strings.ToLower(err.Error())

	suggestions := map[string]string{
		"unknown field":    "Available fields: ip, timestamp, method, url, protocol, status, size, referer, user_agent",
		"unknown function": "Available functions: COUNT, SUM, AVG, MIN, MAX, HOUR, DAY, UPPER, LOWER, etc.",
		"syntax error":     "Check for missing quotes, parentheses, or keywords like SELECT, FROM, WHERE",
		"invalid operator": "Available operators: =, !=, <, >, LIKE, CONTAINS, IN, BETWEEN, IS_BOT, etc.",
	}

	for pattern, suggestion := range suggestions {
		if strings.Contains(errorMsg, pattern) {
			return suggestion
		}
	}

	return "Check the query syntax and available fields/functions"
}