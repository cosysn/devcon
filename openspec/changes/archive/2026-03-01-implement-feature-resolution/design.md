## Context

Current state:
- `devcon build` parses devcontainer.json and extracts features from config
- Features are passed to builder but ignored - builder only uses base image or Dockerfile
- No OCI registry resolution exists for remote features

## Goals / Non-Goals

**Goals:**
- Resolve shorthand feature names (e.g., "docker-in-docker") to full OCI references
- Download features from OCI registry during build
- Generate Dockerfile that installs features using devcontainer-feature scripts
- Handle feature dependencies (dependsOn)

**Non-Goals:**
- Caching layer (may add later)
- Publishing features to registry
- Full feature option validation (basic validation only)

## Decisions

1. **Where to resolve features?**
   - Option A: Resolve in build.go before calling builder
   - Option B: Resolve in builder/docker.go during build
   - **Decision**: Resolve in build.go - separates concerns, easier to test

2. **How to generate Dockerfile with features?**
   - Option A: Generate temp Dockerfile that includes feature install scripts
   - Option B: Use multi-stage build with feature layer
   - **Decision**: Generate temp Dockerfile in build context that:
     - Starts FROM base image
     - Copies feature tarballs
     - Runs each feature's install.sh

3. **Shorthand resolution strategy:**
   - Shorthand names like "docker-in-docker" need default registry
   - Default registry: ghcr.io/devcontainers/features/
   - Full reference: ghcr.io/devcontainers/features/docker-in-docker:latest

## Risks / Trade-offs

- **Risk**: OCI registry may be slow or unavailable
  - **Mitigation**: Show progress during download, allow offline build if features cached

- **Risk**: Feature install scripts may fail
  - **Mitigation**: Pass through install script output, fail build on error

- **Risk**: Large feature download size
  - **Mitigation**: Cache downloaded features locally
