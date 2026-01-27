# Design: Multi-Turn Conversation Support

## Context

AIQ's chat mode currently processes each query independently. Each call to `llmClient.TranslateToSQL()` sends only the current query without any conversation history. This design adds conversation memory and session persistence.

**Constraints:**
- Must maintain backward compatibility with existing single-query flow
- Session files should be human-readable (JSON format)
- Must handle session file corruption gracefully
- Should limit conversation history size to prevent token limit issues

## Goals / Non-Goals

**Goals:**
- Maintain conversation history during chat session
- Send conversation context to LLM for better understanding
- Save sessions to disk on exit
- Restore sessions via command-line flag
- Support conversation history management (clear, view)

**Non-Goals:**
- Real-time session synchronization across multiple instances
- Session encryption (sessions contain only queries and SQL, no sensitive data)
- Automatic session cleanup (users manage their own sessions)
- Session sharing or collaboration features

## Decisions

### 1. Conversation History Storage: In-Memory During Session

**Decision**: Store conversation history in memory as a slice of message pairs (user query + LLM response).

**Rationale:**
- Fast access during conversation
- Simple to implement and maintain
- Persist to disk only on exit or explicit save

**Alternatives Considered:**
- Database storage: Overkill for this use case
- File-based incremental writes: More complex, not needed for session persistence

### 2. Session File Format: JSON

**Decision**: Use JSON format for session files: `session_<timestamp>.json`

**Structure:**
```json
{
  "metadata": {
    "created_at": "2026-01-26T10:00:00Z",
    "last_updated": "2026-01-26T10:15:00Z",
    "data_source": "local-mysql",
    "database_type": "mysql"
  },
  "messages": [
    {
      "role": "user",
      "content": "Show total sales for last week",
      "timestamp": "2026-01-26T10:00:00Z"
    },
    {
      "role": "assistant",
      "content": "SELECT SUM(amount) FROM sales WHERE date >= DATE_SUB(NOW(), INTERVAL 7 DAY)",
      "timestamp": "2026-01-26T10:00:05Z"
    }
  ]
}
```

**Rationale:**
- Human-readable and debuggable
- Easy to parse and modify
- Standard format supported by Go's encoding/json

### 3. LLM Message Format: OpenAI-Compatible Chat Format

**Decision**: Use OpenAI-compatible message format with conversation history.

**Format:**
```go
Messages: []ChatMessage{
    {Role: "system", Content: "You are a SQL expert..."},
    {Role: "user", Content: "First query"},
    {Role: "assistant", Content: "SELECT ..."},
    {Role: "user", Content: "Modify to show only last 3 days"},
    // ... more conversation history
    {Role: "user", Content: "Current query"},
}
```

**Rationale:**
- Compatible with existing LLM client implementation
- Standard format supported by most LLM APIs
- Allows LLM to understand conversation context

### 4. Conversation History Limit: Configurable with Default

**Decision**: Limit conversation history to prevent token limit issues. Default: 20 message pairs (40 messages total).

**Rationale:**
- Prevents hitting LLM token limits
- Keeps context relevant (older messages may be less relevant)
- Configurable for users who need longer conversations

### 5. Session Restoration: Command-Line Flag

**Decision**: Use `-s` or `--session` flag to restore a session.

**Usage:**
```bash
aiq -s ~/.aiqconfig/session_20260126100000.json
```

**Rationale:**
- Simple and intuitive
- Follows common CLI patterns
- Easy to integrate with existing Cobra CLI structure

### 6. Session Save Location: `~/.aiqconfig/`

**Decision**: Save sessions to `~/.aiqconfig/session_<timestamp>.json`

**Rationale:**
- Consistent with existing config directory structure
- Easy to find and manage
- Keeps all AIQ data in one place

## Architecture

### New Components

1. **`internal/session/` package**:
   - `session.go`: Session struct and management
   - `manager.go`: Session save/load operations
   - `history.go`: Conversation history management

2. **Modified Components**:
   - `internal/sql/mode.go`: Integrate session management
   - `internal/llm/client.go`: Support conversation history in TranslateToSQL
   - `cmd/aiq/main.go`: Add `-s/--session` flag

### Data Flow

1. **Start Chat Mode**:
   - If `-s` flag provided: Load session, restore conversation history
   - Otherwise: Create new session

2. **During Conversation**:
   - User enters query
   - Add user message to conversation history
   - Send conversation history + current query to LLM
   - Receive SQL response
   - Add assistant response to conversation history
   - Display results

3. **Exit Chat Mode**:
   - Save session to `~/.aiqconfig/session_<timestamp>.json`
   - Display session file path
   - Show command to resume: `aiq -s ~/.aiqconfig/session_<timestamp>.json`

## Risks

1. **Token Limit**: Long conversations may exceed LLM token limits
   - **Mitigation**: Implement conversation history limit, trim old messages

2. **Session File Corruption**: Invalid JSON or file system errors
   - **Mitigation**: Validate JSON on load, handle errors gracefully, provide fallback

3. **Backward Compatibility**: Existing single-query behavior should still work
   - **Mitigation**: Make conversation history optional, default to empty history

4. **Performance**: Large conversation histories may slow down LLM calls
   - **Mitigation**: Limit history size, optimize message encoding

## Open Questions

1. Should we support conversation history editing (delete specific messages)?
2. Should we add a command to list all saved sessions?
3. Should we support session merging or splitting?
4. Should we add conversation history search functionality?
