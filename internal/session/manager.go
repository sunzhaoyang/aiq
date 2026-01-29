package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aiq/aiq/internal/config"
)

// SaveSession saves the session to a JSON file
func SaveSession(session *Session, filePath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	// Update last updated timestamp
	session.UpdateLastUpdated()

	// Marshal session to JSON
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// LoadSession loads a session from a JSON file
func LoadSession(filePath string) (*Session, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	// Validate JSON format
	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to parse session file (invalid JSON): %w", err)
	}

	// Validate required fields
	if err := validateSession(&session); err != nil {
		return nil, fmt.Errorf("session validation failed: %w", err)
	}

	return &session, nil
}

// validateSession validates that a session has all required fields
func validateSession(session *Session) error {
	if session.Metadata.DataSource == "" {
		return fmt.Errorf("missing data_source in session metadata")
	}
	if session.Metadata.DatabaseType == "" {
		return fmt.Errorf("missing database_type in session metadata")
	}
	if session.Metadata.CreatedAt.IsZero() {
		return fmt.Errorf("missing created_at in session metadata")
	}

	// Validate messages
	for i, msg := range session.Messages {
		if msg.Role != "user" && msg.Role != "assistant" {
			return fmt.Errorf("invalid role '%s' in message %d", msg.Role, i)
		}
		if msg.Content == "" {
			return fmt.Errorf("empty content in message %d", i)
		}
	}

	return nil
}

// GetSessionFilePath generates a session file path with timestamp
// Format: ~/.aiq/sessions/session_YYYYMMDDHHMMSS.json
func GetSessionFilePath(timestamp string) (string, error) {
	sessionsDir, err := config.GetSessionsDir()
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("session_%s.json", timestamp)
	return filepath.Join(sessionsDir, fileName), nil
}
