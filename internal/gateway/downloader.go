// ============================================================
// Gateway Downloader - Download and manage gateway binary
// ============================================================
//
// Downloads the machpay-gateway from GitHub releases.
// Supports:
// - Platform detection (darwin, linux, windows)
// - Architecture detection (amd64, arm64)
// - SHA256 checksum verification
// - Tarball extraction
//
// ============================================================

package gateway

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// GatewayRepo is the GitHub repository for the gateway
	GatewayRepo = "machpay/machpay-gateway"

	// GatewayBinary is the name of the gateway binary
	GatewayBinary = "machpay-gateway"

	// GitHubAPI is the base URL for GitHub API
	GitHubAPI = "https://api.github.com"

	// DownloadTimeout is the timeout for downloading the binary
	DownloadTimeout = 10 * time.Minute
)

// Downloader manages gateway binary downloads
type Downloader struct {
	installDir string
	httpClient *http.Client
}

// NewDownloader creates a new downloader
func NewDownloader() *Downloader {
	home, _ := os.UserHomeDir()
	return &Downloader{
		installDir: filepath.Join(home, ".machpay", "bin"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// BinaryPath returns the full path to the gateway binary
func (d *Downloader) BinaryPath() string {
	binary := GatewayBinary
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	return filepath.Join(d.installDir, binary)
}

// IsInstalled checks if the gateway is installed
func (d *Downloader) IsInstalled() bool {
	_, err := os.Stat(d.BinaryPath())
	return err == nil
}

// InstalledVersion returns the version of the installed gateway
func (d *Downloader) InstalledVersion() (string, error) {
	if !d.IsInstalled() {
		return "", fmt.Errorf("gateway not installed")
	}

	// Run: machpay-gateway --version
	cmd := exec.Command(d.BinaryPath(), "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("get version: %w", err)
	}

	// Parse "machpay-gateway v1.2.0" or "v1.2.0" format
	version := strings.TrimSpace(string(output))
	parts := strings.Fields(version)

	if len(parts) >= 2 {
		return strings.TrimPrefix(parts[1], "v"), nil
	}
	if len(parts) == 1 {
		return strings.TrimPrefix(parts[0], "v"), nil
	}

	return "", fmt.Errorf("unable to parse version from: %s", version)
}

// ============================================================
// GitHub Release Types
// ============================================================

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	PublishedAt string         `json:"published_at"`
	Assets      []GitHubAsset  `json:"assets"`
}

// GitHubAsset represents a release asset
type GitHubAsset struct {
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// ============================================================
// GitHub API Methods
// ============================================================

// GetLatestRelease fetches the latest release info from GitHub
func (d *Downloader) GetLatestRelease() (*GitHubRelease, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/latest", GitHubAPI, GatewayRepo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "machpay-cli")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no releases found for %s", GatewayRepo)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("parse release: %w", err)
	}

	return &release, nil
}

// GetRelease fetches a specific release by tag
func (d *Downloader) GetRelease(tag string) (*GitHubRelease, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/tags/%s", GitHubAPI, GatewayRepo, tag)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "machpay-cli")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("release %s not found", tag)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("parse release: %w", err)
	}

	return &release, nil
}

// ============================================================
// Download Methods
// ============================================================

// Download downloads and installs a specific version
func (d *Downloader) Download(version string, progress io.Writer) error {
	// 1. Normalize version tag
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	// 2. Construct asset name for this platform
	assetName := d.assetName()
	checksumName := "checksums.txt"

	// 3. Get release info
	release, err := d.GetRelease(version)
	if err != nil {
		return fmt.Errorf("get release: %w", err)
	}

	// 4. Find asset URLs
	var assetSize int64
	var checksumURL string
	var downloadURL string

	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			assetSize = asset.Size
		}
		if asset.Name == checksumName {
			checksumURL = asset.BrowserDownloadURL
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no asset found for %s/%s (looking for %s)",
			runtime.GOOS, runtime.GOARCH, assetName)
	}

	// 5. Fetch expected checksum
	expectedChecksum := ""
	if checksumURL != "" {
		checksum, err := d.fetchChecksum(checksumURL, assetName)
		if err != nil {
			fmt.Fprintf(progress, "  Warning: Could not fetch checksum: %v\n", err)
		} else {
			expectedChecksum = checksum
		}
	}

	// 6. Download tarball with progress
	tmpFile, err := d.downloadFile(downloadURL, assetSize, progress)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer os.Remove(tmpFile)

	// 7. Verify checksum
	if expectedChecksum != "" {
		actualChecksum, err := d.fileChecksum(tmpFile)
		if err != nil {
			return fmt.Errorf("compute checksum: %w", err)
		}
		if actualChecksum != expectedChecksum {
			return fmt.Errorf("checksum mismatch:\n  expected: %s\n  got:      %s",
				expectedChecksum, actualChecksum)
		}
		fmt.Fprintf(progress, "  ✓ Checksum verified\n")
	}

	// 8. Extract and install
	if err := d.extractAndInstall(tmpFile, progress); err != nil {
		return fmt.Errorf("install: %w", err)
	}

	return nil
}

// assetName returns the expected asset filename for this platform
func (d *Downloader) assetName() string {
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	return fmt.Sprintf("machpay-gateway_%s_%s.%s", runtime.GOOS, runtime.GOARCH, ext)
}

// downloadFile downloads a file with progress reporting
func (d *Downloader) downloadFile(url string, size int64, progress io.Writer) (string, error) {
	// Create HTTP client with longer timeout for downloads
	client := &http.Client{Timeout: DownloadTimeout}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: %s", resp.Status)
	}

	// Use Content-Length if available, otherwise use provided size
	totalSize := resp.ContentLength
	if totalSize <= 0 {
		totalSize = size
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "machpay-gateway-*.tar.gz")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Wrap reader with progress if size is known
	var reader io.Reader = resp.Body
	if totalSize > 0 && progress != nil {
		reader = NewProgressReader(resp.Body, totalSize, progress)
	}

	// Copy to temp file
	_, err = io.Copy(tmpFile, reader)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("download interrupted: %w", err)
	}

	return tmpFile.Name(), nil
}

// fetchChecksum fetches the expected checksum for an asset
func (d *Downloader) fetchChecksum(checksumURL, assetName string) (string, error) {
	resp, err := d.httpClient.Get(checksumURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetch checksums: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse checksums.txt format: "checksum  filename"
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[1] == assetName {
			return parts[0], nil
		}
	}

	return "", fmt.Errorf("checksum not found for %s", assetName)
}

// fileChecksum computes the SHA256 checksum of a file
func (d *Downloader) fileChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// extractAndInstall extracts the tarball and installs the binary
func (d *Downloader) extractAndInstall(tarPath string, progress io.Writer) error {
	// Ensure install directory exists
	if err := os.MkdirAll(d.installDir, 0755); err != nil {
		return fmt.Errorf("create install dir: %w", err)
	}

	// Open tarball
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("open gzip: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	found := false

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read tar: %w", err)
		}

		// Only extract the binary
		if header.Typeflag != tar.TypeReg {
			continue
		}

		// Match the binary name
		baseName := filepath.Base(header.Name)
		if !strings.HasPrefix(baseName, GatewayBinary) {
			continue
		}

		// Write binary
		outPath := d.BinaryPath()
		outFile, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("create binary: %w", err)
		}

		if _, err := io.Copy(outFile, tr); err != nil {
			outFile.Close()
			return fmt.Errorf("write binary: %w", err)
		}
		outFile.Close()

		// Ensure executable on Unix
		if runtime.GOOS != "windows" {
			if err := os.Chmod(outPath, 0755); err != nil {
				return fmt.Errorf("chmod: %w", err)
			}
		}

		found = true
		fmt.Fprintf(progress, "  ✓ Installed to %s\n", outPath)
		break
	}

	if !found {
		return fmt.Errorf("binary not found in archive")
	}

	return nil
}

// ============================================================
// Update Checking
// ============================================================

// NeedsUpdate checks if an update is available
func (d *Downloader) NeedsUpdate() (bool, string, error) {
	installed, err := d.InstalledVersion()
	if err != nil {
		// Not installed = needs "update" (really, install)
		return true, "", nil
	}

	release, err := d.GetLatestRelease()
	if err != nil {
		return false, "", fmt.Errorf("check for updates: %w", err)
	}

	latest := strings.TrimPrefix(release.TagName, "v")

	// Simple string comparison (works for semver)
	return installed != latest, latest, nil
}

// GetInstallDir returns the installation directory
func (d *Downloader) GetInstallDir() string {
	return d.installDir
}

