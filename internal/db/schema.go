package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// TableInfo represents table information
type TableInfo struct {
	Name    string
	Columns []ColumnInfo
}

// ColumnInfo represents column information
type ColumnInfo struct {
	Name         string
	DataType     string
	IsNullable   string
	ColumnKey    string
	DefaultValue sql.NullString
}

// Schema represents database schema
type Schema struct {
	Tables []TableInfo
}

// GetSchema fetches the database schema
func (c *Connection) GetSchema(ctx context.Context, databaseName string) (*Schema, error) {
	// Get all tables
	tablesQuery := "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME"
	rows, err := c.db.QueryContext(ctx, tablesQuery, databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tableNames = append(tableNames, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tables: %w", err)
	}

	// Get columns for each table
	schema := &Schema{
		Tables: make([]TableInfo, 0, len(tableNames)),
	}

	for _, tableName := range tableNames {
		tableInfo, err := c.getTableInfo(ctx, databaseName, tableName)
		if err != nil {
			return nil, fmt.Errorf("failed to get info for table %s: %w", tableName, err)
		}
		schema.Tables = append(schema.Tables, *tableInfo)
	}

	return schema, nil
}

func (c *Connection) getTableInfo(ctx context.Context, databaseName, tableName string) (*TableInfo, error) {
	query := `
		SELECT 
			COLUMN_NAME,
			DATA_TYPE,
			IS_NULLABLE,
			COLUMN_KEY,
			COLUMN_DEFAULT
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION
	`

	rows, err := c.db.QueryContext(ctx, query, databaseName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tableInfo := &TableInfo{
		Name:    tableName,
		Columns: make([]ColumnInfo, 0),
	}

	for rows.Next() {
		var col ColumnInfo
		var defaultVal sql.NullString
		if err := rows.Scan(&col.Name, &col.DataType, &col.IsNullable, &col.ColumnKey, &defaultVal); err != nil {
			return nil, err
		}
		col.DefaultValue = defaultVal
		tableInfo.Columns = append(tableInfo.Columns, col)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tableInfo, nil
}

// FormatSchema formats schema as a string for LLM context
func (s *Schema) FormatSchema() string {
	var builder strings.Builder

	for _, table := range s.Tables {
		builder.WriteString(fmt.Sprintf("Table: %s\n", table.Name))
		builder.WriteString("Columns:\n")
		for _, col := range table.Columns {
			nullable := "NULL"
			if col.IsNullable == "NO" {
				nullable = "NOT NULL"
			}
			key := ""
			if col.ColumnKey == "PRI" {
				key = " PRIMARY KEY"
			} else if col.ColumnKey == "UNI" {
				key = " UNIQUE"
			}
			builder.WriteString(fmt.Sprintf("  - %s (%s, %s%s)\n", col.Name, col.DataType, nullable, key))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}
