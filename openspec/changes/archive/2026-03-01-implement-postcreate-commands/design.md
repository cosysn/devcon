## Context

Current `devcon up` flow:
1. Build image (if needed)
2. Create container
3. Start container
4. Done - no lifecycle commands executed

## Goals / Non-Goals

**Goals:**
- Execute postCreateCommand after container starts
- Execute onCreateCommand during container creation
- Execute postStartCommand after container starts
- Support both string and array command formats
- Show command output to user

**Non-Goals:**
- Watch for file changes (updateContentCommand)
- Full devcontainer CLI compatibility

## Decisions

1. **Where to execute commands?**
   - Option A: Execute in build.go after Up returns
   - Option B: Execute in builder/docker.go as part of Up
   - **Decision**: Execute in builder/docker.go - cleaner separation

2. **Command format support:**
   - String: `"postCreateCommand": "go mod download"`
   - Array: `"postCreateCommand": ["go", "mod", "download"]`
   - **Decision**: Support both formats

3. **Execution order:**
   - onCreateCommand → create
   - postStartCommand → start
   - postCreateCommand → after start
