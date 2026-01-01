// ============================================================
// Config Package - Configuration Management
// ============================================================

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config represents the CLI configuration
type Config struct {
	Version string `yaml:"version"`
	Role    string `yaml:"role"` // "agent" or "vendor"
	Network string `yaml:"network"` // "mainnet" or "devnet"

	Auth    AuthConfig    `yaml:"auth"`
	Wallet  WalletConfig  `yaml:"wallet"`
	Vendor  VendorConfig  `yaml:"vendor,omitempty"`
	Gateway GatewayConfig `yaml:"gateway,omitempty"`
}

// AuthConfig stores authentication tokens
type AuthConfig struct {
	AccessToken  string `yaml:"access_token,omitempty"`
	RefreshToken string `yaml:"refresh_token,omitempty"`
	UserID       string `yaml:"user_id,omitempty"`
	Email        string `yaml:"email,omitempty"`
}

// WalletConfig stores wallet/keypair info
type WalletConfig struct {
	KeypairPath string `yaml:"keypair_path,omitempty"`
	PublicKey   string `yaml:"public_key,omitempty"`
}

// VendorConfig stores vendor-specific settings
type VendorConfig struct {
	UpstreamURL     string   `yaml:"upstream_url,omitempty"`
	PricePerRequest float64  `yaml:"price_per_request,omitempty"`
	AllowedOrigins  []string `yaml:"allowed_origins,omitempty"`
}

// GatewayConfig stores gateway settings
type GatewayConfig struct {
	Port       int    `yaml:"port,omitempty"`
	BinaryPath string `yaml:"binary_path,omitempty"`
	Version    string `yaml:"version,omitempty"`
}

var (
	configDir  string
	configPath string
	cfg        *Config
)

// Init initializes the configuration
func Init(customPath string) error {
	// Determine config directory
	if customPath != "" {
		configPath = customPath
		configDir = filepath.Dir(customPath)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get home dir: %w", err)
		}
		configDir = filepath.Join(home, ".machpay")
		configPath = filepath.Join(configDir, "config.yaml")
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	// Set up viper
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Load config if it exists
	if _, err := os.Stat(configPath); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("read config: %w", err)
		}
	}

	// Initialize default config
	cfg = &Config{
		Version: "1.0",
		Network: "devnet",
	}

	// Unmarshal into struct
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	return nil
}

// Get returns the current configuration
func Get() *Config {
	if cfg == nil {
		cfg = &Config{Version: "1.0", Network: "devnet"}
	}
	return cfg
}

// Save writes the configuration to disk
func Save() error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	// Create config directory if needed
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	// Write file with secure permissions
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// GetDir returns the config directory path
func GetDir() string {
	return configDir
}

// GetPath returns the config file path
func GetPath() string {
	return configPath
}

// GetConsoleURL returns the MachPay console URL based on network
func GetConsoleURL() string {
	if cfg != nil && cfg.Network == "mainnet" {
		return "https://console.machpay.xyz"
	}
	// Default to devnet console
	return "https://console-dev.machpay.xyz"
}

// Clear removes all auth credentials
func Clear() {
	if cfg != nil {
		cfg.Auth = AuthConfig{}
	}
}

