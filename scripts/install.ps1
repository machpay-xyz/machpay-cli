# ============================================================
# MachPay CLI Installer for Windows
# ============================================================
#
# Usage:
#   iwr machpay.xyz/install.ps1 | iex
#
# Or:
#   iwr -useb https://machpay.xyz/install.ps1 | iex
#
# Environment variables:
#   $env:MACHPAY_INSTALL_DIR - Installation directory
#   $env:MACHPAY_VERSION     - Specific version to install
#
# ============================================================

$ErrorActionPreference = "Stop"

# Configuration
$GitHubRepo = "machpay/machpay-cli"
$BinaryName = "machpay.exe"
$DefaultInstallDir = "$env:LOCALAPPDATA\MachPay"

function Write-Banner {
    Write-Host ""
    Write-Host "███╗   ███╗ █████╗  ██████╗██╗  ██╗██████╗  █████╗ ██╗" -ForegroundColor Blue
    Write-Host "████╗ ████║██╔══██╗██╔════╝██║  ██║██╔══██╗██╔══██╗╚██╗" -ForegroundColor Blue
    Write-Host "██╔████╔██║███████║██║     ███████║██████╔╝███████║ ██║" -ForegroundColor Blue
    Write-Host "██║╚██╔╝██║██╔══██║██║     ██╔══██║██╔═══╝ ██╔══██║ ██║" -ForegroundColor Blue
    Write-Host "██║ ╚═╝ ██║██║  ██║╚██████╗██║  ██║██║     ██║  ██║██╔╝" -ForegroundColor Blue
    Write-Host "╚═╝     ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝     ╚═╝  ╚═╝╚═╝" -ForegroundColor Blue
    Write-Host ""
    Write-Host "CLI Installer for Windows" -ForegroundColor White
    Write-Host ""
}

function Write-Info {
    param([string]$Message)
    Write-Host "→ " -ForegroundColor Cyan -NoNewline
    Write-Host $Message
}

function Write-Success {
    param([string]$Message)
    Write-Host "✓ " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Warn {
    param([string]$Message)
    Write-Host "! " -ForegroundColor Yellow -NoNewline
    Write-Host $Message
}

function Get-Architecture {
    if ([Environment]::Is64BitOperatingSystem) {
        return "amd64"
    } else {
        throw "32-bit Windows is not supported"
    }
}

function Get-LatestVersion {
    Write-Info "Fetching latest version..."
    
    try {
        $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$GitHubRepo/releases/latest" -Headers @{
            "User-Agent" = "MachPay-Installer"
        }
        return $release.tag_name
    } catch {
        throw "Could not fetch latest version: $_"
    }
}

function Get-InstallDir {
    if ($env:MACHPAY_INSTALL_DIR) {
        return $env:MACHPAY_INSTALL_DIR
    }
    return $DefaultInstallDir
}

function Install-MachPay {
    Write-Banner
    
    # Get version
    $version = if ($env:MACHPAY_VERSION) {
        $env:MACHPAY_VERSION
    } else {
        Get-LatestVersion
    }
    Write-Info "Version: $version"
    
    # Get architecture
    $arch = Get-Architecture
    Write-Info "Architecture: windows/$arch"
    
    # Get install directory
    $installDir = Get-InstallDir
    Write-Info "Install directory: $installDir"
    
    # Build download URL
    $zipName = "machpay_windows_$arch.zip"
    $downloadUrl = "https://github.com/$GitHubRepo/releases/download/$version/$zipName"
    $checksumUrl = "https://github.com/$GitHubRepo/releases/download/$version/checksums.txt"
    
    # Create install directory
    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    }
    
    # Create temp directory
    $tempDir = Join-Path $env:TEMP "machpay-install-$(Get-Random)"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    
    try {
        # Download
        Write-Info "Downloading $zipName..."
        $zipPath = Join-Path $tempDir $zipName
        
        try {
            Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -UseBasicParsing
        } catch {
            throw "Download failed: $_"
        }
        
        # Verify checksum
        Write-Info "Verifying checksum..."
        try {
            $checksums = (Invoke-WebRequest -Uri $checksumUrl -UseBasicParsing).Content
            $expectedHash = ($checksums -split "`n" | Where-Object { $_ -match $zipName } | Select-Object -First 1) -replace "\s.*", ""
            $actualHash = (Get-FileHash -Path $zipPath -Algorithm SHA256).Hash.ToLower()
            
            if ($expectedHash -and $expectedHash -ne $actualHash) {
                throw "Checksum verification failed!`n  Expected: $expectedHash`n  Actual:   $actualHash"
            }
            Write-Success "Checksum verified"
        } catch {
            Write-Warn "Could not verify checksum: $_"
        }
        
        # Extract
        Write-Info "Extracting..."
        Expand-Archive -Path $zipPath -DestinationPath $tempDir -Force
        
        # Find binary
        $binaryPath = Get-ChildItem -Path $tempDir -Filter $BinaryName -Recurse | Select-Object -First 1
        if (-not $binaryPath) {
            throw "Binary not found in archive"
        }
        
        # Install
        Write-Info "Installing..."
        Copy-Item -Path $binaryPath.FullName -Destination (Join-Path $installDir $BinaryName) -Force
        
        # Add to PATH
        $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
        if ($userPath -notlike "*$installDir*") {
            Write-Info "Adding to PATH..."
            [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
            $env:Path = "$env:Path;$installDir"
        }
        
        Write-Success "Installed to $installDir\$BinaryName"
        
        # Verify
        Write-Info "Verifying installation..."
        $installedVersion = & (Join-Path $installDir $BinaryName) version 2>&1
        Write-Success "Installation verified: $installedVersion"
        
        # Success message
        Write-Host ""
        Write-Host "╔════════════════════════════════════════════════════════════╗" -ForegroundColor Green
        Write-Host "║          MachPay CLI installed successfully!               ║" -ForegroundColor Green
        Write-Host "╚════════════════════════════════════════════════════════════╝" -ForegroundColor Green
        Write-Host ""
        Write-Host "Restart your terminal, then run:" -ForegroundColor White
        Write-Host ""
        Write-Host "  machpay login     " -ForegroundColor Cyan -NoNewline
        Write-Host "# Link your account"
        Write-Host "  machpay setup     " -ForegroundColor Cyan -NoNewline
        Write-Host "# Configure your node"
        Write-Host "  machpay serve     " -ForegroundColor Cyan -NoNewline
        Write-Host "# Start vendor gateway"
        Write-Host ""
        Write-Host "Documentation: " -NoNewline
        Write-Host "https://docs.machpay.xyz/cli" -ForegroundColor Blue
        Write-Host ""
        
    } finally {
        # Cleanup
        Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}

# Run installer
Install-MachPay

