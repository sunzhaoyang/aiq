package tool

import (
	"context"
	"fmt"

	"github.com/aiq/aiq/internal/db"
)

// ExecuteSQL executes a SQL query and returns results
// This function does NOT print anything - it only returns data
// The LLM will decide how to display the results (via render_table or text description)
func ExecuteSQL(ctx context.Context, conn *db.Connection, sql string) (*db.QueryResult, error) {
	result, err := conn.ExecuteQuery(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	return result, nil
}
