## ADDED Requirements

### Requirement: Lifecycle command support
The devcontainer configuration SHALL support lifecycle commands that execute at various stages of container creation and startup.

#### Scenario: Initialize command runs on host
- **WHEN** devcontainer.json specifies initializeCommand
- **THEN** the command SHALL run on the host machine during initialization

#### Scenario: OnCreateCommand runs in container
- **WHEN** devcontainer.json specifies onCreateCommand
- **THEN** the command SHALL run inside the container after it starts for the first time

#### Scenario: UpdateContentCommand runs for content updates
- **WHEN** devcontainer.json specifies updateContentCommand
- **THEN** the command SHALL run inside the container when new content is available

#### Scenario: PostCreateCommand runs after creation
- **WHEN** devcontainer.json specifies postCreateCommand
- **THEN** the command SHALL run inside the container after container creation is complete

#### Scenario: PostStartCommand runs after container starts
- **WHEN** devcontainer.json specifies postStartCommand
- **THEN** the command SHALL run each time the container is started

#### Scenario: PostAttachCommand runs after tool attaches
- **WHEN** devcontainer.json specifies postAttachCommand
- **THEN** the command SHALL run each time a tool attaches to the container

#### Scenario: WaitFor controls when to connect
- **WHEN** devcontainer.json specifies waitFor
- **THEN** the tool SHALL wait for the specified command to complete before connecting

#### Scenario: Multiple lifecycle commands as array
- **WHEN** lifecycle command is specified as an array
- **THEN** each command in the array SHALL be executed in order without a shell

#### Scenario: Multiple lifecycle commands as string
- **WHEN** lifecycle command is specified as a string
- **THEN** the command SHALL be executed through a shell

### Requirement: Lifecycle command failure handling
If a lifecycle command fails, subsequent commands SHALL NOT be executed.

#### Scenario: PostCreateCommand failure stops execution
- **WHEN** postCreateCommand fails
- **THEN** postStartCommand SHALL NOT be executed
