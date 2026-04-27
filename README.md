# Card Watcher

Your own watcher for ordering cards off [Cardtrader](https://www.cardtrader.com).

Have you ever scrolled through the marketplace for a card over and
over again to check its price debating whether to buy it or not?
Worry not! Cardwatcher can do that for you, and much more!

Cardwatcher is a Go-based application for monitoring trading card prices
with automated notifications.
This project includes a gRPC server, HTTP gateway, CLI client, and
a web frontend.

## Architecture

The cardwatcher project repository follows a monorepo
structure. The main applications included are:

- Server: a gRPC server that wraps Cardwatcher APIs and handles
database interactions;
- CLI: a cli application to quickly interact with the server
- Gateway: a gRPC generated REST-to-gRPC gateway that maps gRPC
endpoints to REST ones;
- Client: a webapp client that interacts with the server through the gateway.

Other supporting applications used are:

- [NTFY](https://ntfy.sh/) to handle mobile notifications;
- [MongoDB](https://www.mongodb.com/) to maintain user data.

## Getting Started

These instructions will get you a copy of the project up
and running on your local machine for development and
testing purposes.

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

`mise run` will provide a list of detailed and descriptive
scripts that you may need to work on the project.

Contributions are always welcome! Once your work is done,
push your branch and open a pull request. For more details about
how to contribute, see **[CONTRIBUTING.md](./CONTRIBUTING.md)**.

### Building and Running

Each component has its own `build` script handled by
mise. These scripts will handle dependencies download
and build the relative binary in the `bin/` folder.

Alternatively, you may directly use the `run` script of
each component to have it automatically built and ran.

A comprehensive docker-compose file is available to
quickly build and spin up the whole application
with `docker compose up`.

## 🤝 Contributing

Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed
contributing guidelines.

## 🔗 Links

- [Cardtrader](https://www.cardtrader.com) and
its [APIs](https://www.cardtrader.com/en/docs/api)
- [Buf Documentation](https://docs.buf.build/)
- [gRPC-Gateway Documentation](https://grpc-ecosystem.github.io/grpc-gateway/)
