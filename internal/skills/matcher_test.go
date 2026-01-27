package skills

import (
	"testing"
)

func TestMatcher_Match_ExactNameMatch(t *testing.T) {
	matcher := NewMatcher()
	
	metadataList := []*Metadata{
		{Name: "seekdb", Description: "Database operations"},
		{Name: "data-analysis", Description: "Data analysis tools"},
	}
	
	// Exact name match should have highest score
	matched := matcher.Match("seekdb", metadataList)
	
	if len(matched) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(matched))
	}
	
	if matched[0].Name != "seekdb" {
		t.Errorf("Expected 'seekdb', got '%s'", matched[0].Name)
	}
}

func TestMatcher_Match_PartialNameMatch(t *testing.T) {
	matcher := NewMatcher()
	
	metadataList := []*Metadata{
		{Name: "seekdb-docs", Description: "SeekDB documentation"},
		{Name: "data-analysis", Description: "Data analysis tools"},
	}
	
	// Partial name match
	matched := matcher.Match("seekdb", metadataList)
	
	if len(matched) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(matched))
	}
	
	if matched[0].Name != "seekdb-docs" {
		t.Errorf("Expected 'seekdb-docs', got '%s'", matched[0].Name)
	}
}

func TestMatcher_Match_DescriptionKeywordMatch(t *testing.T) {
	matcher := NewMatcher()
	
	metadataList := []*Metadata{
		{Name: "database-tool", Description: "Database operations and SQL queries"},
		{Name: "other-tool", Description: "Something else"},
	}
	
	// Description keyword match
	matched := matcher.Match("SQL queries", metadataList)
	
	if len(matched) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(matched))
	}
	
	if matched[0].Name != "database-tool" {
		t.Errorf("Expected 'database-tool', got '%s'", matched[0].Name)
	}
}

func TestMatcher_Match_MultipleMatches(t *testing.T) {
	matcher := NewMatcher()
	
	metadataList := []*Metadata{
		{Name: "seekdb", Description: "SeekDB database operations"},
		{Name: "database-tool", Description: "Database operations"},
		{Name: "other-tool", Description: "Something else"},
	}
	
	// Should return top 3 (default max)
	matched := matcher.Match("database", metadataList)
	
	if len(matched) != 2 {
		t.Fatalf("Expected 2 matches, got %d", len(matched))
	}
	
	// seekdb should be first (name contains "db")
	if matched[0].Name != "seekdb" && matched[0].Name != "database-tool" {
		t.Errorf("Expected 'seekdb' or 'database-tool' as first match, got '%s'", matched[0].Name)
	}
}

func TestMatcher_Match_NoMatches(t *testing.T) {
	matcher := NewMatcher()
	
	metadataList := []*Metadata{
		{Name: "seekdb", Description: "Database operations"},
		{Name: "data-analysis", Description: "Data analysis tools"},
	}
	
	// No matches
	matched := matcher.Match("completely unrelated query", metadataList)
	
	if len(matched) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(matched))
	}
}

func TestMatcher_Match_EmptyMetadataList(t *testing.T) {
	matcher := NewMatcher()
	
	matched := matcher.Match("any query", []*Metadata{})
	
	if len(matched) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(matched))
	}
}

func TestMatcher_Match_CustomMaxSkills(t *testing.T) {
	matcher := NewMatcher()
	matcher.SetMaxSkills(2)
	
	metadataList := []*Metadata{
		{Name: "skill1", Description: "Database operations"},
		{Name: "skill2", Description: "Database queries"},
		{Name: "skill3", Description: "Database tools"},
	}
	
	matched := matcher.Match("database", metadataList)
	
	if len(matched) > 2 {
		t.Errorf("Expected at most 2 matches, got %d", len(matched))
	}
}

func TestMatcher_ScoreSkill_ExactNameMatch(t *testing.T) {
	matcher := NewMatcher()
	
	metadata := &Metadata{
		Name:        "seekdb",
		Description: "Database operations",
	}
	
	score := matcher.scoreSkill([]string{"seekdb"}, metadata)
	
	if score < 100.0 {
		t.Errorf("Expected score >= 100.0 for exact name match, got %f", score)
	}
}

func TestMatcher_ScoreSkill_PartialNameMatch(t *testing.T) {
	matcher := NewMatcher()
	
	metadata := &Metadata{
		Name:        "seekdb-docs",
		Description: "Database operations",
	}
	
	score := matcher.scoreSkill([]string{"seekdb"}, metadata)
	
	if score < 50.0 {
		t.Errorf("Expected score >= 50.0 for partial name match, got %f", score)
	}
}

func TestMatcher_ScoreSkill_DescriptionMatch(t *testing.T) {
	matcher := NewMatcher()
	
	metadata := &Metadata{
		Name:        "other-tool",
		Description: "Database operations and SQL queries",
	}
	
	score := matcher.scoreSkill([]string{"database"}, metadata)
	
	if score < 10.0 {
		t.Errorf("Expected score >= 10.0 for description match, got %f", score)
	}
}

func TestExtractKeywords(t *testing.T) {
	query := "How do I query the database?"
	keywords := extractKeywords(query)
	
	expectedKeywords := []string{"query", "database"}
	
	if len(keywords) != len(expectedKeywords) {
		t.Errorf("Expected %d keywords, got %d: %v", len(expectedKeywords), len(keywords), keywords)
	}
	
	// Check that stop words are filtered
	for _, keyword := range keywords {
		if keyword == "how" || keyword == "do" || keyword == "i" || keyword == "the" {
			t.Errorf("Stop word '%s' should be filtered", keyword)
		}
	}
}

func TestExtractKeywords_EmptyQuery(t *testing.T) {
	keywords := extractKeywords("")
	
	if len(keywords) != 0 {
		t.Errorf("Expected 0 keywords for empty query, got %d", len(keywords))
	}
}

func TestExtractKeywords_OnlyStopWords(t *testing.T) {
	keywords := extractKeywords("the and or")
	
	if len(keywords) != 0 {
		t.Errorf("Expected 0 keywords for only stop words, got %d: %v", len(keywords), keywords)
	}
}

func TestSortByScore(t *testing.T) {
	results := []MatchResult{
		{Metadata: &Metadata{Name: "low"}, Score: 10.0},
		{Metadata: &Metadata{Name: "high"}, Score: 100.0},
		{Metadata: &Metadata{Name: "medium"}, Score: 50.0},
	}
	
	sortByScore(results)
	
	if results[0].Score != 100.0 {
		t.Errorf("Expected highest score first, got %f", results[0].Score)
	}
	
	if results[1].Score != 50.0 {
		t.Errorf("Expected second highest score, got %f", results[1].Score)
	}
	
	if results[2].Score != 10.0 {
		t.Errorf("Expected lowest score last, got %f", results[2].Score)
	}
}
