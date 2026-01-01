package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestFindFreePort(t *testing.T) {
	port, err := FindFreePort()
	if err != nil {
		t.Fatalf("FindFreePort failed: %v", err)
	}

	if port < 1024 || port > 65535 {
		t.Errorf("Port %d is outside valid range", port)
	}
}

func TestStartCallbackServer(t *testing.T) {
	port, err := FindFreePort()
	if err != nil {
		t.Fatalf("FindFreePort failed: %v", err)
	}

	resultChan := make(chan CallbackResult, 1)
	server := StartCallbackServer(port, resultChan)
	defer ShutdownServer(server)

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Test health endpoint
	healthURL := fmt.Sprintf("http://localhost:%d/health", port)
	resp, err := http.Get(healthURL)
	if err != nil {
		t.Logf("Health check failed (expected if server not ready): %v", err)
	} else {
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Health check returned %d, want %d", resp.StatusCode, http.StatusOK)
		}
		resp.Body.Close()
	}

	// Test callback endpoint with token
	testToken := "test-jwt-token"
	callbackURL := fmt.Sprintf("http://localhost:%d/callback?token=%s", port, testToken)
	go func() {
		http.Get(callbackURL)
	}()

	select {
	case result := <-resultChan:
		if result.Error != nil {
			t.Errorf("Unexpected error: %v", result.Error)
		}
		if result.Token != testToken {
			t.Errorf("Token mismatch: got %s, want %s", result.Token, testToken)
		}
	case <-time.After(2 * time.Second):
		t.Error("Timeout waiting for callback result")
	}
}

func TestCallbackServerNoToken(t *testing.T) {
	port, err := FindFreePort()
	if err != nil {
		t.Fatalf("FindFreePort failed: %v", err)
	}

	resultChan := make(chan CallbackResult, 1)
	server := StartCallbackServer(port, resultChan)
	defer ShutdownServer(server)

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Test callback without token
	callbackURL := fmt.Sprintf("http://localhost:%d/callback", port)
	go func() {
		http.Get(callbackURL)
	}()

	select {
	case result := <-resultChan:
		if result.Error == nil {
			t.Error("Expected error for missing token")
		}
	case <-time.After(2 * time.Second):
		t.Error("Timeout waiting for callback result")
	}
}
