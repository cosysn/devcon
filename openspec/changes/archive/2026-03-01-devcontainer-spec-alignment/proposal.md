## Why

The current devcontainer implementation in this project only supports a subset of the official Dev Container specification. Many important properties defined in `devcontainerjson-reference.md` are missing, causing misalignment with the standard and limiting the functionality that can be configured. This prevents users from using advanced features like lifecycle scripts, host requirements, port attributes, and Docker Compose support.

## What Changes

- Add support for all missing properties from the official Dev Container spec
- Implement lifecycle commands: `initializeCommand`, `updateContentCommand`, `postAttachCommand`, `waitFor`
- Add host requirements: `cpus`, `memory`, `storage`
- Add port configuration: `portsAttributes`, `otherPortsAttributes`
- Add security options: `init`, `privileged`, `capAdd`, `securityOpt`
- Add workspace controls: `workspaceMount`, `workspaceFolder`, `runArgs`
- Add environment variables: `remoteEnv`, `updateRemoteUserUID`, `userEnvProbe`
- Add Docker Compose support: `dockerComposeFile`, `service`, `runServices`
- Add command override: `overrideCommand`, `shutdownAction`
- Add `name` property for container naming
- **BREAKING**: Rename `Env` field to `ContainerEnv` for spec compliance

## Capabilities

### New Capabilities
- `devcontainer-lifecycle-commands`: Support for initializeCommand, onCreateCommand, updateContentCommand, postCreateCommand, postStartCommand, postAttachCommand, and waitFor
- `devcontainer-host-requirements`: Support for cpus, memory, and storage requirements
- `devcontainer-port-configuration`: Support for portsAttributes and otherPortsAttributes
- `devcontainer-security-options`: Support for init, privileged, capAdd, and securityOpt
- `devcontainer-workspace-mounting`: Support for workspaceMount, workspaceFolder, and runArgs
- `devcontainer-environment-variables`: Support for remoteEnv, updateRemoteUserUID, and userEnvProbe
- `devcontainer-compose`: Support for dockerComposeFile, service, and runServices
- `devcontainer-command-handling`: Support for overrideCommand and shutdownAction

### Modified Capabilities
- `devcontainer-user-config`: Extend to include the renamed ContainerEnv field and additional user-related properties

## Impact

- Changes to `pkg/config/devcontainer.go` - struct definition updates
- Changes to `pkg/config/devcontainer_test.go` - test updates for new properties
- May affect feature generation in `pkg/feature/generate.go`
