## Context

The current devcontainer implementation in `pkg/config/devcontainer.go` only supports a subset of the official Dev Container specification defined in `devcontainerjson-reference.md`. The existing implementation covers basic properties like image, dockerfile, features, mounts, ports, and customizations, but lacks many important features that users expect from the Dev Container spec.

## Goals / Non-Goals

**Goals:**
- Align the implementation with the official Dev Container specification
- Add support for lifecycle commands (initializeCommand, onCreateCommand, updateContentCommand, postCreateCommand, postStartCommand, postAttachCommand, waitFor)
- Add host requirements support (cpus, memory, storage)
- Add port configuration (portsAttributes, otherPortsAttributes)
- Add security options (init, privileged, capAdd, securityOpt)
- Add workspace controls (workspaceMount, workspaceFolder, runArgs)
- Add environment variable support (remoteEnv, updateRemoteUserUID, userEnvProbe)
- Add Docker Compose support (dockerComposeFile, service, runServices)
- Add command handling (overrideCommand, shutdownAction)
- Rename `Env` field to `ContainerEnv` for spec compliance

**Non-Goals:**
- Implement actual command execution (this is handled by external tools like VS Code)
- Implement Docker Compose container orchestration
- Implement container runtime features (these are handled by the container runtime)
- Add support for deprecated properties

## Decisions

1. **Struct-based configuration**: Use Go structs for parsing devcontainer.json rather than generic map[string]interface{} for better type safety and validation.

2. **Backward compatibility for Env**: Instead of completely removing the `Env` field, add `ContainerEnv` as the spec-compliant name while keeping `Env` as a deprecated alias for backward compatibility.

3. **Slice-based lifecycle commands**: Store lifecycle commands as string slices to support both single command and multiple command formats.

4. **Interface-based feature resolution**: Keep the existing feature resolution pattern but extend it to handle the new configuration options.

## Risks / Trade-offs

- **Risk**: Adding many new fields increases the complexity of the struct
  - **Mitigation**: Use nested structs to group related properties logically

- **Risk**: Breaking changes to the public API
  - **Mitigation**: Use type aliases and deprecated annotations to maintain backward compatibility

- **Risk**: Some spec properties may not be applicable to this project's use case
  - **Mitigation**: Focus on properties that can be useful for Docker build and container setup

## Migration Plan

1. Add new fields to DevcontainerConfig struct in `pkg/config/devcontainer.go`
2. Add corresponding JSON tags for each new field
3. Update the ResolveExtends function to handle merging of new fields
4. Update tests to cover new fields
5. Run existing tests to ensure backward compatibility

## Open Questions

- Should lifecycle commands be stored as strings, arrays, or objects (for parallel execution)?
- How should port attributes be merged in the extends resolution?
- Should we implement validation for host requirements (e.g., memory format validation)?
