# Protocol Buffers

This directory contains Protocol Buffer definitions for the CardWatcher API.

## Directory Structure

```
proto/
└── cardwatcher/
    └── v1/          # v1 API definitions
        ├── entities.proto
        ├── requests.proto
        ├── responses.proto
        └── service.proto
```

## Versioning

API versions are organized under `cardwatcher/vX/` directories. When making breaking changes:

1. Create a new version directory: `proto/cardwatcher/v2/`
2. Copy and modify proto files with new version
3. Update the `package` declaration: `package cardwatcher.v2;`
4. Update the `go_package` option: `option go_package = "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v2;cardwatcherv2";`
5. Update imports in applications

## Files

- **entities.proto**: Shared entity types (Condition, Expansion, Blueprint, Watch)
- **requests.proto**: Request message definitions
- **responses.proto**: Response message definitions
- **service.proto**: gRPC service definitions with HTTP annotations

## Generating Code

From the repository root:

```bash
./scripts/generate-proto.sh
```

Or directly:

```bash
buf generate
```

## Linting

```bash
./scripts/lint-proto.sh
```

## Breaking Change Detection

```bash
./scripts/check-breaking.sh
```

## Guidelines

1. **Never renumber fields** - Always add new fields at the end
2. **Use reserved ranges** - Reserve field numbers for future use
3. **Include HTTP annotations** - All RPC methods should have HTTP mappings for the gateway
4. **Follow naming conventions** - Use snake_case for fields, PascalCase for messages

See [CONTRIBUTING.md](../CONTRIBUTING.md) for detailed guidelines.
