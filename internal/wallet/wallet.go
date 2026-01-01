// ============================================================
// Wallet Package - Solana Keypair Management
// ============================================================
//
// Supports:
// - Generating new Ed25519 keypairs
// - Loading keypairs from Solana CLI format (JSON array)
// - Saving keypairs in Solana CLI format
// - Base58 address encoding
//
// ============================================================

package wallet

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Keypair represents a Solana Ed25519 keypair
type Keypair struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

// Generate creates a new random keypair
func Generate() (*Keypair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate keypair: %w", err)
	}
	return &Keypair{
		PublicKey:  pub,
		PrivateKey: priv,
	}, nil
}

// LoadFromFile loads a keypair from a Solana CLI format file
// The file format is a JSON array of 64 bytes (32 private + 32 public)
func LoadFromFile(path string) (*Keypair, error) {
	// Expand home directory
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("get home dir: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	// Parse as JSON array of bytes
	var bytes []byte
	if err := json.Unmarshal(data, &bytes); err != nil {
		return nil, fmt.Errorf("parse keypair JSON: %w", err)
	}

	if len(bytes) != 64 {
		return nil, fmt.Errorf("invalid keypair length: got %d bytes, want 64", len(bytes))
	}

	// First 32 bytes are the seed (not used directly in ed25519.PrivateKey)
	// The full 64 bytes are the private key in Go's ed25519 format
	privateKey := ed25519.PrivateKey(bytes)
	publicKey := privateKey.Public().(ed25519.PublicKey)

	return &Keypair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

// SaveToFile saves the keypair in Solana CLI format
func (k *Keypair) SaveToFile(path string) error {
	// Expand home directory
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get home dir: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	// Marshal private key as JSON array of bytes
	// ed25519.PrivateKey is 64 bytes (seed + public key)
	data, err := json.Marshal([]byte(k.PrivateKey))
	if err != nil {
		return fmt.Errorf("marshal keypair: %w", err)
	}

	// Write with secure permissions (0600 = owner read/write only)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// PublicKeyBase58 returns the base58-encoded public key (Solana address)
func (k *Keypair) PublicKeyBase58() string {
	return Base58Encode(k.PublicKey)
}

// PublicKeyBytes returns the raw public key bytes
func (k *Keypair) PublicKeyBytes() []byte {
	return k.PublicKey
}

// Sign signs a message with the private key
func (k *Keypair) Sign(message []byte) []byte {
	return ed25519.Sign(k.PrivateKey, message)
}

// Verify verifies a signature against a message
func (k *Keypair) Verify(message, signature []byte) bool {
	return ed25519.Verify(k.PublicKey, message, signature)
}

