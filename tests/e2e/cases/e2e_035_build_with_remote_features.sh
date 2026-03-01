#!/bin/bash
# E2E-035: build_with_remote_features - Build with remote features from OCI registry
# This test verifies that features are actually resolved and installed in the built image.
# Uses 'git' feature which should be available from ghcr.io/devcontainers/features

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE="$SCRIPT_DIR/../fixtures/devcontainer/with-remote-features"

cd "$PROJECT_ROOT"

# Create fixture directory if it doesn't exist
mkdir -p "$FIXTURE/.devcontainer"

# Create devcontainer.json with remote feature (shorthand name)
# Use debian-based image since we use apt-get in the generated Dockerfile
cat > "$FIXTURE/.devcontainer/devcontainer.json" << 'EOF'
{
    "image": "debian:bookworm-slim",
    "features": {
        "git": {}
    }
}
EOF

# Create a simple Dockerfile that will be used if needed
cat > "$FIXTURE/Dockerfile" << 'EOF'
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y git
EOF

# Verify fixture exists
if [ ! -d "$FIXTURE" ]; then
    echo "Error: Fixture not found: $FIXTURE"
    exit 1
fi

# Execute build
OUTPUT=$(./devcon build "$FIXTURE" 2>&1)

# Verify output contains image ID
if ! echo "$OUTPUT" | grep -q "Image built:"; then
    echo "Error: Build did not produce expected output"
    echo "Output: $OUTPUT"
    exit 1
fi

# Get the image name from output
IMAGE_NAME=$(echo "$OUTPUT" | grep "Image built:" | awk '{print $NF}')

# Run the built image and verify git is installed
# This verifies that features were actually injected into the image
CONTAINER_OUTPUT=$(docker run --rm "$IMAGE_NAME" which git 2>&1) || true

if echo "$CONTAINER_OUTPUT" | grep -q "/usr/bin/git" || echo "$CONTAINER_OUTPUT" | grep -q "/bin/git"; then
    echo "Success: git feature is installed in the built image"
    exit 0
else
    echo "Error: git feature was NOT installed in the built image"
    echo "Container output: $CONTAINER_OUTPUT"
    echo "This means features are not being resolved and injected during build"
    exit 1
fi
