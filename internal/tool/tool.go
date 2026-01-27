package tool

import (
	"context"
	"fmt"
)

// Tool represents a callable tool/function
type Tool struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
	Execute     func(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

// ToolCall represents a request to call a tool
type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	ToolName string
	Result   interface{}
	Error    error
}

// Registry manages available tools
type Registry struct {
	tools map[string]*Tool
}

// NewRegistry creates a new tool registry
func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]*Tool),
	}
}

// Register registers a tool in the registry
func (r *Registry) Register(tool *Tool) {
	r.tools[tool.Name] = tool
}

// GetTool retrieves a tool by name
func (r *Registry) GetTool(name string) (*Tool, error) {
	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool '%s' not found", name)
	}
	return tool, nil
}

// ListTools returns all registered tools
func (r *Registry) ListTools() []*Tool {
	tools := make([]*Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// GetToolDescriptions returns descriptions of all tools for LLM context
func (r *Registry) GetToolDescriptions() string {
	if len(r.tools) == 0 {
		return "No tools available."
	}
	
	descriptions := "Available tools:\n"
	for _, tool := range r.tools {
		descriptions += fmt.Sprintf("- %s: %s\n", tool.Name, tool.Description)
	}
	return descriptions
}
