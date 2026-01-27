package sql

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aiq/aiq/internal/db"
	"github.com/aiq/aiq/internal/llm"
	"github.com/aiq/aiq/internal/tool"
	"github.com/aiq/aiq/internal/ui"
)

// ToolHandler handles tool execution and manages tool calling loop
type ToolHandler struct {
	conn *db.Connection
}

// NewToolHandler creates a new tool handler
func NewToolHandler(conn *db.Connection) *ToolHandler {
	return &ToolHandler{
		conn: conn,
	}
}

// ExecuteTool executes a tool call and returns the result
func (h *ToolHandler) ExecuteTool(ctx context.Context, toolCall llm.ToolCall) (json.RawMessage, error) {
	toolName := toolCall.Function.Name
	
	// Parse arguments from JSON string
	args, err := toolCall.ParseArguments()
	if err != nil {
		return nil, fmt.Errorf("failed to parse tool arguments: %w", err)
	}

	switch toolName {
	case "execute_sql":
		sql, ok := args["sql"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid sql parameter")
		}
		
		// Execute SQL - this does NOT print anything, only returns data
		result, err := tool.ExecuteSQL(ctx, h.conn, sql)
		if err != nil {
			errorMsg := fmt.Sprintf(`{"error": "%s"}`, err.Error())
			return json.RawMessage(errorMsg), nil
		}

		// Convert result to JSON and return to LLM
		// LLM will decide how to display this (via render_table or text description)
		resultJSON := map[string]interface{}{
			"columns":   result.Columns,
			"rows":      result.Rows,
			"row_count": len(result.Rows),
		}
		jsonData, err := json.Marshal(resultJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal result: %w", err)
		}
		return json.RawMessage(jsonData), nil

	case "render_table":
		columnsInterface, ok := args["columns"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid columns parameter")
		}
		rowsInterface, ok := args["rows"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid rows parameter")
		}

		// Convert to string slices
		columns := make([]string, len(columnsInterface))
		for i, col := range columnsInterface {
			columns[i] = fmt.Sprintf("%v", col)
		}

		rows := make([][]string, len(rowsInterface))
		for i, rowInterface := range rowsInterface {
			rowArray, ok := rowInterface.([]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid row format")
			}
			rows[i] = make([]string, len(rowArray))
			for j, val := range rowArray {
				rows[i][j] = fmt.Sprintf("%v", val)
			}
		}

		// Format the table as a string (do not print)
		tableOutput, err := tool.RenderTableString(columns, rows)
		if err != nil {
			errorMsg := fmt.Sprintf(`{"error": "%s"}`, err.Error())
			return json.RawMessage(errorMsg), nil
		}

		// Return formatted table to LLM
		resultJSON := map[string]interface{}{
			"status":    "success",
			"format":    "table",
			"output":    tableOutput,
			"row_count": len(rows),
		}
		jsonData, err := json.Marshal(resultJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal result: %w", err)
		}
		return json.RawMessage(jsonData), nil

	case "render_chart":
		columnsInterface, ok := args["columns"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid columns parameter")
		}
		rowsInterface, ok := args["rows"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid rows parameter")
		}
		chartTypeStr, ok := args["chart_type"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid chart_type parameter")
		}

		// Convert to string slices
		columns := make([]string, len(columnsInterface))
		for i, col := range columnsInterface {
			columns[i] = fmt.Sprintf("%v", col)
		}

		rows := make([][]string, len(rowsInterface))
		for i, rowInterface := range rowsInterface {
			rowArray, ok := rowInterface.([]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid row format")
			}
			rows[i] = make([]string, len(rowArray))
			for j, val := range rowArray {
				rows[i][j] = fmt.Sprintf("%v", val)
			}
		}

		// Create QueryResult
		result := &db.QueryResult{
			Columns: columns,
			Rows:    rows,
		}

		chartOutput, err := tool.RenderChartString(result, chartTypeStr)
		if err != nil {
			errorMsg := fmt.Sprintf(`{"error": "%s"}`, err.Error())
			return json.RawMessage(errorMsg), nil
		}

		// Return formatted chart to LLM
		resultJSON := map[string]interface{}{
			"status":     "success",
			"format":     "chart",
			"output":     chartOutput,
			"chart_type": chartTypeStr,
			"row_count":  len(result.Rows),
		}
		jsonData, err := json.Marshal(resultJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal result: %w", err)
		}
		return json.RawMessage(jsonData), nil

	default:
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}
}

// HandleToolCallLoop handles the complete tool calling loop
// Returns the final response content and any query result
func (h *ToolHandler) HandleToolCallLoop(ctx context.Context, llmClient *llm.Client, userInput string, schemaContext string, databaseType string, conversationHistory []llm.ChatMessage, tools []llm.Function) (string, *db.QueryResult, error) {
	// Build initial messages
	messages := []interface{}{
		llm.ChatMessage{
			Role: "system",
			Content: fmt.Sprintf(`You are a helpful AI assistant for database queries. You can have natural conversations with users, or help them query databases using available tools.

IMPORTANT CONTEXT:
- Database engine type: %s (this is the database system type like MySQL/PostgreSQL/seekdb, NOT a schema or database name)
- Current database connection information is provided below

Database connection and schema information:
%s

CRITICAL RULES FOR SQL GENERATION:
1. Database type "%s" is the ENGINE TYPE (like MySQL, PostgreSQL, seekdb), NOT a database name or schema name
2. When user asks "show tables" or "list tables", use execute_sql tool with: SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() OR use SHOW TABLES (for MySQL/seekdb)
3. NEVER use the database engine type (like "%s") as a schema name in WHERE table_schema = '%s'
4. Always use the actual database name from the connection context, or use DATABASE() function to get current database
5. For MySQL/seekdb: Use SHOW TABLES; or SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE();
6. For PostgreSQL: Use SELECT tablename FROM pg_tables WHERE schemaname = 'public'; or SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';

Use the available tools to help users query the database. Follow this flow:
1. Use execute_sql to query the database (this returns data, does NOT display it)
2. Decide how to present results:
   - Default: use render_table to format results as a table string when there are multiple rows/columns.
   - Use render_chart to format results as a chart string only when the user explicitly asks for a chart.
   - Summarize in text only when the user explicitly requests a summary or when data is trivial (e.g., a single value).
3. If you call render_table or render_chart, include the returned output string in your final response.

Remember: tools return data; you decide what to show in the final response.`, databaseType, schemaContext, databaseType, databaseType, databaseType),
		},
	}

	// Add conversation history
	for _, msg := range conversationHistory {
		messages = append(messages, msg)
	}

	// Add user input
	messages = append(messages, llm.ChatMessage{
		Role:    "user",
		Content: userInput,
	})

	var lastQueryResult *db.QueryResult
	var lastFormattedOutput string
	maxIterations := 10 // Prevent infinite loops

	for i := 0; i < maxIterations; i++ {
		// Call LLM
		response, err := llmClient.ChatWithTools(ctx, messages, tools)
		if err != nil {
			return "", nil, fmt.Errorf("LLM call failed: %w", err)
		}

		if len(response.Choices) == 0 {
			return "", nil, fmt.Errorf("no choices in response")
		}

		choice := response.Choices[0]
		message := choice.Message

		// Add assistant message to history
		// If there are tool_calls, content might be empty or null
		assistantMsg := map[string]interface{}{
			"role": "assistant",
		}
		// Only add content if it's not empty
		if message.Content != "" {
			assistantMsg["content"] = message.Content
		}
		if len(message.ToolCalls) > 0 {
			assistantMsg["tool_calls"] = message.ToolCalls
		}
		messages = append(messages, assistantMsg)

		// If no tool calls, return the final response
		if len(message.ToolCalls) == 0 {
			if message.Content == "" {
				if lastFormattedOutput != "" {
					return lastFormattedOutput, lastQueryResult, nil
				}
				return "", lastQueryResult, fmt.Errorf("empty response from LLM")
			}
			return message.Content, lastQueryResult, nil
		}

		// Execute tool calls (show status only)
		for _, toolCall := range message.ToolCalls {
			// For execute_sql, show SQL and ask for confirmation before executing
			if toolCall.Function.Name == "execute_sql" {
				args, err := toolCall.ParseArguments()
				if err != nil {
					ui.ShowError(fmt.Sprintf("Tool [%s] failed: %v", toolCall.Function.Name, err))
					errorMsg := fmt.Sprintf(`{"error": "%s"}`, err.Error())
					toolResult := json.RawMessage(errorMsg)
					toolMsg := map[string]interface{}{
						"role":        "tool",
						"content":     string(toolResult),
						"tool_call_id": toolCall.ID,
					}
					messages = append(messages, toolMsg)
					continue
				}

				sql, ok := args["sql"].(string)
				if !ok {
					err := fmt.Errorf("invalid sql parameter")
					ui.ShowError(fmt.Sprintf("Tool [%s] failed: %v", toolCall.Function.Name, err))
					errorMsg := fmt.Sprintf(`{"error": "%s"}`, err.Error())
					toolResult := json.RawMessage(errorMsg)
					toolMsg := map[string]interface{}{
						"role":        "tool",
						"content":     string(toolResult),
						"tool_call_id": toolCall.ID,
					}
					messages = append(messages, toolMsg)
					continue
				}

				fmt.Println()
				ui.ShowInfo("Generated SQL:")
				fmt.Println(ui.HighlightSQL(sql))
				fmt.Println()

				confirm, err := ui.ShowConfirm("Execute this query?")
				if err != nil {
					fmt.Println()
					// Treat as cancelled
					ui.ShowWarning("Query execution cancelled.")
					toolResult := json.RawMessage(`{"status":"cancelled","message":"query execution cancelled by user"}`)
					toolMsg := map[string]interface{}{
						"role":        "tool",
						"content":     string(toolResult),
						"tool_call_id": toolCall.ID,
					}
					messages = append(messages, toolMsg)
					continue
				}
				if !confirm {
					ui.ShowWarning("Query execution cancelled.")
					toolResult := json.RawMessage(`{"status":"cancelled","message":"query execution cancelled by user"}`)
					toolMsg := map[string]interface{}{
						"role":        "tool",
						"content":     string(toolResult),
						"tool_call_id": toolCall.ID,
					}
					messages = append(messages, toolMsg)
					continue
				}
			}

			ui.ShowInfo(fmt.Sprintf("Calling tool [%s]...", toolCall.Function.Name))
			toolResult, err := h.ExecuteTool(ctx, toolCall)
			if err != nil {
				errorMsg := fmt.Sprintf(`{"error": "%s"}`, err.Error())
				toolResult = json.RawMessage(errorMsg)
				ui.ShowError(fmt.Sprintf("Tool [%s] failed: %v", toolCall.Function.Name, err))
			} else {
				ui.ShowSuccess(fmt.Sprintf("Tool [%s] executed successfully", toolCall.Function.Name))
			}

			// Store query result if it's execute_sql
			if toolCall.Function.Name == "execute_sql" && err == nil {
				var resultData map[string]interface{}
				if err := json.Unmarshal(toolResult, &resultData); err == nil {
					if columns, ok := resultData["columns"].([]interface{}); ok {
						if rows, ok := resultData["rows"].([]interface{}); ok {
							// Convert to QueryResult
							cols := make([]string, len(columns))
							for i, col := range columns {
								cols[i] = fmt.Sprintf("%v", col)
							}
							rowsData := make([][]string, len(rows))
							for i, rowInterface := range rows {
								rowArray, ok := rowInterface.([]interface{})
								if ok {
									rowsData[i] = make([]string, len(rowArray))
									for j, val := range rowArray {
										rowsData[i][j] = fmt.Sprintf("%v", val)
									}
								}
							}
							lastQueryResult = &db.QueryResult{
								Columns: cols,
								Rows:    rowsData,
							}
						}
					}
				}
			}

			// Store last formatted output if tool returned it
			if err == nil {
				var outputData map[string]interface{}
				if jsonErr := json.Unmarshal(toolResult, &outputData); jsonErr == nil {
					if output, ok := outputData["output"].(string); ok && output != "" {
						lastFormattedOutput = output
					}
				}
			}

			// Step 5: Add tool result message back to LLM
			// LLM will decide next step based on tool results
			toolMsg := map[string]interface{}{
				"role":        "tool",
				"content":     string(toolResult), // Convert json.RawMessage to string
				"tool_call_id": toolCall.ID,
			}
			messages = append(messages, toolMsg)
		}
		// Loop continues: LLM will process tool results and decide next action
	}

	return "", nil, fmt.Errorf("max iterations reached")
}
