# MachPay CLI

The unified command-line interface for the MachPay network.

## Overview

MachPay CLI is a lightweight orchestrator that provides:

- **For Agents**: A wallet and payment tool for AI services
- **For Vendors**: An orchestrator that downloads and runs the payment gateway
- **For Everyone**: The onboarding wizard and quick commands

## Installation

### Homebrew (macOS/Linux)

```bash
brew install machpay/tap/machpay
```

### Direct Download

```bash
curl -fsSL https://machpay.xyz/install.sh | sh
```

### From Source

```bash
git clone https://github.com/machpay/machpay-cli.git
cd machpay-cli
go build -o machpay ./cmd/machpay
```

## Quick Start

```bash
# Authenticate with MachPay
machpay login

# Check your status
machpay status

# Interactive setup wizard
machpay setup

# Start vendor gateway (vendors only)
machpay serve
```

## Commands

| Command | Description |
|---------|-------------|
| `machpay login` | Authenticate via browser |
| `machpay logout` | Clear stored credentials |
| `machpay status` | Show current status |
| `machpay setup` | Interactive setup wizard |
| `machpay serve` | Start vendor gateway |
| `machpay open` | Launch web console |
| `machpay version` | Show version info |

## Configuration

Configuration is stored in `~/.machpay/config.yaml`:

```yaml
version: "1.0"
role: "agent"  # or "vendor"
network: "devnet"

auth:
  access_token: "eyJ..."
  user_id: "user_123"
  email: "you@example.com"

wallet:
  keypair_path: "~/.machpay/wallet.json"
  public_key: "7xK9..."
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--config` | Custom config file path |
| `--debug` | Enable debug output |
| `--no-color` | Disable colored output |

## Development

```bash
# Run in development
go run ./cmd/machpay

# Build
go build -o machpay ./cmd/machpay

# Run tests
go test ./...
```

## License

MIT License - see [LICENSE](LICENSE) for details.

