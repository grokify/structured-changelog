# TRD — INIT-SCHANGELOG-001: RMI Trailer Flow — Commit Trailers to Changelog Entries

**Initiative:** INIT-SCHANGELOG-001
**Status:** completed

## T1 — Trailer extraction (gitlog)

- Read git trailers via `git log --format` with `%(trailers:key=Refs,valueonly)`
  (or `Interpret-trailers`-compatible body parsing for log sources that
  lack trailer support).
- Extract IDs matching `^RMI-[A-Z0-9]+-\d+$`; a commit may carry multiple
  `Refs:` trailers.
- Add `RMIs []string` to the parsed-commit type; include in both TOON and
  JSON output of `parse-commits`.
- The existing `#123` issue-ref extraction is untouched.

## T2 — Entry schema (additive, non-breaking)

- Add `RMIs []string \`json:"rmis,omitempty"\`` to `Entry` beside the
  legacy `RMI string`; add `WithRMIs(...string)` builder.
- Readers treat `rmi` as equivalent to a one-element `rmis`; writers prefer
  `rmis`. No deprecation of `rmi` in this initiative.
- Regenerate the JSON Schema (Go-first: structs → invopop/jsonschema →
  schemakit lint → go:embed); commit types and schema together.
- `validate`: warn on strings in `rmis` not matching the RMI pattern; warn
  when both `rmi` and `rmis` are set and disagree.

## T3 — Generation flow

- When building suggested entries from parsed commits, union the commits'
  RMI IDs into the entry's `rmis` (deduplicated, sorted).
- `initiative` is never auto-populated: with an RMI present it is
  derivable downstream (VisionStudio resolves RMI → initiative); guidance
  documents `initiative` only for entries with no RMI.

## T4 — Downstream consumption (informational)

VisionStudio ingest (INIT-VISIONSTUDIO-006, RMI-VISIONSTUDIO-303) gains a
precedence order once this ships:

1. Entry-level `rmis` in CHANGELOG.json — human-curated, highest precision.
2. Trailer-chain walk of the tag range — automatic fallback.
3. Cross-check between the two; disagreement flags curation errors.

No structured-changelog code depends on VisionStudio — the spec stays
tracker-agnostic (RMI IDs are pattern-validated strings).

## Risks

1. **Trailer formats vary** (`Refs:` casing, multiple trailers, squashed
   merges concatenating bodies) — parse leniently, validate strictly on
   output.
2. **Dual singular/plural fields** invite confusion — validate rule (T2)
   plus docs; consider deprecating `rmi` in a later major only.
