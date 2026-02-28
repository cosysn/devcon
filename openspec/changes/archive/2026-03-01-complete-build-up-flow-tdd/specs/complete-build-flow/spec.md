## ADDED Requirements

### Requirement: Build image from base image
The devcon build command SHALL build a devcontainer image from a base image specified in devcontainer.json.

#### Scenario: Build with image field
- **WHEN** devcontainer.json contains `"image": "alpine:latest"`
- **THEN** the build command SHALL build a Docker image from alpine:latest

### Requirement: Build image from Dockerfile
The devcon build command SHALL build a devcontainer image from a Dockerfile specified in devcontainer.json.

#### Scenario: Build with dockerfile
- **WHEN** devcontainer.json contains `"build": {"dockerfile": "Dockerfile"}`
- **THEN** the build command SHALL build a Docker image using the specified Dockerfile

### Requirement: Build with features
The devcon build command SHALL build a devcontainer image with features included.

#### Scenario: Build with local features
- **WHEN** devcontainer.json contains `"features": {"test-feature": {}}`
- **THEN** the build command SHALL include the specified features in the image

### Requirement: Build with extends
The devcon build command SHALL build a devcontainer image that extends a base configuration.

#### Scenario: Build with extends
- **WHEN** devcontainer.json contains `"extends": "./base.json"`
- **THEN** the build command SHALL resolve the base configuration and build the merged image
