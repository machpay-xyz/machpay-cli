# Sprint 3: Distribution Channels - Execution Prompts

**Sprint Goal:** Enable all install methods (Homebrew, Docker, direct download)  
**Estimated Time:** 2-3 hours  
**Prerequisites:** Sprint 2 Complete (v0.1.0 released on GitHub)

---

## Overview

| Task | Description | Priority | Est. Time |
|------|-------------|----------|-----------|
| 3.1 | Update Homebrew Formula with SHA256 | HIGH | 30 min |
| 3.2 | Test Homebrew Installation | HIGH | 15 min |
| 3.3 | Create Install Scripts | MEDIUM | 30 min |
| 3.4 | Enable Docker in GoReleaser | MEDIUM | 30 min |
| 3.5 | Test Docker Image | MEDIUM | 15 min |
| 3.6 | Update Documentation | MEDIUM | 30 min |

---

## Prompt 3.1: Update Homebrew Formula with SHA256

### Objective
Update the Homebrew formula with real SHA256 checksums from the v0.1.0 release.

### Prompt
```
Update the Homebrew formula at homebrew-tap/Formula/machpay.rb with the actual
SHA256 checksums from the v0.1.0 release.

Steps:

1. Get checksums from release:
   curl -sL https://github.com/machpay-xyz/machpay-cli/releases/download/v0.1.0/checksums.txt

2. Extract the SHA256 for each platform:
   - machpay_darwin_arm64.tar.gz
   - machpay_darwin_amd64.tar.gz
   - machpay_linux_arm64.tar.gz
   - machpay_linux_amd64.tar.gz

3. Update Formula/machpay.rb:
   - Replace PLACEHOLDER_SHA256_DARWIN_ARM64 with actual hash
   - Replace PLACEHOLDER_SHA256_DARWIN_AMD64 with actual hash
   - Replace PLACEHOLDER_SHA256_LINUX_ARM64 with actual hash
   - Replace PLACEHOLDER_SHA256_LINUX_AMD64 with actual hash

4. Verify the formula syntax:
   brew audit --strict Formula/machpay.rb

5. Commit and push to homebrew-tap repo:
   git add Formula/machpay.rb
   git commit -m "Update to v0.1.0 with real checksums"
   git push origin main

Expected Formula after update:
```

### Current Formula Location
`/Users/abhishektomar/Desktop/git/homebrew-tap/Formula/machpay.rb`

### Verification Commands
```bash
# Get checksums
curl -sL https://github.com/machpay-xyz/machpay-cli/releases/download/v0.1.0/checksums.txt

# After updating, verify syntax
cd /Users/abhishektomar/Desktop/git/homebrew-tap
brew audit --strict Formula/machpay.rb
```

### Expected Changes
```ruby
# Update these lines with actual hashes:
sha256 "8861644a08b87a29a7f0d7a3fe760deb88d88cadeae46ecd5f696ff4391c2d40" # darwin_arm64
sha256 "ACTUAL_HASH_HERE" # darwin_amd64
sha256 "ACTUAL_HASH_HERE" # linux_arm64
sha256 "ACTUAL_HASH_HERE" # linux_amd64
```

### Verification Checklist
- [ ] checksums.txt downloaded from release
- [ ] All 4 SHA256 hashes extracted
- [ ] Formula updated with real hashes
- [ ] `brew audit` passes
- [ ] Changes pushed to homebrew-tap

---

## Prompt 3.2: Test Homebrew Installation

### Objective
Verify that `brew install machpay/tap/machpay` works correctly.

### Prompt
```
Test the Homebrew installation flow for machpay CLI.

Steps:

1. Add the tap (if not already added):
   brew tap machpay/tap https://github.com/machpay-xyz/homebrew-tap

2. Install machpay:
   brew install machpay/tap/machpay

3. Verify installation:
   which machpay
   machpay version
   machpay --help

4. Test basic commands:
   machpay status
   machpay version --json

5. Uninstall and reinstall to test upgrade path:
   brew uninstall machpay
   brew install machpay/tap/machpay

6. Test upgrade (for future releases):
   # After v0.1.1 is released:
   brew upgrade machpay

If issues occur:
- Check brew audit output
- Verify SHA256 matches
- Check formula URL resolves correctly
```

### Verification Commands
```bash
# Full test sequence
brew tap machpay/tap https://github.com/machpay-xyz/homebrew-tap
brew install machpay/tap/machpay
machpay version
machpay status
machpay --help

# Cleanup
brew uninstall machpay
brew untap machpay/tap
```

### Expected Output
```
==> Downloading https://github.com/machpay-xyz/machpay-cli/releases/download/v0.1.0/machpay_darwin_arm64.tar.gz
######################################################################## 100.0%
==> Installing machpay from machpay/tap
🍺  /opt/homebrew/Cellar/machpay/0.1.0: 3 files, 10MB, built in 2 seconds

$ machpay version
machpay version 0.1.0
```

### Verification Checklist
- [ ] `brew tap` succeeds
- [ ] `brew install` downloads correct binary
- [ ] `machpay version` shows 0.1.0
- [ ] `machpay --help` displays usage
- [ ] `machpay status` runs without errors
- [ ] Uninstall/reinstall works

---

## Prompt 3.3: Create Install Scripts

### Objective
Create cross-platform install scripts for users who don't use Homebrew.

### Prompt
```
Create install scripts for machpay CLI that work across platforms.

1. Create install.sh for Linux/macOS:

scripts/install.sh:
- Detect OS (darwin/linux) and architecture (amd64/arm64)
- Download correct binary from GitHub releases
- Verify checksum
- Install to /usr/local/bin (or ~/.local/bin if no sudo)
- Make executable
- Verify installation

Features:
- --version flag to install specific version (default: latest)
- --no-verify to skip checksum verification
- --dry-run to show what would be done
- Colored output with progress

2. Create install.ps1 for Windows:

scripts/install.ps1:
- Download Windows binary from GitHub releases
- Install to user's PATH
- Verify installation

3. Test the scripts:
   # Linux/macOS:
   curl -fsSL https://raw.githubusercontent.com/machpay-xyz/machpay-cli/main/scripts/install.sh | sh
   
   # Windows PowerShell:
   iwr https://raw.githubusercontent.com/machpay-xyz/machpay-cli/main/scripts/install.ps1 | iex

Provide complete scripts with error handling and user-friendly output.
```

### Expected install.sh
```bash
#!/bin/bash
# MachPay CLI Installer
# Usage: curl -fsSL https://machpay.xyz/install.sh | sh

set -euo pipefail

VERSION="${MACHPAY_VERSION:-latest}"
INSTALL_DIR="${MACHPAY_INSTALL_DIR:-/usr/local/bin}"
GITHUB_REPO="machpay-xyz/machpay-cli"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1" >&2; exit 1; }

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) error "Unsupported architecture: $ARCH" ;;
    esac
    
    case "$OS" in
        darwin|linux) ;;
        *) error "Unsupported OS: $OS" ;;
    esac
    
    echo "${OS}_${ARCH}"
}

# Get latest version from GitHub
get_latest_version() {
    curl -sL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | 
        grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/'
}

# Download and install
main() {
    info "Installing MachPay CLI..."
    
    PLATFORM=$(detect_platform)
    info "Detected platform: $PLATFORM"
    
    if [ "$VERSION" = "latest" ]; then
        VERSION=$(get_latest_version)
    fi
    info "Installing version: $VERSION"
    
    DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/v${VERSION}/machpay_${PLATFORM}.tar.gz"
    
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT
    
    info "Downloading from $DOWNLOAD_URL..."
    curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/machpay.tar.gz"
    
    info "Extracting..."
    tar -xzf "$TMP_DIR/machpay.tar.gz" -C "$TMP_DIR"
    
    info "Installing to $INSTALL_DIR..."
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TMP_DIR/machpay" "$INSTALL_DIR/machpay"
    else
        sudo mv "$TMP_DIR/machpay" "$INSTALL_DIR/machpay"
    fi
    chmod +x "$INSTALL_DIR/machpay"
    
    info "Verifying installation..."
    if command -v machpay &>/dev/null; then
        echo ""
        machpay version
        echo ""
        info "✅ MachPay CLI installed successfully!"
        info "Run 'machpay login' to get started."
    else
        error "Installation failed. Please check your PATH."
    fi
}

main "$@"
```

### Verification Checklist
- [ ] install.sh created and tested on macOS
- [ ] install.sh tested on Linux (if available)
- [ ] install.ps1 created for Windows
- [ ] Scripts handle errors gracefully
- [ ] Scripts show progress and success messages

---

## Prompt 3.4: Enable Docker in GoReleaser

### Objective
Configure GoReleaser to automatically build and push Docker images on release.

### Prompt
```
Enable Docker image publishing in GoReleaser configuration.

1. Update .goreleaser.yaml to add Docker configuration:

dockers:
  - image_templates:
      - "ghcr.io/machpay-xyz/cli:{{ .Tag }}"
      - "ghcr.io/machpay-xyz/cli:v{{ .Major }}"
      - "ghcr.io/machpay-xyz/cli:v{{ .Major }}.{{ .Minor }}"
      - "ghcr.io/machpay-xyz/cli:latest"
    dockerfile: Dockerfile
    use: buildx
    build_flag_templates:
      - "--platform=linux/amd64,linux/arm64"
      - "--label=org.opencontainers.image.source=https://github.com/machpay-xyz/machpay-cli"
      - "--label=org.opencontainers.image.version={{ .Version }}"
      - "--label=org.opencontainers.image.created={{ .Date }}"
    goos: linux
    goarch: amd64

2. Update release.yml workflow to login to GHCR:

- name: Login to GitHub Container Registry
  uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}

3. Verify Dockerfile works with GoReleaser:
   goreleaser release --snapshot --clean --skip=publish

4. Test Docker build locally:
   docker build -t machpay-test .

Note: The Dockerfile expects the binary to be copied from GoReleaser's build.
```

### Current Dockerfile Location
`/Users/abhishektomar/Desktop/git/machpay-cli/Dockerfile`

### GoReleaser Docker Config
```yaml
# Add to .goreleaser.yaml
dockers:
  - image_templates:
      - "ghcr.io/machpay-xyz/cli:{{ .Tag }}"
      - "ghcr.io/machpay-xyz/cli:v{{ .Major }}"
      - "ghcr.io/machpay-xyz/cli:v{{ .Major }}.{{ .Minor }}"
      - "ghcr.io/machpay-xyz/cli:latest"
    dockerfile: Dockerfile
    use: buildx
    build_flag_templates:
      - "--platform=linux/amd64"
      - "--label=org.opencontainers.image.source=https://github.com/machpay-xyz/machpay-cli"
      - "--label=org.opencontainers.image.version={{ .Version }}"
      - "--label=org.opencontainers.image.revision={{ .FullCommit }}"
      - "--label=org.opencontainers.image.created={{ .Date }}"
    goos: linux
    goarch: amd64
```

### Release Workflow Update
```yaml
# Add before GoReleaser step in release.yml
- name: Login to GitHub Container Registry
  uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}

- name: Set up Docker Buildx
  uses: docker/setup-buildx-action@v3
```

### Verification Checklist
- [ ] .goreleaser.yaml updated with docker config
- [ ] release.yml updated with GHCR login
- [ ] `goreleaser check` passes
- [ ] Local Docker build works
- [ ] Snapshot build creates Docker image

---

## Prompt 3.5: Test Docker Image

### Objective
Verify Docker image works correctly after release.

### Prompt
```
Test the Docker image for machpay CLI.

1. Build locally for testing:
   cd /Users/abhishektomar/Desktop/git/machpay-cli
   
   # Build the Go binary first
   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o machpay ./cmd/machpay
   
   # Build Docker image
   docker build -t machpay-test .

2. Test basic commands:
   docker run --rm machpay-test version
   docker run --rm machpay-test --help
   docker run --rm machpay-test status

3. Test with volume mount for config persistence:
   docker run --rm -v ~/.machpay:/home/machpay/.machpay machpay-test status

4. After release with Docker enabled:
   docker pull ghcr.io/machpay-xyz/cli:latest
   docker run --rm ghcr.io/machpay-xyz/cli:latest version
   docker run --rm ghcr.io/machpay-xyz/cli:latest status

5. Test health check:
   docker run -d --name machpay-health ghcr.io/machpay-xyz/cli:latest sleep 3600
   docker inspect --format='{{.State.Health.Status}}' machpay-health
   docker stop machpay-health && docker rm machpay-health

Report any issues with the Docker image.
```

### Verification Commands
```bash
# Local build test
cd /Users/abhishektomar/Desktop/git/machpay-cli
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o machpay ./cmd/machpay
docker build -t machpay-test .

# Test commands
docker run --rm machpay-test version
docker run --rm machpay-test --help
docker run --rm machpay-test status

# Check image size
docker images machpay-test --format "{{.Size}}"
```

### Expected Output
```
$ docker run --rm machpay-test version
machpay version 0.1.0

$ docker run --rm machpay-test --help
MachPay CLI is the unified entry point for the MachPay Network...

$ docker images machpay-test --format "{{.Size}}"
~15MB
```

### Verification Checklist
- [ ] Docker image builds successfully
- [ ] Image size is reasonable (<50MB)
- [ ] `machpay version` works in container
- [ ] `machpay --help` works in container
- [ ] Volume mount for config works
- [ ] Health check passes

---

## Prompt 3.6: Update Documentation

### Objective
Update README and create INSTALL.md with all installation methods.

### Prompt
```
Update documentation with all installation methods.

1. Create INSTALL.md with comprehensive installation guide:

# Installation

## Quick Install

### Homebrew (macOS/Linux)
```bash
brew tap machpay/tap
brew install machpay
```

### Shell Script (macOS/Linux)
```bash
curl -fsSL https://machpay.xyz/install.sh | sh
```

### PowerShell (Windows)
```powershell
iwr https://machpay.xyz/install.ps1 | iex
```

### Docker
```bash
docker pull ghcr.io/machpay-xyz/cli:latest
docker run --rm ghcr.io/machpay-xyz/cli version
```

### Go Install
```bash
go install github.com/machpay-xyz/machpay-cli/cmd/machpay@latest
```

### Manual Download
Download the appropriate binary from the releases page...

## Verification
```bash
machpay version
machpay --help
```

## Upgrading
...

## Uninstalling
...

2. Update README.md installation section to link to INSTALL.md

3. Add badges to README.md:
   - Release version badge
   - Go version badge
   - License badge
   - Docker pulls badge

Provide complete markdown files.
```

### README Badges
```markdown
[![Release](https://img.shields.io/github/v/release/machpay-xyz/machpay-cli)](https://github.com/machpay-xyz/machpay-cli/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/machpay-xyz/machpay-cli)](https://go.dev/)
[![License](https://img.shields.io/github/license/machpay-xyz/machpay-cli)](LICENSE)
[![Docker](https://img.shields.io/docker/pulls/ghcr.io/machpay-xyz/cli)](https://ghcr.io/machpay-xyz/cli)
```

### Verification Checklist
- [ ] INSTALL.md created with all methods
- [ ] README.md updated with badges
- [ ] Installation section links to INSTALL.md
- [ ] All install commands are tested
- [ ] Upgrade/uninstall instructions included

---

## Sprint 3 Completion Checklist

```
Sprint 3 Final Checklist:

Homebrew:
[ ] SHA256 hashes updated in formula
[ ] brew audit passes
[ ] brew install machpay/tap/machpay works
[ ] Formula pushed to homebrew-tap repo

Install Scripts:
[ ] install.sh created and tested
[ ] install.ps1 created (Windows)
[ ] Scripts handle errors gracefully

Docker:
[ ] GoReleaser Docker config added
[ ] release.yml updated with GHCR login
[ ] Docker image builds locally
[ ] Image size < 50MB

Documentation:
[ ] INSTALL.md created
[ ] README.md badges added
[ ] All install methods documented
[ ] Verification steps included
```

---

## Commands Summary

```bash
# Navigate to repos
cd /Users/abhishektomar/Desktop/git/machpay-cli
cd /Users/abhishektomar/Desktop/git/homebrew-tap

# Prompt 3.1 - Update Homebrew formula
curl -sL https://github.com/machpay-xyz/machpay-cli/releases/download/v0.1.0/checksums.txt
# Update Formula/machpay.rb with hashes
brew audit --strict Formula/machpay.rb

# Prompt 3.2 - Test Homebrew
brew tap machpay/tap https://github.com/machpay-xyz/homebrew-tap
brew install machpay/tap/machpay
machpay version

# Prompt 3.3 - Install scripts
# Create scripts/install.sh
chmod +x scripts/install.sh
./scripts/install.sh --dry-run

# Prompt 3.4 - Docker setup
# Update .goreleaser.yaml with docker config
goreleaser check
goreleaser release --snapshot --clean --skip=publish

# Prompt 3.5 - Test Docker
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o machpay ./cmd/machpay
docker build -t machpay-test .
docker run --rm machpay-test version

# Prompt 3.6 - Documentation
# Create INSTALL.md
# Update README.md
```

---

## Next Sprint

After Sprint 3 is complete, proceed to **Sprint 4: Desktop App Fixes**

Sprint 4 will:
- Fix Google OAuth in Tauri app (critical bug)
- Set up macOS code signing
- Create cross-platform installers

---

**Document Version:** 1.0  
**Created:** 2025-01-01


