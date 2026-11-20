# ADR-0002: nana is a platform-agnostic Go CLI, released as checksummed, attested binaries

- Status: Proposed
- Date: 2026-11-20
- Deciders: project author

## Context

`nana` is a command-line tool for the checks a delivery pipeline runs on infrastructure
code and container images: test and lint container images (including suggestions to make
them smaller and safer), and validate Terraform. The first consumer is
[terraform-container-image](https://github.com/eshiji/terraform-container-image).

Pipelines run on different platforms (GitHub Actions today; maybe GitLab CI or Jenkins
later) and on laptops. When the logic lives inside one platform's YAML, it's duplicated
per platform, it can't run locally, and it's hard to test.

Prior art, studied for ideas only, no code copied: CLIs that are released on their own and
wrapped by thin per-platform integrations, such as
[trivy](https://github.com/aquasecurity/trivy) +
[trivy-action](https://github.com/aquasecurity/trivy-action),
[golangci-lint](https://github.com/golangci/golangci-lint) +
[golangci-lint-action](https://github.com/golangci/golangci-lint-action), and
[terraform](https://github.com/hashicorp/terraform) +
[setup-terraform](https://github.com/hashicorp/setup-terraform).

## Decision

1. **A standalone CLI with no platform code.** `nana` doesn't read `GITHUB_*` or any other
   CI variable and doesn't call any CI API. Platform integrations live in their own
   repositories (the first one is
   [github-workflows](https://github.com/eshiji/github-workflows), see its ADR-0001) and
   only install `nana` and pass arguments.
2. **Written in Go.**
   - It ships as one static binary for linux/darwin/windows × amd64/arm64, with no runtime
     to install on the runner.
   - The domain's libraries are Go-native: the
     [BuildKit Dockerfile parser](https://pkg.go.dev/github.com/moby/buildkit/frontend/dockerfile/parser),
     [go-containerregistry](https://github.com/google/go-containerregistry) (inspect images
     without pulling them) and [hcl](https://github.com/hashicorp/hcl). Using them as
     libraries removes most of the work of calling external tools and parsing their output.
3. **Small subcommands grouped by subject, with sensible defaults**, built with
   [cobra](https://github.com/spf13/cobra):
   `nana image lint`, `nana image test`, `nana terraform validate`. No do-everything
   command steered by flags. Output is human-readable by default, with `--json` for
   machines and documented exit codes (0 ok, 1 findings, 2 usage or runtime error).
4. **Layered code, with dependencies pointing inward:**

   | Layer | Path | Contains |
   | --- | --- | --- |
   | Domain | `internal/<subject>` (e.g. `internal/dockerfile`) | Pure rules and types, no I/O |
   | Ports | interfaces in the domain packages | e.g. `ImageInspector`, `CommandRunner` |
   | Adapters | `internal/adapters/...` | registry, `exec`, filesystem implementations |
   | Delivery | `internal/cli`, `cmd/nana` | cobra commands: parse flags, call the domain, print |

5. **Test-first.** Every behaviour starts as a failing test.
   - Unit tests: table-driven tests for the domain.
   - Integration tests: the adapters.
   - End-to-end tests: the CLI, with
     [testscript](https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript).
   - Coverage target: 80% or more, reported in CI.
   - [golangci-lint](https://golangci-lint.run/) and `go vet` gate every PR.
6. **One version, released automatically.**
   [release-please](https://github.com/googleapis/release-please) (release type `go`)
   reads [Conventional Commits](https://www.conventionalcommits.org/) and opens the
   release PR. Merging it creates tag `vX.Y.Z`, and then
   [GoReleaser](https://goreleaser.com/) publishes the binaries, `checksums.txt` and
   [build provenance attestations](https://docs.github.com/en/actions/how-tos/secure-your-work/use-artifact-attestations/use-artifact-attestations).
   Releases are [immutable](https://docs.github.com/en/code-security/concepts/supply-chain-security/immutable-releases):
   a bad release is never deleted or re-tagged. It gets a "do not use" note and a new
   patch release.
7. **Semver honestly.** An intentional change of behaviour is marked breaking (`feat!:` /
   `fix!:`). An accidental break is a bug, fixed by a patch release that starts with a test
   reproducing it.
8. **Independent work.** Built on the author's own time and equipment, from public
   documentation, with no code, names, schemas, configuration or design taken from any
   employer or client.

## Alternatives considered

- **Logic in CI YAML and shell.** It can't run locally, it's hard to test, and it's
  duplicated per platform.
- **Python CLI.** Faster to write glue code in, but it needs a Python runtime and a
  dependency install on every runner, and the container and Terraform libraries are
  Go-first.
- **Container image as the distribution format.** It pins the toolchain, but it needs a
  registry and an image pull per run, and it doesn't run natively on macOS or Windows.
  An image may be published later as an extra artifact, not as the main one.

## Consequences

- Glue code is more verbose in Go. The compiler and explicit error handling make up for
  it, and so do libraries replacing calls to external tools.
- Every platform integration pins one exact `nana` version and is updated by Renovate
  (see github-workflows ADR-0001).
- A release pipeline (release-please + GoReleaser) has to exist before the first consumer
  can install a binary. Until then, consumers build from source with `go install`.

## Open questions

- First command for terraform-container-image: to be decided together with its test
  structure, before any code.
