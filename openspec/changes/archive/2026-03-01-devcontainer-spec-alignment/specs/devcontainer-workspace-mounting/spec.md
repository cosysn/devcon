## ADDED Requirements

### Requirement: Workspace mounting configuration
The devcontainer configuration SHALL support workspace mount and runtime options.

#### Scenario: Workspace mount specified
- **WHEN** devcontainer.json specifies workspaceMount
- **THEN** the configuration SHALL store the mount configuration string

#### Scenario: Workspace folder specified
- **WHEN** devcontainer.json specifies workspaceFolder
- **THEN** the configuration SHALL store the workspace path in the container

#### Scenario: Run arguments specified
- **WHEN** devcontainer.json specifies runArgs
- **THEN** the configuration SHALL store Docker CLI arguments for container runtime

#### Scenario: Variables in workspaceMount
- **WHEN** workspaceMount contains ${localWorkspaceFolder}
- **THEN** the variable SHALL be resolved to the actual workspace path
