## 1. Update DevcontainerConfig struct

- [x] 1.1 Add HostRequirements struct with cpus, memory, storage fields
- [x] 1.2 Add PortAttributes struct with label, protocol, onAutoForward, requireLocalPort, elevateIfNeeded fields
- [x] 1.3 Add SecurityOptions struct with init, privileged, capAdd, securityOpt fields
- [x] 1.4 Add WorkspaceConfig struct with mount, folder, runArgs fields
- [x] 1.5 Add RemoteEnv map for remote environment variables
- [x] 1.6 Add ComposeConfig struct with dockerComposeFile, service, runServices, workspaceFolder fields
- [x] 1.7 Add command handling fields: overrideCommand, shutdownAction, waitFor, appPort
- [x] 1.8 Add lifecycle command fields: initializeCommand, onCreateCommand, updateContentCommand, postCreateCommand, postStartCommand, postAttachCommand
- [x] 1.9 Add new environment fields: updateRemoteUserUID, userEnvProbe, remoteEnv
- [x] 1.10 Add ContainerEnv field (spec-compliant name) alongside existing Env field for backward compatibility
- [x] 1.11 Add Name field for container naming
- [x] 1.12 Add PortsAttributes and OtherPortsAttributes fields

## 2. Update JSON parsing

- [x] 2.1 Update ParseDevcontainer to handle new fields
- [x] 2.2 Add validation for host requirements format (memory/storage units)
- [x] 2.3 Add validation for port attributes values

## 3. Update ResolveExtends function

- [x] 3.1 Update merge logic for HostRequirements
- [x] 3.2 Update merge logic for PortsAttributes
- [x] 3.3 Update merge logic for SecurityOptions
- [x] 3.4 Update merge logic for WorkspaceConfig
- [x] 3.5 Update merge logic for RemoteEnv
- [x] 3.6 Update merge logic for ComposeConfig
- [x] 3.7 Update merge logic for lifecycle commands (array concatenation)
- [x] 3.8 Update merge logic for ContainerEnv (keep existing Env handling for backward compatibility)

## 4. Update tests

- [x] 4.1 Update DevcontainerConfig tests to cover new fields
- [x] 4.2 Add tests for lifecycle command parsing
- [x] 4.3 Add tests for host requirements parsing
- [x] 4.4 Add tests for port attributes parsing
- [x] 4.5 Add tests for security options parsing
- [x] 4.6 Add tests for compose configuration parsing
- [x] 4.7 Add tests for extends resolution of new fields

## 5. Verify backward compatibility

- [x] 5.1 Ensure existing test fixtures still work
- [x] 5.2 Verify Env field still works (deprecated alias)
- [x] 5.3 Run full test suite to confirm no regressions
