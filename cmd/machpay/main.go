// ============================================================
// MachPay CLI - Main Entry Point
// ============================================================
//
// The MachPay CLI is a lightweight orchestrator for the MachPay
// network. It handles authentication, gateway management, and
// provides a unified interface for both agents and vendors.
//
// Usage:
//   machpay login     - Authenticate with MachPay
//   machpay setup     - Interactive setup wizard
//   machpay serve     - Start vendor gateway
//   machpay status    - Show current status
//   machpay open      - Launch web console
//
// ============================================================

package main

import (
	"fmt"
	"os"

	"github.com/machpay-xyz/machpay-cli/internal/cmd"
)

// Version information (set at build time via ldflags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Set version info for the root command
	cmd.SetVersionInfo(version, commit, date)

	// Execute the root command
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

