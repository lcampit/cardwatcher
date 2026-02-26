# Card Watcher

Your own watcher for ordering cards off [Cardtrader](https://www.cardtrader.com).

Have you ever scrolled through the marketplace for a card over and
over again to check its price debating whether to buy it or not?
Worry not! Cardwatcher can do that for you, and much more!

Cardwatcher is a Go-based application for monitoring trading card prices
with automated notifications.
This project includes a gRPC server, HTTP gateway, CLI client, and
a web frontend.

## Getting Started

These instructions will get you a copy of the project up and running on your
local machine for development and testing purposes.
See deployment for notes on how to deploy the project on a live system.
For detailed information about the monorepo structure
see **[MONOREPO.md](./MONOREPO.md)**.

Cardwatcher repo uses [mise-en-place](https://mise.jdx.dev/) to
handle dependencies, tooling, environment and commonly
used commands.
With mise installed, setting up the project on you own machine
is as simple as:

```bash
git clone https://github.com/lcampit/cardwatcher && cd cardwatcher
mise install
```

This will take care of SDKs and environment variables,
as well as multiple tools that I use when working on the
project.

To initialize Go workspaces and start working on the
Go part of the application, run `go work sync`.

`mise run` will provide a list of detailed and descriptive
scripts that the project needs.

Contributions are always welcome! Once your work is done,
push your branch and open a pull request. For more details about
how to contribute, see **[CONTRIBUTING.md](./CONTRIBUTING.md)**.

## Architecture

The cardwatcher project repository follows a monorepo
structure using a multi-module Go workspace layout. The main applications
included are:

- Server: a gRPC server that wraps Cardwatcher's API and handles
database interactions;
- CLI: a cli application to quickly interact with the server
- Gateway: a gRPC generated REST-to-gRPC gateway that maps gRPC
endpoints to REST ones;
- Client: a webapp client that interacts with the server through the gateway.

Cardwatcher uses a document database such as [MongoDB](https://www.mongodb.com/)
to maintain user data.

## MakeFile

- **apps/server** - gRPC server implementing the CardWatcher service
- **apps/gateway** - gRPC-Gateway for HTTP/JSON REST API access
- **apps/cli** - Command-line client for interacting with the service
- **apps/web** - Web frontend (placeholder for future implementation)

Run build make command with tests

## 🚀 Quick Start

```bash
make all
```

### Prerequisites

Build the application

- Go 1.24 or later
- Buf CLI (for Protocol Buffer management)
- Docker and Docker Compose

```bash
make build
```

### Setup

Run the application

1. Clone the repository
2. Install Buf: `brew install bufbuild/buf/buf` (macOS) or see [MONOREPO.md](./MONOREPO.md) for other platforms
3. Sync the Go workspace: `go work sync`
4. Generate protobuf code: `./scripts/generate-proto.sh`
5. Start services with Docker Compose: `docker-compose up -d`

### Building Applications

```bash
make run
```

# Build the gRPC server

cd apps/server && go build -o cardwatcher-server ./cmd/main.go

Create DB container

# Build the gateway

cd apps/gateway && go build -o cardwatcher-gateway ./cmd/main.go

```bash
make docker-run
# Build the CLI
cd apps/cli && go build -o cardwatcher-cli ./cmd/main.go
```

Shutdown DB Container

### Running Services

```bash
make docker-down
# Start all services with Docker Compose
docker-compose up -d

# Or run individual services
cd apps/server && go run ./cmd/main.go
cd apps/gateway && go run ./cmd/main.go
```

DB Integrations Test:

## 📖 API Usage

### gRPC (Direct)

```bash
make itest
# Using the CLI
cd apps/cli
./cardwatcher-cli list-watches
```

Live reload the application:

### HTTP/JSON (via Gateway)

```bash
make watch
# List watches
curl http://localhost:8080/v1/watches

# Create a watch
curl -X POST http://localhost:8080/v1/watches \
  -H "Content-Type: application/json" \
  -d '{"expansion_id": 123, "blueprint_id": 456, "condition": "CONDITION_NEAR_MINT", "foil": false}'
```

Run the test suite:

## 🛠️ Development

### Protocol Buffers

```bash
make test
# Lint proto files
./scripts/lint-proto.sh

# Check for breaking changes
./scripts/check-breaking.sh

# Generate code
./scripts/generate-proto.sh
```

Clean up binary from the last build:
For detailed protobuf guidelines and versioning patterns, see [MONOREPO.md](./MONOREPO.md#protocol-buffers).

```bash
make clean
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

