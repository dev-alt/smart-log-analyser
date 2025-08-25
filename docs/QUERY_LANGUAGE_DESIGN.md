# Advanced Query Language Design

## Overview

The Smart Log Analyser Query Language (SLAQ) provides powerful filtering and analysis capabilities for Nginx access logs. It allows users to construct complex queries using SQL-like syntax with log-specific operators and functions.

## Syntax Overview

### Basic Query Structure
```sql
SELECT [fields] FROM logs WHERE [conditions] [GROUP BY field] [ORDER BY field] [LIMIT number]
```

### Simple Filtering
```sql
-- Filter by status code
WHERE status = 404

-- Filter by IP address
WHERE ip = "192.168.1.1"

-- Filter by URL pattern
WHERE url LIKE "/api/*"
```

### Complex Filtering
```sql
-- Multiple conditions
WHERE status >= 400 AND size > 10000

-- Time range filtering
WHERE timestamp BETWEEN "2024-08-20 00:00:00" AND "2024-08-20 23:59:59"

-- Regular expressions
WHERE user_agent MATCHES ".*bot.*"

-- IP range filtering
WHERE ip IN_RANGE "192.168.1.0/24"
```

## Available Fields

| Field | Type | Description | Examples |
|-------|------|-------------|----------|
| `ip` | string | Client IP address | `"192.168.1.1"`, `"203.0.113.42"` |
| `timestamp` | datetime | Request timestamp | `"2024-08-20 15:30:45"` |
| `method` | string | HTTP method | `"GET"`, `"POST"`, `"PUT"` |
| `url` | string | Request URL/path | `"/index.html"`, `"/api/users"` |
| `protocol` | string | HTTP protocol | `"HTTP/1.1"`, `"HTTP/2.0"` |
| `status` | int | HTTP status code | `200`, `404`, `500` |
| `size` | int | Response size in bytes | `1024`, `4096` |
| `referer` | string | HTTP referer header | `"https://google.com"` |
| `user_agent` | string | User agent string | `"Mozilla/5.0..."` |

## Operators

### Comparison Operators
- `=` - Equals
- `!=` or `<>` - Not equals
- `<` - Less than
- `<=` - Less than or equal
- `>` - Greater than
- `>=` - Greater than or equal

### Logical Operators
- `AND` - Logical AND
- `OR` - Logical OR
- `NOT` - Logical NOT

### String Operators
- `LIKE` - Pattern matching with wildcards (`*` for multiple chars, `?` for single char)
- `MATCHES` - Regular expression matching
- `CONTAINS` - String contains substring
- `STARTS_WITH` - String starts with prefix
- `ENDS_WITH` - String ends with suffix

### Special Operators
- `IN` - Value in list: `status IN (200, 201, 202)`
- `BETWEEN` - Value in range: `size BETWEEN 1000 AND 5000`
- `IN_RANGE` - IP in CIDR range: `ip IN_RANGE "192.168.0.0/16"`
- `IS_BOT` - User agent is a bot/crawler
- `IS_ERROR` - Status code is 4xx or 5xx
- `IS_SUCCESS` - Status code is 2xx

## Functions

### Time Functions
- `HOUR(timestamp)` - Extract hour (0-23)
- `DAY(timestamp)` - Extract day of month
- `WEEKDAY(timestamp)` - Extract weekday (0=Sunday)
- `DATE(timestamp)` - Extract date part
- `TIME_DIFF(timestamp1, timestamp2)` - Time difference in seconds

### String Functions
- `UPPER(field)` - Convert to uppercase
- `LOWER(field)` - Convert to lowercase
- `LENGTH(field)` - String length
- `SUBSTR(field, start, length)` - Substring extraction

### Network Functions
- `IP_TO_INT(ip)` - Convert IP to integer
- `IS_PRIVATE_IP(ip)` - Check if IP is private
- `COUNTRY(ip)` - Get country from IP (requires GeoIP)

### Aggregate Functions (with GROUP BY)
- `COUNT()` - Count records
- `SUM(field)` - Sum values
- `AVG(field)` - Average value
- `MIN(field)` - Minimum value
- `MAX(field)` - Maximum value

## Advanced Examples

### Security Analysis
```sql
-- Find potential brute force attacks
SELECT ip, COUNT() as attempts 
FROM logs 
WHERE status = 401 AND url LIKE "/login*"
GROUP BY ip 
HAVING attempts > 10
ORDER BY attempts DESC

-- Identify suspicious scanning activity
WHERE status = 404 AND url MATCHES ".*\\.(php|asp|jsp)$"

-- Find large file access attempts
WHERE size > 50000000 AND status = 200
```

### Performance Analysis
```sql
-- Slow endpoints (by response size as proxy)
SELECT url, AVG(size) as avg_size, COUNT() as requests
FROM logs 
WHERE status = 200
GROUP BY url
HAVING requests > 100
ORDER BY avg_size DESC
LIMIT 10

-- Peak hour analysis
SELECT HOUR(timestamp) as hour, COUNT() as requests
FROM logs
GROUP BY hour
ORDER BY requests DESC
```

### Traffic Analysis
```sql
-- Bot vs human traffic
SELECT CASE 
  WHEN IS_BOT(user_agent) THEN 'Bot' 
  ELSE 'Human' 
END as traffic_type, COUNT() as requests
FROM logs
GROUP BY traffic_type

-- Geographic distribution
SELECT COUNTRY(ip) as country, COUNT() as requests
FROM logs
WHERE NOT IS_PRIVATE_IP(ip)
GROUP BY country
ORDER BY requests DESC
LIMIT 20
```

### Error Analysis
```sql
-- Error rate by endpoint
SELECT url, 
       COUNT() as total_requests,
       SUM(CASE WHEN IS_ERROR(status) THEN 1 ELSE 0 END) as errors,
       ROUND(100.0 * SUM(CASE WHEN IS_ERROR(status) THEN 1 ELSE 0 END) / COUNT(), 2) as error_rate
FROM logs
GROUP BY url
HAVING total_requests > 50 AND error_rate > 5
ORDER BY error_rate DESC
```

## Query Optimization

### Performance Tips
1. **Index-friendly conditions**: Use equality and range conditions on commonly filtered fields
2. **Limit result sets**: Use `LIMIT` to prevent memory issues
3. **Time range filtering**: Always specify time ranges for large datasets
4. **Field selection**: Use `SELECT` to limit returned fields

### Best Practices
1. **Combine conditions efficiently**: Put most selective conditions first
2. **Use appropriate operators**: `LIKE` is slower than `=` or `STARTS_WITH`
3. **Group by sparingly**: GROUP BY operations are memory intensive
4. **Regular expressions**: Use `MATCHES` only when necessary

## Integration Points

### CLI Integration
```bash
# Query via command line
./smart-log-analyser analyse --query "WHERE status >= 400" logs/*.log

# Save query results
./smart-log-analyser analyse --query "SELECT ip, COUNT() FROM logs GROUP BY ip" --export-csv results.csv
```

### Menu Integration
```
📊 Advanced Query Mode:
1. Guided query builder
2. Manual query entry
3. Saved query templates
4. Query history
```

### API Integration
```json
{
  "query": "SELECT url, COUNT() as hits FROM logs WHERE status = 404 GROUP BY url ORDER BY hits DESC LIMIT 10",
  "parameters": {
    "time_range": "last_24h",
    "format": "json"
  }
}
```

## Implementation Architecture

### Components
1. **Lexer**: Tokenizes query string into tokens
2. **Parser**: Builds Abstract Syntax Tree (AST) from tokens  
3. **Validator**: Validates query syntax and semantics
4. **Executor**: Executes query against log entries
5. **Formatter**: Formats results in requested output format

### Query Execution Pipeline
```
Raw Query → Lexer → Parser → Validator → Executor → Formatter → Results
```

### Error Handling
- **Syntax Errors**: Clear error messages with position information
- **Semantic Errors**: Field validation and type checking
- **Runtime Errors**: Graceful handling of malformed data
- **Performance Warnings**: Alerts for potentially slow queries

## Future Enhancements

### Advanced Features
- **Subqueries**: Nested query support
- **Joins**: Cross-log file correlation
- **Window Functions**: Advanced analytics
- **Saved Views**: Reusable query templates
- **Real-time Queries**: Streaming log analysis

### UI Enhancements
- **Visual Query Builder**: Drag-and-drop interface
- **Query Autocomplete**: Intelligent suggestions
- **Result Visualization**: Charts and graphs
- **Query Performance Metrics**: Execution statistics