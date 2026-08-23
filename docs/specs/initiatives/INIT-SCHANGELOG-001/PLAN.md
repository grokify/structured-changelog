# PLAN — INIT-SCHANGELOG-001: RMI Trailer Flow — Commit Trailers to Changelog Entries

## Sequencing Rationale

Single-phase initiative — the four RMIs form one short pipeline
(parse → schema → generate → document) and each is independently
shippable behind the previous one:

1. Trailer parsing first: it's pure gitlog work, testable against real
   repos immediately (visionstudio and acts commit history carry trailers
   today).
2. Schema second: additive `rmis` field + regenerated JSON Schema; no
   behavior depends on it yet.
3. Generation third: wires 1 into 2 — suggested entries arrive
   pre-populated.
4. Docs/release last: convention (`rmis` yes, `initiative` only without an
   RMI), examples, release notes.

## Working Agreements

- Commits carry `Refs: RMI-SCHANGELOG-10N` — this initiative's own
  changelog entries become the first real `rmis` usage (dogfood).
- Additive only: no breaking changes to Entry or the schema; `rmi`
  singular stays.
- Go-first schema flow: structs → invopop/jsonschema → schemakit lint →
  go:embed; types and schema commit together.

## Dependencies and Coordination

- **INIT-VISIONSTUDIO-006 (RMI-VISIONSTUDIO-303):** consumes entry-level
  `rmis` with precedence over the trailer walk once both ship. No blocking
  dependency either direction; coordinate the precedence/cross-check
  behavior when 303 is implemented.
- **Release skill / org CLAUDE.md:** once shipped, the changelog step of
  the release ritual produces RMI-annotated entries automatically — no
  operating-model doc changes needed beyond an example.
