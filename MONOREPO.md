# CardWatcher Monorepo

A Go-based monorepo for monitoring trading card prices with multiple services and applications.

## Architecture Overview

This monorepo follows a multi-module Go workspace layout with the following applications:

- **apps/server** - gRPC server implementing the CardWatcher service
- **apps/gateway** - gRPC-Gateway for HTTP/JSON REST API access
- **apps/cli** - Command-line client for interacting with the service
- **apps/web** - Web frontend (placeholder for future implementation)

## Project Structure

```
.
├── apps/
│   ├── server/          # gRPC server application
│   ├── gateway/         # gRPC-Gateway (HTTP/JSON to gRPC)
│   ├── cli/             # CLI client application
│   └── web/             # Web frontend (placeholder)
├── proto/
│   └── cardwatcher/
│       └── v1/          # v1 Protocol Buffer definitions
├── gen/
│   ├── go/              # Generated Go code
│   ├── web/             # Generated web/TypeScript code
│   └── openapi/         # Generated OpenAPI/Swagger specs
├── scripts/             # Utility scripts
├── go.work              # Go workspace file
├── buf.yaml             # Buf configuration
├── buf.gen.yaml         # Buf generation configuration
└── buf.work.yaml        # Buf workspace configuration
```

## Getting Started

### Prerequisites

- Go 1.24 or later
- Buf CLI (for Protocol Buffer management)
- Docker and Docker Compose (for running services)
- MongoDB (can run via Docker Compose)
- Ntfy (can run via Docker Compose)

### Installing Buf

```bash
# On macOS
brew install bufbuild/buf/buf

# On Linux
curl -sSL "https://github.com/bufbuild/buf/releases/latest/download/buf-$(uname -s)-$(uname -m)" -o "/usr/local/bin/buf"
chmod +x "/usr/local/bin/buf"

# Or using go install
go install github.com/bufbuild/buf/cmd/buf@latest
```

## Go Workspace

This project uses Go workspaces to manage multiple modules in a single repository. The `go.work` file at the root defines the workspace.

```bash
# Sync the workspace
go work sync

# Build all modules in the workspace
go work build ./...
```

## Protocol Buffers

### Directory Structure

Protocol Buffer definitions are organized under `proto/cardwatcher/v1/` following a versioned package structure:

```
proto/cardwatcher/v1/
├── entities.proto       # Shared entities (Condition, Expansion, Blueprint, Watch)
├── requests.proto       # Request messages
├── responses.proto      # Response messages
└── service.proto        # Service definitions with HTTP annotations
```

### Versioning Pattern

When creating a new major version (e.g., v2), follow this pattern:

```
proto/cardwatcher/v2/
├── entities.proto
├── requests.proto
├── responses.proto
└── service.proto
```

Update the `go_package` option in each proto file:

```protobuf
option go_package = "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v2;cardwatcherv2";
```

### Generating Code

Use the provided script to generate all protobuf artifacts:

```bash
./scripts/generate-proto.sh
```

This will generate:
- **Go code** in `gen/go/cardwatcher/v1/` - gRPC server and client stubs
- **Gateway stubs** in `gen/go/cardwatcher/v1/` - gRPC-Gateway reverse proxy
- **OpenAPI specs** in `gen/openapi/` - Swagger JSON documentation
- **TypeScript types** in `gen/web/` - For frontend integration

Or run buf directly:

```bash
buf generate
```

### Linting Protobuf

Lint your proto files to ensure they follow best practices:

```bash
./scripts/lint-proto.sh
```

Or:

```bash
buf lint
```

### Checking Breaking Changes

Before merging changes, check for breaking changes:

```bash
./scripts/check-breaking.sh
```

Or check against a specific commit:

```bash
./scripts/check-breaking.sh .git#main
```

### Proto Development Guidelines

#### Reserved Fields

Always reserve field numbers for future use when adding new fields. Reserve a range of field numbers (e.g., 100-199) to allow future additions without renumbering:

```protobuf
message MyMessage {
  string name = 1;
  // reserved field numbers 2-9 for future use
  reserved 2 to 9;
}
```

#### Field Renumbering

**Never renumber existing fields**. This is a breaking change. Always use new field numbers when adding fields.

#### Compatibility Guardrails

1. **Use reserved ranges** for future additions
2. **Never delete or reuse field numbers**
3. **Add new fields at the end** of messages
4. **Run buf breaking checks** before committing
5. **Always lint** with `buf lint`

## Building Applications

### Server (gRPC)

```bash
cd apps/server
go build -o cardwatcher-server ./cmd/main.go
```

Build Docker image:

```bash
docker build -t cardwatcher-server:latest apps/server
```

### Gateway (HTTP/JSON)

```bash
cd apps/gateway
go build -o cardwatcher-gateway ./cmd/main.go
```

Build Docker image:

```bash
docker build -t cardwatcher-gateway:latest apps/gateway
```

### CLI Client

```bash
cd apps/cli
go build -o cardwatcher-cli ./cmd/main.go
```

Build Docker image:

```bash
docker build -t cardwatcher-cli:latest apps/cli
```

### Web Frontend

The web frontend is a placeholder. When implemented:

```bash
cd apps/web
# Follow framework-specific build instructions
```

## Running Services

### Using Docker Compose

Start all services:

```bash
docker-compose up -d
```

This will start:
- MongoDB database
- Ntfy notification service
- CardWatcher server
- CardWatcher gateway

### Individual Services

#### Server

```bash
# Set environment variables (see .env.example)
export SERVER_PORT=50051
export MONGO_HOST=localhost
export MONGO_PORT=27017
export MONGO_DATABASE=cardwatcher
export CARDTRADER_ACCESS_TOKEN=your_token
export CARDTRADER_API_BASE_URL=https://api.cardtrader.com
export NTFY_HOST=localhost
export NTFY_PORT=80

cd apps/server
go run ./cmd/main.go
```

#### Gateway

```bash
export GATEWAY_PORT=8080
export GRPC_SERVER=localhost
export GRPC_PORT=50051

cd apps/gateway
go run ./cmd/main.go
```

#### CLI

```bash
cd apps/cli
./cardwatcher-cli --help
```

## API Usage

### gRPC (Direct)

Connect to the gRPC server on port 50051:

```bash
# Example using grpcurl
grpcurl -plaintext localhost:50051 list

# Example using the CLI
cd apps/cli
./cardwatcher-cli list-watches
```

### HTTP/JSON (via Gateway)

The gateway provides RESTful access to all gRPC endpoints:

```bash
# List expansions
curl http://localhost:8080/v1/expansions

# List blueprints
curl http://localhost:8080/v1/blueprints?expansion_id=123

# Create a watch
curl -X POST http://localhost:8080/v1/watches \
  -H "Content-Type: application/json" \
  -d '{
    "expansion_id": 123,
    "blueprint_id": 456,
    "condition": "CONDITION_NEAR_MINT",
    "foil": false
  }'

# List watches
curl http://localhost:8080/v1/watches

# Delete a watch
curl -X DELETE http://localhost:8080/v1/watches/{watch_id}
```

### OpenAPI Documentation

The generated OpenAPI spec is available in `gen/openapi/api.swagger.json`. You can view it using Swagger UI or other OpenAPI tools.

## Development Workflow

### 1. Make Proto Changes

Edit `.proto` files in `proto/cardwatcher/v1/`.

### 2. Lint Protobuf

```bash
./scripts/lint-proto.sh
```

### 3. Check Breaking Changes

```bash
./scripts/check-breaking.sh
```

### 4. Generate Code

```bash
./scripts/generate-proto.sh
```

### 5. Update Imports

If you changed the proto package or added new messages/services, update imports in your Go code.

### 6. Build and Test

```bash
cd apps/server
go build ./...
go test ./...

cd apps/gateway
go build ./...

cd apps/cli
go build ./...
```

## Deployment

### Docker Images

Build all images:

```bash
docker build -t cardwatcher-server:latest apps/server
docker build -t cardwatcher-gateway:latest apps/gateway
docker build -t cardwatcher-cli:latest apps/cli
```

### Push to Registry

```bash
docker tag cardwatcher-server:latest your-registry/cardwatcher-server:latest
docker push your-registry/cardwatcher-server:latest

docker tag cardwatcher-gateway:latest your-registry/cardwatcher-gateway:latest
docker push your-registry/cardwatcher-gateway:latest
```

### Kubernetes Deployment

Example deployment manifests would be placed in a `deploy/` directory:

```
deploy/
├── server-deployment.yaml
├── gateway-deployment.yaml
├── service.yaml
└── ingress.yaml
```

The gateway exposes the HTTP API, while the server communicates internally via gRPC.

## Testing

### Unit Tests

```bash
# Server tests
cd apps/server
go test ./...

# Gateway tests
cd apps/gateway
go test ./...

# CLI tests
cd apps/cli
go test ./...
```

### Integration Tests

```bash
cd apps/server
go test -tags=integration ./...
```

## Troubleshooting

### Proto Generation Issues

If `buf generate` fails:

1. Ensure buf is installed and up to date
2. Check that `buf.yaml` and `buf.gen.yaml` are valid
3. Verify proto file syntax is correct
4. Check for circular dependencies

### Import Path Issues

If you encounter import errors:

1. Ensure `go.work` is properly configured
2. Run `go work sync` to sync dependencies
3. Verify the generated code exists in `gen/go/`

### Module Dependency Issues

If modules can't find each other:

1. Ensure you're working from the repository root
2. Run `go work use ./apps/*` to ensure all modules are in the workspace
3. Check `go.mod` files have the correct `replace` directives

## Contributing

1. Follow the [Proto Development Guidelines](#proto-development-guidelines)
2. Always run `./scripts/lint-proto.sh` before committing
3. Run `./scripts/check-breaking.sh` before proposing changes
4. Generate code with `./scripts/generate-proto.sh` after proto changes
5. Update documentation as needed

## Resources

- [Buf Documentation](https://docs.buf.build/)
- [gRPC-Gateway Documentation](https://grpc-ecosystem.github.io/grpc-gateway/)
- [Go Workspace Documentation](https://go.dev/ref/mod#workspaces)
- [Protocol Buffers Style Guide](https://developers.google.com/protocol-buffers/docs/style)
