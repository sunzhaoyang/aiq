package session

import (
	"testing"
	"time"
)

func TestNewSession(t *testing.T) {
	dataSource := "test_source"
	databaseType := "mysql"
	
	sess := NewSession(dataSource, databaseType)
	
	if sess == nil {
		t.Fatal("NewSession returned nil")
	}
	
	if sess.Metadata.DataSource != dataSource {
		t.Errorf("Expected DataSource %s, got %s", dataSource, sess.Metadata.DataSource)
	}
	
	if sess.Metadata.DatabaseType != databaseType {
		t.Errorf("Expected DatabaseType %s, got %s", databaseType, sess.Metadata.DatabaseType)
	}
	
	if sess.Metadata.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	
	if sess.Metadata.LastUpdated.IsZero() {
		t.Error("LastUpdated should not be zero")
	}
	
	if sess.Messages == nil {
		t.Error("Messages should be initialized as empty slice")
	}
	
	if len(sess.Messages) != 0 {
		t.Errorf("Expected empty messages, got %d messages", len(sess.Messages))
	}
}

func TestUpdateLastUpdated(t *testing.T) {
	sess := NewSession("test", "mysql")
	initialTime := sess.Metadata.LastUpdated
	
	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)
	
	sess.UpdateLastUpdated()
	
	if !sess.Metadata.LastUpdated.After(initialTime) {
		t.Error("LastUpdated should be updated to a later time")
	}
}

func TestGetTimestamp(t *testing.T) {
	timestamp := GetTimestamp()
	
	if len(timestamp) != 14 {
		t.Errorf("Expected timestamp length 14 (YYYYMMDDHHMMSS), got %d", len(timestamp))
	}
	
	// Verify format by parsing
	_, err := time.Parse("20060102150405", timestamp)
	if err != nil {
		t.Errorf("Timestamp format invalid: %v", err)
	}
}
