## 1. Query Type Detection

- [x] 1.1 Add `DetectQueryType(query string)` function in `internal/skills/matcher.go` to classify queries as "action" or "information"
- [x] 1.2 Implement action keyword detection (install, setup, configure, start, stop, run)
- [x] 1.3 Implement information keyword detection (show, list, query, select, display, get)
- [x] 1.4 Add fallback logic: if query type is ambiguous, default to allowing all skills

## 2. Skill Matching Precision

- [x] 2.1 Add `MinRelevanceScore` constant (default: 0.7) to `internal/skills/matcher.go`
- [x] 2.2 Modify `MatchWithLLM()` to apply relevance threshold filtering (Note: LLM returns skill names without scores, so filtering is applied after LLM matching via query type filtering)
- [x] 2.3 Add skill filtering based on query type: exclude setup/installation skills for information queries
- [x] 2.4 Add skill name pattern detection to identify setup/installation skills (e.g., names containing "init", "setup", "install")
- [x] 2.5 Update `Match()` function to pass query type to filtering logic
- [x] 2.6 Return empty match list if all matches are filtered out

## 3. Enhanced System Prompts

- [x] 3.1 Update free mode prompt in `internal/prompt/loader.go` to add explicit input validation guidance
- [x] 3.2 Add validation instruction: "Before executing any tools, validate if the user's request is appropriate for the current mode"
- [x] 3.3 Add free mode guidance: "If user asks about database operations (show tables, query data, etc.), inform them that a database connection is required"
- [x] 3.4 Add database mode guidance: "If user asks about database operations, use execute_sql tool, NOT shell commands"
- [x] 3.5 Add examples of invalid queries and expected responses to system prompt
- [x] 3.6 Add explicit warning: "Do NOT use execute_command with mysql/psql commands to query databases. Use execute_sql tool instead (in database mode)."

## 4. Tool Description Enhancement

- [x] 4.1 Update `execute_sql` tool description in `internal/tool/llm_functions.go` to mention "Available ONLY in database mode when a database source is selected"
- [x] 4.2 Update `execute_command` tool description to clarify "Use for system operations, NOT for database queries"
- [x] 4.3 Add mode-specific tool availability notes in system prompt (already included in prompt files)

## 5. Integration and Testing

- [ ] 5.1 Test query type detection with various query types (requires manual testing)
- [ ] 5.2 Test skill matching precision: verify "show tables" doesn't match setup skills (requires manual testing)
- [ ] 5.3 Test free mode: verify LLM asks clarifying question for database queries instead of executing commands (requires manual testing)
- [ ] 5.4 Test database mode: verify LLM uses execute_sql tool, not execute_command (requires manual testing)
- [ ] 5.5 Test with ambiguous queries to verify fallback behavior (requires manual testing)
- [ ] 5.6 Verify skill loading behavior: fewer irrelevant skills loaded (requires manual testing)
- [ ] 5.7 Test edge cases: queries that could match multiple skill types (requires manual testing)
