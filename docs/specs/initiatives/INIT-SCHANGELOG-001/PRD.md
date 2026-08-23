# PRD — INIT-SCHANGELOG-001: RMI Trailer Flow — Commit Trailers to Changelog Entries

**Initiative:** INIT-SCHANGELOG-001
**Status:** proposed
**Type:** feature
**Workflow:** pbhq-lite

## Problem

The spec already has what we need but nothing feeds it:

1. `Entry.rmi` and `Entry.initiative` exist in the changelog schema
   (entry.go, "PRISM Control metadata") with `WithRMI`/`WithInitiative`
   builders — but **zero** CHANGELOG.json files across the grokify /
   plexusone / ProductBuildersHQ orgs populate them.
2. Commits carry `Refs: RMI-<SLUG>-<NNN>` git trailers per org convention,
   but `schangelog parse-commits` only extracts `#123`-style issue refs
   (gitlog/conventional.go) — trailers are dropped on the floor.
3. A curated entry often summarizes several commits spanning multiple RMIs;
   the single-string `rmi` field cannot express that.

VisionStudio release ingest (INIT-VISIONSTUDIO-006) currently derives
release↔initiative associations by walking git history for trailers.
Entry-level RMIs in CHANGELOG.json would be a higher-precision,
human-curated association source — if the pipeline populated them.

## Goals

1. **G1 — Trailer parsing.** `parse-commits` reads `Refs:` git trailers and
   carries RMI IDs on parsed commits (TOON and JSON output).
2. **G2 — Plural `rmis`.** Additive `rmis []string` alongside the legacy
   `rmi` string; schema regenerated; validate checks RMI ID pattern.
3. **G3 — Pre-populated entries.** The changelog-generation flow suggests
   entries with `rmis` filled from the underlying commits' trailers; the
   human curates rather than transcribes.
4. **G4 — Convention documented.** Populate `rmis`; leave `initiative`
   empty when an RMI is present (derivable via VisionStudio — storing both
   invites drift); use `initiative` only for entries with no RMI.

## Non-Goals

- Changing VisionStudio ingest itself (that is INIT-VISIONSTUDIO-006
  RMI-VISIONSTUDIO-303; it gains entry-level RMIs as a source when this
  ships).
- Removing or deprecating the legacy `rmi` singular field (back-compat).
- Any org-specific coupling in the public spec — RMI IDs remain plain
  strings validated by pattern, usable by any tracker with similar IDs.

## Success Criteria

- `schangelog parse-commits` on a repo with trailered commits shows RMI IDs
  per commit.
- A generated changelog update for a real release carries `rmis` on its
  entries with no hand-editing.
- `schangelog validate` accepts the new field and flags malformed RMI IDs.
- VisionStudio ingest can associate a release from entry-level RMIs alone
  (verified once 006's RMI-VISIONSTUDIO-303 lands).
