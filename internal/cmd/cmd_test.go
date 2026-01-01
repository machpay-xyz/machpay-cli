// ============================================================
// Command Tests
// ============================================================

package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	// Test that root command exists and has expected properties
	if rootCmd.Use != "machpay" {
		t.Errorf("rootCmd.Use = %v, want machpay", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("rootCmd.Short should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("rootCmd.Long should not be empty")
	}
}

func TestVersionCommand(t *testing.T) {
	SetVersionInfo("1.0.0", "abc123", "2024-01-01")

	// Test that version command exists and has expected properties
	if versionCmd.Use != "version" {
		t.Errorf("versionCmd.Use = %v, want version", versionCmd.Use)
	}

	if versionCmd.Short == "" {
		t.Error("versionCmd.Short should not be empty")
	}

	// Version info should be set
	if versionInfo.Version != "1.0.0" {
		t.Errorf("Version = %v, want 1.0.0", versionInfo.Version)
	}
}

func TestSubcommands(t *testing.T) {
	expectedCommands := []string{
		"version",
		"login",
		"logout",
		"status",
		"setup",
		"open",
		"serve",
		"stop",
		"restart",
		"logs",
		"update",
	}

	commands := rootCmd.Commands()
	commandMap := make(map[string]*cobra.Command)
	for _, cmd := range commands {
		commandMap[cmd.Name()] = cmd
	}

	for _, name := range expectedCommands {
		if _, ok := commandMap[name]; !ok {
			t.Errorf("expected command %q not found", name)
		}
	}
}

func TestGlobalFlags(t *testing.T) {
	// Test that global flags are registered
	flags := []string{
		"config",
		"debug",
		"no-color",
	}

	for _, flag := range flags {
		f := rootCmd.PersistentFlags().Lookup(flag)
		if f == nil {
			t.Errorf("expected flag --%s not found", flag)
		}
	}
}

func TestLoginCommandFlags(t *testing.T) {
	f := loginCmd.Flags().Lookup("no-browser")
	if f == nil {
		t.Error("login command should have --no-browser flag")
	}
}

func TestServeCommandFlags(t *testing.T) {
	flags := []string{
		"port",
		"upstream",
		"detach",
		"debug",
	}

	for _, flag := range flags {
		f := serveCmd.Flags().Lookup(flag)
		if f == nil {
			t.Errorf("serve command should have --%s flag", flag)
		}
	}
}

func TestSetupCommandFlags(t *testing.T) {
	f := setupCmd.Flags().Lookup("non-interactive")
	if f == nil {
		t.Error("setup command should have --non-interactive flag")
	}
}

func TestStatusCommandFlags(t *testing.T) {
	flags := []string{
		"json",
		"watch",
	}

	for _, flag := range flags {
		f := statusCmd.Flags().Lookup(flag)
		if f == nil {
			t.Errorf("status command should have --%s flag", flag)
		}
	}
}

func TestSetVersionInfo(t *testing.T) {
	SetVersionInfo("2.0.0", "def456", "2024-12-31")

	if versionInfo.Version != "2.0.0" {
		t.Errorf("Version = %v, want 2.0.0", versionInfo.Version)
	}
	if versionInfo.Commit != "def456" {
		t.Errorf("Commit = %v, want def456", versionInfo.Commit)
	}
	if versionInfo.Date != "2024-12-31" {
		t.Errorf("Date = %v, want 2024-12-31", versionInfo.Date)
	}
}

func TestGetConfigDir(t *testing.T) {
	dir := getConfigDir()
	// Should end with .machpay
	if dir != "" && dir[len(dir)-8:] != ".machpay" {
		t.Errorf("getConfigDir() = %v, should end with .machpay", dir)
	}
}

