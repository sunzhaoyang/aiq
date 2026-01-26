package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Table represents a formatted table
type Table struct {
	headers []string
	rows    [][]string
	style   lipgloss.Style
}

// NewTable creates a new table
func NewTable(headers []string) *Table {
	return &Table{
		headers: headers,
		rows:    make([][]string, 0),
		style:   lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("8")),
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

// Render renders the table as a string
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

	// Add padding
	for i := range widths {
		widths[i] += 2
	}

	// Build table
	var builder strings.Builder
	
	// Header row
	headerCells := make([]string, len(t.headers))
	for i, header := range t.headers {
		headerCells[i] = Highlight.Render(fmt.Sprintf(" %-*s ", widths[i]-2, header))
	}
	builder.WriteString("┌" + strings.Repeat("─", widths[0]))
	for i := 1; i < len(widths); i++ {
		builder.WriteString("┬" + strings.Repeat("─", widths[i]))
	}
	builder.WriteString("┐\n")
	builder.WriteString("│" + strings.Join(headerCells, "│") + "│\n")
	
	// Separator
	builder.WriteString("├" + strings.Repeat("─", widths[0]))
	for i := 1; i < len(widths); i++ {
		builder.WriteString("┼" + strings.Repeat("─", widths[i]))
	}
	builder.WriteString("┤\n")
	
	// Data rows
	for _, row := range t.rows {
		cells := make([]string, len(row))
		for i, cell := range row {
			cells[i] = fmt.Sprintf(" %-*s ", widths[i]-2, cell)
		}
		builder.WriteString("│" + strings.Join(cells, "│") + "│\n")
	}
	
	// Footer
	builder.WriteString("└" + strings.Repeat("─", widths[0]))
	for i := 1; i < len(widths); i++ {
		builder.WriteString("┴" + strings.Repeat("─", widths[i]))
	}
	builder.WriteString("┘\n")

	return t.style.Render(builder.String())
}

// PrintTable prints a table directly
func PrintTable(headers []string, rows [][]string) {
	table := NewTable(headers)
	for _, row := range rows {
		table.AddRow(row)
	}
	fmt.Print(table.Render())
}
