# Alternative Design: LLM-Assisted Risk Assessment

## Option C: LLM Provides Risk Level in Tool Call (Recommended)

### Core Concept
- LLM adds optional `risk_level` field in arguments when calling tools
- Code layer prioritizes LLM-provided risk_level, if not provided makes conservative judgment
- Combines code layer basic rules with LLM's intelligent judgment

### Implementation

#### 1. Tool Definition Extension
Add optional `risk_level` field to parameters of all tools:

```json
{
  "type": "object",
  "properties": {
    "sql": {...},
    "risk_level": {
      "type": "string",
      "enum": ["low", "medium", "high"],
      "description": "Optional: Risk level assessment for this operation. 'low' = safe to execute automatically, 'medium'/'high' = requires user confirmation. If not provided, system will assess risk conservatively."
    }
  }
}
```

#### 2. Prompt Guidance
Guide LLM in prompts:
- For obviously safe operations (SELECT, SHOW, ls, read), set `risk_level: "low"`
- For obviously dangerous operations (DROP, TRUNCATE, rm, write), set `risk_level: "high"`
- For uncertain operations (init, reboot, custom scripts), set `risk_level: "high"` or ask user

#### 3. Code Layer Processing Logic
```go
func assessRisk(toolCall, args) RiskDecision {
    // 1. Check LLM-provided risk_level
    if riskLevel, ok := args["risk_level"].(string); ok {
        switch riskLevel {
        case "low":
            return RiskLow  // Execute directly
        case "medium", "high":
            return RiskHigh  // Require confirmation
        }
    }
    
    // 2. LLM didn't provide risk_level, code layer makes conservative judgment
    // Only execute directly for clearly safe operations (whitelist)
    if isWhitelisted(toolCall, args) {
        return RiskLow
    }
    
    // 3. Other cases default to requiring confirmation (conservative strategy)
    return RiskHigh
}
```

### Advantages
- ✅ LLM can intelligently judge risk of unknown commands
- ✅ Code layer has basic rules guarantee (whitelist)
- ✅ No need to exhaustively list all dangerous commands
- ✅ When LLM judgment is wrong, code layer conservative strategy ensures safety
- ✅ Conforms to Agent flow: LLM makes decisions, code layer executes

### Disadvantages
- ⚠️ Need to modify all tool definitions
- ⚠️ LLM may not always provide risk_level (needs fallback)

---

## Option D: Two-Phase Tool Call (LLM Assesses Risk First)

### Core Concept
- LLM calls an `assess_risk` tool before calling actual tool
- `assess_risk` tool returns risk level
- Code layer decides whether to confirm based on risk level

### Implementation

#### 1. New assess_risk Tool
```go
{
  "name": "assess_risk",
  "description": "Assess the risk level of a tool operation before execution",
  "parameters": {
    "tool_name": "string",
    "tool_args": "object"
  }
}
```

#### 2. Flow
```
1. LLM calls assess_risk tool
2. assess_risk returns risk_level
3. Code layer decides based on risk_level:
   - low → Execute original tool directly
   - high → Require confirmation before execution
4. LLM calls actual tool
```

### Advantages
- ✅ Separation of concerns: Risk assessment and execution separated
- ✅ LLM can fully assess risk

### Disadvantages
- ⚠️ Adds one LLM call, affects performance
- ⚠️ Increases complexity
- ⚠️ May be skipped by LLM (doesn't call assess_risk)

---

## Option E: Code Layer Asynchronously Asks LLM (Not Recommended)

### Core Concept
- Code layer makes basic judgment first
- If uncertain, asynchronously call LLM for quick risk assessment
- Decide whether to confirm based on LLM assessment result

### Disadvantages
- ❌ Need to asynchronously call LLM, affects performance
- ❌ Adds latency
- ❌ Complex implementation

---

## Option F: LLM Returns Text to Ask User (Conforms to Agent Flow)

### Core Concept
- LLM can return text to ask user before calling tool if uncertain about risk
- After user confirms, LLM calls tool
- Code layer only handles clear cases (whitelist executes directly, others require confirmation)

### Implementation

#### Prompt Guidance
```
<RISK_ASSESSMENT>
When calling tools, assess the risk level:
- **Low-risk operations** (SELECT, SHOW, ls, read, GET): Call tool directly with risk_level="low"
- **High-risk operations** (DROP, TRUNCATE, rm, write, DELETE): Call tool with risk_level="high" or ask user first
- **Uncertain operations** (init, reboot, custom scripts): 
  - Option 1: Return text asking user for confirmation before calling tool
  - Option 2: Call tool with risk_level="high" (system will ask for confirmation)
  
If you're uncertain about an operation's risk, you can:
1. Return text asking user: "This operation (init system) may be risky. Proceed?"
2. Or call tool with risk_level="high" and let system handle confirmation
</RISK_ASSESSMENT>
```

#### Code Layer Processing
```go
// If LLM returns text asking user (no tool_calls)
if len(message.ToolCalls) == 0 && message.Content != "" {
    // Check if it's a risk confirmation question
    if isRiskConfirmationQuestion(message.Content) {
        // Display to user, wait for user confirmation
        // After user confirms, LLM can continue calling tool
        return message.Content, nil, nil
    }
}
```

### Advantages
- ✅ Conforms to Agent flow: LLM can make decisions
- ✅ Flexible: LLM can choose to ask or call directly
- ✅ No need to modify tool definitions
- ✅ Code layer remains simple

### Disadvantages
- ⚠️ May violate "must call tool" principle (but can allow asking when uncertain)

---

## Recommended Option: Option C (LLM Provides Risk Level in Tool Call)

Combining advantages of Option C and Option F:
1. **Primary Method**: LLM provides `risk_level` in tool call arguments
2. **Alternative Method**: If LLM is uncertain, can return text to ask user (allow this exception)
3. **Code Layer Guarantee**: Code layer has basic whitelist, can work even if LLM doesn't provide risk_level

### Final Flow

```
1. LLM decides to call tool
   ↓
2. LLM adds risk_level in arguments (optional)
   ↓
3. Code layer assessment:
   - If LLM provided risk_level="low" → Execute directly
   - If LLM provided risk_level="high" → Require confirmation
   - If LLM didn't provide risk_level:
     - Code layer whitelist → Execute directly
     - Others → Require confirmation (conservative strategy)
   ↓
4. If confirmation needed:
   - Display operation content
   - Ask user for confirmation
   - Execute after user confirms
```

### Special Case Handling

If LLM is uncertain about operation risk, can:
- **Method 1**: Set `risk_level="high"` when calling tool, let system require confirmation
- **Method 2**: Return text asking user: "This operation (init system) may be risky. Should I proceed?"
  - After user confirms, LLM calls tool
  - This is part of Agent flow, allows LLM to ask user when uncertain
