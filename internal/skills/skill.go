package skills

// Priority represents the priority level of a Skill
type Priority int

const (
	PriorityInactive Priority = iota // Not matched by current query
	PriorityRelevant                 // Matched but not yet used
	PriorityActive                   // Referenced in current/recent queries
)

// Skill represents a parsed Skill with metadata and content
type Skill struct {
	// Metadata from YAML frontmatter
	Name        string
	Description string

	// Content from Markdown body
	Content string

	// Runtime state
	Priority Priority
	Loaded   bool // Whether full content has been loaded
}

// Metadata represents lightweight Skill metadata (loaded on startup)
type Metadata struct {
	Name        string
	Description string
	Path        string // Path to SKILL.md file
}
