package sql

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aiq/aiq/internal/config"
	"github.com/aiq/aiq/internal/db"
	"github.com/aiq/aiq/internal/llm"
	"github.com/aiq/aiq/internal/source"
	"github.com/aiq/aiq/internal/ui"
)

// RunSQLMode runs the SQL interactive mode
func RunSQLMode() error {
	// Select source
	sources, err := source.LoadSources()
	if err != nil {
		return fmt.Errorf("failed to load sources: %w", err)
	}

	if len(sources) == 0 {
		return fmt.Errorf("no data sources configured. Please add a source first")
	}

	items := make([]ui.MenuItem, 0, len(sources))
	for _, s := range sources {
		label := fmt.Sprintf("%s (%s/%s:%d/%s)", s.Name, s.Type, s.Host, s.Port, s.Database)
		items = append(items, ui.MenuItem{Label: label, Value: s.Name})
	}

	sourceName, err := ui.ShowMenu("Select Data Source", items)
	if err != nil {
		return err
	}

	// Load source
	src, err := source.GetSource(sourceName)
	if err != nil {
		return fmt.Errorf("failed to load source: %w", err)
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create database connection
	conn, err := db.NewConnection(src.DSN(), string(src.Type))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close()

	// Create LLM client
	llmClient := llm.NewClient(cfg.LLM.URL, cfg.LLM.APIKey, cfg.LLM.Model)

	// Fetch schema for context
	ctx := context.Background()
	schema, err := conn.GetSchema(ctx, src.Database)
	if err != nil {
		ui.ShowWarning(fmt.Sprintf("Failed to fetch schema: %v. Continuing without schema context.", err))
		schema = &db.Schema{}
	}

	ui.ShowInfo(fmt.Sprintf("Entering SQL mode. Source: %s (%s/%s:%d/%s)", src.Name, src.Type, src.Host, src.Port, src.Database))
	ui.ShowInfo("Type your question in natural language, or 'exit' to return to main menu.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Show prompt
		fmt.Print(ui.InfoText("aiq> "))

		// Read multi-line input
		var lines []string
		hasInput := false
		
		for scanner.Scan() {
			hasInput = true
			line := strings.TrimSpace(scanner.Text())
			
			// Handle empty lines in multi-line input
			if line == "" && len(lines) == 0 {
				// Empty first line - show prompt again
				fmt.Print(ui.InfoText("aiq> "))
				continue
			}
			
			// Add non-empty lines
			if line != "" {
				lines = append(lines, line)
			}
			
			// Check if line ends with semicolon or is a command
			if line != "" && (strings.HasSuffix(line, ";") || strings.ToLower(line) == "exit" || strings.ToLower(line) == "back") {
				break
			}
		}

		// Check for scanner errors
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}

		// Handle EOF (Ctrl+D) - if no input was received, exit SQL mode
		if !hasInput && len(lines) == 0 {
			fmt.Println()
			ui.ShowInfo("Exiting SQL mode (EOF).")
			return nil
		}

		// If we have some input but scanner stopped (EOF), process what we have
		if len(lines) == 0 {
			fmt.Println()
			continue
		}

		query := strings.Join(lines, " ")
		query = strings.TrimSpace(query)

		// Handle exit command
		if strings.ToLower(query) == "exit" || strings.ToLower(query) == "back" {
			return nil
		}

		// Translate natural language to SQL
		stopLoading := ui.ShowLoading("Translating to SQL...")
		sqlQuery, err := llmClient.TranslateToSQL(ctx, query, schema.FormatSchema(), src.GetDatabaseType())
		stopLoading()

		if err != nil {
			ui.ShowError(fmt.Sprintf("Failed to translate to SQL: %v", err))
			ui.ShowInfo("Please check your LLM configuration and try again.")
			fmt.Println()
			continue
		}

		// Display translated SQL
		fmt.Println()
		ui.ShowInfo("Generated SQL:")
		fmt.Println(ui.HighlightSQL(sqlQuery))
		fmt.Println()

		// Confirm execution
		confirm, err := ui.ShowConfirm("Execute this query?")
		if err != nil {
			fmt.Println()
			continue
		}

		if !confirm {
			ui.ShowInfo("Query cancelled.")
			fmt.Println()
			continue
		}

		// Execute query
		stopLoading = ui.ShowLoading("Executing query...")
		result, err := conn.ExecuteQuery(ctx, sqlQuery)
		stopLoading()

		if err != nil {
			ui.ShowError(fmt.Sprintf("Query execution failed: %v", err))
			ui.ShowInfo("You can modify the query and try again.")
			fmt.Println()
			continue
		}

		// Display results
		fmt.Println()
		if len(result.Rows) == 0 {
			ui.ShowInfo("Query executed successfully. No rows returned.")
		} else {
			ui.ShowSuccess(fmt.Sprintf("Query executed successfully. %d row(s) returned.", len(result.Rows)))
			fmt.Println()
			ui.PrintTable(result.Columns, result.Rows)
		}
		fmt.Println()
	}
}
