// ============================================================
// Update Command - Update CLI and gateway
// ============================================================
//
// Usage: machpay update [gateway|cli|all]
//
// Checks for and installs updates to the CLI and gateway.
//
// ============================================================

package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/machpay/machpay-cli/internal/gateway"
	"github.com/machpay/machpay-cli/internal/tui"
)

var updateCmd = &cobra.Command{
	Use:   "update [target]",
	Short: "Update CLI and gateway",
	Long: `Check for and install updates to the CLI and gateway.

Targets:
  gateway    Update gateway only
  cli        Update CLI only (shows instructions)
  all        Update everything (default)

Examples:
  machpay update           # Update everything
  machpay update gateway   # Update gateway only
  machpay update cli       # Show CLI update instructions`,
	Args: cobra.MaximumNArgs(1),
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	target := "all"
	if len(args) > 0 {
		target = args[0]
	}

	switch target {
	case "gateway":
		return updateGateway()
	case "cli":
		return updateCLI()
	case "all":
		fmt.Println()
		fmt.Println(tui.Bold("Checking for updates..."))
		fmt.Println()

		// Update gateway
		if err := updateGateway(); err != nil {
			tui.PrintWarning(fmt.Sprintf("Gateway: %v", err))
		}

		fmt.Println()

		// Show CLI instructions
		return updateCLI()
	default:
		return fmt.Errorf("unknown target: %s (use 'gateway', 'cli', or 'all')", target)
	}
}

func updateGateway() error {
	dl := gateway.NewDownloader()

	fmt.Println(tui.Bold("Gateway"))

	// Check if installed
	if !dl.IsInstalled() {
		fmt.Println("  " + tui.Muted("Not installed"))
		fmt.Println("  " + tui.Muted("Run 'machpay serve' to install"))
		return nil
	}

	// Get installed version
	installed, err := dl.InstalledVersion()
	if err != nil {
		fmt.Printf("  Installed: %s\n", tui.Muted("unknown"))
	} else {
		fmt.Printf("  Installed: %s\n", tui.Primary(installed))
	}

	// Check for updates
	needsUpdate, latest, err := dl.NeedsUpdate()
	if err != nil {
		return fmt.Errorf("check for updates: %w", err)
	}

	if !needsUpdate {
		fmt.Println("  " + tui.Success("✓ Up to date"))
		return nil
	}

	fmt.Printf("  Latest:    %s\n", tui.Primary(latest))
	fmt.Println()

	// Check if gateway is running
	pm := gateway.NewProcessManager(dl.BinaryPath(), 0, "")
	wasRunning := pm.IsRunning()

	if wasRunning {
		fmt.Println("  Stopping gateway for update...")
		if err := pm.Stop(); err != nil {
			return fmt.Errorf("stop gateway: %w", err)
		}
	}

	// Download update
	fmt.Printf("  Downloading v%s...\n", latest)
	fmt.Println()

	if err := dl.Download("v"+latest, os.Stdout); err != nil {
		return fmt.Errorf("download: %w", err)
	}

	fmt.Println()
	tui.PrintSuccess(fmt.Sprintf("Gateway updated to v%s", latest))

	// Restart if was running
	if wasRunning {
		fmt.Println()
		fmt.Println("  Restarting gateway...")
		if err := pm.Start(); err != nil {
			tui.PrintWarning("Failed to restart gateway")
			fmt.Println("  " + tui.Muted("Run 'machpay serve --detach' to restart manually"))
		} else {
			tui.PrintSuccess("Gateway restarted")
		}
	}

	return nil
}

func updateCLI() error {
	fmt.Println(tui.Bold("CLI"))
	fmt.Printf("  Installed: %s\n", tui.Primary(versionInfo.Version))

	// TODO: Check for CLI updates via GitHub API
	// For now, just show instructions

	fmt.Println()
	fmt.Println("  To update the CLI, run:")
	fmt.Println()

	switch runtime.GOOS {
	case "darwin":
		fmt.Println("    brew upgrade machpay")
	case "linux":
		fmt.Println("    curl -fsSL https://machpay.xyz/install.sh | sh")
	case "windows":
		fmt.Println("    winget upgrade machpay")
		fmt.Println()
		fmt.Println("  Or download from: https://github.com/machpay/machpay-cli/releases")
	default:
		fmt.Println("    curl -fsSL https://machpay.xyz/install.sh | sh")
	}

	return nil
}

