# Release Process

Releases are fully automated via [release-please](https://github.com/googleapis/release-please) and [GoReleaser](https://goreleaser.com). Merging to `main` is all that is needed — no manual tagging or changelog editing.

## How it works

```
commit (conventional) → merge to main
        │
        ▼
release-please opens/updates a release PR
        │
        ▼  (PR merged)
release-please creates a semver tag  e.g. v1.2.0
        │
        ▼
GitHub Actions: run tests
        │
        ▼  (tests pass)
GoReleaser builds binaries for all platforms
and publishes the GitHub release with changelog
```

## Commit message format

Commits must follow the [Conventional Commits](https://www.conventionalcommits.org) specification. The commit type determines which version component is bumped:

| Commit prefix | Version bump | Example |
|---|---|---|
| `feat:` | minor — `0.x.0` | `feat: add clock format option` |
| `fix:` | patch — `0.0.x` | `fix: show indicator without sudo` |
| `feat!:` or `fix!:` | major — `x.0.0` | `feat!: redesign config format` |
| `docs:`, `chore:`, `ci:` | none (excluded from changelog) | — |

The body or footer can also trigger a major bump:

```
feat: redesign config format

BREAKING CHANGE: config keys have been renamed.
```

## Workflows

### `release-please.yml`

Triggers on every push to `main`. Runs three sequential jobs:

1. **release-please** — Maintains a long-running PR that updates `CHANGELOG.md` and bumps the version. When the PR is merged it creates the semver tag and sets `release_created=true`.
2. **test** — Runs `go test ./...`. Only runs when a release was just created.
3. **goreleaser** — Cross-compiles binaries for all platforms and publishes the GitHub release. Only runs after tests pass.

Release artefacts per version:

| File | Description |
|---|---|
| `tztui_linux_amd64.tar.gz` | Linux x86-64 binary + README |
| `tztui_linux_arm64.tar.gz` | Linux ARM64 binary + README |
| `tztui_darwin_amd64.tar.gz` | macOS x86-64 binary + README |
| `tztui_darwin_arm64.tar.gz` | macOS Apple Silicon binary + README |
| `tztui_checksums.txt` | SHA-256 checksums for all archives |

## Required secrets

| Secret | Required | Purpose |
|---|---|---|
| `GITHUB_TOKEN` | Automatic | Create tags, releases, and PRs |
| `HOMEBREW_TAP_GITHUB_TOKEN` | Optional | Push formula updates to `homebrew-tztui` tap |

`GITHUB_TOKEN` is provided automatically by GitHub Actions. `HOMEBREW_TAP_GITHUB_TOKEN` must be added manually under **Settings → Secrets and variables → Actions** if Homebrew tap publishing is desired.

## First release

The manifest is initialised at `0.0.0`. The first merged `feat:` commit will produce a `v0.1.0` release PR.
