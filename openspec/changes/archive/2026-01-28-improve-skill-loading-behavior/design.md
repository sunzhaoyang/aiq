## Context

Current skill matching uses LLM-based semantic matching with keyword fallback, but lacks precision filters. This causes irrelevant skills to be loaded (e.g., "show me all tables" matches MySQL installation skills). Additionally, the system prompt doesn't sufficiently guide LLM to:
1. Recognize invalid queries for current mode (e.g., database queries in free mode)
2. Understand tool availability differences between modes
3. Ask clarifying questions instead of guessing commands when input is ambiguous

The existing codebase has:
- LLM-based semantic matching in `internal/skills/matcher.go` with caching
- Mode-aware system prompts in `internal/sql/tool_handler.go` (free mode vs database mode)
- Tool registration in `internal/tool/llm_functions.go` that conditionally includes `execute_sql` based on connection

## Goals / Non-Goals

**Goals:**
- Improve skill matching precision to avoid loading irrelevant skills for simple queries
- Enhance LLM's mode awareness to recognize invalid queries for current mode
- Add input validation guidance so LLM asks clarifying questions instead of guessing
- Improve tool availability communication in system prompts
- Add minimum relevance threshold for skill matching

**Non-Goals:**
- Changing the fundamental skill matching algorithm (keep LLM semantic matching)
- Removing skill loading entirely (still load relevant skills)
- Changing tool registration logic (keep conditional `execute_sql` registration)
- Adding user-facing skill filtering options (keep automatic matching)

## Decisions

### 1. Skill Matching Precision Improvement
**Decision**: Add minimum relevance threshold and query type filtering to skill matching.

**Rationale**:
- Current LLM matching returns any skills it thinks are relevant, even if relevance is low
- Simple queries like "show tables" shouldn't match installation/setup skills
- Need to filter out skills that are only tangentially related

**Implementation**:
- Add `MinRelevanceScore` threshold (e.g., 0.7) - if LLM returns skills but confidence is low, return empty
- Add query type detection: distinguish between "action queries" (install, setup) vs "information queries" (show, list, query)
- Only match action-oriented skills (like `init-mysql-mac`) for action queries
- For information queries in database mode, don't load setup/installation skills

**Alternatives Considered**:
- User-defined skill filters: Too complex, defeats automatic matching
- Completely disable skills for simple queries: Too restrictive, might miss relevant skills
- Manual skill whitelist: Not scalable, requires maintenance

### 2. Mode-Aware Input Validation
**Decision**: Enhance system prompt to explicitly guide LLM to validate input appropriateness for current mode before executing tools.

**Rationale**:
- Current prompt mentions free mode limitations but doesn't emphasize validation
- LLM needs explicit instruction to check if query is valid for current mode
- Should ask clarifying questions when input is ambiguous or invalid

**Implementation**:
- Add explicit validation step in system prompt: "Before executing any tools, validate if the user's request is appropriate for the current mode"
- For free mode: "If user asks about database operations (show tables, query data, etc.), inform them that a database connection is required"
- For database mode: "If user asks about database operations, use execute_sql tool, NOT shell commands"
- Add examples of invalid queries and expected responses

**Alternatives Considered**:
- Pre-validate queries before sending to LLM: Too restrictive, might block valid queries
- Add separate validation LLM call: Adds latency and complexity
- Only validate after tool execution fails: Too late, wastes API calls

### 3. Tool Availability Communication
**Decision**: Make tool availability more explicit in system prompt and tool descriptions.

**Rationale**:
- Current prompt lists available tools but doesn't emphasize mode differences clearly
- Tool descriptions don't mention mode requirements
- LLM might confuse `execute_sql` tool with shell commands

**Implementation**:
- Enhance system prompt to clearly state: "In FREE MODE: execute_sql tool is NOT available. Use execute_command only for system operations, NOT for database queries."
- Add to tool descriptions: "execute_sql: Available ONLY in database mode when a database source is selected"
- Add explicit warning: "Do NOT use execute_command with mysql/psql commands to query databases. Use execute_sql tool instead (in database mode)."

**Alternatives Considered**:
- Remove execute_command in free mode: Too restrictive, needed for other operations
- Add mode check in tool execution: Already exists, but LLM needs better guidance upfront

### 4. Query Type Detection for Skill Filtering
**Decision**: Detect query type (action vs information) and filter skills accordingly.

**Rationale**:
- Installation/setup skills (like `init-mysql-mac`) should only match action queries
- Information queries (show tables, list data) shouldn't load setup skills
- Reduces irrelevant skill loading

**Implementation**:
- Add query type detection: check for action keywords (install, setup, configure, start) vs information keywords (show, list, query, select)
- For action queries: allow matching setup/installation skills
- For information queries: exclude setup/installation skills from matching
- Pass query type to skill matcher for filtering

**Alternatives Considered**:
- Skill-level metadata for query types: Requires updating all skills, not scalable
- User-specified query type: Too complex, defeats automatic matching
- Always exclude setup skills: Too restrictive, might miss relevant cases

## Risks / Trade-offs

**[Risk] Over-filtering skills might miss relevant skills**
- **Mitigation**: Use conservative thresholds, allow override for high-confidence matches

**[Risk] Query type detection might misclassify queries**
- **Mitigation**: Use simple keyword-based detection with fallback to allow all skills if uncertain

**[Risk] Enhanced prompts might increase token usage**
- **Mitigation**: Keep prompts concise, focus on key points, use compression when needed

**[Risk] LLM might still ignore validation guidance**
- **Mitigation**: Make validation instructions prominent, add examples, test with various queries

**[Trade-off] Precision vs. Recall**
- **Decision**: Favor precision (fewer false positives) over recall (might miss some relevant skills). Better to load fewer skills correctly than load many irrelevant ones.

## Migration Plan

1. **Phase 1: Skill Matching Precision**
   - Add query type detection function
   - Add minimum relevance threshold to LLM matching
   - Add skill filtering based on query type
   - Test with various query types

2. **Phase 2: Enhanced System Prompts**
   - Update free mode prompt with explicit validation guidance
   - Update database mode prompt with tool usage clarification
   - Add examples of invalid queries and expected responses
   - Test LLM responses with invalid queries

3. **Phase 3: Tool Description Enhancement**
   - Update tool descriptions to mention mode requirements
   - Add explicit warnings about tool misuse
   - Test tool selection behavior

4. **Phase 4: Integration and Testing**
   - Test end-to-end with various query types
   - Verify skill loading behavior
   - Verify LLM response quality for invalid queries

## Open Questions

- What should be the minimum relevance threshold? (Currently: 0.7, but needs testing)
- Should we add skill-level metadata for query types? (Currently: No, use query analysis)
- How to handle ambiguous queries that could be either action or information? (Currently: Allow all skills, but prioritize based on query type)
