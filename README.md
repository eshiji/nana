# nana

Pipeline checks for container images and Terraform, the same on a laptop and in CI.

## Install

```sh
go install github.com/eshiji/nana/cmd/nana@latest
```

## Usage

```sh
nana --help
nana version
```

## Development

Requires Go 1.26+ and [golangci-lint](https://golangci-lint.run/).

```sh
make help            # list targets
make test cover lint # what CI runs
make build           # bin/nana
```

See [CONTRIBUTING.md](CONTRIBUTING.md) and the [architecture decisions](docs/adr/).

## License

[Apache-2.0](LICENSE)
