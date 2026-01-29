package ui

import (
	"fmt"
	"strings"
)

// Table represents a formatted table (mysql client style)
type Table struct {
	headers []string
	rows    [][]string
}

// NewTable creates a new table
func NewTable(headers []string) *Table {
	return &Table{
		headers: headers,
		rows:    make([][]string, 0),
	}
}

// AddRow adds a row to the table
func (t *Table) AddRow(row []string) {
	// Pad row to match header length
	for len(row) < len(t.headers) {
		row = append(row, "")
	}
	t.rows = append(t.rows, row[:len(t.headers)])
}

// Render renders the table as a string (mysql client style: +---+---+)
func (t *Table) Render() string {
	if len(t.headers) == 0 {
		return ""
	}

	// Calculate column widths
	widths := make([]int, len(t.headers))
	for i, header := range t.headers {
		widths[i] = len(header)
	}

	for _, row := range t.rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Build separator line: +------+------+
	var sepBuilder strings.Builder
	sepBuilder.WriteString("+")
	for _, w := range widths {
		sepBuilder.WriteString(strings.Repeat("-", w+2))
		sepBuilder.WriteString("+")
	}
	separator := sepBuilder.String()

	// Build table
	var builder strings.Builder

	// Top border
	builder.WriteString(separator)
	builder.WriteString("\n")

	// Header row: | col1 | col2 |
	builder.WriteString("|")
	for i, header := range t.headers {
		builder.WriteString(fmt.Sprintf(" %-*s |", widths[i], header))
	}
	builder.WriteString("\n")

	// Separator after header
	builder.WriteString(separator)
	builder.WriteString("\n")

	// Data rows
	for _, row := range t.rows {
		builder.WriteString("|")
		for i, cell := range row {
			builder.WriteString(fmt.Sprintf(" %-*s |", widths[i], cell))
		}
		builder.WriteString("\n")
	}

	// Bottom border
	builder.WriteString(separator)

	return builder.String()
}

// PrintTable prints a table directly
func PrintTable(headers []string, rows [][]string) {
	table := NewTable(headers)
	for _, row := range rows {
		table.AddRow(row)
	}
	fmt.Println(table.Render())
}
