package query

import (
	"fmt"
	"time"

	"smart-log-analyser/pkg/parser"
)

// TokenType represents different types of tokens in the query language
type TokenType int

const (
	// Literals
	TokenString TokenType = iota
	TokenNumber
	TokenBool
	TokenDate

	// Identifiers
	TokenField
	TokenFunction

	// Operators
	TokenEquals
	TokenNotEquals
	TokenLessThan
	TokenLessThanOrEqual
	TokenGreaterThan
	TokenGreaterThanOrEqual
	TokenLike
	TokenMatches
	TokenContains
	TokenStartsWith
	TokenEndsWith
	TokenIn
	TokenBetween
	TokenInRange
	TokenIsBot
	TokenIsError
	TokenIsSuccess

	// Logical operators
	TokenAnd
	TokenOr
	TokenNot

	// Keywords
	TokenSelect
	TokenFrom
	TokenWhere
	TokenGroup
	TokenBy
	TokenGroupBy
	TokenOrder
	TokenOrderBy
	TokenHaving
	TokenLimit
	TokenAs

	// Punctuation
	TokenLeftParen
	TokenRightParen
	TokenComma
	TokenSemicolon

	// Special
	TokenEOF
	TokenInvalid
)

// Token represents a single token in the query
type Token struct {
	Type     TokenType
	Value    string
	Position int
}

// QueryField represents available log fields for querying
type QueryField string

const (
	FieldIP        QueryField = "ip"
	FieldTimestamp QueryField = "timestamp"
	FieldMethod    QueryField = "method"
	FieldURL       QueryField = "url"
	FieldProtocol  QueryField = "protocol"
	FieldStatus    QueryField = "status"
	FieldSize      QueryField = "size"
	FieldReferer   QueryField = "referer"
	FieldUserAgent QueryField = "user_agent"
)

// Operator represents comparison and logical operators
type Operator string

const (
	OpEquals              Operator = "="
	OpNotEquals           Operator = "!="
	OpLessThan            Operator = "<"
	OpLessThanOrEqual     Operator = "<="
	OpGreaterThan         Operator = ">"
	OpGreaterThanOrEqual  Operator = ">="
	OpLike                Operator = "LIKE"
	OpMatches             Operator = "MATCHES"
	OpContains            Operator = "CONTAINS"
	OpStartsWith          Operator = "STARTS_WITH"
	OpEndsWith            Operator = "ENDS_WITH"
	OpIn                  Operator = "IN"
	OpBetween             Operator = "BETWEEN"
	OpInRange             Operator = "IN_RANGE"
	OpIsBot               Operator = "IS_BOT"
	OpIsError             Operator = "IS_ERROR"
	OpIsSuccess           Operator = "IS_SUCCESS"
	OpAnd                 Operator = "AND"
	OpOr                  Operator = "OR"
	OpNot                 Operator = "NOT"
)

// Value represents a query value with its type
type Value struct {
	Type      ValueType
	StringVal string
	IntVal    int64
	FloatVal  float64
	BoolVal   bool
	TimeVal   time.Time
	ListVal   []Value
}

// ValueType represents the type of a query value
type ValueType int

const (
	ValueString ValueType = iota
	ValueInt
	ValueFloat
	ValueBool
	ValueTime
	ValueList
)

// String returns string representation of Value
func (v Value) String() string {
	switch v.Type {
	case ValueString:
		return fmt.Sprintf("'%s'", v.StringVal)
	case ValueInt:
		return fmt.Sprintf("%d", v.IntVal)
	case ValueFloat:
		return fmt.Sprintf("%.2f", v.FloatVal)
	case ValueBool:
		return fmt.Sprintf("%t", v.BoolVal)
	case ValueTime:
		return fmt.Sprintf("'%s'", v.TimeVal.Format("2006-01-02 15:04:05"))
	case ValueList:
		result := "("
		for i, val := range v.ListVal {
			if i > 0 {
				result += ", "
			}
			result += val.String()
		}
		result += ")"
		return result
	default:
		return "unknown"
	}
}

// AST Node interfaces
type Node interface {
	String() string
}

type Expression interface {
	Node
	Evaluate(entry *parser.LogEntry) (Value, error)
}

type Statement interface {
	Node
}

// SelectStatement represents a SELECT query
type SelectStatement struct {
	Fields   []SelectField
	From     string
	Where    Expression
	GroupBy  []Expression
	OrderBy  []OrderByClause
	Having   Expression
	Limit    *int64
}

func (s SelectStatement) String() string {
	result := "SELECT "
	for i, field := range s.Fields {
		if i > 0 {
			result += ", "
		}
		result += field.String()
	}
	result += " FROM " + s.From
	if s.Where != nil {
		result += " WHERE " + s.Where.String()
	}
	if len(s.GroupBy) > 0 {
		result += " GROUP BY "
		for i, expr := range s.GroupBy {
			if i > 0 {
				result += ", "
			}
			result += expr.String()
		}
	}
	if s.Having != nil {
		result += " HAVING " + s.Having.String()
	}
	if len(s.OrderBy) > 0 {
		result += " ORDER BY "
		for i, clause := range s.OrderBy {
			if i > 0 {
				result += ", "
			}
			result += clause.String()
		}
	}
	if s.Limit != nil {
		result += fmt.Sprintf(" LIMIT %d", *s.Limit)
	}
	return result
}

// SelectField represents a field in SELECT clause
type SelectField struct {
	Expression Expression
	Alias      string
}

func (sf SelectField) String() string {
	result := sf.Expression.String()
	if sf.Alias != "" {
		result += " AS " + sf.Alias
	}
	return result
}

// OrderByClause represents ORDER BY clause
type OrderByClause struct {
	Expression Expression
	Descending bool
}

func (ob OrderByClause) String() string {
	result := ob.Expression.String()
	if ob.Descending {
		result += " DESC"
	}
	return result
}

// FieldExpression represents a field reference
type FieldExpression struct {
	Field QueryField
}

func (fe FieldExpression) String() string {
	return string(fe.Field)
}

func (fe FieldExpression) Evaluate(entry *parser.LogEntry) (Value, error) {
	switch fe.Field {
	case FieldIP:
		return Value{Type: ValueString, StringVal: entry.IP}, nil
	case FieldTimestamp:
		return Value{Type: ValueTime, TimeVal: entry.Timestamp}, nil
	case FieldMethod:
		return Value{Type: ValueString, StringVal: entry.Method}, nil
	case FieldURL:
		return Value{Type: ValueString, StringVal: entry.URL}, nil
	case FieldProtocol:
		return Value{Type: ValueString, StringVal: entry.Protocol}, nil
	case FieldStatus:
		return Value{Type: ValueInt, IntVal: int64(entry.Status)}, nil
	case FieldSize:
		return Value{Type: ValueInt, IntVal: entry.Size}, nil
	case FieldReferer:
		return Value{Type: ValueString, StringVal: entry.Referer}, nil
	case FieldUserAgent:
		return Value{Type: ValueString, StringVal: entry.UserAgent}, nil
	default:
		return Value{}, fmt.Errorf("unknown field: %s", fe.Field)
	}
}

// LiteralExpression represents a literal value
type LiteralExpression struct {
	Value Value
}

func (le LiteralExpression) String() string {
	return le.Value.String()
}

func (le LiteralExpression) Evaluate(entry *parser.LogEntry) (Value, error) {
	return le.Value, nil
}

// BinaryExpression represents binary operations
type BinaryExpression struct {
	Left     Expression
	Operator Operator
	Right    Expression
}

func (be BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", be.Left.String(), be.Operator, be.Right.String())
}

func (be BinaryExpression) Evaluate(entry *parser.LogEntry) (Value, error) {
	left, err := be.Left.Evaluate(entry)
	if err != nil {
		return Value{}, err
	}

	right, err := be.Right.Evaluate(entry)
	if err != nil {
		return Value{}, err
	}

	return evaluateBinaryOperation(left, be.Operator, right)
}

// UnaryExpression represents unary operations
type UnaryExpression struct {
	Operator Operator
	Operand  Expression
}

func (ue UnaryExpression) String() string {
	return fmt.Sprintf("%s %s", ue.Operator, ue.Operand.String())
}

func (ue UnaryExpression) Evaluate(entry *parser.LogEntry) (Value, error) {
	operand, err := ue.Operand.Evaluate(entry)
	if err != nil {
		return Value{}, err
	}

	return evaluateUnaryOperation(ue.Operator, operand)
}

// FunctionExpression represents function calls
type FunctionExpression struct {
	Name      string
	Arguments []Expression
}

func (fe FunctionExpression) String() string {
	result := fe.Name + "("
	for i, arg := range fe.Arguments {
		if i > 0 {
			result += ", "
		}
		result += arg.String()
	}
	result += ")"
	return result
}

func (fe FunctionExpression) Evaluate(entry *parser.LogEntry) (Value, error) {
	// Evaluate arguments
	args := make([]Value, len(fe.Arguments))
	for i, arg := range fe.Arguments {
		val, err := arg.Evaluate(entry)
		if err != nil {
			return Value{}, err
		}
		args[i] = val
	}

	return evaluateFunction(fe.Name, args, entry)
}

// QueryResult represents the result of a query execution
type QueryResult struct {
	Columns []string
	Rows    [][]Value
	Count   int
}

// QueryError represents errors that occur during query processing
type QueryError struct {
	Message  string
	Position int
	Type     string
}

func (e QueryError) Error() string {
	return fmt.Sprintf("Query error at position %d: %s", e.Position, e.Message)
}

// NewQueryError creates a new query error
func NewQueryError(message string, position int, errorType string) *QueryError {
	return &QueryError{
		Message:  message,
		Position: position,
		Type:     errorType,
	}
}