package auth

import (
	"encoding/base64"
	"testing"
)

func TestParseUserFromToken(t *testing.T) {
	// Create a mock JWT token
	// Header: {"alg":"HS256","typ":"JWT"}
	// Payload: {"sub":"user123","email":"test@example.com","name":"Test User"}
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"user123","email":"test@example.com","name":"Test User"}`))
	signature := "signature"

	token := header + "." + payload + "." + signature

	user, err := ParseUserFromToken(token)
	if err != nil {
		t.Fatalf("ParseUserFromToken failed: %v", err)
	}

	if user.ID != "user123" {
		t.Errorf("ID mismatch: got %s, want user123", user.ID)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Email mismatch: got %s, want test@example.com", user.Email)
	}
	if user.Name != "Test User" {
		t.Errorf("Name mismatch: got %s, want Test User", user.Name)
	}
}

func TestParseUserFromToken_InvalidFormat(t *testing.T) {
	_, err := ParseUserFromToken("invalid")
	if err == nil {
		t.Error("Expected error for invalid token format")
	}
}

func TestParseUserFromToken_InvalidBase64(t *testing.T) {
	_, err := ParseUserFromToken("a.!!!invalid!!!.b")
	if err == nil {
		t.Error("Expected error for invalid base64")
	}
}

func TestParseUserFromToken_InvalidJSON(t *testing.T) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`not-json`))
	token := header + "." + payload + ".sig"

	_, err := ParseUserFromToken(token)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

