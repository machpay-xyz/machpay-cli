# ============================================================
# MachPay CLI Installer for Windows
# ============================================================
#
# Usage:
#   iwr https://raw.githubusercontent.com/machpay-xyz/machpay-cli/main/scripts/install.ps1 | iex
#
# Or:
#   .\install.ps1 [-Version "0.1.0"] [-InstallDir "$env:USERPROFILE\.machpay\bin"]
#
# ============================================================

param(
    [string]$Version = "latest",
    [string]$InstallDir = "$env:USERPROFILE\.machpay\bin"
)

$ErrorActionPreference = "Stop"

# Configuration
$GitHubRepo = "machpay-xyz/machpay-cli"
$BinaryName = "machpay"

# Colors
function Write-Info { Write-Host "▸ $args" -ForegroundColor Cyan }
function Write-Success { Write-Host "✓ $args" -ForegroundColor Green }
function Write-Warn { Write-Host "▸ $args" -ForegroundColor Yellow }
function Write-Err { Write-Host "✗ $args" -ForegroundColor Red }

# Print banner
function Show-Banner {
    Write-Host ""
    Write-Host "  __  __            _     ____             " -ForegroundColor Blue
    Write-Host " |  \/  | __ _  ___| |__ |  _ \ __ _ _   _ " -ForegroundColor Blue
    Write-Host " | |\/| |/ _`` |/ __| '_ \| |_) / _`` | | | |" -ForegroundColor Blue
    Write-Host " | |  | | (_| | (__| | | |  __/ (_| | |_| |" -ForegroundColor Blue
    Write-Host " |_|  |_|\__,_|\___|_| |_|_|   \__,_|\__, |" -ForegroundColor Blue
    Write-Host "                                    |___/ " -ForegroundColor Blue
    Write-Host ""
    Write-Host "MachPay CLI Installer for Windows" -ForegroundColor White
    Write-Host ""
}

# Get latest version from GitHub
function Get-LatestVersion {
    try {
        $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$GitHubRepo/releases/latest"
        return $release.tag_name -replace '^v', ''
    }
    catch {
        Write-Err "Failed to fetch latest version: $_"
        exit 1
    }
}

# Main installation
function Install-MachPay {
    Show-Banner

    # Determine version
    if ($Version -eq "latest") {
        Write-Info "Fetching latest version..."
        $Version = Get-LatestVersion
    }
    Write-Info "Installing version: v$Version"

    # Windows is always amd64
    $Platform = "windows_amd64"
    $ArchiveExt = "zip"
    
    # Construct URLs
    $DownloadUrl = "https://github.com/$GitHubRepo/releases/download/v$Version/${BinaryName}_${Platform}.${ArchiveExt}"
    $ChecksumsUrl = "https://github.com/$GitHubRepo/releases/download/v$Version/checksums.txt"

    # Create temp directory
    $TempDir = New-Item -ItemType Directory -Path (Join-Path $env:TEMP "machpay-install-$(Get-Random)")
    
    try {
        # Download archive
        Write-Info "Downloading from GitHub..."
        $ArchivePath = Join-Path $TempDir "$BinaryName.$ArchiveExt"
        Invoke-WebRequest -Uri $DownloadUrl -OutFile $ArchivePath -UseBasicParsing
        Write-Success "Download complete"

        # Download checksums
        Write-Info "Verifying checksum..."
        $ChecksumsPath = Join-Path $TempDir "checksums.txt"
        Invoke-WebRequest -Uri $ChecksumsUrl -OutFile $ChecksumsPath -UseBasicParsing
        
        # Verify checksum
        $ExpectedChecksum = (Get-Content $ChecksumsPath | Where-Object { $_ -match "${BinaryName}_${Platform}.${ArchiveExt}" }) -split '\s+' | Select-Object -First 1
        $ActualChecksum = (Get-FileHash -Path $ArchivePath -Algorithm SHA256).Hash.ToLower()
        
        if ($ExpectedChecksum -ne $ActualChecksum) {
            Write-Err "Checksum verification failed!"
            Write-Host "Expected: $ExpectedChecksum"
            Write-Host "Actual:   $ActualChecksum"
            exit 1
        }
        Write-Success "Checksum verified"

        # Extract
        Write-Info "Extracting..."
        $ExtractPath = Join-Path $TempDir "extracted"
        Expand-Archive -Path $ArchivePath -DestinationPath $ExtractPath -Force
        Write-Success "Extracted successfully"

        # Create install directory
        if (-not (Test-Path $InstallDir)) {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }

        # Install binary
        Write-Info "Installing to $InstallDir..."
        $BinaryPath = Get-ChildItem -Path $ExtractPath -Filter "$BinaryName.exe" -Recurse | Select-Object -First 1
        if (-not $BinaryPath) {
            $BinaryPath = Get-ChildItem -Path $ExtractPath -Filter "$BinaryName" -Recurse | Select-Object -First 1
        }
        Copy-Item -Path $BinaryPath.FullName -Destination (Join-Path $InstallDir "$BinaryName.exe") -Force
        Write-Success "Installed to $InstallDir\$BinaryName.exe"

        # Add to PATH if not already there
        $CurrentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
        if ($CurrentPath -notlike "*$InstallDir*") {
            Write-Info "Adding $InstallDir to PATH..."
            [Environment]::SetEnvironmentVariable("PATH", "$CurrentPath;$InstallDir", "User")
            $env:PATH = "$env:PATH;$InstallDir"
            Write-Success "Added to PATH"
        }

        # Verify installation
        Write-Host ""
        Write-Info "Verifying installation..."
        & "$InstallDir\$BinaryName.exe" version
        
        Write-Host ""
        Write-Host "✅ MachPay CLI installed successfully!" -ForegroundColor Green
        Write-Host ""
        Write-Host "Get started with:"
        Write-Host "  machpay login     # Authenticate" -ForegroundColor White
        Write-Host "  machpay setup     # Configure your node" -ForegroundColor White
        Write-Host "  machpay status    # Check status" -ForegroundColor White
        Write-Host ""
        Write-Host "Documentation: https://docs.machpay.xyz/cli"
        Write-Host ""
        Write-Warn "Note: You may need to restart your terminal for PATH changes to take effect."
    }
    finally {
        # Cleanup
        Remove-Item -Path $TempDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}

# Run installation
Install-MachPay
