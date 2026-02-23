#!/bin/bash
set -e

# Check for breaking changes in protobuf files using buf
echo "Checking for breaking changes in protobuf files..."

# Check if buf is installed
if ! command -v buf &> /dev/null; then
    echo "Error: buf is not installed"
    echo "Please install buf from https://docs.buf.build/installation"
    exit 1
fi

# Check against the latest commit (main branch or specified baseline)
if [ -n "$1" ]; then
    echo "Checking against baseline: $1"
    buf breaking --against "$1"
else
    echo "Checking against .git (HEAD^)"
    buf breaking --against '.git#HEAD^'
fi

echo "No breaking changes detected!"
