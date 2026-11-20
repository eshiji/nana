# Contributing

## Workflow

1. Write a failing test that describes the change.
2. Make it pass with the smallest change, then refactor.
3. Run `make test cover lint` before opening a pull request.

## Commits and pull requests

Pull request titles follow [Conventional Commits](https://www.conventionalcommits.org/)
(`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`), because the squash-merge title
drives the next release version. Mark breaking changes with `!` (for example `feat!: ...`).

## Architecture decisions

Significant decisions (hard to reverse, or affecting many parts of the code) are recorded
as ADRs in [`docs/adr/`](docs/adr/). ADRs are immutable: to change a decision, add a new
ADR that supersedes the old one and update only the old one's status.

## Changing the CLI

Every command needs a `Short` description, an `Example` and usage text for each flag; a
test enforces this. When `--help` output changes on purpose, run `make golden` and review
the diff of the golden files.
