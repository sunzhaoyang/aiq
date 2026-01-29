package cli

import (
	"fmt"
	"strconv"

	"github.com/aiq/aiq/internal/db"
	"github.com/aiq/aiq/internal/source"
	"github.com/aiq/aiq/internal/ui"
)

var activeSourceName string

// GetActiveSource returns the currently active source name
func GetActiveSource() string {
	return activeSourceName
}

// SetActiveSource sets the active source name
func SetActiveSource(name string) {
	activeSourceName = name
}

// RunSourceMenu runs the data source management menu
func RunSourceMenu() error {
	for {
		items := []ui.MenuItem{
			{Label: "add     - Add a new data source", Value: "add"},
			{Label: "list    - List all data sources", Value: "list"},
			{Label: "edit    - Edit a data source", Value: "edit"},
			{Label: "remove  - Remove a data source", Value: "remove"},
			{Label: "back    - Back to main menu", Value: "back"},
		}

		choice, err := ui.ShowMenu("Data Sources", items)
		if err != nil {
			return err
		}

		switch choice {
		case "add":
			if err := addSource(); err != nil {
				ui.ShowError(err.Error())
			} else {
				ui.ShowSuccess("Data source added successfully!")
			}
		case "list":
			if err := listSources(); err != nil {
				ui.ShowError(err.Error())
			}
		case "edit":
			if err := editSource(); err != nil {
				ui.ShowError(err.Error())
			} else {
				ui.ShowSuccess("Data source updated successfully!")
			}
		case "remove":
			if err := removeSource(); err != nil {
				ui.ShowError(err.Error())
			} else {
				ui.ShowSuccess("Data source removed successfully!")
			}
		case "back":
			return nil
		}

		fmt.Println()
	}
}

func addSource() error {
	src := &source.Source{}

	// Select database type first
	fmt.Println()
	fmt.Println("Select Database Type:")
	typeItems := []ui.MenuItem{
		{Label: "seekdb", Value: "seekdb"},
		{Label: "MySQL", Value: "mysql"},
		{Label: "PostgreSQL", Value: "postgresql"},
	}
	
	dbType, err := ui.ShowMenu("Database Type", typeItems)
	if err != nil {
		return fmt.Errorf("failed to select database type: %w", err)
	}
	src.Type = source.DatabaseType(dbType)

	name, err := ui.ShowInput("Enter source name", "")
	if err != nil {
		return fmt.Errorf("failed to get source name: %w", err)
	}
	src.Name = name

	// Set default port based on database type
	defaultPort := "3306"
	if src.Type == source.DatabaseTypePostgreSQL {
		defaultPort = "5432"
	} else if src.Type == source.DatabaseTypeSeekDB {
		defaultPort = "3306" // Adjust as needed
	}

	host, err := ui.ShowInput("Enter host", "localhost")
	if err != nil {
		return fmt.Errorf("failed to get host: %w", err)
	}
	src.Host = host

	portStr, err := ui.ShowInput("Enter port", defaultPort)
	if err != nil {
		return fmt.Errorf("failed to get port: %w", err)
	}

	port, err := source.ParsePort(portStr)
	if err != nil {
		return err
	}
	src.Port = port

	database, err := ui.ShowInput("Enter database name", "")
	if err != nil {
		return fmt.Errorf("failed to get database name: %w", err)
	}
	src.Database = database

	username, err := ui.ShowInput("Enter username", "")
	if err != nil {
		return fmt.Errorf("failed to get username: %w", err)
	}
	src.Username = username

	password, err := ui.ShowPassword("Enter password")
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}
	src.Password = password

	if err := source.Validate(src); err != nil {
		return err
	}

	// Optional: Test connection
	test, err := ui.ShowConfirm("Test connection before saving?")
	if err == nil && test {
		ui.ShowInfo("Testing connection...")
		if err := db.TestConnection(src.DSN(), string(src.Type)); err != nil {
			ui.ShowWarning(fmt.Sprintf("Connection test failed: %v", err))
			proceed, _ := ui.ShowConfirm("Save anyway?")
			if !proceed {
				return fmt.Errorf("connection test failed, not saving")
			}
		} else {
			ui.ShowSuccess("Connection test successful!")
		}
	}

	return source.AddSource(src)
}

func listSources() error {
	sources, err := source.LoadSources()
	if err != nil {
		return err
	}

	if len(sources) == 0 {
		ui.ShowInfo("No data sources configured.")
		ui.ShowInfo("Use 'Add data source' to add your first source.")
		return nil
	}

	fmt.Println()
	ui.ShowInfo("Configured Data Sources:")
	fmt.Println()

	headers := []string{"Name", "Type", "Host", "Port", "Database", "Username"}
	rows := make([][]string, 0, len(sources))

	for _, s := range sources {
		rows = append(rows, []string{
			s.Name,
			string(s.Type),
			s.Host,
			strconv.Itoa(s.Port),
			s.Database,
			s.Username,
		})
	}

	ui.PrintTable(headers, rows)
	fmt.Println()

	return nil
}

func removeSource() error {
	sources, err := source.LoadSources()
	if err != nil {
		return err
	}

	if len(sources) == 0 {
		return fmt.Errorf("no data sources configured")
	}

	items := make([]ui.MenuItem, 0, len(sources))
	for _, s := range sources {
		label := fmt.Sprintf("%s (%s/%s:%d/%s)", s.Name, s.Type, s.Host, s.Port, s.Database)
		items = append(items, ui.MenuItem{Label: label, Value: s.Name})
	}

	selected, err := ui.ShowMenu("Select Source to Remove", items)
	if err != nil {
		return err
	}

	confirm, err := ui.ShowConfirm(fmt.Sprintf("Are you sure you want to remove '%s'?", selected))
	if err != nil {
		return err
	}

	if !confirm {
		ui.ShowInfo("Removal cancelled.")
		return nil
	}

	return source.RemoveSource(selected)
}

func editSource() error {
	sources, err := source.LoadSources()
	if err != nil {
		return err
	}

	if len(sources) == 0 {
		return fmt.Errorf("no data sources configured")
	}

	items := make([]ui.MenuItem, 0, len(sources))
	for _, s := range sources {
		label := fmt.Sprintf("%s (%s/%s:%d/%s)", s.Name, s.Type, s.Host, s.Port, s.Database)
		items = append(items, ui.MenuItem{Label: label, Value: s.Name})
	}

	selected, err := ui.ShowMenu("Select Source to Edit", items)
	if err != nil {
		return err
	}

	// Load the selected source
	oldSource, err := source.GetSource(selected)
	if err != nil {
		return err
	}

	// Create updated source with current values as defaults
	updated := &source.Source{
		Type:     oldSource.Type,
		Name:     oldSource.Name,
		Host:     oldSource.Host,
		Port:     oldSource.Port,
		Database: oldSource.Database,
		Username: oldSource.Username,
		Password: oldSource.Password,
	}

	// Prompt for all fields with current values as defaults
	name, err := ui.ShowInput("Enter source name", oldSource.Name)
	if err != nil {
		return fmt.Errorf("failed to get source name: %w", err)
	}
	updated.Name = name

	host, err := ui.ShowInput("Enter host", oldSource.Host)
	if err != nil {
		return fmt.Errorf("failed to get host: %w", err)
	}
	updated.Host = host

	portStr, err := ui.ShowInput("Enter port", strconv.Itoa(oldSource.Port))
	if err != nil {
		return fmt.Errorf("failed to get port: %w", err)
	}

	port, err := source.ParsePort(portStr)
	if err != nil {
		return err
	}
	updated.Port = port

	database, err := ui.ShowInput("Enter database name", oldSource.Database)
	if err != nil {
		return fmt.Errorf("failed to get database name: %w", err)
	}
	updated.Database = database

	username, err := ui.ShowInput("Enter username", oldSource.Username)
	if err != nil {
		return fmt.Errorf("failed to get username: %w", err)
	}
	updated.Username = username

	// For password, ask if user wants to change it
	changePassword, err := ui.ShowConfirm("Change password?")
	if err != nil {
		return fmt.Errorf("failed to get password change confirmation: %w", err)
	}

	if changePassword {
		password, err := ui.ShowPassword("Enter password")
		if err != nil {
			return fmt.Errorf("failed to get password: %w", err)
		}
		updated.Password = password
	}

	if err := source.Validate(updated); err != nil {
		return err
	}

	// Optional: Test connection
	test, err := ui.ShowConfirm("Test connection before saving?")
	if err == nil && test {
		ui.ShowInfo("Testing connection...")
		if err := db.TestConnection(updated.DSN(), string(updated.Type)); err != nil {
			ui.ShowWarning(fmt.Sprintf("Connection test failed: %v", err))
			proceed, _ := ui.ShowConfirm("Save anyway?")
			if !proceed {
				return fmt.Errorf("connection test failed, not saving")
			}
		} else {
			ui.ShowSuccess("Connection test successful!")
		}
	}

	return source.UpdateSource(selected, updated)
}
