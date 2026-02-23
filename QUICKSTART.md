# CardWatcher Monorepo Quick Start

A quick reference guide for common tasks in the CardWatcher monorepo.

## Setup

```bash
# Clone and enter the directory
git clone https://github.com/lcampit/cardwatcher.git
cd cardwatcher

# Install Buf (required for proto generation)
# macOS: brew install bufbuild/buf/buf
# Linux: See https://docs.buf.build/installation

# Sync Go workspace
go work sync

# Generate protobuf code
./scripts/generate-proto.sh
```

## Common Commands

### Proto Development

```bash
# Generate all protobuf artifacts
./scripts/generate-proto.sh

# Lint proto files
./scripts/lint-proto.sh

# Check for breaking changes
./scripts/check-breaking.sh

# Or use buf directly
buf generate
buf lint
buf breaking --against '.git#main'
```

### Building Applications

```bash
# Server (gRPC)
cd apps/server && go build -o cardwatcher-server ./cmd/main.go

# Gateway (HTTP/JSON)
cd apps/gateway && go build -o cardwatcher-gateway ./cmd/main.go

# CLI
cd apps/cli && go build -o cardwatcher-cli ./cmd/main.go

# Build all from root
go work build ./...
```

### Running Services

```bash
# Using Docker Compose (recommended)
docker-compose up -d

# Individual services
cd apps/server && go run ./cmd/main.go
cd apps/gateway && go run ./cmd/main.go
```

### Testing

```bash
# Run all tests
go work test ./...

# Server tests with integration tests
cd apps/server
go test ./...
go test -tags=integration ./...
```

### Docker Images

```bash
# Build images
docker build -t cardwatcher-server:latest apps/server
docker build -t cardwatcher-gateway:latest apps/gateway
docker build -t cardwatcher-cli:latest apps/cli

# Run server container
docker run -p 50051:50051 cardwatcher-server:latest

# Run gateway container
docker run -p 8080:8080 -e GRPC_SERVER=host.docker.internal cardwatcher-gateway:latest
```

## API Endpoints

### gRPC (Server on port 50051)

```bash
# Using CLI
cd apps/cli
./cardwatcher-cli list-watches
./cardwatcher-cli list-expansions
./cardwatcher-cli list-blueprints --expansion-id 123
./cardwatcher-cli save-watch --expansion-id 123 --blueprint-id 456 --condition NM --foil false
./cardwatcher-cli delete-watch --watch-id <id>
```

### HTTP/JSON (Gateway on port 8080)

```bash
# List expansions
curl http://localhost:8080/v1/expansions

# List blueprints
curl http://localhost:8080/v1/blueprints?expansion_id=123

# Create watch
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

# Delete watch
curl -X DELETE http://localhost:8080/v1/watches/{watch_id}
```

## Project Structure

```
.
├── apps/
│   ├── server/       # gRPC server app
│   ├── gateway/      # HTTP gateway app
│   ├── cli/          # CLI client app
│   └── web/          # Web frontend (placeholder)
├── proto/            # Protocol Buffer definitions
│   └── cardwatcher/v1/
├── gen/              # Generated code
│   ├── go/           # Go generated code
│   ├── web/          # TypeScript types
│   └── openapi/      # OpenAPI specs
├── scripts/          # Utility scripts
├── go.work           # Go workspace
├── buf.yaml          # Buf configuration
└── buf.gen.yaml      # Buf generation config
```

## Adding a New Proto Field

1. Edit `.proto` file in `proto/cardwatcher/v1/`
2. Add field at the end with new number
3. Run `./scripts/lint-proto.sh`
4. Run `./scripts/check-breaking.sh`
5. Run `./scripts/generate-proto.sh`
6. Update Go code to use new field

## Troubleshooting

### Import errors after proto changes

```bash
# Regenerate proto code
./scripts/generate-proto.sh

# Sync workspace
go work sync
```

### Build failures

```bash
# Clean and rebuild
go clean -cache
go work build ./...
```

### Docker build fails

```bash
# Build without cache
docker build --no-cache -t cardwatcher-server:latest apps/server
```

## Documentation

- **[MONOREPO.md](./MONOREPO.md)** - Comprehensive monorepo documentation
- **[CONTRIBUTING.md](./CONTRIBUTING.md)** - Contribution guidelines
- **[README.md](./README.md)** - Project overview

## Support

- Open an issue for bugs or questions
- Start a discussion for proposals
- Check existing documentation first
