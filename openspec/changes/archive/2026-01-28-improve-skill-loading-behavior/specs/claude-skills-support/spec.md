## MODIFIED Requirements

### Requirement: Skills matching algorithm
The system SHALL match user queries to Skills using LLM-based semantic relevance judgment with precision filtering, with keyword-based matching as fallback.

#### Scenario: LLM semantic matching with precision filter
- **WHEN** user submits a query
- **THEN** system sends query and all Skills metadata to LLM for semantic relevance judgment, then applies precision filters (query type, relevance threshold)

#### Scenario: LLM returns relevant Skills with filtering
- **WHEN** LLM processes matching request
- **THEN** system filters LLM results based on query type and relevance threshold, returning only highly relevant Skills

#### Scenario: Filter out irrelevant skills for simple queries
- **WHEN** user submits simple information query (e.g., "show tables") and LLM matches setup/installation skills
- **THEN** system filters out setup/installation skills, returning empty or only directly relevant skills

#### Scenario: Fallback to keyword matching
- **WHEN** LLM semantic matching fails or is unavailable
- **THEN** system falls back to keyword-based matching algorithm with query type filtering

#### Scenario: Select top N Skills after filtering
- **WHEN** filtering completes
- **THEN** system selects top N most relevant Skills (default: 3) from filtered results for loading

#### Scenario: Cache matching results
- **WHEN** same query is matched multiple times
- **THEN** system uses cached results (LLM or keyword) to avoid repeated processing
