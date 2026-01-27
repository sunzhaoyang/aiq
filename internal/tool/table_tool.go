package tool

import (
	"github.com/aiq/aiq/internal/ui"
)

// RenderTableString formats query results as a table string
// This function does NOT print anything - it only returns the formatted string
func RenderTableString(columns []string, rows [][]string) (string, error) {
	table := ui.NewTable(columns)
	for _, row := range rows {
		table.AddRow(row)
	}
	return table.Render(), nil
}
