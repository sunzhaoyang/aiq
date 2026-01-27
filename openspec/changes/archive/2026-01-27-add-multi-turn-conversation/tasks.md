## 1. Session Package Foundation

- [x] 1.1 Create `internal/session/` directory structure
- [x] 1.2 Create `internal/session/session.go` with Session, SessionMetadata, Message structs
- [x] 1.3 Implement session creation and initialization
- [x] 1.4 Add timestamp generation for session files
- [x] 1.5 Implement session metadata management (created_at, last_updated, data_source, database_type)

## 2. Conversation History Management

- [x] 2.1 Create `internal/session/history.go` for conversation history management
- [x] 2.2 Implement AddMessage function (user and assistant messages)
- [x] 2.3 Implement GetHistory function to retrieve conversation history
- [x] 2.4 Implement ClearHistory function
- [x] 2.5 Add conversation history limit (default: 20 message pairs)
- [x] 2.6 Implement automatic trimming of old messages when limit exceeded

## 3. Session Persistence

- [x] 3.1 Create `internal/session/manager.go` for session save/load operations
- [x] 3.2 Implement SaveSession function to write session to JSON file
- [x] 3.3 Implement LoadSession function to read session from JSON file
- [x] 3.4 Add session file validation (JSON format, required fields)
- [x] 3.5 Implement error handling for corrupted session files
- [x] 3.6 Add session file path generation (`~/.aiqconfig/session_<timestamp>.json`)

## 4. LLM Client Enhancement

- [x] 4.1 Modify `internal/llm/client.go` TranslateToSQL to accept conversation history
- [x] 4.2 Update message building logic to include conversation history
- [x] 4.3 Ensure system message is always first in message list
- [x] 4.4 Maintain backward compatibility (empty history works as before)
- [x] 4.5 Test with various conversation history lengths

## 5. Chat Mode Integration

- [x] 5.1 Modify `internal/sql/mode.go` to create session on entry
- [x] 5.2 Integrate conversation history into query flow
- [x] 5.3 Add user message to history before LLM call
- [x] 5.4 Add assistant response to history after LLM call
- [x] 5.5 Implement session save on exit (Ctrl+D or 'exit' command)
- [x] 5.6 Display session save message with resume command
- [x] 5.7 Add `/history` command to view conversation history
- [x] 5.8 Add `/clear` command to clear conversation history

## 6. CLI Flag Support

- [x] 6.1 Modify `cmd/aiq/main.go` to add `-s/--session` flag
- [x] 6.2 Implement session file loading from flag
- [x] 6.3 Pass session to chat mode initialization
- [x] 6.4 Handle session file not found error
- [x] 6.5 Validate session file before loading

## 7. Session Restoration

- [x] 7.1 Implement session restoration in chat mode entry
- [x] 7.2 Restore conversation history from session
- [x] 7.3 Restore data source from session metadata
- [x] 7.4 Handle missing data source gracefully (prompt for selection)
- [x] 7.5 Display restored session info to user

## 8. Testing and Edge Cases

- [x] 8.1 Test conversation history with single query (backward compatibility)
- [x] 8.2 Test multi-turn conversation flow
- [x] 8.3 Test conversation history limit enforcement
- [x] 8.4 Test session save and restore
- [x] 8.5 Test corrupted session file handling
- [x] 8.6 Test missing data source restoration (handled in mode.go, no separate test needed)
- [x] 8.7 Test session file with very long conversation history
- [x] 8.8 Test `/history` and `/clear` commands
- [x] 8.9 Test session file path generation and uniqueness

## 9. Documentation

- [x] 9.1 Update README.md with multi-turn conversation feature
- [x] 9.2 Add session management examples
- [x] 9.3 Document `-s/--session` flag usage
- [x] 9.4 Add conversation history commands documentation
- [x] 9.5 Update README_CN.md with Chinese documentation
