# Resolve — Build Progress Tracker

## Phase 0: Research & Validation
**Target: April – May 2026**

### Spec & Design
- [x] Initial spec written (`unified_linux_package_manager_spec_roadmap.md`)
- [x] Competitive landscape documented (PackageKit, GNOME Software, Flatpak, Snap, Nix, apx, Homebrew)
- [x] Pre-Phase 0 data validation complete (442 HN + StackExchange posts analysed)
- [x] HN thread deep-read complete ("Ask HN: How would you set up a child's first Linux computer?")
- [x] Linux Foundation FAIR project investigated — confirmed not a competitor
- [x] Distro-aware resolution strategy decided: Ubuntu/derivatives → Snap-first; others → Flatpak-first
- [x] Snap transitional stub problem identified and addressed in spec

### Prototype
- [x] Git repo initialised
- [x] Go module initialised (`github.com/pawanjot/resolve`)
- [x] Dependencies added: `cobra` (CLI), `viper` (config)
- [x] Project structure scaffolded (`cmd/`, `internal/adapter/`, `internal/resolver/`, `internal/distro/`)

#### CLI
- [x] `resolve install <app> [--prefer flatpak|snap|native] [--dry-run]`
- [x] `resolve dry-run <app>` — subcommand alias for `install --dry-run`
- [x] `resolve search <app>` — shows results across all available sources with version info
- [x] Config file support (`~/.config/resolve/config.toml`, `preferred_source` key)
- [x] Shared `resolvePreference()` + `printDecision()` helper (`cmd/resolve_helper.go`) — eliminates duplication between `install` and `dry-run`

#### Adapters
- [x] **apt adapter** — search (`apt-cache`), version (`apt-cache policy`), install (`pkexec apt-get`)
- [x] **Snap adapter** — search (exact-match only), version (`snap info`), install (`snap install`), installed version detection (`snap list`)
- [x] **Flatpak adapter** — search, version, install (user-space, no elevation)
- [x] **Flatpak bootstrap** — auto-installs Flatpak + adds Flathub remote on non-Ubuntu distros when resolution picks Flatpak; one-time session restart warning

#### Resolution Engine
- [x] Distro detection via `/etc/os-release` (Ubuntu family, Fedora, Arch, Debian)
- [x] App type classification: GUI vs CLI (known CLI tool list)
- [x] **Ubuntu hierarchy**: Snap → apt → Flatpak (last resort)
- [x] **Non-Ubuntu hierarchy**: Flatpak → native (GUI) / native → Flatpak (CLI)
- [x] Explicit `--prefer` flag overrides algorithm per-command
- [x] Snap transitional stub detection — apt packages that secretly install a Snap are flagged and skipped during resolution

#### Quality
- [x] Snap search restricted to exact name matches (prevents e.g. `git` → `gitkraken`)
- [x] `.gitignore` in place

---

## Phase 0 Exit Criteria (outstanding)
- [ ] Prototype installs ≥ 20 of the top-50 Ubuntu apps successfully with visible source selection rationale
- [ ] Prototype posted publicly (GitHub + Show HN, r/linux, r/linuxquestions)
- [ ] ≥ 20 substantive reactions within 2 weeks of public release
- [ ] No commenter names an existing tool that already does cross-source selection with explanation

---

## Phase 1: MVP
**Target: June – August 2026** — not started

## Phase 2: Intelligent Resolution Engine
**Target: September – October 2026** — not started

## Phase 3: Expanded Distro Support
**Target: November 2026 – January 2027** — not started

## Phase 4: Developer Publishing Platform
**Target: February – April 2027** — not started

## Phase 5: Trust & Verification Layer
**Target: May – June 2027** — not started

## Phase 6: GUI & Ecosystem Layer
**Target: Q3–Q4 2027** — not started
