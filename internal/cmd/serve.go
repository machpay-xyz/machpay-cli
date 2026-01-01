// ============================================================
// Serve Command - Start the payment gateway
// ============================================================
//
// Usage: machpay serve [--port] [--upstream] [--detach] [--debug]
//
// Downloads the gateway if not installed, then starts it.
// Foreground mode shows live output, background mode detaches.
//
// ============================================================

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/machpay-xyz/machpay-cli/internal/auth"
	"github.com/machpay-xyz/machpay-cli/internal/config"
	"github.com/machpay-xyz/machpay-cli/internal/gateway"
	"github.com/machpay-xyz/machpay-cli/internal/tui"
)

var (
	servePort     int
	serveUpstream string
	serveDetach   bool
	serveDebug    bool
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the payment gateway",
	Long: `Start the MachPay payment gateway for your vendor service.

The gateway acts as a reverse proxy in front of your API, handling
x402 payment verification and request forwarding.

If the gateway is not installed, it will be downloaded automatically.

Examples:
  machpay serve                           # Start with config defaults
  machpay serve --port 8402               # Custom port
  machpay serve --upstream http://localhost:11434  # Override upstream
  machpay serve --detach                  # Run in background`,
	RunE: runServe,
}

func init() {
	serveCmd.Flags().IntVar(&servePort, "port", 8402, "Gateway listen port")
	serveCmd.Flags().StringVar(&serveUpstream, "upstream", "", "Upstream API URL (overrides config)")
	serveCmd.Flags().BoolVar(&serveDetach, "detach", false, "Run in background")
	serveCmd.Flags().BoolVar(&serveDebug, "debug", false, "Enable debug logging")
}

func runServe(cmd *cobra.Command, args []string) error {
	// 1. Check prerequisites
	if !auth.IsLoggedIn() {
		tui.PrintError("Not logged in")
		fmt.Println(tui.Muted("  Run 'machpay login' first"))
		return fmt.Errorf("not logged in")
	}

	cfg := config.Get()
	if cfg.Role != "vendor" {
		tui.PrintError("Not configured as vendor")
		fmt.Println(tui.Muted("  Run 'machpay setup' and select 'Vendor'"))
		return fmt.Errorf("not a vendor")
	}

	// Use config values if flags not set
	upstream := serveUpstream
	if upstream == "" {
		upstream = cfg.Vendor.UpstreamURL
	}
	if upstream == "" {
		tui.PrintError("No upstream URL configured")
		fmt.Println(tui.Muted("  Run 'machpay setup' to configure your upstream URL"))
		fmt.Println(tui.Muted("  Or use --upstream flag"))
		return fmt.Errorf("no upstream URL")
	}

	// 2. Ensure gateway is installed
	dl := gateway.NewDownloader()
	if !dl.IsInstalled() {
		if err := downloadGateway(dl); err != nil {
			return err
		}
	} else {
		// Check for updates in background
		go checkForGatewayUpdates(dl)
	}

	// 3. Create process manager
	pm := gateway.NewProcessManager(dl.BinaryPath(), servePort, upstream)
	pm.SetDebug(serveDebug)

	// 4. Handle detach mode
	if serveDetach {
		return runDetached(pm, upstream)
	}

	// 5. Run in foreground
	return runForeground(pm, upstream)
}

func downloadGateway(dl *gateway.Downloader) error {
	fmt.Println()
	tui.PrintWarning("Gateway not found. Downloading...")
	fmt.Println()

	// Get latest release
	release, err := dl.GetLatestRelease()
	if err != nil {
		return fmt.Errorf("fetch latest version: %w", err)
	}

	version := release.TagName
	fmt.Printf("  Downloading %s for %s/%s...\n",
		tui.Primary(version), runtime.GOOS, runtime.GOARCH)
	fmt.Println()

	if err := dl.Download(version, os.Stdout); err != nil {
		return fmt.Errorf("download gateway: %w", err)
	}

	fmt.Println()
	tui.PrintSuccess(fmt.Sprintf("Gateway %s installed", version))
	fmt.Println()

	return nil
}

func runDetached(pm *gateway.ProcessManager, upstream string) error {
	if pm.IsRunning() {
		pid, _ := pm.GetPID()
		fmt.Printf("Gateway already running (PID %d)\n", pid)
		fmt.Println(tui.Muted("  Use 'machpay stop' to stop it"))
		return nil
	}

	fmt.Println()
	fmt.Println(tui.Bold("Starting gateway in background..."))

	if err := pm.Start(); err != nil {
		return fmt.Errorf("start gateway: %w", err)
	}

	// Wait a moment and check health
	time.Sleep(2 * time.Second)

	if err := pm.HealthCheck(); err != nil {
		pid, _ := pm.GetPID()
		tui.PrintWarning(fmt.Sprintf("Gateway started (PID %d) but health check failed", pid))
		fmt.Println(tui.Muted("  Check logs with 'machpay logs'"))
		return nil
	}

	pid, _ := pm.GetPID()
	fmt.Println()
	tui.PrintSuccess(fmt.Sprintf("Gateway started (PID %d)", pid))
	fmt.Println()
	tui.PrintKeyValue("Port", fmt.Sprintf("http://localhost:%d", servePort))
	tui.PrintKeyValue("Upstream", upstream)
	tui.PrintKeyValue("Logs", pm.LogFile())
	fmt.Println()
	fmt.Println(tui.Muted("  Use 'machpay logs -f' to follow logs"))
	fmt.Println(tui.Muted("  Use 'machpay stop' to stop the gateway"))

	return nil
}

func runForeground(pm *gateway.ProcessManager, upstream string) error {
	if pm.IsRunning() {
		pid, _ := pm.GetPID()
		return fmt.Errorf("gateway already running (PID %d). Stop it first with 'machpay stop'", pid)
	}

	// Print status header
	printGatewayHeader(upstream)

	// Handle Ctrl+C
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println()
		fmt.Println(tui.Muted("Shutting down..."))
		cancel()
	}()

	// Start gateway
	err := pm.StartForeground(ctx, os.Stdout, os.Stderr)

	if ctx.Err() != nil {
		// Graceful shutdown via Ctrl+C
		fmt.Println()
		tui.PrintSuccess("Gateway stopped")
		return nil
	}

	return err
}

func printGatewayHeader(upstream string) {
	fmt.Println()
	fmt.Println(tui.Bold("MachPay Gateway"))
	fmt.Println()
	tui.PrintKeyValue("Port", fmt.Sprintf("http://localhost:%d", servePort))
	tui.PrintKeyValue("Upstream", upstream)
	fmt.Println()
	fmt.Println(tui.Muted("Press Ctrl+C to stop"))
	tui.PrintSection()
}

func checkForGatewayUpdates(dl *gateway.Downloader) {
	needsUpdate, latest, err := dl.NeedsUpdate()
	if err != nil || !needsUpdate {
		return
	}

	fmt.Println()
	tui.PrintInfo(fmt.Sprintf("Update available: %s", latest))
	fmt.Println(tui.Muted("  Run 'machpay update' to upgrade"))
}

// ============================================================
// Stop Command
// ============================================================

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the gateway",
	Long: `Stop the MachPay gateway if it's running in the background.

Examples:
  machpay stop           # Stop gracefully
  machpay stop --force   # Force kill`,
	RunE: runStop,
}

var stopForce bool

func init() {
	stopCmd.Flags().BoolVar(&stopForce, "force", false, "Force kill the gateway")
}

func runStop(cmd *cobra.Command, args []string) error {
	dl := gateway.NewDownloader()
	pm := gateway.NewProcessManager(dl.BinaryPath(), 0, "")

	if !pm.IsRunning() {
		fmt.Println(tui.Muted("Gateway is not running"))
		return nil
	}

	pid, _ := pm.GetPID()
	fmt.Printf("Stopping gateway (PID %d)...\n", pid)

	var err error
	if stopForce {
		err = pm.Kill()
	} else {
		err = pm.Stop()
	}

	if err != nil {
		return fmt.Errorf("stop gateway: %w", err)
	}

	tui.PrintSuccess("Gateway stopped")
	return nil
}

// ============================================================
// Restart Command
// ============================================================

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the gateway",
	Long: `Restart the MachPay gateway (stop + start).

Examples:
  machpay restart`,
	RunE: runRestart,
}

func runRestart(cmd *cobra.Command, args []string) error {
	dl := gateway.NewDownloader()
	pm := gateway.NewProcessManager(dl.BinaryPath(), 0, "")

	if pm.IsRunning() {
		pid, _ := pm.GetPID()
		fmt.Printf("Stopping gateway (PID %d)...\n", pid)
		if err := pm.Stop(); err != nil {
			return fmt.Errorf("stop gateway: %w", err)
		}
		tui.PrintSuccess("Gateway stopped")
		fmt.Println()
	}

	// Re-run serve in detached mode
	serveDetach = true
	return runServe(cmd, args)
}

