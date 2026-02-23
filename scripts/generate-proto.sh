#!/bin/bash
set -e

# Generate protobuf files using buf
echo "Generating protobuf files..."

# Check if buf is installed
if ! command -v buf &> /dev/null; then
    echo "Error: buf is not installed"
    echo "Please install buf from https://docs.buf.build/installation"
    exit 1
fi

# Run buf generate
buf generate

echo "Proto generation complete!"
echo "Generated files can be found in:"
echo "  - gen/go/    (Go code)"
echo "  - gen/web/   (Web/TypeScript code)"
echo "  - gen/openapi/ (OpenAPI/Swagger specs)"
