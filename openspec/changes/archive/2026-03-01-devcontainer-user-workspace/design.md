## Context

The current devcontainer configuration in `.devcontainer/devcontainer.json` uses a base image with pre-installed features but lacks support for:
1. Specifying a custom user (defaults to root)
2. Configuring the working directory to a user's workspace
3. Default VSCode settings for custom users

This change needs to support scenarios where developers want to run as a non-root user with their own workspace directory.

## Goals / Non-Goals

**Goals:**
- Add support for specifying a custom user in devcontainer configuration
- Configure working directory to point to the specified user's workspace
- Apply default VSCode customizations
- Maintain backward compatibility with existing configurations

**Non-Goals:**
- User management beyond creation/selection (no complex user permissions)
- Multiple simultaneous users
- Authentication or authorization features

## Decisions

1. **User Configuration via devcontainer.json**
   - Decision: Add `user` and `workspace` properties to devcontainer.json
   - Rationale: Keeps configuration declarative and consistent with devcontainer spec

2. **Default User Selection**
   - Decision: Default to `vscode` user if not specified
   - Rationale: `vscode` is the standard non-root user in Microsoft devcontainer images

3. **Workspace Directory**
   - Decision: Default workspace to `/home/<user>/workspace`
   - Rationale: Follows common convention and provides clear separation

4. **Dockerfile Generation**
   - Decision: Modify generated.Dockerfile to handle user creation/selection
   - Rationale: User must be set during image build, not container runtime

## Risks / Trade-offs

- **Risk**: User permission issues with mounted directories
  - Mitigation: Ensure the specified user has correct ownership of workspace directory

- **Risk**: Breaking existing configurations that rely on root user
  - Mitigation: Make user specification optional, defaulting to current behavior

- **Trade-off**: Simplicity vs Flexibility
  - Decision: Keep configuration simple (user + workspace) rather than supporting complex permission schemes
