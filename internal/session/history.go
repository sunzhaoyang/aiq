package session

import (
	"time"
)

const (
	// DefaultHistoryLimit is the default maximum number of message pairs to keep
	DefaultHistoryLimit = 20
)

// AddMessage adds a message to the conversation history
func (s *Session) AddMessage(role, content string) {
	message := Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now().UTC(),
	}
	
	s.Messages = append(s.Messages, message)
	s.UpdateLastUpdated()
	
	// Trim history if it exceeds the limit
	s.trimHistory(DefaultHistoryLimit)
}

// GetHistory returns all conversation messages
func (s *Session) GetHistory() []Message {
	return s.Messages
}

// ClearHistory clears all conversation messages
func (s *Session) ClearHistory() {
	s.Messages = make([]Message, 0)
	s.UpdateLastUpdated()
}

// trimHistory trims the conversation history to keep only the most recent messages
// Keeps the most recent `limit` message pairs (limit * 2 messages total)
func (s *Session) trimHistory(limit int) {
	maxMessages := limit * 2 // Each pair has user + assistant message
	
	if len(s.Messages) <= maxMessages {
		return
	}
	
	// Keep only the most recent messages
	trimCount := len(s.Messages) - maxMessages
	s.Messages = s.Messages[trimCount:]
}

// GetHistoryLimit returns the current history limit
func GetHistoryLimit() int {
	return DefaultHistoryLimit
}
