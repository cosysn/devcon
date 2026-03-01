## Why

The current devcontainer configuration lacks support for specifying a custom user and configuring the working directory to the user's workspace. This limits flexibility for multi-user environments and makes it difficult to maintain consistent development environments across different users who may have different preferred workspace locations.

## What Changes

- Add support for specifying a custom user in the devcontainer configuration
- Add default VSCode customization settings
- Add configuration for working directory to point to the specified user's workspace
- Add user creation/selection logic in the devcontainer Dockerfile

## Capabilities

### New Capabilities

- `devcontainer-user-config`: Support specifying a custom user and workspace directory in devcontainer configuration
  - `user` parameter: Allow specifying which user to run as (default: root or vscode)
  - `workspace` parameter: Configure the working directory to user's workspace
  - VSCode defaults: Apply default VSCode settings for the specified user

### Modified Capabilities

- None. This is a new capability.

## Impact

- New configuration options in `.devcontainer/devcontainer.json`
- Updates to `.devcontainer/generated.Dockerfile` for user handling
- Potentially affects how devcontainer is built and started
