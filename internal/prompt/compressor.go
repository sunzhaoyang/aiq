package prompt

import (
	"fmt"
	"strings"

	"github.com/aiq/aiq/internal/skills"
)

const (
	// Compression thresholds as percentages of context window
	ThresholdCompressHistory = 0.80 // 80% - start compressing conversation history
	ThresholdEvictSkills     = 0.90 // 90% - start evicting low-priority Skills
	ThresholdAggressive      = 0.95 // 95% - aggressive compression

	// Default context window size (conservative estimate for most models)
	DefaultContextWindow = 100000 // 100k tokens
)

// Compressor manages prompt compression to stay within token limits
type Compressor struct {
	contextWindow int
}

// NewCompressor creates a new prompt compressor
func NewCompressor(contextWindow int) *Compressor {
	if contextWindow <= 0 {
		contextWindow = DefaultContextWindow
	}
	return &Compressor{
		contextWindow: contextWindow,
	}
}

// CompressionResult represents the result of compression
type CompressionResult struct {
	CompressedHistory []string
	RemainingSkills   []*skills.Skill
	Compressed        bool
}

// Compress compresses a prompt if it exceeds token thresholds
func (c *Compressor) Compress(
	conversationHistory []string,
	loadedSkills []*skills.Skill,
	systemPrompt string,
	currentQuery string,
) (*CompressionResult, error) {
	// Estimate total tokens
	totalTokens := EstimatePromptTokens(
		systemPrompt,
		strings.Join(conversationHistory, "\n"),
		currentQuery,
	)

	result := &CompressionResult{
		CompressedHistory: conversationHistory,
		RemainingSkills:   loadedSkills,
		Compressed:        false,
	}

	threshold := float64(totalTokens) / float64(c.contextWindow)

	// 80% threshold: Compress conversation history
	if threshold >= ThresholdCompressHistory {
		result.CompressedHistory = c.compressHistory(conversationHistory, 10) // Keep last 10 messages
		result.Compressed = true

		// Re-estimate after compression
		totalTokens = EstimatePromptTokens(
			systemPrompt,
			strings.Join(result.CompressedHistory, "\n"),
			currentQuery,
		)
		threshold = float64(totalTokens) / float64(c.contextWindow)
	}

	// 90% threshold: Evict low-priority Skills
	if threshold >= ThresholdEvictSkills {
		result.RemainingSkills = c.evictLowPrioritySkills(loadedSkills, skills.PriorityRelevant)
		totalTokens = EstimatePromptTokens(
			systemPrompt,
			strings.Join(result.CompressedHistory, "\n"),
			currentQuery,
		)
		threshold = float64(totalTokens) / float64(c.contextWindow)
	}

	// 95% threshold: Aggressive compression
	if threshold >= ThresholdAggressive {
		result.CompressedHistory = c.compressHistory(conversationHistory, 5) // Keep only last 5 messages
		result.RemainingSkills = c.evictLowPrioritySkills(loadedSkills, skills.PriorityActive) // Keep only active Skills
	}

	return result, nil
}

// compressHistory compresses conversation history, keeping only the last N messages
func (c *Compressor) compressHistory(history []string, keepLast int) []string {
	if len(history) <= keepLast {
		return history
	}

	// Keep last N messages
	kept := history[len(history)-keepLast:]

	// Summarize the rest
	toSummarize := history[:len(history)-keepLast]
	summary := c.summarizeHistory(toSummarize)

	// Return summary + kept messages
	result := []string{summary}
	result = append(result, kept...)

	return result
}

// summarizeHistory creates a summary of conversation history
// Simple implementation: just indicate how many messages were compressed
func (c *Compressor) summarizeHistory(messages []string) string {
	return fmt.Sprintf("[Previous conversation: %d messages compressed]", len(messages))
}

// evictLowPrioritySkills removes Skills with priority lower than the given threshold
func (c *Compressor) evictLowPrioritySkills(skillList []*skills.Skill, minPriority skills.Priority) []*skills.Skill {
	var remaining []*skills.Skill
	for _, skill := range skillList {
		if skill.Priority >= minPriority {
			remaining = append(remaining, skill)
		}
	}
	return remaining
}
