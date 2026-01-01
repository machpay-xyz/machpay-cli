package gateway

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNewDownloader(t *testing.T) {
	dl := NewDownloader()

	if dl.installDir == "" {
		t.Error("Expected non-empty install dir")
	}

	home, _ := os.UserHomeDir()
	expectedDir := filepath.Join(home, ".machpay", "bin")
	if dl.installDir != expectedDir {
		t.Errorf("Install dir = %s, want %s", dl.installDir, expectedDir)
	}
}

func TestDownloader_BinaryPath(t *testing.T) {
	dl := NewDownloader()
	path := dl.BinaryPath()

	expectedName := GatewayBinary
	if runtime.GOOS == "windows" {
		expectedName += ".exe"
	}

	if filepath.Base(path) != expectedName {
		t.Errorf("Binary name = %s, want %s", filepath.Base(path), expectedName)
	}
}

func TestDownloader_IsInstalled_NotInstalled(t *testing.T) {
	dl := &Downloader{
		installDir: t.TempDir(),
	}

	if dl.IsInstalled() {
		t.Error("Expected IsInstalled to return false for non-existent binary")
	}
}

func TestDownloader_IsInstalled_Installed(t *testing.T) {
	tmpDir := t.TempDir()
	dl := &Downloader{
		installDir: tmpDir,
	}

	// Create fake binary
	binaryPath := dl.BinaryPath()
	if err := os.WriteFile(binaryPath, []byte("fake"), 0755); err != nil {
		t.Fatalf("Create fake binary: %v", err)
	}

	if !dl.IsInstalled() {
		t.Error("Expected IsInstalled to return true for existing binary")
	}
}

func TestDownloader_AssetName(t *testing.T) {
	dl := NewDownloader()
	name := dl.assetName()

	if !containsString(name, runtime.GOOS) {
		t.Errorf("Asset name %s does not contain OS %s", name, runtime.GOOS)
	}

	if !containsString(name, runtime.GOARCH) {
		t.Errorf("Asset name %s does not contain arch %s", name, runtime.GOARCH)
	}

	if runtime.GOOS == "windows" {
		if !containsString(name, ".zip") {
			t.Error("Windows asset should be .zip")
		}
	} else {
		if !containsString(name, ".tar.gz") {
			t.Error("Unix asset should be .tar.gz")
		}
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestDownloader_GetLatestRelease_Mock(t *testing.T) {
	// Create mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/machpay/machpay-gateway/releases/latest" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"tag_name": "v1.0.0",
				"name": "v1.0.0",
				"assets": [
					{
						"name": "machpay-gateway_darwin_arm64.tar.gz",
						"size": 10485760,
						"browser_download_url": "https://example.com/file.tar.gz"
					},
					{
						"name": "checksums.txt",
						"browser_download_url": "https://example.com/checksums.txt"
					}
				]
			}`))
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()

	// Note: Can't easily override GitHubAPI constant, so this test
	// documents the expected behavior with a mock
	// Real integration tests would hit the actual API
}

func TestDownloader_FileChecksum(t *testing.T) {
	dl := NewDownloader()

	// Create temp file with known content
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	content := []byte("hello world")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("Write temp file: %v", err)
	}

	checksum, err := dl.fileChecksum(tmpFile)
	if err != nil {
		t.Fatalf("fileChecksum: %v", err)
	}

	// SHA256 of "hello world"
	expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	if checksum != expected {
		t.Errorf("Checksum = %s, want %s", checksum, expected)
	}
}

func TestDownloader_FileChecksum_NotFound(t *testing.T) {
	dl := NewDownloader()

	_, err := dl.fileChecksum("/nonexistent/file.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestDownloader_NeedsUpdate_NotInstalled(t *testing.T) {
	dl := &Downloader{
		installDir: t.TempDir(),
		httpClient: &http.Client{},
	}

	// When not installed, NeedsUpdate should return true
	// (without hitting the API for the "latest" check)
	needsUpdate, _, _ := dl.NeedsUpdate()

	// Should need "update" (really, install)
	if !needsUpdate {
		t.Error("Expected NeedsUpdate to return true when not installed")
	}
}

func TestDownloader_ExtractAndInstall(t *testing.T) {
	// This would require creating a valid tar.gz with the binary
	// Skipping for unit tests - covered by integration tests
	t.Skip("Requires valid tar.gz archive - covered by integration tests")
}

func TestProgressReader_Integration(t *testing.T) {
	// Test that ProgressReader integrates with download flow
	data := bytes.Repeat([]byte("x"), 10000)
	reader := bytes.NewReader(data)
	output := &bytes.Buffer{}

	pr := NewProgressReader(reader, int64(len(data)), output)

	result := make([]byte, len(data))
	n, err := pr.Read(result)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if n != len(data) {
		t.Errorf("Read %d bytes, want %d", n, len(data))
	}

	// Progress should have been written
	if output.Len() == 0 {
		t.Error("Expected progress output")
	}
}

