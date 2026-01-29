package ui

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// RollingOutput manages a rolling window display for command output
// It uses ANSI escape codes to update in-place, showing only the last N lines
type RollingOutput struct {
	windowSize   int      // Number of lines to display (e.g., 3)
	lines        []string // Buffer of all output lines
	printedLines int      // Number of lines currently displayed
	enabled      bool     // Whether ANSI sequences are supported
	mu           sync.Mutex
}

// NewRollingOutput creates a new rolling output display
func NewRollingOutput(windowSize int) *RollingOutput {
	if windowSize < 1 {
		windowSize = 3
	}
	if windowSize > 5 {
		windowSize = 5
	}

	return &RollingOutput{
		windowSize:   windowSize,
		lines:        make([]string, 0),
		printedLines: 0,
		enabled:      isANSISupported(),
	}
}

// isANSISupported checks if the terminal supports ANSI escape sequences
func isANSISupported() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// AddLine adds a new line and updates the rolling display
func (r *RollingOutput) AddLine(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Clean line
	line = strings.TrimRight(line, "\n\r")

	// Add to buffer
	r.lines = append(r.lines, line)

	// Update display
	r.updateDisplay()
}

// updateDisplay refreshes the rolling window
func (r *RollingOutput) updateDisplay() {
	if !r.enabled {
		// Fallback: just print the latest line
		if len(r.lines) > 0 {
			fmt.Printf("  %s\n", HintText(r.lines[len(r.lines)-1]))
		}
		return
	}

	// Calculate which lines to show (last windowSize lines)
	start := len(r.lines) - r.windowSize
	if start < 0 {
		start = 0
	}
	displayLines := r.lines[start:]

	// Clear previously printed lines by moving cursor up and clearing
	if r.printedLines > 0 {
		for i := 0; i < r.printedLines; i++ {
			// Move up one line and clear it
			fmt.Print("\033[A\033[K")
		}
	}

	// Print the new lines
	for _, line := range displayLines {
		fmt.Printf("  %s\n", HintText(line))
	}

	// Update count of printed lines
	r.printedLines = len(displayLines)
}

// Finish finalizes the display
// Shows summary line if there are many lines
func (r *RollingOutput) Finish() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.lines) == 0 {
		return
	}

	// If there are more lines than displayed, show a hint
	if len(r.lines) > r.windowSize {
		fmt.Printf("  %s\n", HintText(fmt.Sprintf("... (%d more lines above)", len(r.lines)-r.windowSize)))
	}
}

// GetTotalLines returns the total number of lines captured
func (r *RollingOutput) GetTotalLines() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.lines)
}

// GetLines returns all buffered lines
func (r *RollingOutput) GetLines() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]string, len(r.lines))
	copy(result, r.lines)
	return result
}
