# Installation Guide

MachPay CLI can be installed using several methods. Choose the one that best fits your workflow.

## Quick Install

### Homebrew (macOS/Linux) - Recommended

```bash
brew tap machpay/tap
brew install machpay
```

### Shell Script (macOS/Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/machpay-xyz/machpay-cli/main/scripts/install.sh | sh
```

### PowerShell (Windows)

```powershell
iwr https://raw.githubusercontent.com/machpay-xyz/machpay-cli/main/scripts/install.ps1 | iex
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

---

## Manual Download

Download the appropriate binary for your platform from the [releases page](https://github.com/machpay-xyz/machpay-cli/releases/latest).

| Platform | Architecture | Download |
|----------|--------------|----------|
| macOS | Apple Silicon (M1/M2/M3) | `machpay_darwin_arm64.tar.gz` |
| macOS | Intel | `machpay_darwin_amd64.tar.gz` |
| Linux | x64 | `machpay_linux_amd64.tar.gz` |
| Linux | ARM64 | `machpay_linux_arm64.tar.gz` |
| Windows | x64 | `machpay_windows_amd64.zip` |

### Extract and Install

**macOS/Linux:**
```bash
# Download (replace URL with your platform)
curl -LO https://github.com/machpay-xyz/machpay-cli/releases/latest/download/machpay_darwin_arm64.tar.gz

# Extract
tar -xzf machpay_darwin_arm64.tar.gz

# Move to PATH
sudo mv machpay /usr/local/bin/

# Verify
machpay version
```

**Windows:**
```powershell
# Extract the zip file
# Move machpay.exe to a directory in your PATH
# Or add the directory to your PATH environment variable
```

---

## Verification

After installation, verify it's working:

```bash
machpay version
machpay --help
```

Expected output:
```
machpay version 0.1.0
```

---

## Upgrading

### Homebrew

```bash
brew upgrade machpay
```

### Shell Script

Re-run the install script:
```bash
curl -fsSL https://raw.githubusercontent.com/machpay-xyz/machpay-cli/main/scripts/install.sh | sh
```

### Docker

```bash
docker pull ghcr.io/machpay-xyz/cli:latest
```

### Go

```bash
go install github.com/machpay-xyz/machpay-cli/cmd/machpay@latest
```

---

## Uninstalling

### Homebrew

```bash
brew uninstall machpay
brew untap machpay/tap
```

### Manual

```bash
# macOS/Linux
sudo rm /usr/local/bin/machpay
rm -rf ~/.machpay

# Windows
# Delete machpay.exe and the ~/.machpay directory
```

### Docker

```bash
docker rmi ghcr.io/machpay-xyz/cli:latest
```

---

## Configuration

After installation, MachPay stores its configuration in:

| Platform | Location |
|----------|----------|
| macOS/Linux | `~/.machpay/config.yaml` |
| Windows | `%USERPROFILE%\.machpay\config.yaml` |

---

## Troubleshooting

### "command not found" after installation

Make sure the installation directory is in your PATH:

```bash
# Check if machpay is in PATH
which machpay

# If not, add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$PATH:/usr/local/bin"
```

### Permission denied on macOS

If you see a security warning on macOS:

1. Go to **System Preferences** → **Security & Privacy** → **General**
2. Click **"Allow Anyway"** next to the machpay message
3. Run `machpay version` again and click **"Open"**

### Docker permission issues

On Linux, you may need to run Docker commands with `sudo` or add your user to the `docker` group:

```bash
sudo usermod -aG docker $USER
# Log out and back in for changes to take effect
```

---

## Getting Help

- **Documentation:** https://docs.machpay.xyz/cli
- **GitHub Issues:** https://github.com/machpay-xyz/machpay-cli/issues
- **Discord:** https://discord.gg/machpay

---

## Quick Start

After installation:

```bash
# Authenticate
machpay login

# Configure your node
machpay setup

# Check status
machpay status

# Start vendor gateway (if vendor role)
machpay serve
```

