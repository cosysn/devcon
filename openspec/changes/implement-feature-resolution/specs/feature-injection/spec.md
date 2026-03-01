## ADDED Requirements

### Requirement: Generate Dockerfile with features
The system SHALL generate a Dockerfile that includes feature installation.

#### Scenario: Build with features
- **WHEN** devcontainer.json has features and build is called
- **THEN** the system SHALL generate a Dockerfile that:
  - Starts FROM the base image
  - Copies downloaded feature tarballs
  - Runs each feature's install.sh script

#### Scenario: Handle missing install.sh
- **WHEN** feature doesn't have install.sh
- **THEN** the system SHALL skip that feature's installation step

### Requirement: Feature dependency resolution
The system SHALL resolve feature dependencies before building.

#### Scenario: Simple dependency
- **WHEN** feature A depends on feature B (dependsOn: ["B"])
- **THEN** the system SHALL ensure feature B is resolved and installed before feature A

#### Scenario: Circular dependency detection
- **WHEN** features have circular dependencies
- **THEN** the system SHALL return error "circular dependency detected"

### Requirement: Pass feature options to install script
The system SHALL pass configured options to feature install scripts.

#### Scenario: Options passed as environment variables
- **WHEN** feature has options configured (e.g., `"option1": "value1"`)
- **THEN** the system SHALL set environment variables FEATURE_OPTION1=value1 before running install.sh
