## Why

The current `devcon up` command creates and starts a container but does not execute the devcontainer.json lifecycle commands (postCreateCommand, onCreateCommand, etc.). This means user-defined setup commands like `go mod download` are not run, requiring manual intervention.

## What Changes

- Implement execution of postCreateCommand after container starts
- Implement execution of onCreateCommand during container creation
- Implement execution of updateContentCommand
- Implement execution of postStartCommand
- Add support for both string and array command formats

## Capabilities

### New Capabilities
- `postcreate-commands`: Execute lifecycle commands in containers

### Modified Capabilities
- (None)

## Impact

- `cmd/devcon/up.go` - Add lifecycle command execution
- `internal/builder/docker.go` - Add Exec method to run commands in container
