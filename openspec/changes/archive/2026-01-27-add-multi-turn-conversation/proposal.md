# Add Multi-Turn Conversation Support

## Why

Currently, the chat mode in AIQ processes each query independently without maintaining conversation context. This limits the user experience when:
- Users need to refine queries based on previous results
- Users want to ask follow-up questions
- Users need to modify previous queries incrementally

Without conversation memory, users must repeat context in each query, making the interaction less natural and efficient.

## What Changes

Add multi-turn conversation support to the chat mode:

1. **Conversation Memory**: Maintain a conversation history during the chat session
2. **Context-Aware Queries**: Send conversation history along with the current query to the LLM
3. **Session Persistence**: Save conversation sessions to disk when exiting
4. **Session Restoration**: Allow users to resume previous conversations via command-line flag

## Capabilities

### New Capability: `multi-turn-conversation`
- Maintain conversation history within a chat session
- Persist sessions to `~/.aiqconfig/session_<timestamp>.json`
- Restore sessions via `aiq -s <session-file>` or `aiq --session <session-file>`

### Modified Capability: `sql-interactive-mode` (chat mode)
- Enhanced to support conversation context
- Messages are encapsulated and sent together to LLM
- Session management integrated into chat flow

## Impact

- **User Experience**: More natural conversation flow, similar to popular chatbots
- **Efficiency**: Users can build upon previous queries without repeating context
- **Persistence**: Conversations can be saved and resumed later
- **LLM Integration**: Better context understanding through conversation history
