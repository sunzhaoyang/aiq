package prompt

// TokenEstimator estimates token count for text content
// Uses a simple approximation: ~4 characters per token (conservative estimate)
const charsPerToken = 4.0

// EstimateTokens estimates the number of tokens in a text string
func EstimateTokens(text string) int {
	if len(text) == 0 {
		return 0
	}
	// Conservative estimate: 4 characters per token
	// This is a rough approximation; actual tokenization varies by model
	return int(float64(len(text)) / charsPerToken)
}

// EstimatePromptTokens estimates total tokens for a prompt with multiple components
func EstimatePromptTokens(components ...string) int {
	total := 0
	for _, component := range components {
		total += EstimateTokens(component)
	}
	return total
}
