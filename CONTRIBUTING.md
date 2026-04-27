# Contributing to CardWatcher

Thank you for your interest in contributing to CardWatcher!
We welcome contributions of all kinds — bug reports,
feature requests, documentation improvements,
and code contributions.

## Project Setup

CardWatcher uses [mise](https://mise.jdx.dev/) to manage development
tools, SDKs, and environment variables.
This ensures consistent tooling across all contributors.

### Getting Started

1. Clone the repository:

   ```bash
   git clone https://github.com/lcampit/cardwatcher && cd cardwatcher
   ```

2. Install all dependencies and tools:

   ```bash
   mise install
   ```

3. Install git hooks with [lefthook](https://lefthook.dev/)

```bash
  lefthook install
```

That's it! All required tools and dependencies will be automatically
configured by mise.

### Available Tasks

Most project operations are handled through mise tasks. To see all available tasks:

```bash
mise run
```

This will display a list of all defined tasks with descriptions.
Common tasks include:

- Building applications components
- Running tests
- Generating protobuf code
- Linting

Each task is defined in a standalone *executable* file in the
`.mise/tasks` directory.
Tasks are grouped based on the component they work on: i.e.
`server:run`, `server:build` and `cli:run`.
When adding a new task, make sure to maintain
this grouping. Tasks that are not in any subdirectory
are considered generic and should do something related
to the application as a whole rather than a single
component.

## Project Layout

```
.
├── apps/                   # Application modules
│   ├── server/             # gRPC server
│   ├── gateway/            # REST-to-gRPC gateway
│   ├── cli/                # Command-line client
│   └── client/             # Web frontend
├── proto/                  # Protocol Buffer definitions
│   └── cardwatcher/v1/     # API definitions
├── gen/                    # Generated code
│   └── go/                 # Go generated code (from proto files)
├── docker-compose.yml
└── mise.toml               # Project configuration
```

When developing a new feature,

### Key Directories

- **`apps/`**: Contains each application as an independent module,
each with its own `Dockerfile` for containerization.
- **`proto/`**: Protocol Buffer definitions that define the API contracts.
- **`gen/`**: Generated code from proto files,
organized by language (e.g., `gen/go/` for Go).
- **`docker-compose.yml`**: Orchestrates all services (server, gateway,
database, notifications).

## Making Changes

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/my-feature`)
3. **Make** your changes
4. **Test** your changes (run `mise run` to see available test commands)
5. **Commit** your changes (`git commit -am 'Add new feature'`)
6. **Push** to the branch (`git push origin feature/my-feature`)
7. **Open** a Pull Request

## Code Style

Github Actions and git hooks are in place to
ensure a standard formatting and style across the
repo.
When not sure, refer to `mise run` commands to run
formatters and linters.

## Questions?

If you have questions about contributing, feel free to open
an issue or reach out. We're happy to help you get started!
