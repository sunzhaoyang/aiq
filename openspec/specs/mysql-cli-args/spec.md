## ADDED Requirements

### Requirement: MySQL-compatible CLI arguments
The system SHALL support MySQL-compatible command-line connection arguments (`-h`, `-u`, `-P`, `-p`, `-D`) for direct database connection.

#### Scenario: Parse MySQL host argument
- **WHEN** user runs `aiq -h 11.124.9.201`
- **THEN** system parses `-h` flag and stores host value

#### Scenario: Parse MySQL user argument
- **WHEN** user runs `aiq -u root`
- **THEN** system parses `-u` flag and stores username value

#### Scenario: Parse MySQL port argument
- **WHEN** user runs `aiq -P 2900`
- **THEN** system parses `-P` flag and stores port value

#### Scenario: Parse MySQL password argument
- **WHEN** user runs `aiq -ptest1111` (no space between `-p` and password)
- **THEN** system parses `-p` flag and stores password value

#### Scenario: Parse MySQL database argument
- **WHEN** user runs `aiq -D aiq`
- **THEN** system parses `-D` flag and stores database name value

#### Scenario: Use -D database for current session only
- **WHEN** user runs `aiq -h host -u user -P port -ppassword -D database1` and source already exists (based on host-port-username)
- **THEN** system uses `database1` for the current chat session, overriding the source's stored database value

#### Scenario: -D parameter does not persist to source
- **WHEN** user runs `aiq -h host -u user -P port -ppassword -D database1` and source already exists
- **THEN** system does not update the source's database field in `sources.yaml`, only uses `database1` for the current session

#### Scenario: Source uniqueness based on host-port-username
- **WHEN** user runs `aiq -h host -u user -P port -ppassword -D database1` multiple times with different `-D` values
- **THEN** system reuses the same source (based on host-port-username) but uses the `-D` value specified in each command for that session

#### Scenario: Parse combined MySQL arguments
- **WHEN** user runs `aiq -h 11.124.9.201 -u root -P 2900 -ptest1111 -D aiq` (note: `-p` and password have no space)
- **THEN** system parses all flags and stores all connection parameters

### Requirement: PostgreSQL-compatible CLI arguments
The system SHALL support PostgreSQL-compatible command-line connection arguments (`-h`, `-U`, `-p`, `-d`, `-W` or `PGPASSWORD` env var) for direct database connection.

#### Scenario: Parse PostgreSQL host argument
- **WHEN** user runs `aiq -h localhost`
- **THEN** system parses `-h` flag and stores host value

#### Scenario: Parse PostgreSQL user argument
- **WHEN** user runs `aiq -U postgres`
- **THEN** system parses `-U` flag and stores username value

#### Scenario: Parse PostgreSQL port argument
- **WHEN** user runs `aiq -p 5432`
- **THEN** system parses `-p` flag and stores port value (PostgreSQL uses lowercase `-p` for port)

#### Scenario: Parse PostgreSQL database argument
- **WHEN** user runs `aiq -d mydb`
- **THEN** system parses `-d` flag and stores database name value

#### Scenario: Parse PostgreSQL password from environment variable
- **WHEN** user sets `PGPASSWORD=secret` and runs `aiq -h host -U user -p port -d db`
- **THEN** system reads password from `PGPASSWORD` environment variable

#### Scenario: Parse PostgreSQL password from flag
- **WHEN** user runs `aiq -W secret` (non-standard but supported)
- **THEN** system parses `-W` flag and stores password value

#### Scenario: Parse combined PostgreSQL arguments
- **WHEN** user sets `PGPASSWORD=secret` and runs `aiq -h localhost -U postgres -p 5432 -d mydb`
- **THEN** system parses all flags, reads password from environment, and stores all connection parameters

### Requirement: Database type auto-detection
The system SHALL automatically detect database type (MySQL vs PostgreSQL) based on argument patterns, resolving `-p` flag ambiguity.

#### Scenario: Detect MySQL from argument pattern
- **WHEN** user provides `-P` (uppercase port) or `-D` (uppercase database) flags
- **THEN** system detects MySQL type, and `-p` flag is interpreted as password (format: `-ppassword`, no space)

#### Scenario: Detect PostgreSQL from argument pattern
- **WHEN** user provides `-U` (uppercase username) or `-d` (lowercase database) flags
- **THEN** system detects PostgreSQL type, and `-p` flag is interpreted as port (format: `-p 5432` or `-p5432`, both supported)

#### Scenario: Resolve -p flag ambiguity for MySQL
- **WHEN** user provides MySQL flags (`-P` or `-D`) along with `-p` flag
- **THEN** system interprets `-p` as password (e.g., `-ppassword`), not port

#### Scenario: Resolve -p flag ambiguity for PostgreSQL
- **WHEN** user provides PostgreSQL flags (`-U` or `-d`) along with `-p` flag
- **THEN** system interprets `-p` as port (e.g., `-p 5432` or `-p5432`), not password

#### Scenario: Explicit engine override
- **WHEN** user provides `--engine mysql` or `-e postgresql` flag
- **THEN** system uses explicitly specified engine type, ignoring auto-detection, and interprets `-p` accordingly

#### Scenario: Reserve -t flag for future use
- **WHEN** user provides `-t` flag
- **THEN** system reserves `-t` for future use (e.g., OceanBase tenant parameter) and does not interpret it as database engine type

#### Scenario: Default to MySQL when ambiguous
- **WHEN** user provides ambiguous flags that don't clearly indicate MySQL or PostgreSQL (e.g., only `-p` without `-P`, `-D`, `-U`, or `-d`)
- **THEN** system defaults to MySQL type, interpreting `-p` as password

### Requirement: Argument validation
The system SHALL validate required arguments and connection parameters.

#### Scenario: Validate required arguments
- **WHEN** user provides CLI args but missing required fields (host, user, port, password, database)
- **THEN** system displays error message indicating missing required arguments and exits with error code 1

#### Scenario: Validate port number range
- **WHEN** user provides invalid port number (e.g., negative, zero, or > 65535)
- **THEN** system displays error message and exits with error code 1

#### Scenario: Validate PostgreSQL password requirement
- **WHEN** user provides PostgreSQL flags but neither `PGPASSWORD` env var nor `-W` flag is set
- **THEN** system displays error message indicating password is required and exits with error code 1

### Requirement: Connection validation before source creation
The system SHALL test database connection immediately when database CLI args are provided, before creating source or entering chat mode.

#### Scenario: Test MySQL connection with valid credentials
- **WHEN** user provides MySQL CLI args with valid connection parameters
- **THEN** system attempts to connect to MySQL database and validates connection succeeds

#### Scenario: Test PostgreSQL connection with valid credentials
- **WHEN** user provides PostgreSQL CLI args with valid connection parameters
- **THEN** system attempts to connect to PostgreSQL database and validates connection succeeds

#### Scenario: Fail on invalid connection
- **WHEN** user provides database CLI args with invalid connection parameters (wrong host, port, credentials, or database)
- **THEN** system displays clear error message indicating connection failure and exits with error code 1 without creating source

#### Scenario: Fail on network error
- **WHEN** user provides database CLI args but database is unreachable
- **THEN** system displays network error message and exits with error code 1

#### Scenario: Fail on authentication error
- **WHEN** user provides database CLI args with incorrect username or password
- **THEN** system displays authentication error message and exits with error code 1
