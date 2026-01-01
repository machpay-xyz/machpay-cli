# Changelog

All notable changes to MachPay CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-01-01

### 🎉 Initial Public Release

First public release of MachPay CLI - the command-line interface for the MachPay AI payment network.

### Added

#### Authentication & Configuration
- Browser-based OAuth authentication with `machpay login`
- Secure token storage in `~/.machpay/config.yaml`
- Headless mode for CI/CD environments (`--no-browser` flag)
- Logout command to clear credentials: `machpay logout`
- Interactive setup wizard with `machpay setup`
- Role selection (Agent or Vendor)
- Network selection (Devnet or Mainnet)
- Non-interactive mode for automation via environment variables

#### Wallet Management
- Ed25519/Solana keypair generation
- Import existing wallet from keypair file
- Secure wallet storage with proper file permissions
- Public key display in Base58 format

#### Gateway Management (Vendors)
- Automatic gateway binary download from GitHub releases
- Process management: `machpay serve`, `machpay stop`, `machpay restart`
- Background mode with `--detach` flag
- Foreground mode for debugging
- Health check monitoring
- Log viewing with follow mode: `machpay logs -f`
- Gateway update command: `machpay update gateway`
- Configurable upstream URL and port

#### Status & Monitoring
- Comprehensive status display: `machpay status`
- JSON output for scripting: `machpay status --json`
- Watch mode for continuous monitoring: `machpay status --watch`
- Gateway status (running/stopped, PID, port)
- Authentication status
- Wallet information display

#### Desktop Integration
- Launch desktop app from CLI: `machpay open`
- Deep linking to specific routes: `machpay open marketplace`
- Supported routes: marketplace, funding, settings, invoices, disputes, etc.
- Automatic browser fallback when desktop app not installed
- Force browser mode: `machpay open --web`

#### Developer Experience
- Beautiful TUI with Lipgloss styling
- Colored output with status icons
- Progress indicators for downloads
- Helpful error messages with suggestions
- Comprehensive `--help` for all commands
- Version information: `machpay version`

#### Cross-Platform Support
- macOS (Intel and Apple Silicon)
- Linux (x64 and ARM64)
- Windows (x64)

### Known Issues

- **Google OAuth in Desktop App**: Authentication via Google in the desktop app may not complete automatically due to cross-origin security restrictions. **Workaround**: Use email OTP or wallet login in the desktop app, or use the web console at https://console.machpay.xyz for Google OAuth. See [KNOWN_ISSUES.md](./KNOWN_ISSUES.md) for technical details.

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
go install github.com/machpay-xyz/machpay-cli/cmd/machpay@v0.1.0
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

- [README](./README.md) - Full documentation and usage guide
- [CLI Status](./CLI_STATUS.md) - Feature implementation status
- [Known Issues](./KNOWN_ISSUES.md) - Known issues and workarounds
- [Contributing](./CONTRIBUTING.md) - How to contribute

---

[Unreleased]: https://github.com/machpay-xyz/machpay-cli/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/machpay-xyz/machpay-cli/releases/tag/v0.1.0
