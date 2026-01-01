#!/bin/sh
# ============================================================
# MachPay CLI Installer
# ============================================================
#
# Usage:
#   curl -fsSL https://machpay.xyz/install.sh | sh
#
# Environment variables:
#   MACHPAY_INSTALL_DIR - Installation directory (default: /usr/local/bin)
#   MACHPAY_VERSION     - Specific version to install (default: latest)
#
# ============================================================

set -e

# Configuration
GITHUB_REPO="machpay/machpay-cli"
BINARY_NAME="machpay"
INSTALL_DIR="${MACHPAY_INSTALL_DIR:-}"

# Colors (only if terminal supports them)
if [ -t 1 ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    BLUE='\033[0;34m'
    CYAN='\033[0;36m'
    BOLD='\033[1m'
    NC='\033[0m'
else
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    CYAN=''
    BOLD=''
    NC=''
fi

# Print functions
info() { printf "${CYAN}→${NC} %s\n" "$1"; }
success() { printf "${GREEN}✓${NC} %s\n" "$1"; }
warn() { printf "${YELLOW}!${NC} %s\n" "$1"; }
error() { printf "${RED}✗${NC} %s\n" "$1" >&2; exit 1; }

# Banner
print_banner() {
    printf "\n"
    printf "${BLUE}███╗   ███╗ █████╗  ██████╗██╗  ██╗██████╗  █████╗ ██╗${NC}\n"
    printf "${BLUE}████╗ ████║██╔══██╗██╔════╝██║  ██║██╔══██╗██╔══██╗╚██╗${NC}\n"
    printf "${BLUE}██╔████╔██║███████║██║     ███████║██████╔╝███████║ ██║${NC}\n"
    printf "${BLUE}██║╚██╔╝██║██╔══██║██║     ██╔══██║██╔═══╝ ██╔══██║ ██║${NC}\n"
    printf "${BLUE}██║ ╚═╝ ██║██║  ██║╚██████╗██║  ██║██║     ██║  ██║██╔╝${NC}\n"
    printf "${BLUE}╚═╝     ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝     ╚═╝  ╚═╝╚═╝${NC}\n"
    printf "\n"
    printf "${BOLD}CLI Installer${NC}\n"
    printf "\n"
}

# Detect platform
detect_platform() {
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    ARCH="$(uname -m)"

    case "$ARCH" in
        x86_64|amd64)   ARCH="amd64" ;;
        aarch64|arm64)  ARCH="arm64" ;;
        armv7l)         ARCH="arm" ;;
        *)              error "Unsupported architecture: $ARCH" ;;
    esac

    case "$OS" in
        darwin)         OS="darwin" ;;
        linux)          OS="linux" ;;
        mingw*|msys*|cygwin*)
            printf "${YELLOW}Windows detected.${NC}\n"
            printf "\nPlease use PowerShell:\n"
            printf "  ${CYAN}iwr machpay.xyz/install.ps1 | iex${NC}\n\n"
            exit 1
            ;;
        *)              error "Unsupported OS: $OS" ;;
    esac

    info "Platform: ${OS}/${ARCH}"
}

# Check for required commands
check_requirements() {
    for cmd in curl tar; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            error "Required command not found: $cmd"
        fi
    done
}

# Determine install directory
determine_install_dir() {
    if [ -n "$INSTALL_DIR" ]; then
        # User specified
        return
    fi

    # Try /usr/local/bin first
    if [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
    elif command -v sudo >/dev/null 2>&1; then
        INSTALL_DIR="/usr/local/bin"
        NEED_SUDO=1
    else
        # Fall back to user directory
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
        ADD_TO_PATH=1
    fi

    info "Install directory: $INSTALL_DIR"
}

# Get latest version
get_version() {
    if [ -n "$MACHPAY_VERSION" ]; then
        VERSION="$MACHPAY_VERSION"
    else
        info "Fetching latest version..."
        VERSION=$(curl -sL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | \
            grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    fi

    if [ -z "$VERSION" ]; then
        error "Could not determine version. Check your internet connection."
    fi

    info "Version: $VERSION"
}

# Download and install
download_and_install() {
    TARBALL="${BINARY_NAME}_${OS}_${ARCH}.tar.gz"
    URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${TARBALL}"
    CHECKSUM_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/checksums.txt"

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    info "Downloading ${TARBALL}..."
    if ! curl -sL --fail -o "$TMP_DIR/$TARBALL" "$URL"; then
        error "Download failed. URL: $URL"
    fi

    # Verify checksum
    info "Verifying checksum..."
    CHECKSUMS=$(curl -sL "$CHECKSUM_URL")
    if [ -z "$CHECKSUMS" ]; then
        warn "Could not fetch checksums, skipping verification"
    else
        EXPECTED=$(echo "$CHECKSUMS" | grep "$TARBALL" | cut -d ' ' -f 1)

        if command -v sha256sum >/dev/null 2>&1; then
            ACTUAL=$(sha256sum "$TMP_DIR/$TARBALL" | cut -d ' ' -f 1)
        elif command -v shasum >/dev/null 2>&1; then
            ACTUAL=$(shasum -a 256 "$TMP_DIR/$TARBALL" | cut -d ' ' -f 1)
        else
            warn "Cannot verify checksum (no sha256sum or shasum)"
            ACTUAL="$EXPECTED"
        fi

        if [ "$EXPECTED" != "$ACTUAL" ]; then
            error "Checksum verification failed!\n  Expected: $EXPECTED\n  Actual:   $ACTUAL"
        fi
        success "Checksum verified"
    fi

    # Extract
    info "Extracting..."
    tar -xzf "$TMP_DIR/$TARBALL" -C "$TMP_DIR"

    # Find binary (might be in a subdirectory)
    BINARY_PATH="$TMP_DIR/$BINARY_NAME"
    if [ ! -f "$BINARY_PATH" ]; then
        BINARY_PATH=$(find "$TMP_DIR" -name "$BINARY_NAME" -type f | head -1)
    fi

    if [ ! -f "$BINARY_PATH" ]; then
        error "Binary not found in archive"
    fi

    # Install
    info "Installing to $INSTALL_DIR..."
    if [ "$NEED_SUDO" = "1" ]; then
        sudo mkdir -p "$INSTALL_DIR"
        sudo mv "$BINARY_PATH" "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        mkdir -p "$INSTALL_DIR"
        mv "$BINARY_PATH" "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi

    success "Installed to $INSTALL_DIR/$BINARY_NAME"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        INSTALLED_VERSION=$("$BINARY_NAME" version 2>/dev/null || echo "unknown")
        success "Installation verified: $INSTALLED_VERSION"
    elif [ -x "$INSTALL_DIR/$BINARY_NAME" ]; then
        INSTALLED_VERSION=$("$INSTALL_DIR/$BINARY_NAME" version 2>/dev/null || echo "unknown")
        success "Installation verified: $INSTALLED_VERSION"
    else
        warn "Could not verify installation"
    fi
}

# Print success message
print_success() {
    printf "\n"
    printf "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}\n"
    printf "${GREEN}║${NC}          ${BOLD}MachPay CLI installed successfully!${NC}            ${GREEN}║${NC}\n"
    printf "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"
    printf "\n"

    if [ "$ADD_TO_PATH" = "1" ]; then
        printf "${YELLOW}Note:${NC} Add this to your shell profile (~/.bashrc, ~/.zshrc):\n"
        printf "\n"
        printf "  ${CYAN}export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}\n"
        printf "\n"
    fi

    printf "Get started:\n"
    printf "  ${CYAN}machpay login${NC}     # Link your account\n"
    printf "  ${CYAN}machpay setup${NC}     # Configure your node\n"
    printf "  ${CYAN}machpay serve${NC}     # Start vendor gateway\n"
    printf "\n"
    printf "Documentation: ${BLUE}https://docs.machpay.xyz/cli${NC}\n"
    printf "\n"
}

# Main
main() {
    print_banner
    check_requirements
    detect_platform
    determine_install_dir
    get_version
    download_and_install
    verify_installation
    print_success
}

main "$@"

