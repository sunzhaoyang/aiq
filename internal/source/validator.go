package source

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	minPort = 1
	maxPort = 65535
)

// Validate validates a source configuration
func Validate(source *Source) error {
	if source == nil {
		return fmt.Errorf("source is nil")
	}

	// Validate name
	if strings.TrimSpace(source.Name) == "" {
		return fmt.Errorf("source name is required")
	}

	// Validate type
	if source.Type == "" {
		return fmt.Errorf("database type is required")
	}
	if source.Type != DatabaseTypeMySQL && source.Type != DatabaseTypePostgreSQL && source.Type != DatabaseTypeSeekDB {
		return fmt.Errorf("invalid database type: %s (must be mysql, postgresql, or seekdb)", source.Type)
	}

	// Validate host
	if strings.TrimSpace(source.Host) == "" {
		return fmt.Errorf("host is required")
	}

	// Validate port
	if source.Port < minPort || source.Port > maxPort {
		return fmt.Errorf("port must be between %d and %d", minPort, maxPort)
	}

	// Validate database
	if strings.TrimSpace(source.Database) == "" {
		return fmt.Errorf("database name is required")
	}

	// Validate username
	if strings.TrimSpace(source.Username) == "" {
		return fmt.Errorf("username is required")
	}

	// Validate password (can be empty, but warn)
	if source.Password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

// ValidateHost validates host format
func ValidateHost(host string) error {
	if strings.TrimSpace(host) == "" {
		return fmt.Errorf("host cannot be empty")
	}

	// Try to parse as IP address
	if ip := net.ParseIP(host); ip != nil {
		return nil
	}

	// Try to parse as hostname
	if _, err := net.LookupHost(host); err == nil {
		return nil
	}

	// If it's localhost or 127.0.0.1, that's fine
	if host == "localhost" || host == "127.0.0.1" {
		return nil
	}

	// Basic hostname validation (contains at least one dot or is localhost)
	if strings.Contains(host, ".") || host == "localhost" {
		return nil
	}

	return fmt.Errorf("invalid host format: %s", host)
}

// ValidatePort validates port number
func ValidatePort(port int) error {
	if port < minPort || port > maxPort {
		return fmt.Errorf("port must be between %d and %d", minPort, maxPort)
	}
	return nil
}

// ParsePort parses a port string to int
func ParsePort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("invalid port format: %w", err)
	}

	if err := ValidatePort(port); err != nil {
		return 0, err
	}

	return port, nil
}
