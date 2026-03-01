## Why

The current `devcon build` command does not actually resolve and inject devcontainer features into the built image. When users specify features like `"docker-in-docker": {}` or `"git": {}`, these are ignored and only the base image is used. This breaks the core functionality users expect from devcontainer.

## What Changes

- Implement feature resolution from OCI registry (e.g., ghcr.io)
- Implement shorthand feature name resolution (e.g., "docker-in-docker" → full OCI reference)
- Implement feature download and injection into Dockerfile
- Implement feature dependency resolution (dependsOn)
- Add feature caching to avoid re-downloading

## Capabilities

### New Capabilities
- `feature-oci-resolution`: Resolve features from OCI registries (ghcr.io, etc.)
- `feature-shorthand-resolution`: Convert shorthand names to full OCI references
- `feature-injection`: Inject resolved features into Dockerfile during build

### Modified Capabilities
- (None - this is new functionality)

## Impact

- `cmd/devcon/build.go` - Will call feature resolution before building
- `internal/builder/docker.go` - Needs to accept resolved features and generate Dockerfile
- `pkg/feature/` - Needs OCI resolution logic
- `pkg/config/feature.go` - May need additional resolution functions
