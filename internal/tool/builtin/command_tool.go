package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// CommandTool handles command execution
type CommandTool struct {
	blockedCommands    map[string]bool // Blocklist of dangerous commands
	interactiveCommands map[string]bool // Commands that require interactive input
}

// NewCommandTool creates a new command tool with security blocklist
// By default, all commands are allowed except those in the blocklist
func NewCommandTool() *CommandTool {
	// Default blocklist: dangerous commands that should be blocked
	blocked := map[string]bool{
		"rm":       true, // Remove files/directories - dangerous
		"sudo":     true, // Execute as superuser - dangerous
		"dd":       true, // Disk dump - can destroy data
		"mkfs":     true, // Make filesystem - can destroy data
		"fdisk":    true, // Disk partitioning - can destroy data
		"shutdown": true, // Shutdown system - dangerous
		"reboot":   true, // Reboot system - dangerous
		"halt":     true, // Halt system - dangerous
		"poweroff": true, // Power off system - dangerous
		"init":     true, // Init system - dangerous
		"killall":  true, // Kill all processes - dangerous
		"kill":     true, // Kill processes - dangerous (especially kill -9)
	}

	// Commands that require interactive input (cannot be executed non-interactively)
	interactive := map[string]bool{
		"mysql_secure_installation": true,
		"passwd":                   true,
		"ssh":                      true, // Without -o BatchMode=yes
		"ftp":                      true,
		"telnet":                   true,
		"less":                     true,
		"more":                     true,
		"vi":                       true,
		"vim":                      true,
		"nano":                     true,
		"emacs":                    true,
	}

	return &CommandTool{
		blockedCommands:    blocked,
		interactiveCommands: interactive,
	}
}

// SetBlockedCommands sets the blocklist of blocked commands
func (t *CommandTool) SetBlockedCommands(commands []string) {
	t.blockedCommands = make(map[string]bool)
	for _, cmd := range commands {
		t.blockedCommands[cmd] = true
	}
}

// CommandParams represents parameters for command execution
type CommandParams struct {
	Command     string `json:"command"`
	Args        []string `json:"args,omitempty"`
	WorkingDir  string `json:"working_dir,omitempty"`
	Timeout     int    `json:"timeout,omitempty"` // Timeout in seconds, default 30
}

// CommandResult represents command execution result
type CommandResult struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
}

// Execute executes a shell command
func (t *CommandTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Parse parameters
	var cmdParams CommandParams
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}
	if err := json.Unmarshal(paramsJSON, &cmdParams); err != nil {
		return nil, fmt.Errorf("failed to parse params: %w", err)
	}

	// Validate command
	if cmdParams.Command == "" {
		return nil, fmt.Errorf("command is required")
	}

	// Parse command and args, handling environment variables (e.g., "LC_ALL=C command")
	parts := strings.Fields(cmdParams.Command)
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid command")
	}

	// Extract environment variables and find the actual command name
	var cmdName string
	var cmdArgs []string
	var envVars []string
	cmdStartIdx := 0
	
	for i, part := range parts {
		// Check if this part is an environment variable assignment (VAR=value format)
		if strings.Contains(part, "=") && !strings.HasPrefix(part, "-") {
			// This is an environment variable, collect it
			envVars = append(envVars, part)
			continue
		}
		// Found the actual command
		cmdName = part
		cmdStartIdx = i
		break
	}

	if cmdName == "" {
		return nil, fmt.Errorf("no valid command found in: %s", cmdParams.Command)
	}

	// Check if command is blocked (blacklist approach: allow all except blocked)
	if t.blockedCommands[cmdName] {
		return nil, fmt.Errorf("command '%s' is blocked for security reasons. Blocked commands: %v", cmdName, t.getBlockedCommandList())
	}

	// Check if command requires sudo (commands with sudo need user interaction for password)
	if strings.Contains(cmdParams.Command, "sudo ") || strings.HasPrefix(cmdParams.Command, "sudo ") {
		return nil, fmt.Errorf("command requires sudo privileges and cannot be executed automatically. Please run this command manually in your terminal: %s", cmdParams.Command)
	}

	// Check if command requires interactive input
	if t.interactiveCommands[cmdName] {
		return nil, fmt.Errorf("command '%s' requires interactive input and cannot be executed non-interactively. Please run this command manually in your terminal, or use a non-interactive alternative if available", cmdName)
	}

	// Get command arguments (everything after the command name)
	cmdArgs = parts[cmdStartIdx+1:]
	if len(cmdParams.Args) > 0 {
		cmdArgs = cmdParams.Args
	}

	// Set timeout
	timeout := 30 * time.Second
	if cmdParams.Timeout > 0 {
		timeout = time.Duration(cmdParams.Timeout) * time.Second
	}

	// Create context with timeout
	cmdCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(cmdCtx, cmdName, cmdArgs...)

	// Set environment variables if any were specified
	if len(envVars) > 0 {
		// Start with current environment
		cmd.Env = os.Environ()
		// Add or override with specified environment variables
		for _, envVar := range envVars {
			cmd.Env = append(cmd.Env, envVar)
		}
	}

	// Set working directory
	if cmdParams.WorkingDir != "" {
		cmd.Dir = cmdParams.WorkingDir
	}

	// Execute command
	output, err := cmd.CombinedOutput()
	stdout := string(output)
	stderr := ""
	exitCode := 0

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
			stderr = exitError.Error()
		} else {
			return nil, fmt.Errorf("command execution failed: %w", err)
		}
	}

	// Split stdout and stderr if possible
	// Note: CombinedOutput doesn't separate them, so we use the error for stderr
	if err == nil {
		stdout = string(output)
		stderr = ""
	} else {
		// Try to extract stderr from error message
		stdout = string(output)
		stderr = err.Error()
	}

	result := CommandResult{
		Stdout:   stdout,
		Stderr:   stderr,
		ExitCode: exitCode,
	}

	return result, nil
}

// getBlockedCommandList returns list of blocked commands
func (t *CommandTool) getBlockedCommandList() []string {
	var list []string
	for cmd := range t.blockedCommands {
		list = append(list, cmd)
	}
	return list
}

// GetDefinition returns the tool definition for LLM
func (t *CommandTool) GetDefinition() map[string]interface{} {
	return map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "execute_command",
			"description": "Execute shell commands. Most commands are allowed, but dangerous commands (like rm, sudo, dd) are blocked for security. Use this to run scripts or system commands.",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "Command to execute (e.g., 'ls -la')",
					},
					"args": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "string"},
						"description": "Command arguments",
					},
					"working_dir": map[string]interface{}{
						"type":        "string",
						"description": "Working directory for command execution",
					},
					"timeout": map[string]interface{}{
						"type":        "integer",
						"description": "Timeout in seconds (default: 30)",
					},
				},
				"required": []string{"command"},
			},
		},
	}
}
