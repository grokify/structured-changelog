# RMI Trailer Flow — Commit Trailers to Changelog Entries — Roadmap

**Initiative:** `INIT-SCHANGELOG-001`
**Repository:** `github.com/grokify/structured-changelog`

> RMI IDs are stable and permanent. Commits implementing an item carry the trailer `Refs: RMI-<REPOSLUG>-<NNN>`. Phase status is derived from member RMIs — a phase is complete only when all its required RMIs are complete.

This initiative uses the RMI-SCHANGELOG-1xx block (0xx carries cross-repo items from other initiatives, e.g. INIT-DEVFOLIOQR-001).

## Phase 1 — Trailer-to-Entry Pipeline

**Theme:** Refs trailers flow from commits through parse-commits into pre-populated changelog entries; additive rmis field; convention documented

- [x] `RMI-SCHANGELOG-101` Refs trailer parsing in gitlog
  - git log trailers key=Refs valueonly; pattern ^RMI-[A-Z0-9]+-\d+$; multiple trailers per commit; RMIs []string on parsed commits in TOON + JSON output; existing #123 issue-ref extraction untouched
- [x] `RMI-SCHANGELOG-102` Plural rmis field on Entry
  - Additive rmis []string beside legacy rmi; WithRMIs builder; validate warns on malformed IDs (W006) and rmi/rmis disagreement (W007). Schema: the roadmap block (rmis, rmi, initiative) was hand-added to changelog-v1.schema.json and lint-verified — the hand-crafted draft-07 published schema is richer than current invopop output, so wholesale generator adoption (with the new `schemakit generate --check` drift guard) is deferred to a follow-up
- [x] `RMI-SCHANGELOG-103` Pre-populated entries in generation flow
  - Suggested entries union underlying commits' RMI IDs (deduped, sorted) into rmis; initiative never auto-populated (derivable downstream when RMI present)
- [x] `RMI-SCHANGELOG-104` Convention docs, examples, and release
  - Document rmis-yes / initiative-only-without-RMI convention; real example from own release (dogfood); CHANGELOG + semver release (v0.17.0)
