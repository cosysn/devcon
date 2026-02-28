## ADDED Requirements

### Requirement: Up builds image if needed
The devcon up command SHALL build the devcontainer image if it doesn't exist.

#### Scenario: Up builds missing image
- **WHEN** the image doesn't exist and devcontainer.json specifies an image or Dockerfile
- **THEN** the up command SHALL build the image before starting the container

### Requirement: Up starts container
The devcon up command SHALL start a container from the built devcontainer image.

#### Scenario: Up starts container from image
- **WHEN** devcontainer.json contains `"image": "alpine:latest"`
- **THEN** the up command SHALL start a container running alpine:latest

### Requirement: Up starts container from Dockerfile
The devcon up command SHALL build from Dockerfile if needed and start a container.

#### Scenario: Up builds and starts from Dockerfile
- **WHEN** devcontainer.json contains `"build": {"dockerfile": "Dockerfile"}`
- **THEN** the up command SHALL build the image and start a container

### Requirement: Up with features
The devcon up command SHALL build an image with features and start a container.

#### Scenario: Up with features
- **WHEN** devcontainer.json contains features
- **THEN** the up command SHALL build the image with features and start a container
