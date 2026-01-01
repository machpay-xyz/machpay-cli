# MachPay CLI - Launch Roadmap

**Document Purpose:** Multi-phase execution plan for all pending items  
**Created:** 2025-01-01  
**Target:** Production-ready v0.1.0 release

---

## 📊 Overview

This roadmap consolidates all pending work from Phase 5 (Distribution) and Phase 6 (Polish) into a logical, dependency-aware execution plan.

### Current Status
- **Phases 1-4:** ✅ Complete (Core CLI, Gateway, Desktop App)
- **Phase 5:** ⚠️ 40% complete (Distribution infrastructure ready)
- **Phase 6:** ⚠️ 33% complete (Tests exist, needs polish)

### Total Pending Items: 13
- Phase 5 items: 6
- Phase 6 items: 6
- Critical Bug Fix: 1 (Google OAuth in Desktop App)

---

## 🚀 Execution Phases

### Sprint 1: Foundation & CI/CD (Est: 2-3 hours)
**Goal:** Enable automated releases

| # | Task | Type | Priority | Blocked By |
|---|------|------|----------|------------|
| 1.1 | Create GitHub Actions Release Workflow | P5 | HIGH | - |
| 1.2 | Create PR Build/Test Workflow | P5 | HIGH | - |
| 1.3 | Verify GoReleaser Configuration | P5 | HIGH | - |
| 1.4 | Test Coverage Verification (>80%) | P6 | HIGH | - |

**Deliverables:**
- `.github/workflows/release.yml`
- `.github/workflows/ci.yml`
- Verified test coverage report
- All CI checks passing

---

### Sprint 2: First Release (Est: 1-2 hours)
**Goal:** Publish v0.1.0 to GitHub

| # | Task | Type | Priority | Blocked By |
|---|------|------|----------|------------|
| 2.1 | Update Version Numbers | P5 | HIGH | Sprint 1 |
| 2.2 | Create CHANGELOG.md Release Notes | P5 | HIGH | Sprint 1 |
| 2.3 | Create and Push v0.1.0 Tag | P5 | HIGH | 2.1, 2.2 |
| 2.4 | Verify Release Artifacts | P5 | HIGH | 2.3 |

**Deliverables:**
- GitHub Release v0.1.0
- CLI binaries (macOS, Linux, Windows)
- Checksums file
- Release notes

---

### Sprint 3: Distribution Channels (Est: 2-3 hours)
**Goal:** Enable all install methods

| # | Task | Type | Priority | Blocked By |
|---|------|------|----------|------------|
| 3.1 | Update Homebrew Formula with SHA256 | P5 | MEDIUM | Sprint 2 |
| 3.2 | Test Homebrew Installation | P5 | MEDIUM | 3.1 |
| 3.3 | Optimize Dockerfile | P5 | MEDIUM | - |
| 3.4 | Add Docker to GoReleaser | P5 | MEDIUM | 3.3 |
| 3.5 | Test Docker Image | P5 | MEDIUM | 3.4 |

**Deliverables:**
- Working `brew install machpay/tap/machpay`
- Docker image at `ghcr.io/machpay/cli:latest`
- Verified install scripts

---

### Sprint 4: Desktop App Fixes (Est: 3-4 hours)
**Goal:** Fix Google OAuth and package desktop app

| # | Task | Type | Priority | Blocked By |
|---|------|------|----------|------------|
| 4.1 | **Fix Google OAuth in Tauri App** | BUG | HIGH | - |
| 4.2 | macOS DMG Code Signing | P5 | MEDIUM | 4.1 |
| 4.3 | macOS Notarization | P5 | MEDIUM | 4.2 |
| 4.4 | Windows MSI Installer | P5 | LOW | 4.1 |
| 4.5 | Linux AppImage/DEB | P5 | LOW | 4.1 |

**Deliverables:**
- Working Google OAuth in desktop app
- Signed macOS DMG
- Notarized for Gatekeeper
- Cross-platform installers

---

### Sprint 5: Documentation & Polish (Est: 2-3 hours)
**Goal:** Complete documentation and launch prep

| # | Task | Type | Priority | Blocked By |
|---|------|------|----------|------------|
| 5.1 | Create INSTALL.md Guide | P5 | MEDIUM | Sprint 2 |
| 5.2 | Update README.md Install Section | P5 | MEDIUM | 5.1 |
| 5.3 | Run Performance Benchmarks | P6 | MEDIUM | - |
| 5.4 | Complete Launch Checklist | P6 | HIGH | Sprint 1-4 |
| 5.5 | Generate Man Pages | P6 | LOW | - |

**Deliverables:**
- Comprehensive INSTALL.md
- Updated README with badges
- Performance baseline documented
- Launch checklist verified

---

### Sprint 6: Nice-to-Have (Est: 4-6 hours) [Optional]
**Goal:** Enhanced observability (can defer to v0.2.0)

| # | Task | Type | Priority | Blocked By |
|---|------|------|----------|------------|
| 6.1 | Implement Telemetry System | P6 | LOW | - |
| 6.2 | Implement Error Reporting (Sentry) | P6 | LOW | - |
| 6.3 | Create Demo GIFs/Videos | P5 | LOW | Sprint 2 |

**Deliverables:**
- Optional telemetry with opt-out
- Sentry integration for crash reporting
- Installation demo assets

---

## 🐛 Critical Bug: Google OAuth in Desktop App

### Problem Statement
Google OAuth login in the Tauri desktop app gets stuck in "Connecting..." state after browser completes authentication.

### Root Cause
Cross-origin security restrictions prevent communication between:
- Browser callback (`https://console.machpay.xyz`)
- Desktop app (`tauri://localhost`)
- Local backend (`http://localhost:8081`)

### Attempted Solutions (Failed)
1. ❌ Clipboard relay - Browser blocks automatic writes
2. ❌ LocalStorage relay - Different origins don't share storage
3. ❌ Local HTTP server - Added complexity, Google redirect issues
4. ❌ Backend polling - HTTPS→HTTP Mixed Content blocked

### Proposed Solution: Custom URL Scheme

**Implementation Plan:**

```
Step 1: Register machpay:// URL scheme in Tauri
Step 2: Update Google OAuth redirect to use web callback
Step 3: Web callback exchanges code for token
Step 4: Web callback redirects to machpay://auth?token=xxx
Step 5: Desktop app captures URL and logs in user
```

**Technical Details:**

1. **Tauri Configuration** (`tauri.conf.json`):
```json
{
  "tauri": {
    "allowlist": {
      "protocol": {
        "asset": true,
        "assetScope": ["**"]
      }
    }
  }
}
```

2. **Register URL Handler** (`main.rs`):
```rust
// Register machpay:// protocol handler
tauri_plugin_deep_link::register("machpay", |request| {
    // Parse token from URL and emit to frontend
});
```

3. **Web Callback Flow**:
```
Google → https://console.machpay.xyz/auth/callback?code=xxx
         ↓
         Exchange code for token (server-side)
         ↓
         Redirect to: machpay://auth?token=xxx
         ↓
Desktop App captures and logs in
```

4. **Fallback**: If deep link fails, show token on screen for manual copy

---

## 📋 Detailed Task Prompts

### Sprint 1 Prompts

#### Prompt 1.1: GitHub Actions Release Workflow
```
Create a GitHub Actions workflow at machpay-cli/.github/workflows/release.yml that:

1. Triggers on tag push matching pattern v*.*.*
2. Sets up Go 1.21+ environment
3. Runs tests before release
4. Uses GoReleaser action to:
   - Build binaries for macOS (arm64, amd64), Linux (amd64, arm64), Windows (amd64)
   - Create GitHub Release with artifacts
   - Generate checksums
5. Uploads release artifacts

Include:
- Caching for Go modules
- Matrix strategy for test platforms
- Proper permissions for GITHUB_TOKEN

Reference the existing .goreleaser.yaml configuration.
```

#### Prompt 1.2: PR Build/Test Workflow
```
Create a GitHub Actions workflow at machpay-cli/.github/workflows/ci.yml that:

1. Triggers on pull requests to main branch
2. Runs on ubuntu-latest, macos-latest, windows-latest
3. Steps:
   - Checkout code
   - Setup Go 1.21+
   - Cache Go modules
   - Run go vet
   - Run golangci-lint
   - Run tests with coverage
   - Upload coverage report
4. Fail if coverage < 80%

Use actions/setup-go@v4 and golangci/golangci-lint-action@v3
```

#### Prompt 1.3: Verify GoReleaser Config
```
Review and update machpay-cli/.goreleaser.yaml to ensure:

1. Correct binary name (machpay)
2. All target platforms configured:
   - darwin/amd64, darwin/arm64
   - linux/amd64, linux/arm64
   - windows/amd64
3. Archive format: tar.gz for Unix, zip for Windows
4. Checksum file generation
5. Changelog generation from git commits
6. Release notes template

Test locally with: goreleaser release --snapshot --clean --skip-publish
```

#### Prompt 1.4: Test Coverage Verification
```
Check test coverage for machpay-cli:

1. Run: go test -coverprofile=coverage.out ./...
2. Generate report: go tool cover -html=coverage.out -o coverage.html
3. Check percentage: go tool cover -func=coverage.out | tail -1

If coverage < 80%:
- Identify packages with low coverage
- Add tests for critical paths:
  - internal/auth/
  - internal/cmd/
  - internal/gateway/
  - internal/config/

Target: Minimum 80% overall coverage
```

---

### Sprint 2 Prompts

#### Prompt 2.1: Update Version Numbers
```
Update version to 0.1.0 in all relevant files:

1. machpay-cli/cmd/machpay/main.go - version constant
2. machpay-console/package.json - version field
3. machpay-console/src-tauri/tauri.conf.json - package.version
4. machpay-console/src-tauri/Cargo.toml - package.version

Create a script at scripts/bump-version.sh that:
- Takes version as argument
- Updates all files automatically
- Creates a git commit with message "chore: bump version to X.X.X"
```

#### Prompt 2.2: Create Release Notes
```
Create comprehensive release notes in CHANGELOG.md for v0.1.0:

# Changelog

## [0.1.0] - 2025-01-XX

### 🎉 Initial Public Release

First public release of MachPay CLI - the command-line interface for the MachPay AI payment network.

### ✨ Features

#### Authentication
- Browser-based OAuth login
- Headless mode for CI/CD environments
- Secure token storage

#### Setup & Configuration
- Interactive setup wizard
- Role selection (Agent/Vendor)
- Network selection (Devnet/Mainnet)
- Wallet generation (Ed25519/Solana)

#### Gateway Management
- Automatic gateway download
- Process management (start/stop/restart)
- Health monitoring
- Log viewing with follow mode

#### Desktop Integration
- Launch desktop app from CLI
- Deep linking to specific routes
- Browser fallback when app not installed

### 📦 Installation

[Include installation methods]

### 🐛 Known Issues

- Google OAuth in desktop app requires web fallback (see KNOWN_ISSUES.md)

### 📝 Documentation

- Full README with usage examples
- Command reference with --help
- Configuration guide

---

Make it professional, informative, and include all features from Phases 1-4.
```

#### Prompt 2.3: Create and Push Tag
```
Guide me through creating the v0.1.0 release:

Pre-flight checks:
1. Ensure all tests pass: go test ./...
2. Ensure main branch is up to date: git pull origin main
3. Ensure no uncommitted changes: git status
4. Verify version numbers updated (Prompt 2.1)
5. Verify CHANGELOG.md updated (Prompt 2.2)

Create release:
1. Create annotated tag:
   git tag -a v0.1.0 -m "Release v0.1.0 - Initial public release"

2. Push tag:
   git push origin v0.1.0

3. Monitor GitHub Actions workflow

4. Verify release at:
   https://github.com/machpay-xyz/machpay-cli/releases/tag/v0.1.0
```

---

### Sprint 4 Prompts (Google OAuth Fix)

#### Prompt 4.1: Fix Google OAuth in Tauri App
```
Implement the custom URL scheme solution for Google OAuth in the desktop app:

Part 1: Tauri Configuration
1. Add tauri-plugin-deep-link to Cargo.toml
2. Register machpay:// URL scheme in tauri.conf.json
3. Update main.rs to handle deep links and emit events

Part 2: Web Callback Updates
1. Update AuthCallback.jsx to detect desktop source
2. After code exchange, redirect to machpay://auth?token=xxx
3. Add fallback UI showing "Return to app" if redirect fails

Part 3: Desktop App Handler
1. Listen for deep-link events in App.jsx
2. Parse token from machpay://auth?token=xxx URL
3. Call login() with token to authenticate user
4. Navigate to appropriate dashboard

Part 4: Google Cloud Console
1. Keep existing redirect URIs (no changes needed)
2. Web callback handles code exchange
3. Desktop app receives token via URL scheme

Test the full flow:
1. Click Google login in desktop app
2. Browser opens, complete Google auth
3. Browser redirects to console.machpay.xyz/auth/callback
4. Callback exchanges code, gets token
5. Callback redirects to machpay://auth?token=xxx
6. Desktop app captures, logs in user

Provide complete code changes for all files.
```

---

## ✅ Completion Checklist

### Sprint 1: Foundation & CI/CD
- [ ] `.github/workflows/release.yml` created and working
- [ ] `.github/workflows/ci.yml` created and working
- [ ] GoReleaser config verified
- [ ] Test coverage ≥80%
- [ ] All CI checks passing on main branch

### Sprint 2: First Release
- [ ] Version numbers updated to 0.1.0
- [ ] CHANGELOG.md has release notes
- [ ] Git tag v0.1.0 created and pushed
- [ ] GitHub Release published
- [ ] All binary artifacts available
- [ ] Checksums verified

### Sprint 3: Distribution
- [ ] Homebrew formula updated with SHA256
- [ ] `brew install machpay/tap/machpay` works
- [ ] Docker image builds successfully
- [ ] Docker image pushed to GHCR
- [ ] Install scripts tested

### Sprint 4: Desktop App
- [ ] **Google OAuth fixed and working**
- [ ] macOS DMG signed
- [ ] macOS app notarized
- [ ] Windows MSI created (optional)
- [ ] Linux packages created (optional)

### Sprint 5: Documentation
- [ ] INSTALL.md created
- [ ] README.md updated
- [ ] Performance benchmarks documented
- [ ] Launch checklist complete
- [ ] Man pages generated (optional)

### Sprint 6: Nice-to-Have (Optional)
- [ ] Telemetry system
- [ ] Error reporting
- [ ] Demo videos/GIFs

---

## 📅 Estimated Timeline

| Sprint | Duration | Cumulative |
|--------|----------|------------|
| Sprint 1 | 2-3 hours | 2-3 hours |
| Sprint 2 | 1-2 hours | 3-5 hours |
| Sprint 3 | 2-3 hours | 5-8 hours |
| Sprint 4 | 3-4 hours | 8-12 hours |
| Sprint 5 | 2-3 hours | 10-15 hours |
| Sprint 6 | 4-6 hours | 14-21 hours |

**Minimum Viable Launch:** Sprints 1-2 (3-5 hours)  
**Full Launch:** Sprints 1-5 (10-15 hours)  
**Complete:** All Sprints (14-21 hours)

---

## 🎯 Quick Start

To begin execution, start with:

```
"Execute Sprint 1, Prompt 1.1: Create GitHub Actions Release Workflow"
```

Progress through each sprint sequentially. Mark tasks complete in this document as you go.

---

**Document Version:** 1.0  
**Last Updated:** 2025-01-01

