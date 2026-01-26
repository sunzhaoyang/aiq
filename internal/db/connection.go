package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Connection represents a database connection
type Connection struct {
	db *sql.DB
}

// NewConnection creates a new database connection
func NewConnection(dsn string, dbType string) (*Connection, error) {
	driverName := "mysql"
	if dbType == "postgresql" {
		driverName = "postgres"
	} else if dbType == "seekdb" {
		// SeekDB might use MySQL driver or custom driver
		driverName = "mysql"
	}
	
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Connection{db: db}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// GetDB returns the underlying sql.DB instance
func (c *Connection) GetDB() *sql.DB {
	return c.db
}

// Ping tests the database connection
func (c *Connection) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}
