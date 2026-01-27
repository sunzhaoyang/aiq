package db

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// QueryResult represents a query result
type QueryResult struct {
	Columns []string
	Rows    [][]string
}

// ExecuteQuery executes a SQL query and returns the results
func (c *Connection) ExecuteQuery(ctx context.Context, sqlQuery string) (*QueryResult, error) {
	// Set timeout
	queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	rows, err := c.db.QueryContext(queryCtx, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	// Read rows
	result := &QueryResult{
		Columns: columns,
		Rows:    make([][]string, 0),
	}

	for rows.Next() {
		// Create slice to hold column values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert values to strings
		row := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				row[i] = "NULL"
			} else {
				// Handle different types properly
				switch v := val.(type) {
				case []byte:
					// Convert byte slice to string
					row[i] = string(v)
				case string:
					row[i] = v
				case int64:
					row[i] = fmt.Sprintf("%d", v)
				case float64:
					row[i] = fmt.Sprintf("%g", v)
				case bool:
					if v {
						row[i] = "true"
					} else {
						row[i] = "false"
					}
				case time.Time:
					row[i] = v.Format("2006-01-02 15:04:05")
				default:
					// Fallback to string representation
					row[i] = fmt.Sprintf("%v", val)
				}
			}
		}

		// Check if row contains error information (for CALL statements and stored procedures)
		// MySQL may return error messages in result sets, especially for stored procedures
		// MySQL error format: "ERROR <code> (<SQLSTATE>): <message>"
		// Example: "ERROR 11114 (HY000): The param 'provider' is empty or null"
		for _, cell := range row {
			if cell != "" {
				cellUpper := strings.ToUpper(cell)
				// Check for MySQL error format: ERROR followed by number and SQLSTATE
				if strings.HasPrefix(cellUpper, "ERROR") {
					// Check if it matches MySQL error pattern: ERROR <number> (<SQLSTATE>)
					// This is a MySQL error message returned in result set
					return nil, fmt.Errorf("query execution failed: %s", cell)
				}
			}
		}

		result.Rows = append(result.Rows, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return result, nil
}

// ExecuteNonQuery executes a non-query SQL statement (INSERT, UPDATE, DELETE, etc.)
func (c *Connection) ExecuteNonQuery(ctx context.Context, sqlQuery string) (int64, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	result, err := c.db.ExecContext(queryCtx, sqlQuery)
	if err != nil {
		return 0, fmt.Errorf("query execution failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// TestConnection tests the database connection
func TestConnection(dsn string, dbType string) error {
	conn, err := NewConnection(dsn, dbType)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return conn.Ping(ctx)
}
