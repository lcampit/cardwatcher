# Contributing to CardWatcher

Thank you for your interest in contributing to CardWatcher! This document provides guidelines for contributing to the monorepo.

## Development Workflow

### 1. Fork and Clone

Fork the repository and clone your fork locally:

```bash
git clone https://github.com/your-username/cardwatcher.git
cd cardwatcher
```

### 2. Set Up Go Workspace

Ensure you have Go 1.24 or later installed:

```bash
go version
```

Sync the Go workspace:

```bash
go work sync
```

### 3. Install Buf

Follow the instructions in [MONOREPO.md](./MONOREPO.md) to install the Buf CLI.

### 4. Generate Proto Code

Generate the protobuf code:

```bash
./scripts/generate-proto.sh
```

## Making Changes

### Protocol Buffer Changes

If you need to modify protobuf definitions:

1. **Make changes** to `.proto` files in `proto/cardwatcher/v1/`
2. **Lint** the changes: `./scripts/lint-proto.sh`
3. **Check for breaking changes**: `./scripts/check-breaking.sh`
4. **Generate code**: `./scripts/generate-proto.sh`
5. **Update imports** in Go code if needed
6. **Test** your changes

### Go Code Changes

For Go code changes:

1. Make your changes in the appropriate app directory
2. Run tests: `go test ./...`
3. Ensure code passes linters
4. Update documentation if needed

### Creating a New Proto Version

When making breaking changes that require a new version:

1. Create `proto/cardwatcher/v2/` directory
2. Copy and modify proto files with new version
3. Update `go_package` option to point to `gen/go/cardwatcher/v2`
4. Update `buf.gen.yaml` to generate v2 code
5. Update imports in applications to use v2
6. Update documentation

## Proto Development Guidelines

### Field Numbering

1. **Never renumber existing fields** - This is a breaking change
2. **Add new fields at the end** of messages
3. **Use reserved ranges** for future additions

Example:

```protobuf
message MyMessage {
  string name = 1;
  uint64 id = 2;

  // Reserved field numbers 3-9 for future use
  reserved 3 to 9;

  // New field added at the end
  string description = 10;
}
```

### Field Naming

- Use `snake_case` for field names
- Use descriptive names that convey the purpose
- Avoid abbreviations unless widely understood

### Enum Values

- Use `ENUM_NAME_VALUE_FORMAT` for enum values
- Always include a `_UNSPECIFIED` value with value 0
- Don't remove enum values; deprecate them instead

Example:

```protobuf
enum Condition {
  CONDITION_UNSPECIFIED = 0;
  CONDITION_NEAR_MINT = 1;
  CONDITION_SLIGHTLY_PLAYED = 2;
  // Don't remove values when deprecating
}
```

### Service Methods

- Use descriptive RPC names that indicate the action
- Keep request/response messages in separate files
- Include `google/api/annotations.proto` for HTTP mapping

Example:

```protobuf
service CardWatcher {
  rpc GetWatch(GetWatchRequest) returns (GetWatchResponse) {
    option (google.api.http) = {
      get: "/v1/watches/{watch_id}"
    };
  }
}
```

### Compatibility Checks

Before committing:

```bash
# Lint
buf lint

# Check breaking changes against main
buf breaking --against '.git#main'

# Or check against previous commit
buf breaking --against '.git#HEAD^'
```

## Testing

### Unit Tests

Write unit tests for new functionality:

```go
func TestMyFunction(t *testing.T) {
    result := MyFunction(input)
    expected := "expected value"

    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}
```

### Integration Tests

Integration tests use testcontainers:

```bash
cd apps/server
go test -tags=integration ./...
```

### Running All Tests

```bash
# From repository root
go work test ./...
```

## Code Style

### Go Code

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Run `gofmt` on all Go files
- Use `golangci-lint` for additional checks

### Proto Code

- Follow [Google API Style Guide](https://cloud.google.com/apis/design)
- Run `buf lint` before committing
- Use clear, descriptive comments for messages and fields

## Commit Messages

Follow conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `proto`: Protocol buffer changes

Examples:

```
feat(server): add watch deletion endpoint

Implement DeleteWatch RPC to allow users to remove
watches by ID.

Closes #123
```

```
fix(proto): correct field type in Watch message

Change condition field from string to enum to match
API specification.

Fixes #456
```

## Pull Request Process

1. **Update documentation** if needed
2. **Add tests** for new functionality
3. **Ensure all tests pass**
4. **Run linting**: `buf lint` and `golangci-lint`
5. **Update CHANGELOG** if applicable
6. **Create pull request** with clear description
7. **Address review feedback**

## Pull Request Checklist

- [ ] Code follows project style guidelines
- [ ] Tests pass locally
- [ ] Proto changes are linted (`buf lint`)
- [ ] No breaking changes (or documented with migration guide)
- [ ] Documentation is updated
- [ ] Commit messages follow conventional commits
- [ ] PR description is clear and complete

## Questions?

- Check [MONOREPO.md](./MONOREPO.md) for detailed project documentation
- Open an issue for questions or bugs
- Start a discussion for proposals and ideas

Thank you for contributing! 🎉
