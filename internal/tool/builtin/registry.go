package builtin

import (
	"context"
	"fmt"

	"github.com/aiq/aiq/internal/db"
)

// GetBuiltinToolDefinitions returns LLM function definitions for all built-in tools
func GetBuiltinToolDefinitions(dbConn *db.Connection) []map[string]interface{} {
	definitions := []map[string]interface{}{}

	// HTTP tool
	httpTool := NewHTTPTool()
	definitions = append(definitions, httpTool.GetDefinition())

	// Command tool
	commandTool := NewCommandTool()
	definitions = append(definitions, commandTool.GetDefinition())

	// File tool
	fileTool, err := NewFileTool()
	if err == nil {
		definitions = append(definitions, fileTool.GetDefinition())
	}

	// Note: Database query tool is already available as "execute_sql" in the main tool set
	// No need to add a duplicate "query_database" tool definition

	return definitions
}

// ExecuteBuiltinTool executes a built-in tool by name
func ExecuteBuiltinTool(ctx context.Context, name string, params map[string]interface{}, dbConn *db.Connection) (interface{}, error) {
	return ExecuteBuiltinToolWithCallback(ctx, name, params, dbConn, nil)
}

// ExecuteBuiltinToolWithCallback executes a built-in tool with optional output callback
func ExecuteBuiltinToolWithCallback(ctx context.Context, name string, params map[string]interface{}, dbConn *db.Connection, callback OutputCallback) (interface{}, error) {
	switch name {
	case "http_request":
		httpTool := NewHTTPTool()
		return httpTool.Execute(ctx, params)
	case "execute_command":
		commandTool := NewCommandTool()
		return commandTool.ExecuteWithCallback(ctx, params, callback)
	case "file_operations":
		fileTool, err := NewFileTool()
		if err != nil {
			return nil, err
		}
		return fileTool.Execute(ctx, params)
	// Note: Database queries use "execute_sql" tool, handled separately
	default:
		return nil, fmt.Errorf("unknown built-in tool: %s", name)
	}
}
