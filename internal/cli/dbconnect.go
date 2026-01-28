package cli

import (
	"fmt"
	"strings"

	"github.com/aiq/aiq/internal/db"
	"github.com/aiq/aiq/internal/source"
)

// ValidateConnection validates a database connection with the given parameters
func ValidateConnection(args *DatabaseArgs) error {
	// Create source temporarily for DSN generation
	tempSource := &source.Source{
		Type:     args.Engine,
		Host:     args.Host,
		Port:     args.Port,
		Database: args.Database,
		Username: args.Username,
		Password: args.Password,
	}

	dsn := tempSource.DSN()
	dbType := string(args.Engine)

	conn, err := db.NewConnection(dsn, dbType)
	if err != nil {
		// Provide clearer error messages
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "no such host") {
			return fmt.Errorf("cannot connect to database at %s:%d: %w", args.Host, args.Port, err)
		}
		if strings.Contains(err.Error(), "access denied") || strings.Contains(err.Error(), "authentication") {
			return fmt.Errorf("authentication failed for user '%s': %w", args.Username, err)
		}
		if strings.Contains(err.Error(), "Unknown database") || strings.Contains(err.Error(), "database") {
			return fmt.Errorf("database '%s' does not exist: %w", args.Database, err)
		}
		return fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()

	return nil
}
