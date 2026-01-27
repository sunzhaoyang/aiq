package tool

import (
	"encoding/json"
)

// GetToolDefinitions returns the list of available tools as LLM Function definitions
func GetToolDefinitions() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "execute_sql",
				"description": "Execute a SQL query against the database and return the results. Use this when the user wants to query the database.",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"sql": map[string]interface{}{
							"type":        "string",
							"description": "The SQL query to execute",
						},
					},
					"required": []string{"sql"},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "render_table",
				"description": "Display query results as a formatted table. Use this to show data in a tabular format.",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"columns": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Column names",
						},
						"rows": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
							"description": "Row data, each row is an array of string values",
						},
					},
					"required": []string{"columns", "rows"},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "render_chart",
				"description": "Display query results as a chart (bar, line, pie, or scatter). Use this when the user wants to visualize data.",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"columns": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Column names",
						},
						"rows": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
							"description": "Row data, each row is an array of string values",
						},
						"chart_type": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"bar", "line", "pie", "scatter"},
							"description": "Type of chart to render",
						},
					},
					"required": []string{"columns", "rows", "chart_type"},
				},
			},
		},
	}
}

// ConvertToLLMFunctions converts tool definitions to LLM Function format
func ConvertToLLMFunctions() ([]map[string]interface{}, error) {
	return GetToolDefinitions(), nil
}

// ToolCallResult represents the result of a tool execution for LLM
type ToolCallResult struct {
	ToolCallID string          `json:"tool_call_id"`
	Role      string          `json:"role"` // "tool"
	Name      string          `json:"name"`
	Content   json.RawMessage `json:"content"`
}
