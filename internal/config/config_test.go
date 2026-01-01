// ============================================================
// Config Package Tests
// ============================================================

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	// Create temp directory for tests
	tmpDir, err := os.MkdirTemp("", "machpay-config-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name       string
		customPath string
		setup      func() error
		wantErr    bool
	}{
		{
			name:       "new config with custom path",
			customPath: filepath.Join(tmpDir, "custom", "config.yaml"),
			wantErr:    false,
		},
		{
			name:       "existing config",
			customPath: filepath.Join(tmpDir, "existing", "config.yaml"),
			setup: func() error {
				dir := filepath.Join(tmpDir, "existing")
				if err := os.MkdirAll(dir, 0700); err != nil {
					return err
				}
				content := "version: \"1.0\"\nrole: agent\nnetwork: devnet\n"
				return os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(content), 0600)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			cfg = nil
			configDir = ""
			configPath = ""

			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			err := Init(tt.customPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && cfg == nil {
				t.Error("Init() should set cfg")
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name      string
		setup     func()
		wantRole  string
		wantEmpty bool
	}{
		{
			name: "returns existing config",
			setup: func() {
				cfg = &Config{
					Role:    "vendor",
					Network: "mainnet",
				}
			},
			wantRole:  "vendor",
			wantEmpty: false,
		},
		{
			name: "returns default if nil",
			setup: func() {
				cfg = nil
			},
			wantRole:  "",
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := Get()
			if got == nil {
				t.Fatal("Get() returned nil")
			}
			if got.Role != tt.wantRole {
				t.Errorf("Get().Role = %v, want %v", got.Role, tt.wantRole)
			}
		})
	}
}

func TestSave(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "machpay-config-save")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "saves config successfully",
			setup: func() {
				configDir = tmpDir
				configPath = filepath.Join(tmpDir, "config.yaml")
				cfg = &Config{
					Version: "1.0",
					Role:    "agent",
					Network: "devnet",
					Auth: AuthConfig{
						AccessToken: "test-token",
						Email:       "test@example.com",
					},
				}
			},
			wantErr: false,
		},
		{
			name: "fails when cfg is nil",
			setup: func() {
				cfg = nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := Save()
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify file was created
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					t.Error("Save() did not create config file")
				}
			}
		})
	}
}

func TestGetDir(t *testing.T) {
	expected := "/test/dir"
	configDir = expected
	if got := GetDir(); got != expected {
		t.Errorf("GetDir() = %v, want %v", got, expected)
	}
}

func TestGetPath(t *testing.T) {
	expected := "/test/path/config.yaml"
	configPath = expected
	if got := GetPath(); got != expected {
		t.Errorf("GetPath() = %v, want %v", got, expected)
	}
}

func TestGetConsoleURL(t *testing.T) {
	tests := []struct {
		name    string
		network string
		want    string
	}{
		{
			name:    "mainnet",
			network: "mainnet",
			want:    "https://console.machpay.xyz",
		},
		{
			name:    "devnet",
			network: "devnet",
			want:    "https://console-dev.machpay.xyz",
		},
		{
			name:    "empty defaults to devnet",
			network: "",
			want:    "https://console-dev.machpay.xyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg = &Config{Network: tt.network}
			if got := GetConsoleURL(); got != tt.want {
				t.Errorf("GetConsoleURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClear(t *testing.T) {
	cfg = &Config{
		Auth: AuthConfig{
			AccessToken:  "test-token",
			RefreshToken: "refresh-token",
			Email:        "test@example.com",
		},
	}

	Clear()

	if cfg.Auth.AccessToken != "" {
		t.Error("Clear() did not clear AccessToken")
	}
	if cfg.Auth.RefreshToken != "" {
		t.Error("Clear() did not clear RefreshToken")
	}
	if cfg.Auth.Email != "" {
		t.Error("Clear() did not clear Email")
	}
}

func TestConfigStruct(t *testing.T) {
	// Test default values
	c := &Config{
		Version: "1.0",
		Role:    "agent",
		Network: "devnet",
		Auth: AuthConfig{
			AccessToken: "token",
			UserID:      "user123",
			Email:       "test@example.com",
		},
		Wallet: WalletConfig{
			KeypairPath: "~/.machpay/wallet.json",
			PublicKey:   "abc123",
		},
		Vendor: VendorConfig{
			UpstreamURL:     "http://localhost:8080",
			PricePerRequest: 0.001,
			AllowedOrigins:  []string{"*"},
		},
		Gateway: GatewayConfig{
			Port:       8402,
			BinaryPath: "/usr/local/bin/machpay-gateway",
			Version:    "v1.0.0",
		},
	}

	// Verify struct fields
	if c.Version != "1.0" {
		t.Errorf("Version = %v, want 1.0", c.Version)
	}
	if c.Role != "agent" {
		t.Errorf("Role = %v, want agent", c.Role)
	}
	if c.Network != "devnet" {
		t.Errorf("Network = %v, want devnet", c.Network)
	}
	if c.Vendor.PricePerRequest != 0.001 {
		t.Errorf("PricePerRequest = %v, want 0.001", c.Vendor.PricePerRequest)
	}
	if c.Gateway.Port != 8402 {
		t.Errorf("Gateway.Port = %v, want 8402", c.Gateway.Port)
	}
}

