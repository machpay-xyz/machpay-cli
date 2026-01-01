# Phase 5: Distribution - Task Prompts

**Document Purpose:** Step-by-step prompts to complete Phase 5 (Distribution)  
**Created:** 2025-01-01  
**Status:** Ready to execute

---

## Task 1: GitHub Actions Release Workflow

### Prompt 1.1: Create Release Workflow
```
Create a GitHub Actions workflow file at .github/workflows/release.yml that:
1. Triggers on tag push (pattern: v*.*.*)
2. Checks out code with full git history
3. Sets up Go 1.21+
4. Runs tests before release
5. Runs GoReleaser with GITHUB_TOKEN
6. Uploads artifacts to GitHub Releases
7. Includes job for building desktop app (Tauri)

Use best practices for GitHub Actions and GoReleaser integration.
Include caching for Go modules and build artifacts.
Add status badges that can be used in README.
```

### Prompt 1.2: Create Pre-release Workflow
```
Create a GitHub Actions workflow file at .github/workflows/pr-build.yml that:
1. Triggers on pull requests
2. Runs tests on multiple platforms (ubuntu, macos, windows)
3. Runs GoReleaser with --snapshot --skip-publish
4. Validates that binaries build successfully
5. Runs linter (golangci-lint)

This ensures PRs don't break releases.
```

### Prompt 1.3: Test Workflow Locally
```
Install act (GitHub Actions local runner) and test the release workflow:
1. Install: brew install act
2. Create .secrets file with GITHUB_TOKEN
3. Run: act -j release --secret-file .secrets
4. Verify artifacts are generated correctly

Document any issues and create fixes.
```

---

## Task 2: First GitHub Release (v0.1.0)

### Prompt 2.1: Prepare Release Checklist
```
Create a pre-release checklist at docs/RELEASE_CHECKLIST.md that includes:
1. Run full test suite (make test)
2. Update CHANGELOG.md with release notes
3. Bump version in relevant files (main.go, package.json)
4. Verify all platforms build (make build-all)
5. Test install scripts locally
6. Update documentation with new version
7. Create git tag
8. Push tag to trigger release

Mark each step as [ ] for checkbox tracking.
```

### Prompt 2.2: Update Version Numbers
```
Update version numbers to 0.1.0 in:
1. cmd/machpay/main.go (version constant)
2. machpay-console/package.json
3. machpay-console/src-tauri/tauri.conf.json
4. .goreleaser.yaml (if version is hardcoded)

Verify all version references are consistent.
Create a script scripts/bump-version.sh for future use.
```

### Prompt 2.3: Create Release Notes
```
Create comprehensive release notes for v0.1.0 in CHANGELOG.md:

Format:
# v0.1.0 - Initial Public Release (2025-01-XX)

## 🎉 Highlights
- First public release of MachPay CLI
- Desktop app with Tauri
- Cross-platform support

## ✨ Features
[List all major features from Phase 1-4]

## 📦 Installation
[Brief install instructions]

## 🐛 Known Issues
[Reference KNOWN_ISSUES.md]

## 📝 Documentation
[Links to docs]

Include emoji, be professional but exciting.
```

### Prompt 2.4: Create and Push Git Tag
```
Guide me through creating and pushing the v0.1.0 tag:

1. Verify all changes are committed
2. Create annotated tag: git tag -a v0.1.0 -m "Release v0.1.0"
3. Verify tag: git tag -l -n9 v0.1.0
4. Push tag: git push origin v0.1.0
5. Monitor GitHub Actions workflow
6. Verify release artifacts on GitHub Releases page

Provide the exact commands with safety checks.
```

---

## Task 3: Homebrew Tap Setup

### Prompt 3.1: Verify Homebrew Tap Repository
```
Check the homebrew-tap repository:
1. Verify repo exists at github.com/machpay/homebrew-tap
2. Check Formula/machpay.rb file structure
3. Verify it follows Homebrew formula standards
4. Ensure description, homepage, license are correct

If repo doesn't exist, guide me through creating it.
```

### Prompt 3.2: Update Homebrew Formula Post-Release
```
After v0.1.0 is released, update Formula/machpay.rb:
1. Update version number to "0.1.0"
2. Update download URL to actual GitHub release tarball
3. Calculate and update SHA256 checksum
4. Test formula locally: brew install --build-from-source ./Formula/machpay.rb
5. Verify installed binary works
6. Commit and push to homebrew-tap repo

Provide the complete updated formula with correct URLs.
```

### Prompt 3.3: Test Homebrew Installation
```
Test the Homebrew installation end-to-end:
1. Uninstall any existing machpay: brew uninstall machpay
2. Add tap: brew tap machpay/tap
3. Install: brew install machpay
4. Verify: machpay version
5. Test commands: machpay status, machpay setup --help
6. Check binary location: which machpay

Document any issues and create fixes.
```

---

## Task 4: Docker Image

### Prompt 4.1: Review and Optimize Dockerfile
```
Review the existing Dockerfile at machpay-cli/Dockerfile:
1. Ensure it uses multi-stage build
2. Use minimal base image (alpine or distroless)
3. Copy only necessary binaries
4. Set proper ENTRYPOINT and CMD
5. Add HEALTHCHECK if applicable
6. Optimize layer caching
7. Add labels (org.opencontainers.image.*)

Provide the optimized Dockerfile.
```

### Prompt 4.2: Add Docker Build to GoReleaser
```
Update .goreleaser.yaml to include Docker image builds:
1. Add docker section
2. Configure image_templates for multiple tags
3. Set up multi-arch builds (amd64, arm64)
4. Push to GitHub Container Registry (ghcr.io)
5. Optionally push to Docker Hub

Provide the complete docker configuration block.
```

### Prompt 4.3: Test Docker Build Locally
```
Test the Docker build locally:
1. Build image: docker build -t machpay:test .
2. Run container: docker run -it machpay:test version
3. Test commands: docker run -it machpay:test status
4. Check image size: docker images machpay:test
5. Scan for vulnerabilities: docker scan machpay:test (or trivy)

Verify the image works correctly and document the size.
```

### Prompt 4.4: Create Docker Documentation
```
Create docs/DOCKER.md with:
1. How to pull the image
2. How to run the CLI in Docker
3. Volume mounting for config persistence
4. Environment variable configuration
5. Docker Compose example for full stack
6. Common use cases and examples

Make it beginner-friendly with copy-paste examples.
```

---

## Task 5: Desktop App Distribution Packages

### Prompt 5.1: macOS DMG Creation
```
Set up macOS DMG creation for the desktop app:
1. Update machpay-console/src-tauri/tauri.conf.json bundle settings
2. Configure DMG options (window size, icon positioning)
3. Add code signing configuration (for future signing)
4. Build DMG: cd machpay-console && npm run tauri:build
5. Test DMG installation on clean macOS system
6. Document notarization steps for future

Provide the complete DMG configuration.
```

### Prompt 5.2: Windows MSI Installer
```
Set up Windows MSI installer for the desktop app:
1. Update tauri.conf.json with Windows bundle settings
2. Configure WiX installer options
3. Add custom install directory option
4. Include uninstaller
5. Add Start Menu shortcuts
6. Test on clean Windows system

Provide the Windows-specific configuration.
```

### Prompt 5.3: Linux Packages (DEB/RPM)
```
Set up Linux packages for the desktop app:
1. Configure Debian package in tauri.conf.json
2. Set up dependencies in deb section
3. Add AppImage configuration as alternative
4. Test on Ubuntu/Debian system
5. Test on Fedora/RHEL system (RPM)
6. Verify desktop file and icon installation

Provide the Linux-specific configuration.
```

### Prompt 5.4: Update Release Workflow for Desktop Packages
```
Update .github/workflows/release.yml to build desktop packages:
1. Add job for building Tauri app
2. Build for macOS (DMG)
3. Build for Windows (MSI)
4. Build for Linux (DEB, AppImage)
5. Upload all packages to GitHub Releases
6. Add checksums file

Provide the complete job configuration.
```

---

## Task 6: Installation Documentation

### Prompt 6.1: Create INSTALL.md
```
Create a comprehensive INSTALL.md guide:

Sections:
1. Quick Install (one-liner for each platform)
2. Homebrew (macOS/Linux)
3. Manual binary download
4. Docker installation
5. Desktop app installation (DMG/MSI/DEB)
6. Building from source
7. Verification steps
8. Troubleshooting common issues
9. Uninstallation

Use clear headings, code blocks, and platform badges.
```

### Prompt 6.2: Update README.md Installation Section
```
Update the main README.md:
1. Add prominent "Installation" section near the top
2. Show one-liner for each platform
3. Link to detailed INSTALL.md
4. Add badges for:
   - GitHub Release
   - Homebrew
   - Docker Pulls
   - Platform support (macOS, Linux, Windows)
5. Add "Supported Platforms" table

Keep it concise, link to INSTALL.md for details.
```

### Prompt 6.3: Create Platform-Specific Quick Start Guides
```
Create quick start guides:
1. docs/quickstart/MACOS.md
2. docs/quickstart/LINUX.md
3. docs/quickstart/WINDOWS.md

Each should include:
- Installation command
- First-time setup
- Basic commands to try
- Where to get help
- Platform-specific tips

Maximum 1 page per platform, focus on getting started fast.
```

### Prompt 6.4: Add Installation Videos/GIFs
```
Create installation demo assets:
1. Record GIF of Homebrew install process
2. Record GIF of DMG installation (macOS)
3. Record GIF of first run and setup wizard
4. Add to docs/assets/ folder
5. Reference in INSTALL.md and README.md

Use tools like:
- macOS: use QuickTime or Gifski
- Terminal: use asciinema + agg

Max 10MB per GIF, optimize with gifsicle.
```

---

## Verification Checklist

After completing all tasks, verify:

```
Phase 5 Completion Checklist:

[ ] GitHub Actions workflow triggers on tag push
[ ] v0.1.0 release published with all artifacts
[ ] Homebrew formula works: brew install machpay/tap/machpay
[ ] Docker image available: docker pull ghcr.io/machpay/machpay:latest
[ ] macOS DMG downloads and installs correctly
[ ] Windows MSI downloads and installs correctly
[ ] Linux DEB/AppImage works on Ubuntu
[ ] INSTALL.md is comprehensive and clear
[ ] README.md has updated install section
[ ] All platform-specific docs are complete
[ ] Installation process tested on all 3 platforms
```

---

## Notes

- Execute prompts in order (dependencies exist)
- Test each step before moving to next
- Document any issues in KNOWN_ISSUES.md
- Update CLI_STATUS.md as tasks complete
- Keep CHANGELOG.md updated

**Estimated Time:** 4-6 hours for all tasks

---

**Next:** Start with Prompt 1.1 (GitHub Actions Release Workflow)

