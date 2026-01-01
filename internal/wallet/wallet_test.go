package wallet

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerate(t *testing.T) {
	kp, err := Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(kp.PublicKey) != 32 {
		t.Errorf("PublicKey length = %d, want 32", len(kp.PublicKey))
	}

	if len(kp.PrivateKey) != 64 {
		t.Errorf("PrivateKey length = %d, want 64", len(kp.PrivateKey))
	}
}

func TestPublicKeyBase58(t *testing.T) {
	kp, err := Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	addr := kp.PublicKeyBase58()
	if len(addr) < 32 || len(addr) > 44 {
		t.Errorf("Address length = %d, expected 32-44 chars", len(addr))
	}

	// Should only contain base58 characters
	for _, c := range addr {
		if !isBase58Char(c) {
			t.Errorf("Address contains invalid character: %c", c)
		}
	}
}

func isBase58Char(c rune) bool {
	// Base58 alphabet excludes 0, O, I, l
	return (c >= '1' && c <= '9') ||
		(c >= 'A' && c <= 'H') ||
		(c >= 'J' && c <= 'N') ||
		(c >= 'P' && c <= 'Z') ||
		(c >= 'a' && c <= 'k') ||
		(c >= 'm' && c <= 'z')
}

func TestSaveAndLoad(t *testing.T) {
	// Generate keypair
	kp, err := Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Create temp file
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test-wallet.json")

	// Save
	if err := kp.SaveToFile(path); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	// Verify file permissions
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("File permissions = %o, want 0600", info.Mode().Perm())
	}

	// Load
	loaded, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	// Compare
	if kp.PublicKeyBase58() != loaded.PublicKeyBase58() {
		t.Errorf("Public keys don't match")
	}
}

func TestLoadFromFile_NotFound(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/path/wallet.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestLoadFromFile_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "invalid.json")

	if err := os.WriteFile(path, []byte("not json"), 0600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	_, err := LoadFromFile(path)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestLoadFromFile_WrongLength(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "short.json")

	// Write a valid JSON array but wrong length
	if err := os.WriteFile(path, []byte("[1,2,3]"), 0600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	_, err := LoadFromFile(path)
	if err == nil {
		t.Error("Expected error for wrong length")
	}
}

func TestSignAndVerify(t *testing.T) {
	kp, err := Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	message := []byte("hello world")
	signature := kp.Sign(message)

	if len(signature) != 64 {
		t.Errorf("Signature length = %d, want 64", len(signature))
	}

	if !kp.Verify(message, signature) {
		t.Error("Signature verification failed")
	}

	// Tampered message should fail
	tamperedMessage := []byte("hello world!")
	if kp.Verify(tamperedMessage, signature) {
		t.Error("Tampered message should not verify")
	}
}

func TestBase58Encode(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte{}, ""},
		{[]byte{0}, "1"},
		{[]byte{0, 0}, "11"},
		{[]byte{1}, "2"},
		{[]byte{255}, "5Q"},
	}

	for _, tt := range tests {
		result := Base58Encode(tt.input)
		if result != tt.expected {
			t.Errorf("Base58Encode(%v) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestBase58Decode(t *testing.T) {
	// Encode then decode should give same result
	original := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	encoded := Base58Encode(original)
	decoded, err := Base58Decode(encoded)
	if err != nil {
		t.Fatalf("Base58Decode failed: %v", err)
	}

	if len(decoded) != len(original) {
		t.Fatalf("Length mismatch: got %d, want %d", len(decoded), len(original))
	}

	for i := range original {
		if decoded[i] != original[i] {
			t.Errorf("Byte %d: got %d, want %d", i, decoded[i], original[i])
		}
	}
}

func TestBase58Decode_InvalidChar(t *testing.T) {
	_, err := Base58Decode("0OIl") // Invalid base58 characters
	if err == nil {
		t.Error("Expected error for invalid characters")
	}
}

