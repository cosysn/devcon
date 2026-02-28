#!/bin/bash
# E2E-034: help_all_commands - help 命令

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

cd "$PROJECT_ROOT"

# 测试 devcon --help
OUTPUT1=$(./devcon --help 2>&1)
if ! echo "$OUTPUT1" | grep -q "devcontainer"; then
    echo "Error: devcon --help failed"
    exit 1
fi

# 测试 devcon build --help
OUTPUT2=$(./devcon build --help 2>&1)
if ! echo "$OUTPUT2" | grep -q "build"; then
    echo "Error: devcon build --help failed"
    exit 1
fi

# 测试 devcon features --help
OUTPUT3=$(./devcon features --help 2>&1)
if ! echo "$OUTPUT3" | grep -q "features"; then
    echo "Error: devcon features --help failed"
    exit 1
fi

# 测试 devcon up --help
OUTPUT4=$(./devcon up --help 2>&1)
if ! echo "$OUTPUT4" | grep -q "up"; then
    echo "Error: devcon up --help failed"
    exit 1
fi

# 测试 devcon config --help
OUTPUT5=$(./devcon config --help 2>&1)
if ! echo "$OUTPUT5" | grep -q "config"; then
    echo "Error: devcon config --help failed"
    exit 1
fi

# 测试 devcon inspect --help
OUTPUT6=$(./devcon inspect --help 2>&1)
if ! echo "$OUTPUT6" | grep -q "inspect"; then
    echo "Error: devcon inspect --help failed"
    exit 1
fi

exit 0
