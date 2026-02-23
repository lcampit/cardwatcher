# Card Watcher

Your own watcher for ordering cards off [Cardtrader](https://www.cardtrader.com)

A Go-based monorepo for monitoring trading card prices with automated notifications. This project includes a gRPC server, HTTP gateway, CLI client, and a placeholder web frontend.

## 📚 Documentation

For detailed information about the monorepo structure, build steps, and deployment, see **[MONOREPO.md](./MONOREPO.md)**.

## 🏗️ Architecture

This monorepo follows a multi-module Go workspace layout with the following applications:

- **apps/server** - gRPC server implementing the CardWatcher service
- **apps/gateway** - gRPC-Gateway for HTTP/JSON REST API access
- **apps/cli** - Command-line client for interacting with the service
- **apps/web** - Web frontend (placeholder for future implementation)

## 🚀 Quick Start

### Prerequisites

- Go 1.24 or later
- Buf CLI (for Protocol Buffer management)
- Docker and Docker Compose

### Setup

1. Clone the repository
2. Install Buf: `brew install bufbuild/buf/buf` (macOS) or see [MONOREPO.md](./MONOREPO.md) for other platforms
3. Sync the Go workspace: `go work sync`
4. Generate protobuf code: `./scripts/generate-proto.sh`
5. Start services with Docker Compose: `docker-compose up -d`

### Building Applications

```bash
# Build the gRPC server
cd apps/server && go build -o cardwatcher-server ./cmd/main.go

# Build the gateway
cd apps/gateway && go build -o cardwatcher-gateway ./cmd/main.go

# Build the CLI
cd apps/cli && go build -o cardwatcher-cli ./cmd/main.go
```

### Running Services

```bash
# Start all services with Docker Compose
docker-compose up -d

# Or run individual services
cd apps/server && go run ./cmd/main.go
cd apps/gateway && go run ./cmd/main.go
```

## 📖 API Usage

### gRPC (Direct)

```bash
# Using the CLI
cd apps/cli
./cardwatcher-cli list-watches
```

### HTTP/JSON (via Gateway)

```bash
# List watches
curl http://localhost:8080/v1/watches

# Create a watch
curl -X POST http://localhost:8080/v1/watches \
  -H "Content-Type: application/json" \
  -d '{"expansion_id": 123, "blueprint_id": 456, "condition": "CONDITION_NEAR_MINT", "foil": false}'
```

## 🛠️ Development

### Protocol Buffers

```bash
# Lint proto files
./scripts/lint-proto.sh

# Check for breaking changes
./scripts/check-breaking.sh

# Generate code
./scripts/generate-proto.sh
```

For detailed protobuf guidelines and versioning patterns, see [MONOREPO.md](./MONOREPO.md#protocol-buffers).

## 📦 Tools Used

- **mise** - Development environment management
- **Buf** - Protocol Buffer management and linting
- **Air** - Live reload for Go applications
- **go-simpler/env** - Environment variable management
- **MongoDB** - Database
- **Ntfy** - Push notifications
- **gRPC-Gateway** - HTTP/JSON to gRPC reverse proxy

## 📝 Project Structure

```
.
├── apps/              # Application modules
│   ├── server/       # gRPC server
│   ├── gateway/      # HTTP gateway
│   ├── cli/          # CLI client
│   └── web/          # Web frontend
├── proto/            # Protocol Buffer definitions
│   └── cardwatcher/v1/
├── gen/              # Generated code
│   ├── go/           # Go generated code
│   ├── web/          # TypeScript types
│   └── openapi/      # OpenAPI specs
└── scripts/          # Utility scripts
```

## 🤝 Contributing

Please read [MONOREPO.md](./MONOREPO.md) for detailed contributing guidelines, especially regarding:
- Proto development guidelines
- Breaking change checks
- Code generation workflow

## 📄 License

[Add your license here]

## 🔗 Links

- [Cardtrader](https://www.cardtrader.com)
- [Buf Documentation](https://docs.buf.build/)
- [gRPC-Gateway Documentation](https://grpc-ecosystem.github.io/grpc-gateway/)

