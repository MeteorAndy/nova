# Domain Docs

This repository uses a single-context domain documentation layout.

Before related code exploration, read the root `CONTEXT.md` and relevant ADRs
under `docs/adr/` when they exist. If they do not exist, proceed silently;
domain-modeling skills create them lazily when terminology or decisions are
actually resolved.

Use terminology defined in `CONTEXT.md` consistently in issues, proposals,
tests, and code. When a required concept is missing, reconsider whether it is
new vocabulary or a genuine domain-model gap.

Surface conflicts with existing ADRs explicitly instead of silently overriding
recorded decisions.
