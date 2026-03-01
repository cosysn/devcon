## ADDED Requirements

### Requirement: Environment variable configuration
The devcontainer configuration SHALL support environment variables for both container and remote processes.

#### Scenario: Container environment variables
- **WHEN** devcontainer.json specifies containerEnv
- **THEN** the configuration SHALL store environment variables that apply to the container

#### Scenario: Remote environment variables
- **WHEN** devcontainer.json specifies remoteEnv
- **THEN** the configuration SHALL store environment variables for devcontainer tool processes

#### Scenario: Update user UID enabled
- **WHEN** devcontainer.json specifies updateRemoteUserUID as true
- **THEN** the configuration SHALL indicate that user UID/GID should match local user

#### Scenario: User environment probe specified
- **WHEN** devcontainer.json specifies userEnvProbe
- **THEN** the configuration SHALL store the probe type (none, interactiveShell, loginShell, loginInteractiveShell)

#### Scenario: Variables reference container env
- **WHEN** remoteEnv references ${containerEnv:VARIABLE_NAME}
- **THEN** the variable SHALL be resolved to the container environment value
