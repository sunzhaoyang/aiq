## MODIFIED Requirements

### Requirement: MySQL-compatible CLI arguments
The system SHALL support MySQL-compatible command-line connection arguments (`-h`, `-u`, `-P`, `-p`, `-D`) for direct database connection, with `-D` database parameter used for the current session only without persisting to source configuration.

#### Scenario: Parse MySQL database argument
- **WHEN** user runs `aiq -h host -u user -P port -ppassword -D database`
- **THEN** system parses `-D` flag and stores database name value

#### Scenario: Use -D database for current session
- **WHEN** user runs `aiq -h host -u user -P port -ppassword -D database1` and source already exists (based on host-port-username)
- **THEN** system uses `database1` for the current chat session, overriding the source's stored database value

#### Scenario: -D parameter does not persist to source
- **WHEN** user runs `aiq -h host -u user -P port -ppassword -D database1` and source already exists
- **THEN** system does not update the source's database field in `sources.yaml`, only uses `database1` for the current session

#### Scenario: Source uniqueness based on host-port-username
- **WHEN** user runs `aiq -h host -u user -P port -ppassword -D database1` multiple times with different `-D` values
- **THEN** system reuses the same source (based on host-port-username) but uses the `-D` value specified in each command for that session

#### Scenario: Create source without -D database
- **WHEN** user runs `aiq -h host -u user -P port -ppassword` without `-D` flag and source does not exist
- **THEN** system creates source with empty database field (or prompts for database if required)

#### Scenario: Create source with -D database
- **WHEN** user runs `aiq -h host -u user -P port -ppassword -D database1` and source does not exist
- **THEN** system creates source with `database1` as the initial database value
