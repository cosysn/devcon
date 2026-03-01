## ADDED Requirements

### Requirement: Enhanced error messages with context
The CLI SHALL provide enhanced error messages that include context about what operation failed and why.

#### Scenario: Config parse error shows file path
- **WHEN** devcontainer.json parsing fails
- **THEN** error message includes the file path and specific line/field with the issue

#### Scenario: Build error shows details
- **WHEN** build operation fails
- **THEN** error message includes the underlying docker error and suggestions for resolution

### Requirement: Actionable error suggestions
The CLI SHALL provide actionable suggestions in error messages to help users resolve issues.

#### Scenario: Missing image or Dockerfile suggestion
- **WHEN** user runs build/up without image or dockerfile specified
- **THEN** error message suggests adding one to devcontainer.json

#### Scenario: Docker not available suggestion
- **WHEN** docker is not available or not running
- **THEN** error message suggests checking if docker is installed and running

### Requirement: Error codes for programmatic handling
The CLI SHALL include error codes in error output for programmatic error handling.

#### Scenario: Error includes error code
- **WHEN** command fails
- **THEN** error output includes a machine-readable error code
