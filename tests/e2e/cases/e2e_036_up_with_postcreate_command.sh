#!/bin/bash
# E2E-036: up_with_postcreate_command - Test postCreateCommand execution

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE="$SCRIPT_DIR/../fixtures/devcontainer/with-postcreate-command"

cd "$PROJECT_ROOT"

# Create fixture directory
mkdir -p "$FIXTURE/.devcontainer"

# Create devcontainer.json with postCreateCommand
cat > "$FIXTURE/.devcontainer/devcontainer.json" << 'EOF'
{
    "image": "alpine:latest",
    "postCreateCommand": "echo 'POST_CREATE_EXECUTED' > /tmp/postcreate.txt"
}
EOF

# Execute up
OUTPUT=$(./devcon up "$FIXTURE" 2>&1) || true

# Check if container started
if ! echo "$OUTPUT" | grep -q "Container started:"; then
    echo "Error: Container did not start"
    echo "Output: $OUTPUT"
    exit 1
fi

# Get container ID from output
CONTAINER_ID=$(echo "$OUTPUT" | grep "Container started:" | awk '{print $NF}')

# Wait a bit for command to execute
sleep 2

# Check if postCreateCommand was executed
CONTAINER_OUTPUT=$(docker exec "$CONTAINER_ID" cat /tmp/postcreate.txt 2>&1) || true

# Clean up
docker kill "$CONTAINER_ID" > /dev/null 2>&1 || true
docker rm "$CONTAINER_ID" > /dev/null 2>&1 || true

if echo "$CONTAINER_OUTPUT" | grep -q "POST_CREATE_EXECUTED"; then
    echo "Success: postCreateCommand was executed"
    exit 0
else
    echo "Error: postCreateCommand was NOT executed"
    echo "Container output: $CONTAINER_OUTPUT"
    exit 1
fi
