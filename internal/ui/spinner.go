package ui

import (
	"fmt"
	"time"
)

// Spinner represents a loading spinner
type Spinner struct {
	frames []string
	index  int
	active bool
	done   chan bool
}

// NewSpinner creates a new spinner with default frames
func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		index:  0,
		active: false,
		done:   make(chan bool),
	}
}

// Start starts the spinner animation
func (s *Spinner) Start(message string) {
	if s.active {
		return
	}
	s.active = true
	
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				fmt.Printf("\r%s %s", s.frames[s.index], message)
				s.index = (s.index + 1) % len(s.frames)
			case <-s.done:
				return
			}
		}
	}()
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	if !s.active {
		return
	}
	s.active = false
	s.done <- true
	fmt.Print("\r\033[K") // Clear the line
}

// ShowLoading displays a loading message with spinner
func ShowLoading(message string) func() {
	spinner := NewSpinner()
	spinner.Start(message)
	return spinner.Stop
}
