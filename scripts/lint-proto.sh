#!/bin/bash
set -e

# Lint protobuf files using buf
echo "Linting protobuf files..."

# Check if buf is installed
if ! command -v buf &> /dev/null; then
    echo "Error: buf is not installed"
    echo "Please install buf from https://docs.buf.build/installation"
    exit 1
fi

# Run buf lint
buf lint

echo "Proto linting passed!"
