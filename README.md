# MachPay CLI

[![Build](https://github.com/machpay/machpay-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/machpay/machpay-cli/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/machpay/machpay-cli)](https://github.com/machpay/machpay-cli/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/machpay/machpay-cli)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

The unified command-line interface for the **MachPay AI Payment Network**.

```
███╗   ███╗ █████╗  ██████╗██╗  ██╗██████╗  █████╗ ██╗
████╗ ████║██╔══██╗██╔════╝██║  ██║██╔══██╗██╔══██╗╚██╗
██╔████╔██║███████║██║     ███████║██████╔╝███████║ ██║
██║╚██╔╝██║██╔══██║██║     ██╔══██║██╔═══╝ ██╔══██║ ██║
██║ ╚═╝ ██║██║  ██║╚██████╗██║  ██║██║     ██║  ██║██╔╝
╚═╝     ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝     ╚═╝  ╚═╝╚═╝
```

## What is MachPay?

MachPay is the payment network for AI agents. It enables:

- **Agents** to pay for API services automatically
- **Vendors** to monetize their APIs with zero integration overhead
- **Everyone** to build the future of autonomous commerce

## Quick Start

### Installation

**Homebrew (macOS/Linux):**
```bash
brew install machpay/tap/machpay
```

**Curl (Linux/macOS):**
```bash
curl -fsSL https://machpay.xyz/install.sh | sh
```

**PowerShell (Windows):**
```powershell
iwr machpay.xyz/install.ps1 | iex
```

**Docker:**
```bash
docker pull ghcr.io/machpay/cli:latest
```

### First Steps

```bash
# 1. Authenticate with MachPay
machpay login

# 2. Run the interactive setup wizard
machpay setup

# 3. Check your status
machpay status
```

## Commands

| Command | Description |
|---------|-------------|
| `login` | Authenticate with MachPay |
| `logout` | Clear stored credentials |
| `setup` | Interactive setup wizard |
| `status` | Show current status |
| `serve` | Start vendor gateway |
| `stop` | Stop vendor gateway |
| `restart` | Restart vendor gateway |
| `logs` | View gateway logs |
| `open` | Launch web console |
| `update` | Update CLI and gateway |
| `version` | Show version info |

---

### `machpay login`

Link your CLI to your MachPay account via browser authentication.

```bash
# Standard login (opens browser)
machpay login

# Headless mode - prints URL instead of opening browser
machpay login --no-browser
```

---

### `machpay setup`

Interactive wizard to configure your node as an Agent or Vendor.

```bash
# Interactive mode
machpay setup

# Non-interactive mode for CI/CD
machpay setup --non-interactive
```

**Environment Variables (non-interactive mode):**

| Variable | Description | Values |
|----------|-------------|--------|
| `MACHPAY_ROLE` | Node role | `agent`, `vendor` |
| `MACHPAY_NETWORK` | Network | `mainnet`, `devnet` |
| `MACHPAY_UPSTREAM` | Upstream URL | URL (vendor only) |

---

### `machpay status`

Display current configuration and gateway status.

```bash
# Human-readable output
machpay status

# JSON output for scripting
machpay status --json

# Live updates
machpay status --watch
```

**Example Output:**
```
MachPay Status
══════════════════════════════════════════════════════

  Account:    user@example.com ✓
  Role:       Vendor
  Network:    Mainnet

  Gateway:
    Status:   ● Running (PID 12345)
    Version:  v1.2.0
    Port:     8402

  Wallet:
    Address:  7xK9...3nP
    Balance:  $125.50 USDC

══════════════════════════════════════════════════════
```

---

### `machpay serve`

Start the vendor payment gateway.

```bash
# Start in foreground
machpay serve

# Start in background (daemon mode)
machpay serve --detach

# Custom port
machpay serve --port 9000

# Custom upstream
machpay serve --upstream http://localhost:11434

# Debug mode
machpay serve --debug
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | 8402 | Gateway listen port |
| `--upstream` | from config | Upstream service URL |
| `--detach` | false | Run in background |
| `--debug` | false | Enable debug logging |

---

### `machpay logs`

View gateway logs.

```bash
# Show recent logs
machpay logs

# Follow logs in real-time
machpay logs --follow
```

---

### `machpay open`

Launch the MachPay Console.

```bash
# Open default view
machpay open

# Open specific page
machpay open marketplace
machpay open funding

# Force web browser (even if desktop app installed)
machpay open --web
```

---

### `machpay update`

Update CLI and/or gateway to the latest version.

```bash
# Update everything
machpay update

# Update CLI only
machpay update cli

# Update gateway only
machpay update gateway

# Check for updates without installing
machpay update --check
```

---

## Configuration

Configuration is stored in `~/.machpay/config.yaml`:

```yaml
# Authentication
auth:
  token: "eyJ..."

# Node configuration
role: vendor
network: mainnet

# Vendor settings
vendor:
  upstream: http://localhost:11434
  port: 8402
  price_per_request: 0.001

# Wallet
wallet:
  path: ~/.machpay/wallet.json

# Optional: Telemetry (opt-in)
telemetry:
  enabled: true
```

### Configuration Precedence

1. Command-line flags (highest)
2. Environment variables
3. Config file
4. Defaults (lowest)

---

## For Agents

AI agents use MachPay to pay for API services. Quick setup:

```bash
# 1. Login and setup
machpay login
machpay setup  # Select "Agent"

# 2. Fund your wallet
machpay open funding

# 3. Use the Python SDK
pip install machpay
```

```python
from machpay import MachPay

client = MachPay()  # Uses ~/.machpay config
response = await client.call("weather-api", "/v1/forecast", {"city": "SF"})
```

---

## For Vendors

Vendors monetize their APIs through MachPay. Quick setup:

```bash
# 1. Login and setup
machpay login
machpay setup  # Select "Vendor"

# 2. Start the gateway
machpay serve

# Your API is now available at:
# https://your-service.machpay.network
```

The gateway handles:
- Payment verification
- Request proxying
- Usage tracking
- Settlement

---

## Troubleshooting

### "command not found: machpay"

Add the install directory to your PATH:

```bash
# For ~/.local/bin installs
export PATH="$HOME/.local/bin:$PATH"

# Add to ~/.bashrc or ~/.zshrc for persistence
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
```

### "not logged in"

Run `machpay login` to authenticate with your MachPay account.

### Gateway won't start

1. **Check if port is in use:**
   ```bash
   lsof -i :8402
   ```

2. **Check logs:**
   ```bash
   machpay logs
   ```

3. **Restart:**
   ```bash
   machpay restart
   ```

4. **Check upstream is reachable:**
   ```bash
   curl http://localhost:11434/health
   ```

### "gateway not found"

The gateway binary is downloaded automatically. Force a re-download:

```bash
rm -rf ~/.machpay/bin/machpay-gateway
machpay serve  # Will download fresh
```

### Config issues

Reset configuration:

```bash
rm ~/.machpay/config.yaml
machpay setup
```

---

## Development

### Build from Source

```bash
# Clone the repository
git clone https://github.com/machpay/machpay-cli.git
cd machpay-cli

# Build
go build -o machpay ./cmd/machpay

# Run
./machpay version
```

### Run Tests

```bash
# All tests
go test -v ./...

# With coverage
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Race detection
go test -race ./...
```

### Build for All Platforms

```bash
# Using GoReleaser
goreleaser build --snapshot --clean
```

---

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

---

## Security

Found a security issue? Please report it privately to security@machpay.xyz.

Do NOT open a public issue for security vulnerabilities.

---

## License

MIT - see [LICENSE](LICENSE) for details.

---

## Links

- **Website:** [machpay.xyz](https://machpay.xyz)
- **Documentation:** [docs.machpay.xyz](https://docs.machpay.xyz)
- **Console:** [console.machpay.xyz](https://console.machpay.xyz)
- **Discord:** [discord.gg/machpay](https://discord.gg/machpay)
- **Twitter:** [@machpay](https://twitter.com/machpay)

---

Built with ❤️ by the MachPay Team
