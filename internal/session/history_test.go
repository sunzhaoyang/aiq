package session

import (
	"testing"
)

func TestAddMessage(t *testing.T) {
	sess := NewSession("test", "mysql")
	
	// Add user message
	sess.AddMessage("user", "What is the total revenue?")
	
	if len(sess.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(sess.Messages))
	}
	
	msg := sess.Messages[0]
	if msg.Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", msg.Role)
	}
	
	if msg.Content != "What is the total revenue?" {
		t.Errorf("Expected content 'What is the total revenue?', got '%s'", msg.Content)
	}
	
	if msg.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
	
	// Add assistant message
	sess.AddMessage("assistant", "SELECT SUM(revenue) FROM sales")
	
	if len(sess.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(sess.Messages))
	}
	
	msg2 := sess.Messages[1]
	if msg2.Role != "assistant" {
		t.Errorf("Expected role 'assistant', got '%s'", msg2.Role)
	}
}

func TestGetHistory(t *testing.T) {
	sess := NewSession("test", "mysql")
	
	// Initially empty
	history := sess.GetHistory()
	if len(history) != 0 {
		t.Errorf("Expected empty history, got %d messages", len(history))
	}
	
	// Add messages
	sess.AddMessage("user", "Query 1")
	sess.AddMessage("assistant", "Response 1")
	sess.AddMessage("user", "Query 2")
	
	history = sess.GetHistory()
	if len(history) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(history))
	}
}

func TestClearHistory(t *testing.T) {
	sess := NewSession("test", "mysql")
	
	// Add some messages
	sess.AddMessage("user", "Query 1")
	sess.AddMessage("assistant", "Response 1")
	sess.AddMessage("user", "Query 2")
	
	if len(sess.Messages) != 3 {
		t.Errorf("Expected 3 messages before clear, got %d", len(sess.Messages))
	}
	
	sess.ClearHistory()
	
	if len(sess.Messages) != 0 {
		t.Errorf("Expected empty messages after clear, got %d messages", len(sess.Messages))
	}
	
	// Verify last updated was updated
	if sess.Metadata.LastUpdated.IsZero() {
		t.Error("LastUpdated should be updated after clear")
	}
}

func TestHistoryLimit(t *testing.T) {
	sess := NewSession("test", "mysql")
	
	// Add messages beyond the limit
	// DefaultHistoryLimit is 20, so we need 21 pairs = 42 messages
	for i := 0; i < 42; i++ {
		if i%2 == 0 {
			sess.AddMessage("user", "Query")
		} else {
			sess.AddMessage("assistant", "Response")
		}
	}
	
	// Should be trimmed to 20 pairs = 40 messages
	if len(sess.Messages) != 40 {
		t.Errorf("Expected 40 messages (20 pairs), got %d", len(sess.Messages))
	}
	
	// Verify oldest messages were removed
	// First message should be from later in the sequence
	firstMsg := sess.Messages[0]
	if firstMsg.Content != "Query" {
		t.Errorf("Expected first message content 'Query', got '%s'", firstMsg.Content)
	}
}

func TestHistoryLimitExact(t *testing.T) {
	sess := NewSession("test", "mysql")
	
	// Add exactly 20 pairs (40 messages)
	for i := 0; i < 40; i++ {
		if i%2 == 0 {
			sess.AddMessage("user", "Query")
		} else {
			sess.AddMessage("assistant", "Response")
		}
	}
	
	// Should still be 40 messages (no trimming needed)
	if len(sess.Messages) != 40 {
		t.Errorf("Expected 40 messages, got %d", len(sess.Messages))
	}
	
	// Add one more pair
	sess.AddMessage("user", "New Query")
	sess.AddMessage("assistant", "New Response")
	
	// Should still be 40 messages (oldest pair removed)
	if len(sess.Messages) != 40 {
		t.Errorf("Expected 40 messages after limit, got %d", len(sess.Messages))
	}
}

func TestGetHistoryLimit(t *testing.T) {
	limit := GetHistoryLimit()
	if limit != DefaultHistoryLimit {
		t.Errorf("Expected history limit %d, got %d", DefaultHistoryLimit, limit)
	}
}
