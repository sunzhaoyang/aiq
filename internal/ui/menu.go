package ui

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
)

// MenuItem represents a menu option
type MenuItem struct {
	Label string
	Value string
}

// ShowMenu displays an interactive menu and returns the selected value
func ShowMenu(label string, items []MenuItem) (string, error) {
	searcher := func(input string, index int) bool {
		item := items[index]
		name := strings.Replace(strings.ToLower(item.Label), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:    label,
		Items:    items,
		Searcher: searcher,
		Size:     10,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   Highlight.Render("▸ {{ .Label }}"),
			Inactive: "  {{ .Label }}",
			Selected: Success.Render("✓ {{ .Label }}"),
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return items[index].Value, nil
}

// ShowConfirm displays a confirmation prompt
func ShowConfirm(message string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
		Default:   "n",
	}

	result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return strings.ToLower(result) == "y" || strings.ToLower(result) == "yes", nil
}

// ShowInput displays an input prompt
// If defaultValue is provided, it will be shown as a hint in the label
// User can directly type a new value, or press Enter to use the default
// Pressing Tab will auto-fill the default value
func ShowInput(label string, defaultValue string) (string, error) {
	// If default value is provided, use readline for Tab completion support
	if defaultValue != "" {
		return showInputWithTabCompletion(label, defaultValue)
	}

	// For inputs without default, use standard promptui
	prompt := promptui.Prompt{
		Label: label,
	}

	return prompt.Run()
}

// showInputWithTabCompletion uses readline for better input handling
// Default value is shown as hint, user can type directly or press Enter to use default
func showInputWithTabCompletion(label string, defaultValue string) (string, error) {
	// Format label with hint (default value shown in gray)
	displayLabel := fmt.Sprintf("%s (%s): ", label, InfoText(defaultValue))

	// Create readline instance for better input experience
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          displayLabel,
		HistoryFile:      "",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		// Fallback to promptui if readline fails
		return showInputFallback(label, defaultValue)
	}
	defer rl.Close()

	// Read input
	line, err := rl.Readline()
	if err != nil {
		if err == readline.ErrInterrupt {
			return "", fmt.Errorf("interrupted")
		}
		return "", err
	}

	result := strings.TrimSpace(line)

	// If input is empty, use default value
	if result == "" {
		return defaultValue, nil
	}

	return result, nil
}

// showInputFallback is a fallback implementation using promptui
func showInputFallback(label string, defaultValue string) (string, error) {
	displayLabel := fmt.Sprintf("%s (%s)", label, InfoText(defaultValue))

	prompt := promptui.Prompt{
		Label: displayLabel,
		Validate: func(input string) error {
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	if result == "" && defaultValue != "" {
		return defaultValue, nil
	}

	return result, nil
}

// ShowPassword displays a password input prompt (masked)
func ShowPassword(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
	}

	return prompt.Run()
}

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

// ShowMessage displays a message with optional styling
func ShowMessage(message string, style func(string) string) {
	if style != nil {
		fmt.Println(style(message))
	} else {
		fmt.Println(message)
	}
}

// ShowSuccess displays a success message
func ShowSuccess(message string) {
	fmt.Println(SuccessText("✓ " + message))
}

// ShowError displays an error message
func ShowError(message string) {
	fmt.Println(ErrorText("✗ " + message))
}

// ShowInfo displays an info message
func ShowInfo(message string) {
	fmt.Println(InfoText("ℹ " + message))
}

// ShowWarning displays a warning message
func ShowWarning(message string) {
	fmt.Println(WarningText("⚠ " + message))
}
