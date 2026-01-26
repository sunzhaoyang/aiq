# AIQ

AIQ (AI Query): An intelligent SQL client that translates your natural language questions into precise SQL queries for MySQL, SeekDB, and other databases.

## Features

- 🎯 **Natural Language to SQL**: Ask questions in plain English, get precise SQL queries
- 🔌 **Multiple Database Support**: Currently supports MySQL, with more databases coming soon
- ⚙️ **Easy Configuration**: Guided setup wizard for LLM and database connections
- 🎨 **Beautiful CLI**: Interactive menus with smooth transitions and color-coded output
- 🔒 **Secure**: Local configuration storage, no cloud sync required

## Installation

### Prerequisites

- Go 1.21 or later
- MySQL database (for database queries)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/aiq/aiq.git
cd aiq

# Build the binary
go build -o aiq cmd/aiq/main.go

# Install (optional)
sudo mv aiq /usr/local/bin/
```

## Quick Start

1. **Run AIQ**:
   ```bash
   aiq
   ```

2. **First Run Setup**:
   - On first launch, you'll be guided through LLM configuration
   - Enter your LLM API URL and API Key
   - Configuration is saved to `~/.config/aiq/config.yaml`

3. **Add a Database Source**:
   - Select `source` from the main menu
   - Choose `add` to add a new MySQL connection
   - Enter connection details (host, port, database, username, password)

4. **Query Your Database**:
   - Select `sql` from the main menu
   - Enter your question in natural language
   - Review the generated SQL and confirm execution
   - View results in a formatted table

## Usage

### Main Menu

- **config**: Manage LLM and tool configuration
- **source**: Manage database connection configurations
- **sql**: Enter SQL interactive mode (requires a selected source)
- **exit**: Exit the application

### Configuration Management

- View current LLM configuration
- Update LLM API URL
- Update LLM API Key

### Data Source Management

- Add new database connections
- List all configured sources
- Select active source for queries
- Remove database connections

## Configuration

Configuration files are stored in `~/.config/aiq/`:

- `config.yaml`: LLM configuration (URL, API Key)
- `sources.yaml`: Database connection configurations

## Development

### Project Structure

```
cmd/aiq/          # Main entry point
internal/
  cli/            # CLI commands and menu system
  config/         # Configuration management
  source/         # Data source management
  sql/            # SQL interactive mode
  llm/            # LLM client integration
  db/             # Database connection and query execution
  ui/             # UI components (prompts, colors, loading)
```

### Building

```bash
go build -o aiq cmd/aiq/main.go
```

### Running Tests

```bash
go test ./...
```

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
