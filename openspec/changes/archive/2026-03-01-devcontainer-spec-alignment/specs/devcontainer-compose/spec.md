## ADDED Requirements

### Requirement: Docker Compose configuration
The devcontainer configuration SHALL support Docker Compose for multi-container scenarios.

#### Scenario: Docker Compose file specified
- **WHEN** devcontainer.json specifies dockerComposeFile
- **THEN** the configuration SHALL store the path to the Docker Compose file(s)

#### Scenario: Service name specified
- **WHEN** devcontainer.json specifies service
- **THEN** the configuration SHALL store the name of the service to connect to

#### Scenario: Run services specified
- **WHEN** devcontainer.json specifies runServices
- **THEN** the configuration SHALL store the array of services to start

#### Scenario: Multiple compose files
- **WHEN** dockerComposeFile is an array
- **THEN** each file in the array SHALL be stored in order

#### Scenario: Compose workspace folder
- **WHEN** using Docker Compose with workspaceFolder
- **THEN** workspaceFolder SHALL default to "/" if not specified
