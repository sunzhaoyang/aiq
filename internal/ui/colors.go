package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Success color (green)
	Success = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	
	// Error color (red)
	Error = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	
	// Info color (blue)
	Info = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	
	// Warning color (yellow)
	Warning = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	
	// Primary text color
	Primary = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	
	// Secondary text color (dimmed)
	Secondary = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	
	// Highlight color for selected items
	Highlight = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	
	// SQL keyword color
	SQLKeyword = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
)

// SuccessText returns text styled as success (green)
func SuccessText(text string) string {
	return Success.Render(text)
}

// ErrorText returns text styled as error (red)
func ErrorText(text string) string {
	return Error.Render(text)
}

// InfoText returns text styled as info (blue)
func InfoText(text string) string {
	return Info.Render(text)
}

// WarningText returns text styled as warning (yellow)
func WarningText(text string) string {
	return Warning.Render(text)
}

// HighlightText returns text styled as highlighted
func HighlightText(text string) string {
	return Highlight.Render(text)
}
