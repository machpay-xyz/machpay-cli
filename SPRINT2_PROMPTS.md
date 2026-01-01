# Sprint 2: First Release (v0.1.0) - Execution Prompts

**Sprint Goal:** Publish the first official release to GitHub  
**Estimated Time:** 1-2 hours  
**Prerequisites:** Sprint 1 Complete (CI/CD workflows ready)

---

## Overview

| Task | Description | Priority | Est. Time |
|------|-------------|----------|-----------|
| 2.1 | Update Version Numbers | HIGH | 15 min |
| 2.2 | Create CHANGELOG Release Notes | HIGH | 30 min |
| 2.3 | Pre-Release Verification | HIGH | 15 min |
| 2.4 | Create and Push Git Tag | HIGH | 15 min |
| 2.5 | Verify Release Artifacts | HIGH | 15 min |

---

## Prompt 2.1: Update Version Numbers

### Objective
Update version from "dev" to "0.1.0" across all relevant files.

### Prompt
```
Update version numbers to 0.1.0 in the machpay-cli repository.

Files to update:

1. cmd/machpay/main.go
   - Change: version = "dev" → version = "0.1.0"
   - Leave commit and date as placeholders (GoReleaser will inject)

2. README.md (if version is mentioned)
   - Update any version badges or references

3. Verify .goreleaser.yaml doesn't have hardcoded version

Additionally, create a version bump script for future releases:

scripts/bump-version.sh
- Takes new version as argument
- Updates all version references
- Creates a git commit

After changes:
- Run: go build ./cmd/machpay && ./machpay version
- Verify output shows "0.1.0"
```

### Expected Changes

**cmd/machpay/main.go:**
```go
var (
    version = "0.1.0"  // Changed from "dev"
    commit  = "none"   // GoReleaser injects at build time
    date    = "unknown" // GoReleaser injects at build time
)
```

**scripts/bump-version.sh:**
```bash
#!/bin/bash
# Bump version across all files
# Usage: ./scripts/bump-version.sh 0.2.0

set -e

NEW_VERSION=$1

if [ -z "$NEW_VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 0.2.0"
    exit 1
fi

echo "Bumping version to $NEW_VERSION..."

# Update main.go
sed -i '' "s/version = \".*\"/version = \"$NEW_VERSION\"/" cmd/machpay/main.go

# Verify
echo "Updated files:"
grep -n "version = " cmd/machpay/main.go

echo ""
echo "✅ Version bumped to $NEW_VERSION"
echo "Run: git add . && git commit -m 'chore: bump version to $NEW_VERSION'"
```

### Verification
- [ ] `go build ./cmd/machpay && ./machpay version` shows "0.1.0"
- [ ] No hardcoded versions in .goreleaser.yaml
- [ ] bump-version.sh script created and executable

---

## Prompt 2.2: Create CHANGELOG Release Notes

### Objective
Create comprehensive release notes documenting all features in v0.1.0.

### Prompt
```
Create or update CHANGELOG.md with comprehensive release notes for v0.1.0.

The release notes should document all features from Phases 1-4:

Phase 1 (Foundation):
- CLI repository structure
- Browser-based authentication
- Configuration management
- Status and version commands
- Cross-platform support

Phase 2 (CLI Core):
- Interactive setup wizard
- TUI with Lipgloss styling
- Wallet generation (Ed25519/Solana)
- Role and network selection
- JSON output mode
- Watch mode

Phase 3 (Gateway Integration):
- Gateway auto-download from GitHub
- Process management (start/stop/restart)
- Health checks
- Log viewing with follow mode
- Update command

Phase 4 (Desktop Integration):
- Desktop app launching via CLI
- Deep linking to specific routes
- Browser fallback

Known Issues:
- Google OAuth in desktop app (documented workaround)

Format the changelog using Keep a Changelog format:
https://keepachangelog.com/en/1.1.0/

Include:
- Version header with date
- Categorized changes (Added, Changed, Fixed, etc.)
- Links to documentation
- Installation instructions summary
```

### Expected CHANGELOG.md
```markdown
# Changelog

All notable changes to MachPay CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-01-XX

### 🎉 Initial Public Release

First public release of MachPay CLI - the command-line interface for the MachPay AI payment network.

### Added

#### Authentication & Configuration
- Browser-based OAuth authentication with `machpay login`
- Secure token storage in `~/.machpay/config.yaml`
- Headless mode for CI/CD environments (`--no-browser` flag)
- Interactive setup wizard with `machpay setup`
- Role selection (Agent or Vendor)
- Network selection (Devnet or Mainnet)

#### Wallet Management
- Ed25519/Solana keypair generation
- Import existing wallet from keypair file
- Secure wallet storage with permissions
- Public key display in Base58 format

#### Gateway Management (Vendors)
- Automatic gateway binary download from GitHub releases
- Process management: `machpay serve`, `machpay stop`, `machpay restart`
- Background/foreground operation modes
- Health check monitoring
- Log viewing with follow mode: `machpay logs -f`
- Gateway update command: `machpay update gateway`

#### Status & Monitoring
- Comprehensive status display: `machpay status`
- JSON output for scripting: `machpay status --json`
- Watch mode for continuous monitoring: `machpay status --watch`
- Gateway status (running/stopped, PID, port)

#### Desktop Integration
- Launch desktop app from CLI: `machpay open`
- Deep linking to specific routes: `machpay open marketplace`
- Automatic browser fallback when desktop app not installed
- Route support: marketplace, funding, settings, invoices, etc.

#### Developer Experience
- Beautiful TUI with Lipgloss styling
- Colored output with icons
- Progress indicators for downloads
- Helpful error messages with suggestions
- Cross-platform support (macOS, Linux, Windows)

### Known Issues

- **Google OAuth in Desktop App**: Authentication via Google in the desktop app may not complete automatically. Workaround: Use email OTP or wallet login in the desktop app, or use the web console for Google OAuth. See [KNOWN_ISSUES.md](./KNOWN_ISSUES.md) for details.

### Installation

**Homebrew (macOS/Linux):**
```bash
brew install machpay/tap/machpay
```

**Quick Install (Linux/macOS):**
```bash
curl -fsSL https://machpay.xyz/install.sh | sh
```

**Windows PowerShell:**
```powershell
iwr machpay.xyz/install.ps1 | iex
```

**From Source:**
```bash
go install github.com/machpay-xyz/machpay-cli/cmd/machpay@latest
```

### Quick Start

```bash
# Authenticate
machpay login

# Configure your node
machpay setup

# Check status
machpay status

# Start vendor gateway (if vendor role)
machpay serve

# Open MachPay console
machpay open
```

### Documentation

- [README](./README.md) - Full documentation
- [CLI Status](./CLI_STATUS.md) - Feature status tracking
- [Known Issues](./KNOWN_ISSUES.md) - Known issues and workarounds

---

[Unreleased]: https://github.com/machpay-xyz/machpay-cli/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/machpay-xyz/machpay-cli/releases/tag/v0.1.0
```

### Verification
- [ ] CHANGELOG.md created with all features documented
- [ ] Follows Keep a Changelog format
- [ ] Includes installation instructions
- [ ] Includes known issues
- [ ] Date placeholder ready to be updated

---

## Prompt 2.3: Pre-Release Verification

### Objective
Run comprehensive checks before creating the release tag.

### Prompt
```
Run pre-release verification checks for machpay-cli v0.1.0.

Execute the following checks in order:

1. Clean Build Test
   cd /path/to/machpay-cli
   rm -rf dist/
   go clean -cache
   go build -v ./cmd/machpay
   ./machpay version

2. Full Test Suite
   go test -v ./...

3. Linter Check
   golangci-lint run

4. GoReleaser Dry Run
   goreleaser release --snapshot --clean --skip=publish

5. Binary Verification
   # Test each platform binary from dist/
   ./dist/machpay_darwin_arm64_v8.0/machpay version
   ./dist/machpay_darwin_arm64_v8.0/machpay --help
   ./dist/machpay_darwin_arm64_v8.0/machpay status

6. Git Status Check
   git status
   # Ensure all changes are committed
   # Ensure on main branch
   # Ensure up to date with remote

7. Documentation Check
   # Verify README.md is current
   # Verify CHANGELOG.md has v0.1.0 section
   # Verify CLI_STATUS.md is updated

Report any failures and fix before proceeding.
```

### Pre-Release Checklist
```
Pre-Release Verification Checklist v0.1.0
==========================================

Build & Test:
[ ] go build ./cmd/machpay - SUCCESS
[ ] ./machpay version shows "0.1.0"
[ ] go test ./... - ALL PASS
[ ] golangci-lint run - NO ERRORS

GoReleaser:
[ ] goreleaser check - VALID
[ ] goreleaser release --snapshot - SUCCESS
[ ] All 5 platform binaries generated

Binary Test:
[ ] machpay version - works
[ ] machpay --help - works
[ ] machpay status - works (shows not logged in)

Git Status:
[ ] On main branch
[ ] All changes committed
[ ] Up to date with origin/main
[ ] No uncommitted changes

Documentation:
[ ] README.md is current
[ ] CHANGELOG.md has v0.1.0 entry
[ ] Date updated in CHANGELOG.md

Ready for Release:
[ ] All checks passed
[ ] Team approval (if applicable)
```

### Verification
- [ ] All pre-release checks pass
- [ ] Checklist completed
- [ ] Ready to create tag

---

## Prompt 2.4: Create and Push Git Tag

### Objective
Create the v0.1.0 annotated tag and push to trigger the release workflow.

### Prompt
```
Create and push the v0.1.0 release tag for machpay-cli.

IMPORTANT: This will trigger the GitHub Actions release workflow
and create a public release. Ensure all pre-release checks pass first.

Steps:

1. Final Git Status Check
   git status
   git log --oneline -5
   # Verify you're on main with all changes committed

2. Create Annotated Tag
   git tag -a v0.1.0 -m "Release v0.1.0 - Initial public release

   First public release of MachPay CLI.

   Features:
   - Browser-based authentication
   - Interactive setup wizard
   - Wallet generation (Ed25519/Solana)
   - Gateway management (start/stop/restart)
   - Desktop app integration
   - Cross-platform support (macOS, Linux, Windows)

   See CHANGELOG.md for full details."

3. Verify Tag
   git tag -l -n9 v0.1.0
   git show v0.1.0 --quiet

4. Push Tag (THIS TRIGGERS THE RELEASE)
   git push origin v0.1.0

5. Monitor Release
   # Open GitHub Actions to watch the workflow
   open https://github.com/machpay-xyz/machpay-cli/actions

   # Or watch from CLI
   gh run watch

6. Verify Release Page
   open https://github.com/machpay-xyz/machpay-cli/releases/tag/v0.1.0
```

### Tag Message Template
```
Release v0.1.0 - Initial public release

First public release of MachPay CLI - the command-line interface 
for the MachPay AI payment network.

Highlights:
• Browser-based OAuth authentication
• Interactive setup wizard with TUI
• Ed25519/Solana wallet generation
• Vendor gateway management
• Desktop app integration with deep linking
• Cross-platform support (macOS, Linux, Windows)

Installation:
  brew install machpay/tap/machpay
  # or
  curl -fsSL https://machpay.xyz/install.sh | sh

Quick Start:
  machpay login
  machpay setup
  machpay serve

See CHANGELOG.md for full release notes.
```

### Verification
- [ ] Tag created locally
- [ ] Tag pushed to GitHub
- [ ] GitHub Actions workflow triggered
- [ ] Workflow completes successfully

---

## Prompt 2.5: Verify Release Artifacts

### Objective
Verify that all release artifacts are correctly generated and downloadable.

### Prompt
```
Verify the v0.1.0 release artifacts on GitHub.

1. Check Release Page
   Navigate to: https://github.com/machpay-xyz/machpay-cli/releases/tag/v0.1.0

   Verify the following assets exist:
   [ ] machpay_darwin_amd64.tar.gz (macOS Intel)
   [ ] machpay_darwin_arm64.tar.gz (macOS Apple Silicon)
   [ ] machpay_linux_amd64.tar.gz (Linux x64)
   [ ] machpay_linux_arm64.tar.gz (Linux ARM)
   [ ] machpay_windows_amd64.zip (Windows x64)
   [ ] checksums.txt
   [ ] Source code (zip)
   [ ] Source code (tar.gz)

2. Download and Verify Checksums
   # Download checksums
   curl -LO https://github.com/machpay-xyz/machpay-cli/releases/download/v0.1.0/checksums.txt
   
   # Download your platform binary
   curl -LO https://github.com/machpay-xyz/machpay-cli/releases/download/v0.1.0/machpay_darwin_arm64.tar.gz
   
   # Verify checksum
   shasum -a 256 machpay_darwin_arm64.tar.gz
   grep machpay_darwin_arm64.tar.gz checksums.txt
   # Should match!

3. Test Downloaded Binary
   tar -xzf machpay_darwin_arm64.tar.gz
   ./machpay version
   # Should show: machpay version 0.1.0 (commit: XXXXXXX, built: 2025-XX-XX)

4. Verify Release Notes
   - Check that release notes match CHANGELOG.md
   - Verify installation instructions are present
   - Verify links work

5. Test Install Script (if hosted)
   # Only if install script URL is live
   curl -fsSL https://machpay.xyz/install.sh | sh -s -- --dry-run

6. Update CLI_STATUS.md
   - Mark Phase 5 GitHub Release as complete
   - Update "Next Steps" section
```

### Release Verification Checklist
```
Release Artifacts Verification v0.1.0
=====================================

Assets Present:
[ ] machpay_darwin_amd64.tar.gz
[ ] machpay_darwin_arm64.tar.gz
[ ] machpay_linux_amd64.tar.gz
[ ] machpay_linux_arm64.tar.gz
[ ] machpay_windows_amd64.zip
[ ] checksums.txt

Checksum Verification:
[ ] Downloaded binary checksum matches checksums.txt

Binary Test:
[ ] Extracted successfully
[ ] ./machpay version shows correct version
[ ] Version includes commit hash and build date

Release Notes:
[ ] Title: "MachPay CLI v0.1.0"
[ ] Installation instructions present
[ ] Features documented
[ ] Known issues mentioned

Post-Release:
[ ] CLI_STATUS.md updated
[ ] Team notified (if applicable)
[ ] Social media announcement (if applicable)
```

### Verification
- [ ] All artifacts present and downloadable
- [ ] Checksums verify correctly
- [ ] Binary works after download
- [ ] Release notes are complete
- [ ] Documentation updated

---

## Sprint 2 Completion Checklist

After completing all prompts, verify:

```
Sprint 2 Final Checklist:

Version Update:
[ ] main.go version = "0.1.0"
[ ] bump-version.sh script created

Documentation:
[ ] CHANGELOG.md complete with v0.1.0 section
[ ] Date filled in CHANGELOG.md

Pre-Release:
[ ] All tests pass
[ ] Linter passes
[ ] GoReleaser snapshot works
[ ] All changes committed

Release:
[ ] v0.1.0 tag created
[ ] Tag pushed to GitHub
[ ] GitHub Actions workflow succeeded
[ ] Release page shows all artifacts

Verification:
[ ] Downloaded binary works
[ ] Checksum matches
[ ] Release notes correct

Post-Release:
[ ] CLI_STATUS.md updated
[ ] Announced to team/community
```

### Commands Summary

```bash
# Navigate to repo
cd /Users/abhishektomar/Desktop/git/machpay-cli

# Update version (Prompt 2.1)
# Edit cmd/machpay/main.go: version = "0.1.0"
go build ./cmd/machpay && ./machpay version

# Create CHANGELOG (Prompt 2.2)
# Edit CHANGELOG.md with release notes

# Pre-release checks (Prompt 2.3)
go test -v ./...
golangci-lint run
goreleaser release --snapshot --clean --skip=publish

# Create and push tag (Prompt 2.4)
git add .
git commit -m "chore: prepare v0.1.0 release"
git tag -a v0.1.0 -m "Release v0.1.0 - Initial public release"
git push origin main
git push origin v0.1.0

# Verify release (Prompt 2.5)
open https://github.com/machpay-xyz/machpay-cli/releases/tag/v0.1.0
```

---

## Rollback Plan

If the release has issues:

```bash
# Delete the tag locally
git tag -d v0.1.0

# Delete the tag from GitHub
git push origin :refs/tags/v0.1.0

# Delete the GitHub release (via web UI or gh CLI)
gh release delete v0.1.0 --yes

# Fix issues, then re-tag
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

---

## Next Sprint

After Sprint 2 is complete, proceed to **Sprint 3: Distribution Channels**

Sprint 3 will:
- Update Homebrew formula with real checksums
- Test `brew install machpay/tap/machpay`
- Set up Docker image publishing

---

**Document Version:** 1.0  
**Created:** 2025-01-01

