package skills

import (
	"fmt"
	"log"
	"sync"
)

// Manager manages Skills lifecycle: loading, caching, and eviction
type Manager struct {
	// Metadata loaded on startup (lightweight)
	metadataList []*Metadata

	// Cache of loaded Skills (full content)
	cache map[string]*Skill

	// Mutex for thread-safe access
	mu sync.RWMutex
}

// NewManager creates a new Skills manager
func NewManager() *Manager {
	return &Manager{
		metadataList: []*Metadata{},
		cache:        make(map[string]*Skill),
	}
}

// Initialize loads Skills metadata on startup
func (m *Manager) Initialize() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	metadata, err := LoadSkillsMetadata()
	if err != nil {
		return fmt.Errorf("failed to load skills metadata: %w", err)
	}

	m.metadataList = metadata
	if len(metadata) > 0 {
		log.Printf("Loaded %d skill(s) metadata (progressive loading enabled)", len(metadata))
		for _, md := range metadata {
			log.Printf("  - %s: %s", md.Name, md.Description)
		}
	}

	return nil
}

// GetMetadata returns all Skills metadata
func (m *Manager) GetMetadata() []*Metadata {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make([]*Metadata, len(m.metadataList))
	copy(result, m.metadataList)
	return result
}

// LoadSkill loads a Skill by name, using cache if available
func (m *Manager) LoadSkill(name string) (*Skill, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check cache first
	if skill, exists := m.cache[name]; exists && skill.Loaded {
		return skill, nil
	}

	// Find metadata
	var metadata *Metadata
	for _, md := range m.metadataList {
		if md.Name == name {
			metadata = md
			break
		}
	}

	if metadata == nil {
		return nil, fmt.Errorf("skill '%s' not found", name)
	}

	// Load full content (progressive loading: only now do we read the full file)
	log.Printf("Loading full content for skill: %s", name)
	skill, err := LoadSkillContent(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to load skill content: %w", err)
	}

	// Cache it
	m.cache[name] = skill
	log.Printf("Skill '%s' loaded and cached (content length: %d chars)", name, len(skill.Content))

	return skill, nil
}

// LoadSkills loads multiple Skills by name
func (m *Manager) LoadSkills(names []string) ([]*Skill, error) {
	var skills []*Skill
	var errors []error

	for _, name := range names {
		skill, err := m.LoadSkill(name)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		skills = append(skills, skill)
	}

	if len(errors) > 0 {
		return skills, fmt.Errorf("failed to load some skills: %v", errors)
	}

	return skills, nil
}

// SetPriority sets the priority of a Skill
func (m *Manager) SetPriority(name string, priority Priority) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if skill, exists := m.cache[name]; exists {
		skill.Priority = priority
	}
}

// EvictSkill removes a Skill from cache (for prompt compression)
func (m *Manager) EvictSkill(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.cache, name)
}

// EvictLowPrioritySkills removes Skills with priority lower than the given priority
func (m *Manager) EvictLowPrioritySkills(minPriority Priority) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, skill := range m.cache {
		if skill.Priority < minPriority {
			delete(m.cache, name)
		}
	}
}

// GetCachedSkills returns all currently cached Skills
func (m *Manager) GetCachedSkills() []*Skill {
	m.mu.RLock()
	defer m.mu.RUnlock()

	skills := make([]*Skill, 0, len(m.cache))
	for _, skill := range m.cache {
		skills = append(skills, skill)
	}
	return skills
}
