package prompt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aiq/aiq/internal/skills"
)

func TestCompressor_Compress_NoCompression(t *testing.T) {
	compressor := NewCompressor(100000)
	
	history := []string{"user: Hello", "assistant: Hi"}
	skillsList := []*skills.Skill{
		{Name: "skill1", Priority: skills.PriorityActive, Loaded: true},
	}
	systemPrompt := "System prompt"
	query := "Query"
	
	result, err := compressor.Compress(history, skillsList, systemPrompt, query)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if result.Compressed {
		t.Error("Should not compress when under threshold")
	}
	
	if len(result.CompressedHistory) != len(history) {
		t.Errorf("Expected history length %d, got %d", len(history), len(result.CompressedHistory))
	}
	
	if len(result.RemainingSkills) != len(skillsList) {
		t.Errorf("Expected skills length %d, got %d", len(skillsList), len(result.RemainingSkills))
	}
}

func TestCompressor_CompressHistory_KeepLastN(t *testing.T) {
	compressor := NewCompressor(100000)
	
	history := make([]string, 20)
	for i := 0; i < 20; i++ {
		history[i] = fmt.Sprintf("message %d", i)
	}
	
	compressed := compressor.compressHistory(history, 10)
	
	// Should have summary + 10 kept messages
	if len(compressed) != 11 {
		t.Errorf("Expected 11 items (1 summary + 10 messages), got %d", len(compressed))
	}
	
	// First item should be summary
	if !strings.Contains(compressed[0], "Previous conversation") {
		t.Error("First item should be summary")
	}
	
	// Last 10 should be kept
	if compressed[1] != "message 10" {
		t.Errorf("Expected 'message 10' as first kept message, got '%s'", compressed[1])
	}
	
	if compressed[10] != "message 19" {
		t.Errorf("Expected 'message 19' as last kept message, got '%s'", compressed[10])
	}
}

func TestCompressor_CompressHistory_ShortHistory(t *testing.T) {
	compressor := NewCompressor(100000)
	
	history := []string{"message 1", "message 2", "message 3"}
	
	compressed := compressor.compressHistory(history, 10)
	
	// Should not compress if history is shorter than keepLast
	if len(compressed) != len(history) {
		t.Errorf("Expected history length %d, got %d", len(history), len(compressed))
	}
}

func TestCompressor_EvictLowPrioritySkills(t *testing.T) {
	compressor := NewCompressor(100000)
	
	skillsList := []*skills.Skill{
		{Name: "active", Priority: skills.PriorityActive, Loaded: true},
		{Name: "relevant", Priority: skills.PriorityRelevant, Loaded: true},
		{Name: "inactive", Priority: skills.PriorityInactive, Loaded: true},
	}
	
	// Evict below Relevant priority
	remaining := compressor.evictLowPrioritySkills(skillsList, skills.PriorityRelevant)
	
	if len(remaining) != 2 {
		t.Errorf("Expected 2 remaining skills, got %d", len(remaining))
	}
	
	// Check that inactive is removed
	for _, skill := range remaining {
		if skill.Name == "inactive" {
			t.Error("Inactive skill should be evicted")
		}
	}
}

func TestCompressor_EvictLowPrioritySkills_KeepOnlyActive(t *testing.T) {
	compressor := NewCompressor(100000)
	
	skillsList := []*skills.Skill{
		{Name: "active", Priority: skills.PriorityActive, Loaded: true},
		{Name: "relevant", Priority: skills.PriorityRelevant, Loaded: true},
		{Name: "inactive", Priority: skills.PriorityInactive, Loaded: true},
	}
	
	// Evict below Active priority
	remaining := compressor.evictLowPrioritySkills(skillsList, skills.PriorityActive)
	
	if len(remaining) != 1 {
		t.Errorf("Expected 1 remaining skill, got %d", len(remaining))
	}
	
	if remaining[0].Name != "active" {
		t.Errorf("Expected 'active' skill, got '%s'", remaining[0].Name)
	}
}

func TestCompressor_SummarizeHistory(t *testing.T) {
	compressor := NewCompressor(100000)
	
	messages := []string{"msg1", "msg2", "msg3"}
	summary := compressor.summarizeHistory(messages)
	
	if !strings.Contains(summary, "Previous conversation") {
		t.Error("Summary should contain 'Previous conversation'")
	}
	
	if !strings.Contains(summary, "3") {
		t.Error("Summary should contain message count")
	}
}

func TestCompressor_NewCompressor_DefaultWindow(t *testing.T) {
	compressor := NewCompressor(0)
	
	if compressor.contextWindow != DefaultContextWindow {
		t.Errorf("Expected default context window %d, got %d", DefaultContextWindow, compressor.contextWindow)
	}
}

func TestCompressor_NewCompressor_CustomWindow(t *testing.T) {
	customWindow := 50000
	compressor := NewCompressor(customWindow)
	
	if compressor.contextWindow != customWindow {
		t.Errorf("Expected custom context window %d, got %d", customWindow, compressor.contextWindow)
	}
}
