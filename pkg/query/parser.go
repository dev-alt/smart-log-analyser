package query

import (
	"strconv"
	"strings"
	"time"
)

// Parser parses tokens into an Abstract Syntax Tree
type Parser struct {
	tokens   []Token
	current  int
	position int
}

// NewParser creates a new parser instance
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:   tokens,
		current:  0,
		position: 0,
	}
}

// ParseQuery parses a complete query
func (p *Parser) ParseQuery() (*SelectStatement, error) {
	// Validate tokens first
	if err := ValidateTokens(p.tokens); err != nil {
		return nil, err
	}

	return p.parseSelectStatement()
}

// parseSelectStatement parses a SELECT statement
func (p *Parser) parseSelectStatement() (*SelectStatement, error) {
	stmt := &SelectStatement{}

	// Parse SELECT clause
	if !p.expectToken(TokenSelect) {
		return nil, p.error("Expected SELECT")
	}
	p.advance()

	fields, err := p.parseSelectFields()
	if err != nil {
		return nil, err
	}
	stmt.Fields = fields

	// Parse FROM clause
	if !p.expectToken(TokenFrom) {
		return nil, p.error("Expected FROM")
	}
	p.advance()

	if !p.expectToken(TokenField) {
		return nil, p.error("Expected table name after FROM")
	}
	stmt.From = p.currentToken().Value
	p.advance()

	// Parse optional clauses
	for !p.isAtEnd() && p.currentToken().Type != TokenEOF {
		switch p.currentToken().Type {
		case TokenWhere:
			p.advance()
			where, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			stmt.Where = where

		case TokenGroup:
			p.advance()
			if !p.expectToken(TokenBy) {
				return nil, p.error("Expected BY after GROUP")
			}
			p.advance()
			groupBy, err := p.parseExpressionList()
			if err != nil {
				return nil, err
			}
			stmt.GroupBy = groupBy

		case TokenOrder:
			p.advance()
			if !p.expectToken(TokenBy) {
				return nil, p.error("Expected BY after ORDER")
			}
			p.advance()
			orderBy, err := p.parseOrderByClause()
			if err != nil {
				return nil, err
			}
			stmt.OrderBy = orderBy

		case TokenHaving:
			p.advance()
			having, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			stmt.Having = having

		case TokenLimit:
			p.advance()
			if !p.expectToken(TokenNumber) {
				return nil, p.error("Expected number after LIMIT")
			}
			limit, err := strconv.ParseInt(p.currentToken().Value, 10, 64)
			if err != nil {
				return nil, p.error("Invalid LIMIT value")
			}
			stmt.Limit = &limit
			p.advance()

		default:
			return nil, p.error("Unexpected token: " + p.currentToken().Value)
		}
	}

	return stmt, nil
}

// parseSelectFields parses the field list in SELECT clause
func (p *Parser) parseSelectFields() ([]SelectField, error) {
	var fields []SelectField

	// Handle SELECT *
	if p.currentToken().Value == "*" {
		fields = append(fields, SelectField{
			Expression: &FieldExpression{Field: "*"},
		})
		p.advance()
		return fields, nil
	}

	for {
		field, err := p.parseSelectField()
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)

		if p.currentToken().Type != TokenComma {
			break
		}
		p.advance() // Skip comma
	}

	return fields, nil
}

// parseSelectField parses a single field in SELECT clause
func (p *Parser) parseSelectField() (SelectField, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return SelectField{}, err
	}

	field := SelectField{Expression: expr}

	// Check for AS alias
	if p.currentToken().Type == TokenAs {
		p.advance()
		if !p.expectToken(TokenField) {
			return SelectField{}, p.error("Expected alias after AS")
		}
		field.Alias = p.currentToken().Value
		p.advance()
	}

	return field, nil
}

// parseOrderByClause parses ORDER BY clause
func (p *Parser) parseOrderByClause() ([]OrderByClause, error) {
	var clauses []OrderByClause

	for {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		clause := OrderByClause{Expression: expr}

		// Check for DESC
		if p.currentToken().Value == "DESC" {
			clause.Descending = true
			p.advance()
		} else if p.currentToken().Value == "ASC" {
			p.advance() // Skip ASC (default)
		}

		clauses = append(clauses, clause)

		if p.currentToken().Type != TokenComma {
			break
		}
		p.advance() // Skip comma
	}

	return clauses, nil
}

// parseExpressionList parses a comma-separated list of expressions
func (p *Parser) parseExpressionList() ([]Expression, error) {
	var expressions []Expression

	for {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)

		if p.currentToken().Type != TokenComma {
			break
		}
		p.advance() // Skip comma
	}

	return expressions, nil
}

// parseExpression parses a complete expression with precedence
func (p *Parser) parseExpression() (Expression, error) {
	return p.parseLogicalOr()
}

// parseLogicalOr parses OR expressions
func (p *Parser) parseLogicalOr() (Expression, error) {
	left, err := p.parseLogicalAnd()
	if err != nil {
		return nil, err
	}

	for p.currentToken().Type == TokenOr {
		op := OpOr
		p.advance()
		right, err := p.parseLogicalAnd()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Operator: op, Right: right}
	}

	return left, nil
}

// parseLogicalAnd parses AND expressions
func (p *Parser) parseLogicalAnd() (Expression, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.currentToken().Type == TokenAnd {
		op := OpAnd
		p.advance()
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Operator: op, Right: right}
	}

	return left, nil
}

// parseComparison parses comparison expressions
func (p *Parser) parseComparison() (Expression, error) {
	// Handle NOT operator
	if p.currentToken().Type == TokenNot {
		p.advance()
		operand, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		return &UnaryExpression{Operator: OpNot, Operand: operand}, nil
	}

	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	if IsComparisonOperator(p.currentToken().Type) {
		op := p.tokenToOperator(p.currentToken().Type)
		p.advance()

		// Handle special cases
		switch op {
		case OpBetween:
			return p.parseBetweenExpression(left)
		case OpIn:
			return p.parseInExpression(left)
		case OpIsBot, OpIsError, OpIsSuccess:
			// These are unary operators applied to fields
			return &UnaryExpression{Operator: op, Operand: left}, nil
		default:
			right, err := p.parsePrimary()
			if err != nil {
				return nil, err
			}
			return &BinaryExpression{Left: left, Operator: op, Right: right}, nil
		}
	}

	return left, nil
}

// parseBetweenExpression parses BETWEEN expressions
func (p *Parser) parseBetweenExpression(left Expression) (Expression, error) {
	min, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	if !p.expectToken(TokenAnd) {
		return nil, p.error("Expected AND in BETWEEN expression")
	}
	p.advance()

	max, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	// Create a compound expression: left >= min AND left <= max
	minExpr := &BinaryExpression{Left: left, Operator: OpGreaterThanOrEqual, Right: min}
	maxExpr := &BinaryExpression{Left: left, Operator: OpLessThanOrEqual, Right: max}
	return &BinaryExpression{Left: minExpr, Operator: OpAnd, Right: maxExpr}, nil
}

// parseInExpression parses IN expressions
func (p *Parser) parseInExpression(left Expression) (Expression, error) {
	if !p.expectToken(TokenLeftParen) {
		return nil, p.error("Expected '(' after IN")
	}
	p.advance()

	var values []Value
	for {
		value, err := p.parseLiteral()
		if err != nil {
			return nil, err
		}
		values = append(values, value)

		if p.currentToken().Type == TokenRightParen {
			break
		}
		if p.currentToken().Type != TokenComma {
			return nil, p.error("Expected ',' or ')' in IN list")
		}
		p.advance() // Skip comma
	}

	if !p.expectToken(TokenRightParen) {
		return nil, p.error("Expected ')' after IN list")
	}
	p.advance()

	listValue := Value{Type: ValueList, ListVal: values}
	return &BinaryExpression{
		Left:     left,
		Operator: OpIn,
		Right:    &LiteralExpression{Value: listValue},
	}, nil
}

// parsePrimary parses primary expressions (fields, literals, functions, parentheses)
func (p *Parser) parsePrimary() (Expression, error) {
	token := p.currentToken()

	switch token.Type {
	case TokenField:
		field := p.tokenToField(token.Value)
		p.advance()
		return &FieldExpression{Field: field}, nil

	case TokenFunction:
		return p.parseFunctionCall()

	case TokenString, TokenNumber, TokenBool, TokenDate:
		value, err := p.parseLiteral()
		if err != nil {
			return nil, err
		}
		p.advance()
		return &LiteralExpression{Value: value}, nil

	case TokenLeftParen:
		p.advance()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if !p.expectToken(TokenRightParen) {
			return nil, p.error("Expected ')'")
		}
		p.advance()
		return expr, nil

	default:
		return nil, p.error("Unexpected token in expression: " + token.Value)
	}
}

// parseFunctionCall parses function call expressions
func (p *Parser) parseFunctionCall() (Expression, error) {
	funcName := p.currentToken().Value
	p.advance()

	if !p.expectToken(TokenLeftParen) {
		return nil, p.error("Expected '(' after function name")
	}
	p.advance()

	var args []Expression
	if p.currentToken().Type != TokenRightParen {
		for {
			arg, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)

			if p.currentToken().Type == TokenRightParen {
				break
			}
			if p.currentToken().Type != TokenComma {
				return nil, p.error("Expected ',' or ')' in function arguments")
			}
			p.advance() // Skip comma
		}
	}

	if !p.expectToken(TokenRightParen) {
		return nil, p.error("Expected ')' after function arguments")
	}
	p.advance()

	return &FunctionExpression{Name: funcName, Arguments: args}, nil
}

// parseLiteral parses literal values
func (p *Parser) parseLiteral() (Value, error) {
	token := p.currentToken()

	switch token.Type {
	case TokenString:
		return Value{Type: ValueString, StringVal: token.Value}, nil

	case TokenNumber:
		if strings.Contains(token.Value, ".") {
			val, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return Value{}, p.error("Invalid number: " + token.Value)
			}
			return Value{Type: ValueFloat, FloatVal: val}, nil
		} else {
			val, err := strconv.ParseInt(token.Value, 10, 64)
			if err != nil {
				return Value{}, p.error("Invalid number: " + token.Value)
			}
			return Value{Type: ValueInt, IntVal: val}, nil
		}

	case TokenBool:
		val := strings.ToUpper(token.Value) == "TRUE"
		return Value{Type: ValueBool, BoolVal: val}, nil

	case TokenDate:
		// Try multiple date formats
		formats := []string{
			"2006-01-02 15:04:05",
			"2006-01-02",
			"15:04:05",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, token.Value); err == nil {
				return Value{Type: ValueTime, TimeVal: t}, nil
			}
		}
		return Value{}, p.error("Invalid date format: " + token.Value)

	default:
		return Value{}, p.error("Expected literal value")
	}
}

// Helper methods
func (p *Parser) currentToken() Token {
	if p.current >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.current]
}

func (p *Parser) advance() {
	if p.current < len(p.tokens) {
		p.current++
	}
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens) || p.currentToken().Type == TokenEOF
}

func (p *Parser) expectToken(tokenType TokenType) bool {
	return p.currentToken().Type == tokenType
}

func (p *Parser) error(message string) *QueryError {
	position := 0
	if p.current < len(p.tokens) {
		position = p.tokens[p.current].Position
	}
	return NewQueryError(message, position, "parser")
}

func (p *Parser) tokenToOperator(tokenType TokenType) Operator {
	mapping := map[TokenType]Operator{
		TokenEquals:              OpEquals,
		TokenNotEquals:           OpNotEquals,
		TokenLessThan:            OpLessThan,
		TokenLessThanOrEqual:     OpLessThanOrEqual,
		TokenGreaterThan:         OpGreaterThan,
		TokenGreaterThanOrEqual:  OpGreaterThanOrEqual,
		TokenLike:                OpLike,
		TokenMatches:             OpMatches,
		TokenContains:            OpContains,
		TokenStartsWith:          OpStartsWith,
		TokenEndsWith:            OpEndsWith,
		TokenIn:                  OpIn,
		TokenBetween:             OpBetween,
		TokenInRange:             OpInRange,
		TokenIsBot:               OpIsBot,
		TokenIsError:             OpIsError,
		TokenIsSuccess:           OpIsSuccess,
		TokenAnd:                 OpAnd,
		TokenOr:                  OpOr,
		TokenNot:                 OpNot,
	}
	return mapping[tokenType]
}

func (p *Parser) tokenToField(value string) QueryField {
	mapping := map[string]QueryField{
		"IP":         FieldIP,
		"TIMESTAMP":  FieldTimestamp,
		"METHOD":     FieldMethod,
		"URL":        FieldURL,
		"PROTOCOL":   FieldProtocol,
		"STATUS":     FieldStatus,
		"SIZE":       FieldSize,
		"REFERER":    FieldReferer,
		"USER_AGENT": FieldUserAgent,
		"*":          "*", // Special case for SELECT *
	}
	
	if field, ok := mapping[strings.ToUpper(value)]; ok {
		return field
	}
	return QueryField(value) // Return as-is for unknown fields
}

// ParseQuery is a convenience function to parse a query string
func ParseQuery(query string) (*SelectStatement, error) {
	tokens, err := TokenizeQuery(query)
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	return parser.ParseQuery()
}