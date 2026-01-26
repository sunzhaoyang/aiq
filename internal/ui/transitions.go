package ui

import (
	"fmt"
	"time"
)

// FadeTransition creates a fade transition effect
func FadeTransition(message string, duration time.Duration) {
	steps := 10
	delay := duration / time.Duration(steps)
	
	for i := 0; i <= steps; i++ {
		alpha := float64(i) / float64(steps)
		// Simple fade effect using spaces and clearing
		if i == 0 {
			fmt.Print("\r")
		} else {
			fmt.Print("\r\033[K")
		}
		if alpha > 0.5 {
			fmt.Print(message)
		}
		time.Sleep(delay)
	}
	fmt.Print("\r\033[K")
}

// SlideTransition creates a slide transition effect (simple version)
func SlideTransition(message string) {
	fmt.Print("\033[2K\r") // Clear line
	fmt.Print(message)
	time.Sleep(100 * time.Millisecond)
}

// ClearLine clears the current line
func ClearLine() {
	fmt.Print("\r\033[K")
}

// ShowTransition shows a transition message
func ShowTransition(message string) {
	ClearLine()
	fmt.Print(InfoText(message))
	time.Sleep(200 * time.Millisecond)
	ClearLine()
}
