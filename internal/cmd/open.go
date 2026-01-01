// ============================================================
// Open Command - Launch MachPay Console
// ============================================================
//
// Usage: machpay open [route] [--web]
//
// Opens the MachPay console in your browser or desktop app.
// Supports specific routes like marketplace, funding, etc.
//
// ============================================================

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	"github.com/machpay-xyz/machpay-cli/internal/config"
	"github.com/machpay-xyz/machpay-cli/internal/tui"
)

var openWeb bool

var openCmd = &cobra.Command{
	Use:   "open [route]",
	Short: "Launch MachPay Console",
	Long: `Open the MachPay Console in your default browser or desktop app.

Available routes:
  (none)       Home / Explorer
  marketplace  API Marketplace
  funding      Fund your wallet
  settings     Account settings
  analytics    Usage analytics
  endpoints    Vendor endpoints

Examples:
  machpay open              # Opens console home
  machpay open marketplace  # Opens marketplace
  machpay open funding      # Opens funding page
  machpay open --web        # Force browser (not desktop app)`,
	Args: cobra.MaximumNArgs(1),
	RunE: runOpen,
}

func init() {
	openCmd.Flags().BoolVar(&openWeb, "web", false, "Force open in browser (not desktop app)")
}

// Route mappings
var routes = map[string]string{
	"":            "/explorer",
	"home":        "/explorer",
	"explorer":    "/explorer",
	"marketplace": "/marketplace",
	"funding":     "/agent/finance",
	"finance":     "/agent/finance",
	"settings":    "/settings",
	"analytics":   "/agent/analytics",
	"endpoints":   "/vendor/endpoints",
	"integration": "/vendor/integration",
	"payouts":     "/vendor/payouts",
}

func runOpen(cmd *cobra.Command, args []string) error {
	// Determine route
	route := ""
	if len(args) > 0 {
		route = args[0]
	}

	path, ok := routes[route]
	if !ok {
		fmt.Printf("%s Unknown route: %s\n", tui.ErrorIcon(), route)
		fmt.Println()
		fmt.Println("Available routes:")
		for name, p := range routes {
			if name != "" {
				fmt.Printf("  %-12s → %s\n", tui.Primary(name), tui.Muted(p))
			}
		}
		return fmt.Errorf("unknown route: %s", route)
	}

	// Build full URL
	consoleURL := config.GetConsoleURL() + path

	// Try desktop app first (unless --web)
	if !openWeb {
		appPath := getDesktopAppPath()
		if appPath != "" {
			if _, err := os.Stat(appPath); err == nil {
				if launchDesktopApp(appPath, path) {
					tui.PrintSuccess("Opened MachPay Console")
					return nil
				}
			}
		}
	}

	// Fall back to browser
	fmt.Printf("Opening %s...\n", tui.Primary(consoleURL))
	if err := browser.OpenURL(consoleURL); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	tui.PrintSuccess("Opened in browser")
	return nil
}

// getDesktopAppPath returns the path to the desktop app if installed
func getDesktopAppPath() string {
	switch runtime.GOOS {
	case "darwin":
		return "/Applications/MachPay.app"
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData != "" {
			return localAppData + "\\MachPay\\MachPay.exe"
		}
		return ""
	case "linux":
		home, _ := os.UserHomeDir()
		return home + "/.local/bin/machpay-console"
	default:
		return ""
	}
}

// launchDesktopApp attempts to launch the desktop app
func launchDesktopApp(appPath, route string) bool {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// macOS: open -a "MachPay" --args --route=/marketplace
		cmd = exec.Command("open", "-a", appPath, "--args", "--route="+route)
	case "windows":
		// Windows: Start app with arguments
		cmd = exec.Command("cmd", "/c", "start", "", appPath, "--route="+route)
	case "linux":
		// Linux: Run the binary directly
		cmd = exec.Command(appPath, "--route="+route)
	default:
		return false
	}

	if cmd == nil {
		return false
	}

	return cmd.Start() == nil
}

