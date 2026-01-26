## Why

AIQ (AI Query) aims to bridge the gap between natural language questions and precise SQL queries for databases like MySQL and SeekDB. Currently, users need deep SQL knowledge to query databases effectively. This change establishes the foundational architecture and minimal demo for AIQ, enabling users to interact with databases through natural language via an intuitive CLI tool. Starting with MySQL as the initial target, we'll build a solid foundation that can be extended to support additional databases and features.

## What Changes

- **New CLI Application**: Create a Go-based CLI tool (`aiq`) that launches an interactive command-line interface
- **First-Run Configuration Wizard**: Implement guided setup flow for LLM configuration (URL, API Key) on first launch, saving to configuration file
- **Interactive Menu System**: Build a menu-driven interface with four core functions:
  - `config`: Manage LLM and tool configuration settings
  - `source`: Manage database connection configurations
  - `sql`: Enter SQL interactive mode (requires source selection)
  - `exit`: Exit the application
- **Project Structure**: Initialize project architecture following best practices from popular open-source Go CLI projects
- **User Experience Enhancements**: Implement smooth transitions and visual feedback for LLM interactions and database operations
- **Color Scheme**: Apply industry-standard color palette with syntax highlighting for keywords and output

## Capabilities

### New Capabilities

- `cli-application`: Core CLI application framework with interactive menu system and command routing
- `configuration-management`: Configuration file management for LLM settings (URL, API Key) with first-run wizard
- `data-source-management`: Database connection configuration management (add, list, select, remove sources)
- `sql-interactive-mode`: SQL query interface that integrates with selected data source and LLM for natural language to SQL translation
- `user-experience`: Visual enhancements including loading indicators, transitions, and color-coded output

### Modified Capabilities

<!-- No existing capabilities to modify - this is a new project -->

## Impact

- **New Dependencies**: Go standard library and CLI frameworks (e.g., `cobra`, `promptui`), database drivers (MySQL), LLM client libraries
- **Project Structure**: Establishes directory structure, configuration file locations, and code organization patterns
- **Configuration Files**: Creates user configuration directory and file format (e.g., `~/.aiq/config.yaml` or `~/.config/aiq/config.yaml`)
- **Build System**: Sets up Go module, build scripts, and installation process
- **Documentation**: Requires README, installation instructions, and usage documentation
