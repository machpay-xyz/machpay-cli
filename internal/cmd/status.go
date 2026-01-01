// ============================================================
// Status Command - Show Current State
// ============================================================
//
// Usage: machpay status [--json] [--watch]
//
// ============================================================

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/machpay-xyz/machpay-cli/internal/auth"
	"github.com/machpay-xyz/machpay-cli/internal/config"
	"github.com/machpay-xyz/machpay-cli/internal/tui"
)

var (
	statusJSON  bool
	statusWatch bool
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication and configuration status",
	Long: `Display the current status of your MachPay CLI setup.

Shows:
  - Authentication status
  - Configured role (agent/vendor)
  - Network (mainnet/devnet)
  - Wallet address
  - Gateway status (if vendor)

Flags:
  --json   Output as JSON for scripting
  --watch  Continuously update every 5 seconds`,
	RunE: runStatus,
}

func init() {
	statusCmd.Flags().BoolVar(&statusJSON, "json", false, "Output as JSON")
	statusCmd.Flags().BoolVar(&statusWatch, "watch", false, "Continuous monitoring (every 5s)")
}

// StatusOutput is the JSON-serializable status structure
type StatusOutput struct {
	Auth struct {
		LoggedIn bool   `json:"logged_in"`
		Email    string `json:"email,omitempty"`
		UserID   string `json:"user_id,omitempty"`
	} `json:"auth"`
	Config struct {
		Role    string `json:"role"`
		Network string `json:"network"`
		Path    string `json:"config_path"`
	} `json:"config"`
	Wallet struct {
		Address     string `json:"address,omitempty"`
		KeypairPath string `json:"keypair_path,omitempty"`
	} `json:"wallet,omitempty"`
	Gateway struct {
		Installed bool   `json:"installed"`
		Running   bool   `json:"running"`
		Version   string `json:"version,omitempty"`
		Port      int    `json:"port,omitempty"`
	} `json:"gateway,omitempty"`
	Vendor struct {
		UpstreamURL     string  `json:"upstream_url,omitempty"`
		PricePerRequest float64 `json:"price_per_request,omitempty"`
	} `json:"vendor,omitempty"`
}

func runStatus(cmd *cobra.Command, args []string) error {
	if statusWatch {
		return runStatusWatch()
	}
	return runStatusOnce()
}

func runStatusOnce() error {
	status := gatherStatus()

	if statusJSON {
		return outputJSON(status)
	}
	return outputHuman(status)
}

func runStatusWatch() error {
	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Initial display
	clearScreen()
	status := gatherStatus()
	if statusJSON {
		outputJSON(status)
	} else {
		outputHuman(status)
	}
	fmt.Println(tui.Muted("Watching... (Ctrl+C to exit)"))

	for {
		select {
		case <-ticker.C:
			clearScreen()
			status := gatherStatus()
			if statusJSON {
				outputJSON(status)
			} else {
				outputHuman(status)
			}
			fmt.Println(tui.Muted("Watching... (Ctrl+C to exit)"))

		case <-sigChan:
			fmt.Println()
			return nil
		}
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func gatherStatus() *StatusOutput {
	status := &StatusOutput{}
	cfg := config.Get()

	// Auth
	status.Auth.LoggedIn = auth.IsLoggedIn()
	if user := auth.GetUser(); user != nil {
		status.Auth.Email = user.Email
		status.Auth.UserID = user.ID
	}

	// Config
	status.Config.Role = cfg.Role
	status.Config.Network = cfg.Network
	status.Config.Path = config.GetPath()

	// Wallet
	if cfg.Wallet.PublicKey != "" {
		status.Wallet.Address = cfg.Wallet.PublicKey
		status.Wallet.KeypairPath = cfg.Wallet.KeypairPath
	}

	// Gateway (if vendor)
	if cfg.Role == "vendor" {
		status.Gateway.Installed = cfg.Gateway.Version != ""
		status.Gateway.Version = cfg.Gateway.Version
		status.Gateway.Port = cfg.Gateway.Port
		// TODO: Actually check if running via PID file
		status.Gateway.Running = false

		status.Vendor.UpstreamURL = cfg.Vendor.UpstreamURL
		status.Vendor.PricePerRequest = cfg.Vendor.PricePerRequest
	}

	return status
}

func outputJSON(status *StatusOutput) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(status)
}

func outputHuman(status *StatusOutput) error {
	fmt.Println()
	fmt.Println(tui.Header("MachPay Status"))
	fmt.Println()

	// Authentication
	fmt.Println(tui.Bold("Authentication"))
	if status.Auth.LoggedIn {
		fmt.Printf("  Status:  %s\n", tui.Success("● Logged in"))
		if status.Auth.Email != "" {
			fmt.Printf("  Account: %s\n", tui.Primary(status.Auth.Email))
		}
	} else {
		fmt.Printf("  Status:  %s\n", tui.Muted("○ Not logged in"))
		fmt.Printf("  %s\n", tui.Muted("Run 'machpay login' to authenticate"))
	}
	fmt.Println()

	// Configuration
	fmt.Println(tui.Bold("Configuration"))
	if status.Config.Role != "" {
		fmt.Printf("  Role:    %s\n", tui.Primary(status.Config.Role))
	} else {
		fmt.Printf("  Role:    %s\n", tui.Muted("Not configured"))
		fmt.Printf("  %s\n", tui.Muted("Run 'machpay setup' to configure"))
	}
	fmt.Printf("  Network: %s\n", tui.Primary(status.Config.Network))
	fmt.Printf("  Config:  %s\n", tui.Muted(status.Config.Path))
	fmt.Println()

	// Wallet (if configured)
	if status.Wallet.Address != "" {
		fmt.Println(tui.Bold("Wallet"))
		fmt.Printf("  Address: %s\n", tui.Primary(truncateAddress(status.Wallet.Address)))
		// TODO: Add balance fetching
		fmt.Println()
	}

	// Gateway (if vendor)
	if status.Config.Role == "vendor" {
		fmt.Println(tui.Bold("Gateway"))
		if status.Gateway.Installed {
			fmt.Printf("  Version: %s\n", status.Gateway.Version)
			fmt.Printf("  Port:    %d\n", status.Gateway.Port)
			if status.Gateway.Running {
				fmt.Printf("  Status:  %s\n", tui.Success("● Running"))
			} else {
				fmt.Printf("  Status:  %s\n", tui.Muted("○ Not running"))
			}
		} else {
			fmt.Printf("  Status:  %s\n", tui.Muted("Not installed"))
			fmt.Printf("  %s\n", tui.Muted("Run 'machpay serve' to download and start"))
		}
		fmt.Println()

		// Vendor config
		if status.Vendor.UpstreamURL != "" {
			fmt.Println(tui.Bold("Vendor Config"))
			fmt.Printf("  Upstream: %s\n", tui.Muted(status.Vendor.UpstreamURL))
			fmt.Printf("  Price:    %s USDC/req\n", tui.Primary(fmt.Sprintf("%.4f", status.Vendor.PricePerRequest)))
			fmt.Println()
		}
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

