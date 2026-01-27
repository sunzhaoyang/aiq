package skills

import (
	"strings"
)

const (
	// DefaultMaxSkills is the default maximum number of Skills to load per query
	DefaultMaxSkills = 3
)

// MatchResult represents a match between a query and a Skill
type MatchResult struct {
	Metadata *Metadata
	Score    float64
}

// Matcher matches user queries to Skills based on relevance
type Matcher struct {
	maxSkills int
}

// NewMatcher creates a new Skills matcher
func NewMatcher() *Matcher {
	return &Matcher{
		maxSkills: DefaultMaxSkills,
	}
}

// SetMaxSkills sets the maximum number of Skills to return
func (m *Matcher) SetMaxSkills(max int) {
	m.maxSkills = max
}

// Match matches a query against Skills metadata and returns the most relevant Skills
func (m *Matcher) Match(query string, metadataList []*Metadata) []*Metadata {
	if len(metadataList) == 0 {
		return []*Metadata{}
	}

	// Extract keywords from query
	keywords := extractKeywords(query)

	// Score each Skill
	results := make([]MatchResult, 0, len(metadataList))
	for _, md := range metadataList {
		score := m.scoreSkill(keywords, md)
		if score > 0 {
			results = append(results, MatchResult{
				Metadata: md,
				Score:    score,
			})
		}
	}

	// Sort by score (descending)
	sortByScore(results)

	// Return top N
	max := m.maxSkills
	if len(results) < max {
		max = len(results)
	}

	matched := make([]*Metadata, 0, max)
	for i := 0; i < max; i++ {
		matched = append(matched, results[i].Metadata)
	}

	return matched
}

// scoreSkill scores a Skill based on how well it matches the query keywords
// Scoring priority: exact name match > partial name match > description keyword match
func (m *Matcher) scoreSkill(keywords []string, metadata *Metadata) float64 {
	var score float64

	queryLower := strings.ToLower(strings.Join(keywords, " "))
	nameLower := strings.ToLower(metadata.Name)
	descLower := strings.ToLower(metadata.Description)

	// Exact name match (highest priority)
	if queryLower == nameLower {
		score += 100.0
	}

	// Partial name match
	if strings.Contains(nameLower, queryLower) || strings.Contains(queryLower, nameLower) {
		score += 50.0
	}

	// Check if any keyword matches name
	for _, keyword := range keywords {
		keywordLower := strings.ToLower(keyword)
		if strings.Contains(nameLower, keywordLower) {
			score += 30.0
		}
	}

	// Description keyword match
	for _, keyword := range keywords {
		keywordLower := strings.ToLower(keyword)
		if strings.Contains(descLower, keywordLower) {
			score += 10.0
		}
	}

	return score
}

// extractKeywords extracts keywords from a query
// Simple implementation: split by spaces and filter common words
func extractKeywords(query string) []string {
	// Common words to filter out
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true, "from": true,
		"is": true, "are": true, "was": true, "were": true, "be": true,
		"been": true, "have": true, "has": true, "had": true, "do": true,
		"does": true, "did": true, "will": true, "would": true, "could": true,
		"should": true, "may": true, "might": true, "must": true, "can": true,
		"this": true, "that": true, "these": true, "those": true,
		"i": true, "you": true, "he": true, "she": true, "it": true,
		"we": true, "they": true, "what": true, "which": true, "who": true,
		"when": true, "where": true, "why": true, "how": true,
	}

	words := strings.Fields(strings.ToLower(query))
	keywords := make([]string, 0, len(words))

	for _, word := range words {
		// Remove punctuation
		word = strings.Trim(word, ".,!?;:()[]{}'\"")
		if word == "" {
			continue
		}

		// Filter stop words
		if !stopWords[word] && len(word) > 1 {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

// sortByScore sorts match results by score in descending order
func sortByScore(results []MatchResult) {
	// Simple bubble sort (fine for small lists)
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Score < results[j].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}
