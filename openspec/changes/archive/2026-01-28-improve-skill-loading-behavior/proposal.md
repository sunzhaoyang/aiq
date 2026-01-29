## Why

Current skill matching and loading behavior causes several issues: (1) Skills are matched too loosely, loading irrelevant skills for simple queries (e.g., "show me all tables" incorrectly matches MySQL installation skills); (2) In free mode, LLM fails to recognize that database-related queries are invalid without a database connection, leading to incorrect command execution attempts instead of clarifying questions; (3) LLM doesn't properly understand available tools in different modes, causing confusion between SQL execution via `execute_sql` tool vs shell commands. These issues degrade user experience and reduce system reliability.

## What Changes

- **Improve skill matching precision**: Add stricter relevance thresholds and context-aware filtering to prevent irrelevant skill loading
- **Enhance mode awareness**: Strengthen LLM's understanding of free mode vs database mode, and what queries are valid in each mode
- **Add input validation guidance**: Instruct LLM to validate user input appropriateness for current mode before executing tools
- **Improve tool availability communication**: Make it clearer to LLM which tools are available in each mode and their proper usage
- **Add fallback behavior**: When user input is invalid for current mode, LLM should ask clarifying questions instead of guessing commands

## Capabilities

### New Capabilities
- `skill-matching-precision`: Improve skill matching to avoid loading irrelevant skills for simple queries

### Modified Capabilities
- `claude-skills-support`: Improve skill matching precision and add context-aware filtering
- `free-chat-mode`: Enhance LLM's understanding of valid queries in free mode and improve error handling for invalid database queries
- `sql-interactive-mode`: Improve LLM's mode awareness and tool availability understanding

## Impact

- **Skills matching**: Changes to `internal/skills/matcher.go` to add stricter matching criteria
- **System prompts**: Updates to `internal/sql/tool_handler.go` to improve mode awareness and input validation guidance
- **Tool registration**: Potential changes to `internal/tool/llm_functions.go` to better communicate tool availability
- **User experience**: More accurate responses, fewer irrelevant skill loads, better error handling
