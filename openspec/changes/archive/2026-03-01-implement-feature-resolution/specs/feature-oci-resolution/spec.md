## ADDED Requirements

### Requirement: Resolve feature from OCI registry
The system SHALL resolve features from OCI registries using go-containerregistry.

#### Scenario: Resolve from ghcr.io
- **WHEN** devcontainer.json specifies `"features": {"ghcr.io/user/feature": {}}`
- **THEN** the system SHALL download the feature tarball from the registry

#### Scenario: Resolve with version tag
- **WHEN** devcontainer.json specifies `"features": {"ghcr.io/user/feature:1.0.0": {}}`
- **THEN** the system SHALL download the specific version from the registry

### Requirement: Parse feature definition from OCI image
The system SHALL parse devcontainer-feature.json from the OCI image labels.

#### Scenario: Parse valid feature image
- **WHEN** feature image contains devcontainer-feature label
- **THEN** the system SHALL parse the feature definition including id, name, version, options

### Requirement: Handle OCI resolution errors
The system SHALL provide clear error messages when OCI resolution fails.

#### Scenario: Registry not found
- **WHEN** registry is unreachable or doesn't exist
- **THEN** the system SHALL return error "failed to resolve feature X: not found"

#### Scenario: Feature label missing
- **WHEN** OCI image doesn't contain devcontainer-feature label
- **THEN** the system SHALL return error "feature image missing required label"
