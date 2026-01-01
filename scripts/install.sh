#!/bin/bash
# ============================================================
# MachPay CLI Installer
# ============================================================
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/machpay-xyz/machpay-cli/main/scripts/install.sh | sh
#
# Environment Variables:
#   MACHPAY_VERSION    - Version to install (default: latest)
#   MACHPAY_INSTALL_DIR - Installation directory (default: /usr/local/bin)
#
# ============================================================

set -euo pipefail

# Configuration
VERSION="${MACHPAY_VERSION:-latest}"
INSTALL_DIR="${MACHPAY_INSTALL_DIR:-/usr/local/bin}"
GITHUB_REPO="machpay-xyz/machpay-cli"
BINARY_NAME="machpay"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Logging functions
info() { echo -e "${GREEN}▸${NC} $1"; }
warn() { echo -e "${YELLOW}▸${NC} $1"; }
error() { echo -e "${RED}✗${NC} $1" >&2; exit 1; }
success() { echo -e "${GREEN}✓${NC} $1"; }

# Print banner
print_banner() {
    echo -e "${BLUE}"
    echo "  __  __            _     ____             "
    echo " |  \/  | __ _  ___| |__ |  _ \ __ _ _   _ "
    echo " | |\/| |/ _\` |/ __| '_ \| |_) / _\` | | | |"
    echo " | |  | | (_| | (__| | | |  __/ (_| | |_| |"
    echo " |_|  |_|\__,_|\___|_| |_|_|   \__,_|\__, |"
    echo "                                    |___/ "
    echo -e "${NC}"
    echo -e "${BOLD}MachPay CLI Installer${NC}"
    echo ""
}

# Detect operating system
detect_os() {
    local os
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    
    case "$os" in
        darwin) echo "darwin" ;;
        linux) echo "linux" ;;
        mingw*|msys*|cygwin*) error "Windows detected. Please use install.ps1 instead." ;;
        *) error "Unsupported operating system: $os" ;;
    esac
}

# Detect CPU architecture
detect_arch() {
    local arch
    arch=$(uname -m)
    
    case "$arch" in
        x86_64|amd64) echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        armv7l) error "ARM 32-bit is not supported. Please use ARM64." ;;
        *) error "Unsupported architecture: $arch" ;;
    esac
}

# Get latest version from GitHub API
get_latest_version() {
    local version
    version=$(curl -fsSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" 2>/dev/null | 
        grep '"tag_name"' | 
        sed -E 's/.*"v([^"]+)".*/\1/' || true)
    
    if [ -z "$version" ]; then
        error "Failed to fetch latest version. Check your internet connection."
    fi
    
    echo "$version"
}

# Download file with progress
download_file() {
    local url="$1"
    local output="$2"
    
    if command -v curl &>/dev/null; then
        curl -fsSL --progress-bar "$url" -o "$output"
    elif command -v wget &>/dev/null; then
        wget -q --show-progress "$url" -O "$output"
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
}

# Verify checksum
verify_checksum() {
    local file="$1"
    local expected="$2"
    
    local actual
    if command -v sha256sum &>/dev/null; then
        actual=$(sha256sum "$file" | cut -d' ' -f1)
    elif command -v shasum &>/dev/null; then
        actual=$(shasum -a 256 "$file" | cut -d' ' -f1)
    else
        warn "Cannot verify checksum (no sha256sum or shasum found)"
        return 0
    fi
    
    if [ "$actual" != "$expected" ]; then
        error "Checksum verification failed!\nExpected: $expected\nActual:   $actual"
    fi
    
    success "Checksum verified"
}

# Main installation function
main() {
    print_banner
    
    # Detect platform
    local os arch platform
    os=$(detect_os)
    arch=$(detect_arch)
    platform="${os}_${arch}"
    info "Detected platform: ${BOLD}$platform${NC}"
    
    # Get version
    if [ "$VERSION" = "latest" ]; then
        info "Fetching latest version..."
        VERSION=$(get_latest_version)
    fi
    info "Installing version: ${BOLD}v$VERSION${NC}"
    
    # Construct download URL
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/v${VERSION}/${BINARY_NAME}_${platform}.tar.gz"
    local checksums_url="https://github.com/${GITHUB_REPO}/releases/download/v${VERSION}/checksums.txt"
    
    # Create temporary directory
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap "rm -rf '$tmp_dir'" EXIT
    
    # Download binary
    info "Downloading from GitHub..."
    download_file "$download_url" "$tmp_dir/${BINARY_NAME}.tar.gz"
    success "Download complete"
    
    # Download and verify checksum
    info "Verifying checksum..."
    download_file "$checksums_url" "$tmp_dir/checksums.txt"
    local expected_checksum
    expected_checksum=$(grep "${BINARY_NAME}_${platform}.tar.gz" "$tmp_dir/checksums.txt" | cut -d' ' -f1)
    verify_checksum "$tmp_dir/${BINARY_NAME}.tar.gz" "$expected_checksum"
    
    # Extract
    info "Extracting..."
    tar -xzf "$tmp_dir/${BINARY_NAME}.tar.gz" -C "$tmp_dir"
    success "Extracted successfully"
    
    # Install
    info "Installing to $INSTALL_DIR..."
    if [ -w "$INSTALL_DIR" ]; then
        mv "$tmp_dir/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    else
        info "Requesting sudo access to install to $INSTALL_DIR..."
        sudo mv "$tmp_dir/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    fi
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    success "Installed to $INSTALL_DIR/$BINARY_NAME"
    
    # Verify installation
    echo ""
    if command -v "$BINARY_NAME" &>/dev/null; then
        info "Verifying installation..."
        "$BINARY_NAME" version
        echo ""
        echo -e "${GREEN}${BOLD}✅ MachPay CLI installed successfully!${NC}"
        echo ""
        echo "Get started with:"
        echo -e "  ${BOLD}machpay login${NC}     # Authenticate"
        echo -e "  ${BOLD}machpay setup${NC}     # Configure your node"
        echo -e "  ${BOLD}machpay status${NC}    # Check status"
        echo ""
        echo "Documentation: https://docs.machpay.xyz/cli"
    else
        warn "Installation complete but 'machpay' not found in PATH"
        echo ""
        echo "Add $INSTALL_DIR to your PATH:"
        echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
        echo ""
        echo "Or run directly:"
        echo "  $INSTALL_DIR/$BINARY_NAME version"
    fi
}

# Run main function
main "$@"
