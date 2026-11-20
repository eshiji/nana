# ADR-0001: Record architecture decisions

- Status: Accepted
- Date: 2026-11-20

## Context

Decisions about this project's architecture should be explained where the code lives, so
that future contributors (people or agents) understand why things are the way they are.

## Decision

Significant decisions are recorded as Architecture Decision Records in `docs/adr/`,
following [Michael Nygard's format](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions):
status, context, decision and consequences, plus the alternatives considered.

- Files are numbered: `NNNN-short-title.md`.
- ADRs are immutable. A changed decision gets a new ADR that supersedes the old one; the
  old one only has its status updated.
- A pull request that changes a decision includes the new ADR.

## Consequences

- The reasoning behind the design is versioned together with the code.
- Pull requests stay short: they link the ADR instead of re-explaining the design.
