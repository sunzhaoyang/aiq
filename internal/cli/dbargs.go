package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aiq/aiq/internal/source"
)

// DatabaseArgs represents parsed database connection arguments
type DatabaseArgs struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Engine   source.DatabaseType // mysql, postgresql, seekdb
}

// ParseDatabaseArgs parses and validates database CLI arguments from os.Args
// Returns nil if no database args are provided (backward compatibility)
// This function manually parses args to handle MySQL's -ppassword format (no space)
func ParseDatabaseArgs() (*DatabaseArgs, error) {
	args := make(map[string]string)
	var engine string

	// Check if any database-related flags are present
	hasDBArgs := false
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-h") || strings.HasPrefix(arg, "-u") ||
			strings.HasPrefix(arg, "-U") || strings.HasPrefix(arg, "-P") ||
			strings.HasPrefix(arg, "-p") || strings.HasPrefix(arg, "-D") ||
			strings.HasPrefix(arg, "-d") || strings.HasPrefix(arg, "-W") ||
			strings.HasPrefix(arg, "--engine") || strings.HasPrefix(arg, "-e") {
			hasDBArgs = true
			break
		}
	}

	if !hasDBArgs {
		return nil, nil // No database args, backward compatible
	}

	// Manually parse arguments to handle -ppassword format
	i := 1
	for i < len(os.Args) {
		arg := os.Args[i]

		// Handle --engine or -e
		if arg == "--engine" && i+1 < len(os.Args) {
			engine = os.Args[i+1]
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "--engine=") {
			engine = strings.TrimPrefix(arg, "--engine=")
			i++
			continue
		}
		if arg == "-e" && i+1 < len(os.Args) {
			engine = os.Args[i+1]
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "-e") && len(arg) > 2 {
			engine = arg[2:]
			i++
			continue
		}

		// Handle -h (host, shared)
		if arg == "-h" && i+1 < len(os.Args) {
			args["host"] = os.Args[i+1]
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "-h") && len(arg) > 2 {
			args["host"] = arg[2:]
			i++
			continue
		}

		// Handle -u (MySQL username)
		if arg == "-u" && i+1 < len(os.Args) {
			args["mysql_user"] = os.Args[i+1]
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "-u") && len(arg) > 2 {
			args["mysql_user"] = arg[2:]
			i++
			continue
		}

		// Handle -U (PostgreSQL username)
		if arg == "-U" && i+1 < len(os.Args) {
			args["pg_user"] = os.Args[i+1]
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "-U") && len(arg) > 2 {
			args["pg_user"] = arg[2:]
			i++
			continue
		}

		// Handle -P (MySQL port)
		if arg == "-P" && i+1 < len(os.Args) {
			args["mysql_port"] = os.Args[i+1]
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "-P") && len(arg) > 2 {
			args["mysql_port"] = arg[2:]
			i++
			continue
		}

		// Handle -p (MySQL password or PostgreSQL port - need context)
		// We'll determine this based on other flags during type detection
		if arg == "-p" && i+1 < len(os.Args) {
			nextArg := os.Args[i+1]
			// Store both possibilities, will resolve during type detection
			if _, err := strconv.Atoi(nextArg); err == nil {
				// Could be PostgreSQL port
				args["pg_port"] = nextArg
			} else {
				// Could be MySQL password
				args["mysql_password"] = nextArg
			}
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "-p") && len(arg) > 2 {
			value := arg[2:]
			// Store both possibilities, will resolve during type detection
			if _, err := strconv.Atoi(value); err == nil {
				// Could be PostgreSQL port
				args["pg_port"] = value
			} else {
				// Could be MySQL password
				args["mysql_password"] = value
			}
			i++
			continue
		}

		// Handle -D (MySQL database)
		if arg == "-D" && i+1 < len(os.Args) {
			args["mysql_db"] = os.Args[i+1]
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "-D") && len(arg) > 2 {
			args["mysql_db"] = arg[2:]
			i++
			continue
		}

		// Handle -d (PostgreSQL database)
		if arg == "-d" && i+1 < len(os.Args) {
			args["pg_db"] = os.Args[i+1]
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "-d") && len(arg) > 2 {
			args["pg_db"] = arg[2:]
			i++
			continue
		}

		// Handle -W (PostgreSQL password, non-standard)
		if arg == "-W" && i+1 < len(os.Args) {
			args["pg_password"] = os.Args[i+1]
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "-W") && len(arg) > 2 {
			args["pg_password"] = arg[2:]
			i++
			continue
		}

		// Skip unknown flags (like -s for session)
		i++
	}

	// Determine database type
	dbType := detectDatabaseType(args, engine)

	// Build unified args based on detected type
	result := &DatabaseArgs{
		Engine: dbType,
		Host:   args["host"],
	}

	if dbType == source.DatabaseTypePostgreSQL {
		// PostgreSQL mode
		result.Username = args["pg_user"]
		result.Database = args["pg_db"]
		// -p is port in PostgreSQL mode
		if portStr := args["pg_port"]; portStr != "" {
			if port, err := strconv.Atoi(portStr); err == nil {
				result.Port = port
			}
		}
		// PostgreSQL password from env var or -W flag
		if pwd := args["pg_password"]; pwd != "" {
			result.Password = pwd
		} else if pwd := os.Getenv("PGPASSWORD"); pwd != "" {
			result.Password = pwd
		}
	} else {
		// MySQL mode (default)
		result.Username = args["mysql_user"]
		result.Database = args["mysql_db"]
		// -P is port in MySQL mode
		if portStr := args["mysql_port"]; portStr != "" {
			if port, err := strconv.Atoi(portStr); err == nil {
				result.Port = port
			}
		}
		// -p is password in MySQL mode
		result.Password = args["mysql_password"]
	}

	// Validate required fields
	if err := validateDatabaseArgs(result); err != nil {
		return nil, err
	}

	return result, nil
}

// detectDatabaseType detects database type based on argument patterns
func detectDatabaseType(args map[string]string, explicitEngine string) source.DatabaseType {
	// Explicit engine override takes precedence
	if explicitEngine != "" {
		switch strings.ToLower(explicitEngine) {
		case "mysql":
			return source.DatabaseTypeMySQL
		case "postgresql", "postgres":
			return source.DatabaseTypePostgreSQL
		case "seekdb":
			return source.DatabaseTypeSeekDB
		}
	}

	// Auto-detect based on flags
	hasMySQLPattern := args["mysql_port"] != "" || args["mysql_db"] != ""
	hasPostgreSQLPattern := args["pg_user"] != "" || args["pg_db"] != ""

	if hasPostgreSQLPattern {
		return source.DatabaseTypePostgreSQL
	}
	if hasMySQLPattern {
		return source.DatabaseTypeMySQL
	}

	// Default to MySQL
	return source.DatabaseTypeMySQL
}

// validateDatabaseArgs validates that all required fields are present
func validateDatabaseArgs(args *DatabaseArgs) error {
	if args.Host == "" {
		return fmt.Errorf("host is required (use -h)")
	}
	if args.Username == "" {
		if args.Engine == source.DatabaseTypePostgreSQL {
			return fmt.Errorf("username is required (use -U)")
		}
		return fmt.Errorf("username is required (use -u)")
	}
	if args.Port <= 0 || args.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if args.Database == "" {
		if args.Engine == source.DatabaseTypePostgreSQL {
			return fmt.Errorf("database name is required (use -d)")
		}
		return fmt.Errorf("database name is required (use -D)")
	}
	if args.Password == "" {
		if args.Engine == source.DatabaseTypePostgreSQL {
			return fmt.Errorf("password is required (set PGPASSWORD environment variable or use -W)")
		}
		return fmt.Errorf("password is required (use -ppassword)")
	}
	return nil
}
