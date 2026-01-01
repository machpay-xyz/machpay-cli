// ============================================================
// Token Management - JWT Parsing and Credential Storage
// ============================================================

package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/machpay/machpay-cli/internal/config"
)

// User represents user info from JWT claims
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// IsLoggedIn checks if the user has valid credentials
func IsLoggedIn() bool {
	cfg := config.Get()
	return cfg.Auth.AccessToken != ""
}

// SaveToken saves the JWT token to config
func SaveToken(token string) error {
	cfg := config.Get()

	// Parse user info from token
	user, err := ParseUserFromToken(token)
	if err != nil {
		// Token is valid even if we can't parse claims
		cfg.Auth.AccessToken = token
	} else {
		cfg.Auth.AccessToken = token
		cfg.Auth.UserID = user.ID
		cfg.Auth.Email = user.Email
	}

	return config.Save()
}

// GetToken returns the stored access token
func GetToken() string {
	return config.Get().Auth.AccessToken
}

// GetUser returns the stored user info
func GetUser() *User {
	cfg := config.Get()
	if cfg.Auth.Email == "" {
		return nil
	}
	return &User{
		ID:    cfg.Auth.UserID,
		Email: cfg.Auth.Email,
	}
}

// ClearCredentials removes all stored auth credentials
func ClearCredentials() error {
	config.Clear()
	return config.Save()
}

// ParseUserFromToken extracts user info from a JWT
func ParseUserFromToken(token string) (*User, error) {
	// JWT format: header.payload.signature
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Decode payload (base64url)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}

	// Parse claims
	var claims struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("parse claims: %w", err)
	}

	return &User{
		ID:    claims.Sub,
		Email: claims.Email,
		Name:  claims.Name,
	}, nil
}

