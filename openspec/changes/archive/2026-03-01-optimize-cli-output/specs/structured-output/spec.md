## ADDED Requirements

### Requirement: JSON output mode for programmatic access
The CLI SHALL accept an `--output json` flag that enables JSON formatted output for programmatic consumption.

#### Scenario: JSON output on build command
- **WHEN** user runs `devcon build --output json .`
- **THEN** command outputs JSON with success status, message, and any relevant data

#### Scenario: JSON output on up command
- **WHEN** user runs `devcon up --output json .`
- **THEN** command outputs JSON with success status, container details, and any relevant data

### Requirement: JSON output structure
JSON output SHALL use a consistent envelope format that includes success status, command name, message, and data fields.

#### Scenario: JSON success response
- **WHEN** command succeeds with `--output json`
- **THEN** output is valid JSON with `"success": true`

#### Scenario: JSON error response
- **WHEN** command fails with `--output json`
- **THEN** output is valid JSON with `"success": false` and error details in `"error"` field

#### Scenario: JSON output includes command metadata
- **WHEN** command outputs JSON
- **THEN** output includes `"command": "<command-name>"` and `"timestamp": "<iso8601>"`
