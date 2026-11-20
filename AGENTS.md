# AGENTS.md

Instructions for AI coding agents (and humans) working on `nana`.

## Project

Pipeline checks for container images and Terraform, the same on a laptop and in CI.

A Go command-line tool. `cmd/nana/main.go` only calls
`internal/cli.Execute`; all behaviour lives under `internal/`.

| Path | Holds |
| --- | --- |
| `internal/cli` | cobra commands: parse flags, call the application, print. No business logic. |
| `internal/<subject>` | Domain logic and its ports (interfaces). No I/O. |
| `internal/adapters/...` | Implementations of ports (filesystem, network, external tools). |
| `internal/buildinfo` | The version the binary was built from. |
| `docs/adr` | Architecture decisions (immutable; supersede, don't edit). |

## Commands

| Command | Does |
| --- | --- |
| `make test` | All tests with the race detector |
| `make cover` | Tests, failing below 80% coverage |
| `make lint` | golangci-lint (config in `.golangci.yml`) |
| `make fmt` | Format the code |
| `make build` | Build `bin/nana` (`VERSION=vX.Y.Z` stamps a version) |
| `make golden` | Regenerate `--help` golden files after an intended CLI change |

## Rules

- **Test first.** Write a failing test, make it pass, refactor. Table-driven tests for
  logic, adapter tests for I/O, testscript (`cmd/*/testdata/script/*.txtar`) for the binary.
- **CLI contract.** Every command has `Short`, `Example` and usage text for every flag
  (enforced by `internal/cli/contract_test.go`). Exit codes: 0 ok, 1 findings, 2 error.
- **Dependencies point inward:** `cli` → domain ← adapters. The domain imports no
  framework or I/O package.
- **Errors** are returned and wrapped with context (`fmt.Errorf("doing x: %w", err)`); no
  `panic` outside `main`.
- **Security:** validate input at the CLI boundary; run external tools with
  `exec.Command(name, args...)`, never through a shell or with string-built commands.
- **No logic in the Makefile or shell scripts.** Make targets only call Go tooling.
- **English** for code, comments, docs and commit messages.
- **Commits / PR titles** follow Conventional Commits; breaking changes use `!`.

## Template-managed files

This project was generated from a cookiecutter template and is kept up to date with
[cruft](https://cruft.github.io/cruft/) (`.cruft.json`). Files outside `cmd/`,
`internal/`, `docs/`, `README.md`, `CHANGELOG.md`, `go.mod` and `go.sum` belong to the
template: change them in the template, then run `cruft update` here.
