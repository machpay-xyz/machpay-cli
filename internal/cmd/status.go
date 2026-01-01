// ============================================================
// Status Command - Show Current State
// ============================================================

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/machpay/machpay-cli/internal/auth"
	"github.com/machpay/machpay-cli/internal/config"
	"github.com/machpay/machpay-cli/internal/tui"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication and configuration status",
	Long: `Display the current status of your MachPay CLI setup.

Shows:
  - Authentication status
  - Configured role (agent/vendor)
  - Network (mainnet/devnet)
  - Gateway status (if running)`,
	RunE: runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfg := config.Get()

	fmt.Println()
	fmt.Println(tui.Header("MachPay Status"))
	fmt.Println()

	// Authentication
	fmt.Println(tui.Bold("Authentication"))
	if auth.IsLoggedIn() {
		user := auth.GetUser()
		if user != nil {
			fmt.Printf("  Status:  %s\n", tui.Success("● Logged in"))
			fmt.Printf("  Account: %s\n", tui.Primary(user.Email))
		} else {
			fmt.Printf("  Status:  %s\n", tui.Success("● Logged in"))
		}
	} else {
		fmt.Printf("  Status:  %s\n", tui.Muted("○ Not logged in"))
		fmt.Printf("  %s\n", tui.Muted("Run 'machpay login' to authenticate"))
	}
	fmt.Println()

	// Configuration
	fmt.Println(tui.Bold("Configuration"))
	if cfg.Role != "" {
		fmt.Printf("  Role:    %s\n", tui.Primary(cfg.Role))
	} else {
		fmt.Printf("  Role:    %s\n", tui.Muted("Not configured"))
	}
	fmt.Printf("  Network: %s\n", tui.Primary(cfg.Network))
	fmt.Printf("  Config:  %s\n", tui.Muted(config.GetPath()))
	fmt.Println()

	// Wallet (if configured)
	if cfg.Wallet.PublicKey != "" {
		fmt.Println(tui.Bold("Wallet"))
		fmt.Printf("  Address: %s\n", tui.Primary(truncateAddress(cfg.Wallet.PublicKey)))
		fmt.Println()
	}

	// Gateway (if vendor)
	if cfg.Role == "vendor" {
		fmt.Println(tui.Bold("Gateway"))
		if cfg.Gateway.Version != "" {
			fmt.Printf("  Version:  %s\n", cfg.Gateway.Version)
			fmt.Printf("  Port:     %d\n", cfg.Gateway.Port)
			// TODO: Check if actually running
			fmt.Printf("  Status:   %s\n", tui.Muted("○ Not running"))
		} else {
			fmt.Printf("  Status:   %s\n", tui.Muted("Not installed"))
			fmt.Printf("  %s\n", tui.Muted("Run 'machpay serve' to download and start"))
		}
		fmt.Println()
	}

	return nil
}

// truncateAddress shortens a Solana address for display
func truncateAddress(addr string) string {
	if len(addr) <= 12 {
		return addr
	}
	return addr[:6] + "..." + addr[len(addr)-4:]
}

