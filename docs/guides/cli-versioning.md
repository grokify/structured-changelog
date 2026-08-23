# CLI Versioning & Build Provenance

A convention for making Go CLI binaries self-identify — including when they were built **between releases** or from an **uncommitted** working tree. This is the reference implementation used by `schangelog`, `schemakit`, and `stasks`.

## The Two Signals

A binary's provenance has two independent parts, from two sources:

| Signal | Source | Available when |
|--------|--------|----------------|
| **Release tag** (`v0.16.0`) | ldflags (`-X main.version=…`) or `debug.ReadBuildInfo().Main.Version` | ldflags: release builds. Build info: `go install …@vX.Y.Z`, or (Go 1.24+) derived from the git tag of a local checkout. |
| **Commit + dirty flag** | build info `vcs.revision` / `vcs.time` / `vcs.modified` | Always — the Go toolchain stamps these automatically from git. |
| **Commits past last tag** (`v0.16.0-3-gabc1234`) | `git describe --tags` via ldflags | Only when stamped at build time. Build info alone does not embed the base tag. |

**Key point:** "is this uncommitted?" is answerable at runtime for free (build info `vcs.modified`). "How far past the last release is it?" needs a `git describe` stamp at build time.

## Inspecting Any Binary (no code required)

The Go toolchain embeds build info in every binary. Read it with:

```bash
go version -m $(which schangelog)
```

Look for `mod` (module version — `(devel)` or a pseudo-version means a local/between-releases build), and the `build vcs.revision` / `vcs.time` / `vcs.modified` lines (`vcs.modified=true` means the tree had uncommitted changes).

## The `version` Command

Every CLI exposes both a `version` subcommand and a `--version` flag, reporting the same multi-line output:

```
schangelog v0.16.0-3-gabc1234-dirty  [development build]
  commit: abc1234def56 (modified)
  built:  2026-08-23T17:02:32Z
  go:     go1.26.5 darwin/arm64
```

- **Release build:** `schangelog v0.16.0` — clean, no dev marker.
- **Local `go install`/`go build`:** module version or `(devel)` + commit + `[development build]` if the tree was modified.
- **`make build`:** the `git describe` string `vX.Y.Z-N-gSHA-dirty` — the precise "N commits past last tag, uncommitted" signal.

The version token already carries its own dirty marker (Go's `+dirty` pseudo-version suffix, or git-describe's `-dirty`); the report does not double it.

## Reference Implementation

The pattern: keep ldflags variables with sentinel defaults, and fall back to `debug.ReadBuildInfo()` when they aren't stamped.

```go
var (
	version = "dev"     // -X main.version
	commit  = "none"    // -X main.commit
	date    = "unknown" // -X main.date
)

func resolveVersion() (ver, rev, when string, modified, devel bool) {
	ver, rev, when = version, commit, date
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ver, rev, when, false, ver == "dev"
	}
	if ver == "dev" { // not stamped by ldflags
		if bi.Main.Version != "" && bi.Main.Version != "(devel)" {
			ver = bi.Main.Version
		} else {
			ver, devel = "(devel)", true
		}
	}
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			if rev == "none" { rev = s.Value }
		case "vcs.time":
			if when == "unknown" { when = s.Value }
		case "vcs.modified":
			modified = s.Value == "true"
		}
	}
	return ver, rev, when, modified, devel
}
```

Wire both entry points in `init()`:

```go
rootCmd.Version = versionString()
rootCmd.SetVersionTemplate("{{.Version}}\n")
rootCmd.AddCommand(versionCmd) // versionCmd.Run prints versionString()
```

## Build-Time Stamping

**Release builds (goreleaser)** already stamp the ldflags:

```yaml
ldflags:
  - -X main.version={{.Version}}
  - -X main.commit={{.Commit}}
  - -X main.date={{.Date}}
```

**Local `make build`** should stamp `git describe` for the readable "commits-past-tag" string:

```makefile
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null)
COMMIT  := $(shell git rev-parse HEAD 2>/dev/null)
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/mytool ./cmd/mytool
```

A plain `go install`/`go build` without these still self-identifies via build info — the stamp only upgrades the version string to the friendlier git-describe form.

## Why This Lives Here

Provenance is a versioning concern, and structured-changelog is the versioning tool in the toolkit — so the convention is documented here and adopted by the sibling CLIs.
