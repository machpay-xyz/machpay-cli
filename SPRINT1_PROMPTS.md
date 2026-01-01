# Sprint 1: Foundation & CI/CD - Execution Prompts

**Sprint Goal:** Enable automated releases with proper CI/CD  
**Estimated Time:** 2-3 hours  
**Prerequisites:** None (this is the starting sprint)

---

## Overview

| Task | Description | Priority | Est. Time |
|------|-------------|----------|-----------|
| 1.1 | GitHub Actions Release Workflow | HIGH | 45 min |
| 1.2 | PR Build/Test Workflow | HIGH | 30 min |
| 1.3 | Verify GoReleaser Configuration | HIGH | 30 min |
| 1.4 | Test Coverage Verification | HIGH | 45 min |

---

## Prompt 1.1: GitHub Actions Release Workflow

### Objective
Create a GitHub Actions workflow that automatically builds and releases CLI binaries when a version tag is pushed.

### Prompt
```
Create a GitHub Actions workflow file at machpay-cli/.github/workflows/release.yml

Requirements:
1. Trigger on push of tags matching pattern: v*.*.*
2. Single job called "release" running on ubuntu-latest
3. Steps:
   a. Checkout code with fetch-depth: 0 (needed for GoReleaser changelog)
   b. Setup Go 1.21 using actions/setup-go@v5
   c. Cache Go modules using actions/cache@v4
   d. Run tests: go test -v ./...
   e. Run GoReleaser using goreleaser/goreleaser-action@v5
      - args: release --clean
      - Set GITHUB_TOKEN env var from secrets.GITHUB_TOKEN

4. Permissions needed:
   - contents: write (for creating releases)
   - packages: write (for GHCR if using Docker)

5. Add concurrency group to prevent duplicate runs

The workflow should create:
- Binaries for macOS (arm64, amd64), Linux (amd64, arm64), Windows (amd64)
- Checksums file
- GitHub Release with auto-generated changelog

Reference the existing .goreleaser.yaml in the repo.
```

### Expected Output
```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*.*.*'

permissions:
  contents: write
  packages: write

concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: true

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache: true

      - name: Run Tests
        run: go test -v ./...

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Verification
- [ ] File created at `.github/workflows/release.yml`
- [ ] YAML syntax is valid
- [ ] Workflow appears in GitHub Actions tab (after push)

---

## Prompt 1.2: PR Build/Test Workflow

### Objective
Create a CI workflow that runs on every pull request to catch issues before merge.

### Prompt
```
Create a GitHub Actions workflow file at machpay-cli/.github/workflows/ci.yml

Requirements:
1. Trigger on:
   - Pull requests to main branch
   - Push to main branch

2. Job matrix testing on:
   - ubuntu-latest
   - macos-latest
   - windows-latest

3. Steps for each platform:
   a. Checkout code
   b. Setup Go 1.21 with module caching
   c. Download dependencies: go mod download
   d. Run go vet: go vet ./...
   e. Run golangci-lint using golangci/golangci-lint-action@v4
   f. Run tests with coverage: go test -race -coverprofile=coverage.out -covermode=atomic ./...
   g. Upload coverage to Codecov (optional)

4. Add a separate job to check that GoReleaser config is valid:
   - Run: goreleaser check
   - Run: goreleaser release --snapshot --clean --skip-publish

5. Fail the workflow if:
   - Any test fails
   - Linter finds issues
   - GoReleaser config is invalid

Include timeout-minutes: 15 for each job.
```

### Expected Output
```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    name: Test (${{ matrix.os }})
    runs-on: ${{ matrix.os }}
    timeout-minutes: 15
    strategy:
      fail-fast: false
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache: true

      - name: Download Dependencies
        run: go mod download

      - name: Run go vet
        run: go vet ./...

      - name: Run Tests
        run: go test -race -coverprofile=coverage.out -covermode=atomic -v ./...

      - name: Upload Coverage
        if: matrix.os == 'ubuntu-latest'
        uses: codecov/codecov-action@v4
        with:
          files: coverage.out
          fail_ci_if_error: false

  lint:
    name: Lint
    runs-on: ubuntu-latest
    timeout-minutes: 10
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache: true

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
          args: --timeout=5m

  goreleaser-check:
    name: GoReleaser Check
    runs-on: ubuntu-latest
    timeout-minutes: 10
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache: true

      - name: Check GoReleaser Config
        uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: check

      - name: Test GoReleaser Build
        uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: release --snapshot --clean --skip-publish
```

### Verification
- [ ] File created at `.github/workflows/ci.yml`
- [ ] Workflow runs on PR creation
- [ ] All 3 platforms pass tests
- [ ] Linter runs without errors
- [ ] GoReleaser check passes

---

## Prompt 1.3: Verify GoReleaser Configuration

### Objective
Review and update the existing GoReleaser configuration to ensure it's production-ready.

### Prompt
```
Review and update the existing machpay-cli/.goreleaser.yaml file.

Verify and ensure the following configuration:

1. Project name: machpay

2. Builds section:
   - Binary name: machpay
   - Main package: ./cmd/machpay
   - CGO disabled (env: CGO_ENABLED=0)
   - ldflags for version injection:
     -s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}}
   - Target platforms:
     - darwin/amd64
     - darwin/arm64
     - linux/amd64
     - linux/arm64
     - windows/amd64

3. Archives section:
   - Format: tar.gz for Unix, zip for Windows
   - Name template: {{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}
   - Include files: LICENSE, README.md

4. Checksum:
   - Name template: checksums.txt
   - Algorithm: sha256

5. Changelog:
   - Sort by: asc
   - Filters to exclude: docs, test, chore

6. Release section:
   - GitHub release enabled
   - Draft: false
   - Prerelease: auto (based on tag)

Test the configuration:
1. Run: goreleaser check
2. Run: goreleaser release --snapshot --clean --skip-publish
3. Verify artifacts in dist/ folder

If the main.go doesn't have version variables, add them.
```

### Expected .goreleaser.yaml
```yaml
# .goreleaser.yaml
project_name: machpay

before:
  hooks:
    - go mod tidy
    - go generate ./...

builds:
  - id: machpay
    binary: machpay
    main: ./cmd/machpay
    env:
      - CGO_ENABLED=0
    goos:
      - darwin
      - linux
      - windows
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: windows
        goarch: arm64
    ldflags:
      - -s -w
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}

archives:
  - id: default
    name_template: >-
      {{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}
    format_overrides:
      - goos: windows
        format: zip
    files:
      - LICENSE
      - README.md

checksum:
  name_template: 'checksums.txt'
  algorithm: sha256

changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'
      - '^chore:'
      - '^ci:'
      - Merge pull request
      - Merge branch

release:
  github:
    owner: machpay-xyz
    name: machpay-cli
  draft: false
  prerelease: auto
  name_template: "{{.Tag}}"
```

### Version Variables in main.go
```go
// cmd/machpay/main.go
package main

var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

// Use in version command:
// fmt.Printf("machpay version %s (commit: %s, built: %s)\n", version, commit, date)
```

### Verification
- [ ] `goreleaser check` passes
- [ ] `goreleaser release --snapshot --clean --skip-publish` succeeds
- [ ] `dist/` folder contains binaries for all platforms
- [ ] Version variables are injected correctly

---

## Prompt 1.4: Test Coverage Verification

### Objective
Verify test coverage meets the 80% threshold and add tests if needed.

### Prompt
```
Check and improve test coverage for machpay-cli to meet the >80% threshold.

Step 1: Generate coverage report
cd machpay-cli
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | tail -20
go tool cover -html=coverage.out -o coverage.html

Step 2: Analyze coverage by package
go tool cover -func=coverage.out | grep -E "^[a-z]" | sort -t$'\t' -k3 -n

Step 3: Identify packages below 80%
Look for packages with coverage < 80% and prioritize:
- internal/auth/ (critical for security)
- internal/cmd/ (user-facing commands)
- internal/config/ (configuration handling)
- internal/gateway/ (gateway management)

Step 4: For each low-coverage package, add tests for:
- Happy path scenarios
- Error handling paths
- Edge cases
- Input validation

Step 5: Create test utilities if needed
- internal/testutil/testutil.go for shared test helpers

Step 6: Re-run coverage and verify >80%
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total

Report findings:
- Current coverage percentage
- Packages below 80%
- Recommended tests to add
- Final coverage after improvements
```

### Test File Templates

#### Auth Tests (internal/auth/auth_test.go)
```go
package auth

import (
    "testing"
)

func TestTokenStorage(t *testing.T) {
    // Test saving and loading tokens
}

func TestTokenValidation(t *testing.T) {
    // Test token validation logic
}

func TestTokenExpiration(t *testing.T) {
    // Test expired token handling
}
```

#### Config Tests (internal/config/config_test.go)
```go
package config

import (
    "testing"
    "os"
)

func TestLoadConfig(t *testing.T) {
    // Test loading config from file
}

func TestSaveConfig(t *testing.T) {
    // Test saving config to file
}

func TestDefaultConfig(t *testing.T) {
    // Test default config values
}

func TestConfigValidation(t *testing.T) {
    // Test config validation rules
}
```

#### Gateway Tests (internal/gateway/gateway_test.go)
```go
package gateway

import (
    "testing"
)

func TestDownloadGateway(t *testing.T) {
    // Test gateway download (mock HTTP)
}

func TestProcessManagement(t *testing.T) {
    // Test start/stop/status
}

func TestHealthCheck(t *testing.T) {
    // Test health check logic
}
```

### Verification
- [ ] Coverage report generated
- [ ] All packages >80% coverage (or documented exceptions)
- [ ] Critical packages (auth, config, gateway) well tested
- [ ] No test failures
- [ ] Coverage badge can be added to README

---

## Sprint 1 Completion Checklist

After completing all prompts, verify:

```
Sprint 1 Checklist:

[ ] .github/workflows/release.yml created
[ ] .github/workflows/ci.yml created
[ ] Both workflows have valid YAML syntax
[ ] GoReleaser config verified with `goreleaser check`
[ ] Snapshot build works: `goreleaser release --snapshot --clean --skip-publish`
[ ] Test coverage >= 80%
[ ] All tests passing on all platforms
[ ] Linter passes without errors
[ ] Commit all changes to main branch
```

### Commands Summary

```bash
# Navigate to repo
cd /Users/abhishektomar/Desktop/git/machpay-cli

# Create workflows directory
mkdir -p .github/workflows

# Test GoReleaser
goreleaser check
goreleaser release --snapshot --clean --skip-publish

# Test coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total

# Run full test suite
go test -v ./...

# Run linter
golangci-lint run

# Commit changes
git add .
git commit -m "ci: add GitHub Actions workflows and improve test coverage"
git push origin main
```

---

## Next Sprint

After Sprint 1 is complete, proceed to **Sprint 2: First Release (v0.1.0)**

Sprint 2 will use the CI/CD infrastructure set up in Sprint 1 to create the first official release.

---

**Document Version:** 1.0  
**Created:** 2025-01-01

