## ADDED Requirements

### Requirement: Convert shorthand name to OCI reference
The system SHALL convert shorthand feature names to full OCI references.

#### Scenario: Basic shorthand conversion
- **WHEN** devcontainer.json specifies `"features": {"docker-in-docker": {}}`
- **THEN** the system SHALL convert to `ghcr.io/devcontainers/features/docker-in-docker:latest`

#### Scenario: Shorthand with options
- **WHEN** devcontainer.json specifies `"features": {"docker-in-docker": {"option1": "value1"}}`
- **THEN** the system SHALL convert to full OCI reference with options preserved

#### Scenario: Full reference passthrough
- **WHEN** devcontainer.json specifies `"features": {"ghcr.io/user/feature": {}}`
- **THEN** the system SHALL use the reference as-is without conversion

### Requirement: Default registry configuration
The system SHALL use configurable default registry for shorthand names.

#### Scenario: Default registry
- **WHEN** shorthand feature name is used without explicit registry
- **THEN** the system SHALL use ghcr.io/devcontainers/features/ as default
