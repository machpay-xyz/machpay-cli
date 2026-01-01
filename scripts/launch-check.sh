#!/bin/bash
# ============================================================
# MachPay CLI - Launch Readiness Check
# ============================================================
#
# This script verifies that all components are ready for launch.
# Run before tagging a release.
#
# Usage: ./scripts/launch-check.sh [version]
#
# ============================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

VERSION="${1:-latest}"
PASS=0
FAIL=0
WARN=0

# ============================================================
# Helper Functions
# ============================================================

print_header() {
    echo ""
    echo -e "${BLUE}╔══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║${NC}  🚀 MachPay CLI - Launch Readiness Check                  ${BLUE}║${NC}"
    echo -e "${BLUE}║${NC}     Version: ${YELLOW}${VERSION}${NC}                                          ${BLUE}║${NC}"
    echo -e "${BLUE}╚══════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

check() {
    local name="$1"
    local cmd="$2"
    local critical="${3:-true}"

    printf "  %-45s " "$name"
    
    if eval "$cmd" >/dev/null 2>&1; then
        echo -e "[${GREEN}✓${NC}]"
        ((PASS++))
        return 0
    else
        if [ "$critical" = "true" ]; then
            echo -e "[${RED}✗${NC}]"
            ((FAIL++))
        else
            echo -e "[${YELLOW}!${NC}]"
            ((WARN++))
        fi
        return 1
    fi
}

section() {
    echo ""
    echo -e "${BLUE}## $1${NC}"
    echo ""
}

# ============================================================
# Main Checks
# ============================================================

print_header

# ------------------------------------------------------------
# Local Build Checks
# ------------------------------------------------------------
section "Local Build"

check "Go modules" "go mod verify"
check "Go build" "go build -o /tmp/machpay-check ./cmd/machpay"
check "Go tests" "go test ./..."
check "Go vet" "go vet ./..."
check "GoReleaser config" "goreleaser check" "false"

# ------------------------------------------------------------
# Code Quality
# ------------------------------------------------------------
section "Code Quality"

check "No blocking TODOs" "! grep -r 'TODO:BLOCK' --include='*.go' internal/ cmd/ 2>/dev/null" "false"
check "No FIXME comments" "! grep -r 'FIXME' --include='*.go' internal/ cmd/ 2>/dev/null" "false"
check "README exists" "test -f README.md"
check "LICENSE exists" "test -f LICENSE"
check "CONTRIBUTING exists" "test -f CONTRIBUTING.md"

# ------------------------------------------------------------
# Version Checks
# ------------------------------------------------------------
section "Version"

if [ "$VERSION" != "latest" ]; then
    check "Version in code" "grep -q '$VERSION' cmd/machpay/main.go" "false"
fi
check "Changelog updated" "test -f CHANGELOG.md" "false"

# ------------------------------------------------------------
# Distribution (only if version is specified)
# ------------------------------------------------------------
if [ "$VERSION" != "latest" ]; then
    section "Distribution (GitHub Release)"

    RELEASE_URL="https://github.com/machpay/machpay-cli/releases/download/$VERSION"

    check "Darwin/amd64 binary" "curl -sLI ${RELEASE_URL}/machpay_darwin_amd64.tar.gz | grep -q '200'"
    check "Darwin/arm64 binary" "curl -sLI ${RELEASE_URL}/machpay_darwin_arm64.tar.gz | grep -q '200'"
    check "Linux/amd64 binary" "curl -sLI ${RELEASE_URL}/machpay_linux_amd64.tar.gz | grep -q '200'"
    check "Linux/arm64 binary" "curl -sLI ${RELEASE_URL}/machpay_linux_arm64.tar.gz | grep -q '200'"
    check "Windows/amd64 binary" "curl -sLI ${RELEASE_URL}/machpay_windows_amd64.zip | grep -q '200'"
    check "Checksums file" "curl -sL ${RELEASE_URL}/checksums.txt | grep -q 'machpay'"
fi

# ------------------------------------------------------------
# Install Scripts
# ------------------------------------------------------------
section "Install Scripts"

check "install.sh exists" "test -f scripts/install.sh"
check "install.sh executable" "test -x scripts/install.sh"
check "install.sh shebang" "head -1 scripts/install.sh | grep -q '#!/bin/bash'"
check "install.ps1 exists" "test -f scripts/install.ps1"

# ------------------------------------------------------------
# Docker (only if version is specified)
# ------------------------------------------------------------
if [ "$VERSION" != "latest" ]; then
    section "Docker"
    check "Docker image exists" "docker manifest inspect ghcr.io/machpay/cli:$VERSION 2>/dev/null" "false"
fi

# ------------------------------------------------------------
# Documentation
# ------------------------------------------------------------
section "Documentation"

check "README has installation" "grep -q 'Installation' README.md"
check "README has quick start" "grep -qi 'quick start' README.md"
check "README has commands" "grep -q 'Commands' README.md"
check "Help text works" "/tmp/machpay-check --help | grep -q 'machpay'"
check "Version flag works" "/tmp/machpay-check version | grep -q 'version'"

# ------------------------------------------------------------
# Summary
# ------------------------------------------------------------
echo ""
echo -e "${BLUE}══════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "  Results:"
echo -e "    ${GREEN}Passed:${NC}  $PASS"
echo -e "    ${RED}Failed:${NC}  $FAIL"
echo -e "    ${YELLOW}Warnings:${NC} $WARN"
echo ""

# Cleanup
rm -f /tmp/machpay-check

if [ $FAIL -gt 0 ]; then
    echo -e "  ${RED}❌ Launch check FAILED${NC}"
    echo ""
    echo "  Fix the failing checks before release."
    echo ""
    exit 1
else
    if [ $WARN -gt 0 ]; then
        echo -e "  ${YELLOW}⚠️  Launch check PASSED with warnings${NC}"
    else
        echo -e "  ${GREEN}✅ Launch check PASSED${NC}"
    fi
    echo ""
    echo "  Ready for release!"
    echo ""
    exit 0
fi

