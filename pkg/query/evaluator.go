package query

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"smart-log-analyser/pkg/parser"
)

// evaluateBinaryOperation performs binary operations
func evaluateBinaryOperation(left Value, op Operator, right Value) (Value, error) {
	switch op {
	case OpEquals:
		return compareValues(left, right, func(cmp int) bool { return cmp == 0 })
	case OpNotEquals:
		return compareValues(left, right, func(cmp int) bool { return cmp != 0 })
	case OpLessThan:
		return compareValues(left, right, func(cmp int) bool { return cmp < 0 })
	case OpLessThanOrEqual:
		return compareValues(left, right, func(cmp int) bool { return cmp <= 0 })
	case OpGreaterThan:
		return compareValues(left, right, func(cmp int) bool { return cmp > 0 })
	case OpGreaterThanOrEqual:
		return compareValues(left, right, func(cmp int) bool { return cmp >= 0 })
	case OpLike:
		return evaluateLike(left, right)
	case OpMatches:
		return evaluateMatches(left, right)
	case OpContains:
		return evaluateContains(left, right)
	case OpStartsWith:
		return evaluateStartsWith(left, right)
	case OpEndsWith:
		return evaluateEndsWith(left, right)
	case OpIn:
		return evaluateIn(left, right)
	case OpInRange:
		return evaluateInRange(left, right)
	case OpAnd:
		return evaluateLogicalAnd(left, right)
	case OpOr:
		return evaluateLogicalOr(left, right)
	default:
		return Value{}, fmt.Errorf("unsupported binary operator: %s", op)
	}
}

// evaluateUnaryOperation performs unary operations
func evaluateUnaryOperation(op Operator, operand Value) (Value, error) {
	switch op {
	case OpNot:
		return evaluateLogicalNot(operand)
	case OpIsBot:
		return evaluateIsBot(operand)
	case OpIsError:
		return evaluateIsError(operand)
	case OpIsSuccess:
		return evaluateIsSuccess(operand)
	default:
		return Value{}, fmt.Errorf("unsupported unary operator: %s", op)
	}
}

// compareValues compares two values and applies a comparison function
func compareValues(left, right Value, cmp func(int) bool) (Value, error) {
	// Type coercion if needed
	if left.Type != right.Type {
		var err error
		left, right, err = coerceValues(left, right)
		if err != nil {
			return Value{}, err
		}
	}

	var result int
	switch left.Type {
	case ValueString:
		result = strings.Compare(left.StringVal, right.StringVal)
	case ValueInt:
		if left.IntVal == right.IntVal {
			result = 0
		} else if left.IntVal < right.IntVal {
			result = -1
		} else {
			result = 1
		}
	case ValueFloat:
		if left.FloatVal == right.FloatVal {
			result = 0
		} else if left.FloatVal < right.FloatVal {
			result = -1
		} else {
			result = 1
		}
	case ValueTime:
		if left.TimeVal.Equal(right.TimeVal) {
			result = 0
		} else if left.TimeVal.Before(right.TimeVal) {
			result = -1
		} else {
			result = 1
		}
	case ValueBool:
		if left.BoolVal == right.BoolVal {
			result = 0
		} else if !left.BoolVal && right.BoolVal {
			result = -1
		} else {
			result = 1
		}
	default:
		return Value{}, fmt.Errorf("cannot compare values of type %d", left.Type)
	}

	return Value{Type: ValueBool, BoolVal: cmp(result)}, nil
}

// coerceValues attempts to convert values to the same type
func coerceValues(left, right Value) (Value, Value, error) {
	// String to number coercion
	if left.Type == ValueString && (right.Type == ValueInt || right.Type == ValueFloat) {
		if val, err := strconv.ParseInt(left.StringVal, 10, 64); err == nil {
			left = Value{Type: ValueInt, IntVal: val}
		} else if val, err := strconv.ParseFloat(left.StringVal, 64); err == nil {
			left = Value{Type: ValueFloat, FloatVal: val}
		}
	}
	if right.Type == ValueString && (left.Type == ValueInt || left.Type == ValueFloat) {
		if val, err := strconv.ParseInt(right.StringVal, 10, 64); err == nil {
			right = Value{Type: ValueInt, IntVal: val}
		} else if val, err := strconv.ParseFloat(right.StringVal, 64); err == nil {
			right = Value{Type: ValueFloat, FloatVal: val}
		}
	}

	// Int to float coercion
	if left.Type == ValueInt && right.Type == ValueFloat {
		left = Value{Type: ValueFloat, FloatVal: float64(left.IntVal)}
	}
	if right.Type == ValueInt && left.Type == ValueFloat {
		right = Value{Type: ValueFloat, FloatVal: float64(right.IntVal)}
	}

	return left, right, nil
}

// evaluateLike implements LIKE pattern matching
func evaluateLike(left, right Value) (Value, error) {
	if left.Type != ValueString || right.Type != ValueString {
		return Value{}, fmt.Errorf("LIKE operator requires string operands")
	}

	pattern := right.StringVal
	// Convert SQL LIKE pattern to Go regexp
	pattern = strings.ReplaceAll(pattern, "*", ".*")
	pattern = strings.ReplaceAll(pattern, "?", ".")
	pattern = "^" + pattern + "$"

	matched, err := regexp.MatchString(pattern, left.StringVal)
	if err != nil {
		return Value{}, fmt.Errorf("invalid LIKE pattern: %v", err)
	}

	return Value{Type: ValueBool, BoolVal: matched}, nil
}

// evaluateMatches implements regular expression matching
func evaluateMatches(left, right Value) (Value, error) {
	if left.Type != ValueString || right.Type != ValueString {
		return Value{}, fmt.Errorf("MATCHES operator requires string operands")
	}

	matched, err := regexp.MatchString(right.StringVal, left.StringVal)
	if err != nil {
		return Value{}, fmt.Errorf("invalid regex pattern: %v", err)
	}

	return Value{Type: ValueBool, BoolVal: matched}, nil
}

// evaluateContains implements string contains checking
func evaluateContains(left, right Value) (Value, error) {
	if left.Type != ValueString || right.Type != ValueString {
		return Value{}, fmt.Errorf("CONTAINS operator requires string operands")
	}

	contains := strings.Contains(left.StringVal, right.StringVal)
	return Value{Type: ValueBool, BoolVal: contains}, nil
}

// evaluateStartsWith implements string prefix checking
func evaluateStartsWith(left, right Value) (Value, error) {
	if left.Type != ValueString || right.Type != ValueString {
		return Value{}, fmt.Errorf("STARTS_WITH operator requires string operands")
	}

	startsWith := strings.HasPrefix(left.StringVal, right.StringVal)
	return Value{Type: ValueBool, BoolVal: startsWith}, nil
}

// evaluateEndsWith implements string suffix checking
func evaluateEndsWith(left, right Value) (Value, error) {
	if left.Type != ValueString || right.Type != ValueString {
		return Value{}, fmt.Errorf("ENDS_WITH operator requires string operands")
	}

	endsWith := strings.HasSuffix(left.StringVal, right.StringVal)
	return Value{Type: ValueBool, BoolVal: endsWith}, nil
}

// evaluateIn implements IN list checking
func evaluateIn(left, right Value) (Value, error) {
	if right.Type != ValueList {
		return Value{}, fmt.Errorf("IN operator requires a list on the right side")
	}

	for _, val := range right.ListVal {
		result, err := compareValues(left, val, func(cmp int) bool { return cmp == 0 })
		if err != nil {
			continue // Skip incompatible types
		}
		if result.BoolVal {
			return Value{Type: ValueBool, BoolVal: true}, nil
		}
	}

	return Value{Type: ValueBool, BoolVal: false}, nil
}

// evaluateInRange implements IP range checking
func evaluateInRange(left, right Value) (Value, error) {
	if left.Type != ValueString || right.Type != ValueString {
		return Value{}, fmt.Errorf("IN_RANGE operator requires string operands (IP and CIDR)")
	}

	ip := net.ParseIP(left.StringVal)
	if ip == nil {
		return Value{}, fmt.Errorf("invalid IP address: %s", left.StringVal)
	}

	_, cidr, err := net.ParseCIDR(right.StringVal)
	if err != nil {
		return Value{}, fmt.Errorf("invalid CIDR range: %s", right.StringVal)
	}

	inRange := cidr.Contains(ip)
	return Value{Type: ValueBool, BoolVal: inRange}, nil
}

// evaluateLogicalAnd implements logical AND
func evaluateLogicalAnd(left, right Value) (Value, error) {
	leftBool, err := toBool(left)
	if err != nil {
		return Value{}, err
	}
	rightBool, err := toBool(right)
	if err != nil {
		return Value{}, err
	}

	return Value{Type: ValueBool, BoolVal: leftBool && rightBool}, nil
}

// evaluateLogicalOr implements logical OR
func evaluateLogicalOr(left, right Value) (Value, error) {
	leftBool, err := toBool(left)
	if err != nil {
		return Value{}, err
	}
	rightBool, err := toBool(right)
	if err != nil {
		return Value{}, err
	}

	return Value{Type: ValueBool, BoolVal: leftBool || rightBool}, nil
}

// evaluateLogicalNot implements logical NOT
func evaluateLogicalNot(operand Value) (Value, error) {
	boolVal, err := toBool(operand)
	if err != nil {
		return Value{}, err
	}

	return Value{Type: ValueBool, BoolVal: !boolVal}, nil
}

// evaluateIsBot checks if user agent indicates a bot
func evaluateIsBot(operand Value) (Value, error) {
	if operand.Type != ValueString {
		return Value{}, fmt.Errorf("IS_BOT operator requires string operand")
	}

	userAgent := strings.ToLower(operand.StringVal)
	botPatterns := []string{
		"bot", "crawler", "spider", "scraper", "crawl",
		"googlebot", "bingbot", "slurp", "facebookexternalhit",
		"twitterbot", "whatsapp", "telegram", "curl", "wget",
		"postman", "httpie", "python-requests", "monitoring",
	}

	isBot := false
	for _, pattern := range botPatterns {
		if strings.Contains(userAgent, pattern) {
			isBot = true
			break
		}
	}

	return Value{Type: ValueBool, BoolVal: isBot}, nil
}

// evaluateIsError checks if status code is an error (4xx or 5xx)
func evaluateIsError(operand Value) (Value, error) {
	if operand.Type != ValueInt {
		return Value{}, fmt.Errorf("IS_ERROR operator requires integer operand")
	}

	status := operand.IntVal
	isError := status >= 400 && status <= 599
	return Value{Type: ValueBool, BoolVal: isError}, nil
}

// evaluateIsSuccess checks if status code is success (2xx)
func evaluateIsSuccess(operand Value) (Value, error) {
	if operand.Type != ValueInt {
		return Value{}, fmt.Errorf("IS_SUCCESS operator requires integer operand")
	}

	status := operand.IntVal
	isSuccess := status >= 200 && status <= 299
	return Value{Type: ValueBool, BoolVal: isSuccess}, nil
}

// evaluateFunction evaluates function calls
func evaluateFunction(name string, args []Value, entry *parser.LogEntry) (Value, error) {
	switch strings.ToUpper(name) {
	case "COUNT":
		return Value{Type: ValueInt, IntVal: 1}, nil // Will be aggregated later

	case "SUM":
		if len(args) != 1 {
			return Value{}, fmt.Errorf("SUM function requires exactly 1 argument")
		}
		return args[0], nil // Will be aggregated later

	case "AVG":
		if len(args) != 1 {
			return Value{}, fmt.Errorf("AVG function requires exactly 1 argument")
		}
		return args[0], nil // Will be aggregated later

	case "MIN":
		if len(args) != 1 {
			return Value{}, fmt.Errorf("MIN function requires exactly 1 argument")
		}
		return args[0], nil // Will be aggregated later

	case "MAX":
		if len(args) != 1 {
			return Value{}, fmt.Errorf("MAX function requires exactly 1 argument")
		}
		return args[0], nil // Will be aggregated later

	case "HOUR":
		if len(args) != 1 {
			return Value{}, fmt.Errorf("HOUR function requires exactly 1 argument")
		}
		if args[0].Type != ValueTime {
			return Value{}, fmt.Errorf("HOUR function requires time argument")
		}
		hour := args[0].TimeVal.Hour()
		return Value{Type: ValueInt, IntVal: int64(hour)}, nil

	case "DAY":
		if len(args) != 1 {
			return Value{}, fmt.Errorf("DAY function requires exactly 1 argument")
		}
		if args[0].Type != ValueTime {
			return Value{}, fmt.Errorf("DAY function requires time argument")
		}
		day := args[0].TimeVal.Day()
		return Value{Type: ValueInt, IntVal: int64(day)}, nil

	case "WEEKDAY":
		if len(args) != 1 {
			return Value{}, fmt.Errorf("WEEKDAY function requires exactly 1 argument")
		}
		if args[0].Type != ValueTime {
			return Value{}, fmt.Errorf("WEEKDAY function requires time argument")
		}
		weekday := int(args[0].TimeVal.Weekday())
		return Value{Type: ValueInt, IntVal: int64(weekday)}, nil

	case "DATE":
		if len(args) != 1 {
			return Value{}, fmt.Errorf("DATE function requires exactly 1 argument")
		}
		if args[0].Type != ValueTime {
			return Value{}, fmt.Errorf("DATE function requires time argument")
		}
		dateStr := args[0].TimeVal.Format("2006-01-02")
		return Value{Type: ValueString, StringVal: dateStr}, nil

	case "UPPER":
		if len(args) != 1 {
			return Value{}, fmt.Errorf("UPPER function requires exactly 1 argument")
		}
		if args[0].Type != ValueString {
			return Value{}, fmt.Errorf("UPPER function requires string argument")
		}
		return Value{Type: ValueString, StringVal: strings.ToUpper(args[0].StringVal)}, nil

	case "LOWER":
		if len(args) != 1 {
			return Value{}, fmt.Errorf("LOWER function requires exactly 1 argument")
		}
		if args[0].Type != ValueString {
			return Value{}, fmt.Errorf("LOWER function requires string argument")
		}
		return Value{Type: ValueString, StringVal: strings.ToLower(args[0].StringVal)}, nil

	case "LENGTH":
		if len(args) != 1 {
			return Value{}, fmt.Errorf("LENGTH function requires exactly 1 argument")
		}
		if args[0].Type != ValueString {
			return Value{}, fmt.Errorf("LENGTH function requires string argument")
		}
		return Value{Type: ValueInt, IntVal: int64(len(args[0].StringVal))}, nil

	case "SUBSTR":
		if len(args) < 2 || len(args) > 3 {
			return Value{}, fmt.Errorf("SUBSTR function requires 2 or 3 arguments")
		}
		if args[0].Type != ValueString || args[1].Type != ValueInt {
			return Value{}, fmt.Errorf("SUBSTR function requires string and integer arguments")
		}
		str := args[0].StringVal
		start := int(args[1].IntVal)
		if start < 0 || start >= len(str) {
			return Value{Type: ValueString, StringVal: ""}, nil
		}
		end := len(str)
		if len(args) == 3 {
			if args[2].Type != ValueInt {
				return Value{}, fmt.Errorf("SUBSTR length must be integer")
			}
			length := int(args[2].IntVal)
			if start+length < end {
				end = start + length
			}
		}
		return Value{Type: ValueString, StringVal: str[start:end]}, nil

	case "IS_PRIVATE_IP":
		if len(args) != 1 {
			return Value{}, fmt.Errorf("IS_PRIVATE_IP function requires exactly 1 argument")
		}
		if args[0].Type != ValueString {
			return Value{}, fmt.Errorf("IS_PRIVATE_IP function requires string argument")
		}
		ip := net.ParseIP(args[0].StringVal)
		if ip == nil {
			return Value{Type: ValueBool, BoolVal: false}, nil
		}
		isPrivate := isPrivateIP(ip)
		return Value{Type: ValueBool, BoolVal: isPrivate}, nil

	case "COUNTRY":
		if len(args) != 1 {
			return Value{}, fmt.Errorf("COUNTRY function requires exactly 1 argument")
		}
		if args[0].Type != ValueString {
			return Value{}, fmt.Errorf("COUNTRY function requires string argument")
		}
		// Simplified country detection (would need GeoIP database for real implementation)
		country := detectCountryFromIP(args[0].StringVal)
		return Value{Type: ValueString, StringVal: country}, nil

	default:
		return Value{}, fmt.Errorf("unknown function: %s", name)
	}
}

// toBool converts a value to boolean
func toBool(value Value) (bool, error) {
	switch value.Type {
	case ValueBool:
		return value.BoolVal, nil
	case ValueInt:
		return value.IntVal != 0, nil
	case ValueFloat:
		return value.FloatVal != 0.0, nil
	case ValueString:
		return value.StringVal != "", nil
	default:
		return false, fmt.Errorf("cannot convert value to boolean")
	}
}

// isPrivateIP checks if an IP address is private
func isPrivateIP(ip net.IP) bool {
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}

	for _, rangeStr := range privateRanges {
		_, cidr, err := net.ParseCIDR(rangeStr)
		if err != nil {
			continue
		}
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

// detectCountryFromIP provides basic country detection
func detectCountryFromIP(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "Unknown"
	}

	if isPrivateIP(ip) {
		return "Private"
	}

	// Simple heuristic based on IP ranges (would use GeoIP in real implementation)
	if ip.To4() != nil {
		// Very basic detection for demonstration
		parts := strings.Split(ipStr, ".")
		if len(parts) == 4 {
			first, _ := strconv.Atoi(parts[0])
			switch {
			case first >= 1 && first <= 126:
				return "US/International"
			case first >= 128 && first <= 191:
				return "International"
			case first >= 192 && first <= 223:
				return "International"
			}
		}
	}

	return "Unknown"
}