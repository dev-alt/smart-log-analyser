package query

import (
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Lexer tokenizes query strings
type Lexer struct {
	input    string
	position int
	current  rune
}

// NewLexer creates a new lexer instance
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:    input,
		position: 0,
	}
	l.readChar()
	return l
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	var token Token
	token.Position = l.position

	switch l.current {
	case 0:
		token.Type = TokenEOF
		token.Value = ""
	case '(':
		token.Type = TokenLeftParen
		token.Value = "("
	case ')':
		token.Type = TokenRightParen
		token.Value = ")"
	case ',':
		token.Type = TokenComma
		token.Value = ","
	case ';':
		token.Type = TokenSemicolon
		token.Value = ";"
	case '=':
		token.Type = TokenEquals
		token.Value = "="
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			token.Type = TokenLessThanOrEqual
			token.Value = "<="
		} else if l.peekChar() == '>' {
			l.readChar()
			token.Type = TokenNotEquals
			token.Value = "<>"
		} else {
			token.Type = TokenLessThan
			token.Value = "<"
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			token.Type = TokenGreaterThanOrEqual
			token.Value = ">="
		} else {
			token.Type = TokenGreaterThan
			token.Value = ">"
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			token.Type = TokenNotEquals
			token.Value = "!="
		} else {
			token.Type = TokenInvalid
			token.Value = "!"
		}
	case '"', '\'':
		quote := l.current
		token.Value = l.readString(quote)
		token.Type = l.determineStringTokenType(token.Value)
	default:
		if unicode.IsLetter(l.current) || l.current == '_' {
			token.Value = l.readIdentifier()
			token.Type = l.determineKeywordTokenType(token.Value)
		} else if unicode.IsDigit(l.current) {
			token.Value = l.readNumber()
			token.Type = TokenNumber
		} else {
			token.Type = TokenInvalid
			token.Value = string(l.current)
		}
	}

	l.readChar()
	return token
}

// readChar advances to the next character
func (l *Lexer) readChar() {
	if l.position >= len(l.input) {
		l.current = 0 // ASCII NUL character represents EOF
	} else {
		l.current = rune(l.input[l.position])
	}
	l.position++
}

// peekChar returns the next character without advancing
func (l *Lexer) peekChar() rune {
	if l.position >= len(l.input) {
		return 0
	}
	return rune(l.input[l.position])
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.current == ' ' || l.current == '\t' || l.current == '\n' || l.current == '\r' {
		l.readChar()
	}
}

// readIdentifier reads an identifier (field name, function name, keyword)
func (l *Lexer) readIdentifier() string {
	position := l.position - 1
	for unicode.IsLetter(l.current) || unicode.IsDigit(l.current) || l.current == '_' {
		l.readChar()
	}
	l.position-- // Step back one character
	return l.input[position : l.position]
}

// readNumber reads a numeric literal
func (l *Lexer) readNumber() string {
	position := l.position - 1
	for unicode.IsDigit(l.current) || l.current == '.' {
		l.readChar()
	}
	l.position-- // Step back one character
	return l.input[position : l.position]
}

// readString reads a string literal
func (l *Lexer) readString(quote rune) string {
	position := l.position
	for {
		l.readChar()
		if l.current == quote || l.current == 0 {
			break
		}
	}
	return l.input[position : l.position-1]
}

// determineKeywordTokenType determines if an identifier is a keyword
func (l *Lexer) determineKeywordTokenType(literal string) TokenType {
	keywords := map[string]TokenType{
		"SELECT":      TokenSelect,
		"FROM":        TokenFrom,
		"WHERE":       TokenWhere,
		"GROUP":       TokenGroup,
		"BY":          TokenBy,
		"ORDER":       TokenOrder,
		"HAVING":      TokenHaving,
		"LIMIT":       TokenLimit,
		"AS":          TokenAs,
		"AND":         TokenAnd,
		"OR":          TokenOr,
		"NOT":         TokenNot,
		"LIKE":        TokenLike,
		"MATCHES":     TokenMatches,
		"CONTAINS":    TokenContains,
		"STARTS_WITH": TokenStartsWith,
		"ENDS_WITH":   TokenEndsWith,
		"IN":          TokenIn,
		"BETWEEN":     TokenBetween,
		"IN_RANGE":    TokenInRange,
		"IS_BOT":      TokenIsBot,
		"IS_ERROR":    TokenIsError,
		"IS_SUCCESS":  TokenIsSuccess,
	}

	// Handle compound keywords
	upper := strings.ToUpper(literal)
	if tokenType, ok := keywords[upper]; ok {
		return tokenType
	}

	// Check if it's a field name
	fields := map[string]bool{
		"IP":         true,
		"TIMESTAMP":  true,
		"METHOD":     true,
		"URL":        true,
		"PROTOCOL":   true,
		"STATUS":     true,
		"SIZE":       true,
		"REFERER":    true,
		"USER_AGENT": true,
	}

	if _, ok := fields[upper]; ok {
		return TokenField
	}

	// Check if it's a function
	functions := map[string]bool{
		"COUNT":         true,
		"SUM":           true,
		"AVG":           true,
		"MIN":           true,
		"MAX":           true,
		"HOUR":          true,
		"DAY":           true,
		"WEEKDAY":       true,
		"DATE":          true,
		"TIME_DIFF":     true,
		"UPPER":         true,
		"LOWER":         true,
		"LENGTH":        true,
		"SUBSTR":        true,
		"IP_TO_INT":     true,
		"IS_PRIVATE_IP": true,
		"COUNTRY":       true,
	}

	if _, ok := functions[upper]; ok {
		return TokenFunction
	}

	return TokenField // Default to field for unknown identifiers
}

// determineStringTokenType determines the type of a string token
func (l *Lexer) determineStringTokenType(literal string) TokenType {
	// Check if it's a boolean
	upper := strings.ToUpper(literal)
	if upper == "TRUE" || upper == "FALSE" {
		return TokenBool
	}

	// Check if it's a date/time
	dateFormats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		"15:04:05",
		"2006/01/02 15:04:05",
		"2006/01/02",
	}

	for _, format := range dateFormats {
		if _, err := time.Parse(format, literal); err == nil {
			return TokenDate
		}
	}

	// Check if it's a number
	if _, err := strconv.ParseInt(literal, 10, 64); err == nil {
		return TokenNumber
	}
	if _, err := strconv.ParseFloat(literal, 64); err == nil {
		return TokenNumber
	}

	return TokenString
}

// TokenizeQuery tokenizes a complete query string
func TokenizeQuery(query string) ([]Token, error) {
	lexer := NewLexer(query)
	var tokens []Token

	for {
		token := lexer.NextToken()
		tokens = append(tokens, token)

		if token.Type == TokenEOF {
			break
		}

		if token.Type == TokenInvalid {
			return nil, NewQueryError("Invalid token: "+token.Value, token.Position, "lexer")
		}
	}

	return tokens, nil
}

// ValidateTokens performs basic token validation
func ValidateTokens(tokens []Token) error {
	if len(tokens) == 0 {
		return NewQueryError("Empty query", 0, "validation")
	}

	// Check for balanced parentheses
	parenCount := 0
	for _, token := range tokens {
		if token.Type == TokenLeftParen {
			parenCount++
		} else if token.Type == TokenRightParen {
			parenCount--
			if parenCount < 0 {
				return NewQueryError("Unmatched closing parenthesis", token.Position, "validation")
			}
		}
	}

	if parenCount > 0 {
		return NewQueryError("Unmatched opening parenthesis", tokens[len(tokens)-1].Position, "validation")
	}

	return nil
}

// TokensToString converts tokens back to string (for debugging)
func TokensToString(tokens []Token) string {
	var parts []string
	for _, token := range tokens {
		switch token.Type {
		case TokenString, TokenDate:
			parts = append(parts, "'"+token.Value+"'")
		default:
			parts = append(parts, token.Value)
		}
	}
	return strings.Join(parts, " ")
}

// IsComparisonOperator checks if a token is a comparison operator
func IsComparisonOperator(tokenType TokenType) bool {
	switch tokenType {
	case TokenEquals, TokenNotEquals, TokenLessThan, TokenLessThanOrEqual,
		TokenGreaterThan, TokenGreaterThanOrEqual, TokenLike, TokenMatches,
		TokenContains, TokenStartsWith, TokenEndsWith, TokenIn, TokenBetween,
		TokenInRange, TokenIsBot, TokenIsError, TokenIsSuccess:
		return true
	}
	return false
}

// IsLogicalOperator checks if a token is a logical operator
func IsLogicalOperator(tokenType TokenType) bool {
	switch tokenType {
	case TokenAnd, TokenOr, TokenNot:
		return true
	}
	return false
}

// GetOperatorPrecedence returns the precedence of logical operators
func GetOperatorPrecedence(tokenType TokenType) int {
	switch tokenType {
	case TokenOr:
		return 1
	case TokenAnd:
		return 2
	case TokenNot:
		return 3
	}
	return 0
}