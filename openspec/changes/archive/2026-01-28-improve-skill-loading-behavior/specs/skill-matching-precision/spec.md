## ADDED Requirements

### Requirement: Query type detection for skill filtering
The system SHALL detect query type (action vs information) to filter skills more precisely.

#### Scenario: Detect action query
- **WHEN** user query contains action keywords (install, setup, configure, start, stop, run)
- **THEN** system classifies query as "action" type

#### Scenario: Detect information query
- **WHEN** user query contains information keywords (show, list, query, select, display, get)
- **THEN** system classifies query as "information" type

#### Scenario: Filter setup skills for information queries
- **WHEN** query is classified as information type and matches setup/installation skills
- **THEN** system excludes setup/installation skills from matching results

#### Scenario: Allow setup skills for action queries
- **WHEN** query is classified as action type
- **THEN** system allows matching setup/installation skills

### Requirement: Minimum relevance threshold for skill matching
The system SHALL apply a minimum relevance threshold to skill matching results to avoid loading irrelevant skills.

#### Scenario: Apply relevance threshold to LLM matching
- **WHEN** LLM returns skill matches with low confidence
- **THEN** system filters out skills below minimum relevance threshold (default: 0.7)

#### Scenario: Return empty if all matches below threshold
- **WHEN** all matched skills are below relevance threshold
- **THEN** system returns empty match list instead of loading low-confidence skills

#### Scenario: Allow high-confidence matches even if below threshold
- **WHEN** LLM explicitly indicates high confidence for a skill (e.g., explicit mention in response)
- **THEN** system may include the skill even if below threshold (implementation detail)

### Requirement: Context-aware skill filtering
The system SHALL filter skills based on current mode and query context to avoid irrelevant skill loading.

#### Scenario: Filter database setup skills for simple SQL queries
- **WHEN** user asks simple database query (e.g., "show tables") in database mode
- **THEN** system excludes database setup/installation skills from matching

#### Scenario: Filter skills for free mode database queries
- **WHEN** user asks database-related query in free mode
- **THEN** system excludes all skills and informs user that database connection is required

#### Scenario: Allow relevant skills for complex queries
- **WHEN** user query requires setup or configuration (e.g., "install MySQL", "configure database")
- **THEN** system allows matching setup/configuration skills
