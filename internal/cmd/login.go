// ============================================================
// Login Command - Browser Redirect Authentication
// ============================================================
//
// Usage: machpay login [--no-browser]
//
// Flow:
// 1. Start local callback server on random port
// 2. Open browser to console.machpay.xyz/auth/cli?port=PORT
// 3. User logs in via Google/Wallet/Email
// 4. Console redirects to localhost:PORT/callback?token=JWT
// 5. CLI receives token, saves to config
//
// ============================================================

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	"github.com/machpay-xyz/machpay-cli/internal/auth"
	"github.com/machpay-xyz/machpay-cli/internal/config"
	"github.com/machpay-xyz/machpay-cli/internal/tui"
)

var (
	loginNoBrowser bool
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with MachPay",
	Long: `Link your CLI to your MachPay account via browser login.

This command opens your default browser to the MachPay login page.
After you sign in (using Google, Wallet, or Email), your CLI will
be automatically authenticated.

If you're on a headless system without a browser, use --no-browser
to get a URL you can open on another device.`,
	Example: `  # Standard login (opens browser)
  machpay login

  # Headless mode for SSH/CI environments
  machpay login --no-browser`,
	RunE: runLogin,
}

func init() {
	loginCmd.Flags().BoolVar(&loginNoBrowser, "no-browser", false, "Print URL instead of opening browser")
}

func runLogin(cmd *cobra.Command, args []string) error {
	// Check if already logged in
	if auth.IsLoggedIn() {
		user := auth.GetUser()
		if user != nil {
			fmt.Printf("Already logged in as %s\n", tui.Primary(user.Email))
		} else {
			fmt.Println("Already logged in.")
		}
		fmt.Println(tui.Muted("Use 'machpay logout' to sign out first."))
		return nil
	}

	// Find a free port for the callback server
	port, err := auth.FindFreePort()
	if err != nil {
		return fmt.Errorf("failed to find free port: %w", err)
	}

	// Create channel for callback result
	resultChan := make(chan auth.CallbackResult, 1)

	// Start callback server
	server := auth.StartCallbackServer(port, resultChan)
	defer auth.ShutdownServer(server)

	// Build login URL
	consoleURL := config.GetConsoleURL()
	loginURL := fmt.Sprintf("%s/auth/cli?port=%d", consoleURL, port)

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Open browser or print URL
	fmt.Println()
	if loginNoBrowser {
		fmt.Println(tui.Bold("Open this URL in your browser:"))
		fmt.Println()
		fmt.Printf("  %s\n", tui.Primary(loginURL))
		fmt.Println()
	} else {
		fmt.Println(tui.Info("🌐") + " Opening browser for authentication...")
		if err := browser.OpenURL(loginURL); err != nil {
			// Fall back to printing URL
			fmt.Println(tui.Warning("Could not open browser. Please open this URL manually:"))
			fmt.Println()
			fmt.Printf("  %s\n", tui.Primary(loginURL))
			fmt.Println()
		}
	}

	fmt.Println(tui.Muted("⏳ Waiting for login (press Ctrl+C to cancel)"))
	fmt.Println()

	// Wait for result, timeout, or interrupt
	ctx, cancel := context.WithTimeout(context.Background(), auth.CallbackTimeout)
	defer cancel()

	select {
	case result := <-resultChan:
		if result.Error != nil {
			return fmt.Errorf("login failed: %w", result.Error)
		}

		// Save token
		if err := auth.SaveToken(result.Token); err != nil {
			return fmt.Errorf("failed to save credentials: %w", err)
		}

		// Get user info
		user := auth.GetUser()
		if user != nil {
			tui.PrintSuccess(fmt.Sprintf("Logged in as %s", tui.Bold(user.Email)))
		} else {
			tui.PrintSuccess("Logged in successfully")
		}

		fmt.Println()
		fmt.Printf("  Credentials saved to: %s\n", tui.Muted(config.GetPath()))
		fmt.Println()
		fmt.Println(tui.Muted("Run 'machpay status' to verify your setup."))

		return nil

	case <-sigChan:
		fmt.Println()
		fmt.Println(tui.Warning("Login cancelled"))
		return nil

	case <-ctx.Done():
		return fmt.Errorf("login timed out after %v", auth.CallbackTimeout)
	}
}

