package tool

import "github.com/aiq/aiq/internal/llm"

// GetLLMFunctions converts tool definitions to LLM Function format
func GetLLMFunctions() []llm.Function {
	return []llm.Function{
		{
			Name:        "execute_sql",
			Description: "Execute a SQL query against the database and return the results. Use this when the user wants to query the database.",
			Parameters: map[string]interface{}{
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
		{
			Name:        "render_table",
			Description: "Format query results as a table string. Use this when you want to show data in a tabular format.",
			Parameters: map[string]interface{}{
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
		{
			Name:        "render_chart",
			Description: "Format query results as a chart string (bar, line, pie, or scatter). Use this when the user wants to visualize data.",
			Parameters: map[string]interface{}{
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
	}
}
