// ============================================================
// Logout Command - Clear Credentials
// ============================================================

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/machpay/machpay-cli/internal/auth"
	"github.com/machpay/machpay-cli/internal/tui"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Sign out and clear stored credentials",
	Long: `Sign out of MachPay and remove all stored credentials.

This will remove your access token from the local configuration.
You will need to run 'machpay login' again to authenticate.`,
	RunE: runLogout,
}

func runLogout(cmd *cobra.Command, args []string) error {
	// Check if logged in
	if !auth.IsLoggedIn() {
		fmt.Println(tui.Muted("Not currently logged in."))
		return nil
	}

	// Get user for confirmation message
	user := auth.GetUser()

	// Clear credentials
	if err := auth.ClearCredentials(); err != nil {
		return fmt.Errorf("failed to clear credentials: %w", err)
	}

	// Success message
	if user != nil {
		tui.PrintSuccess(fmt.Sprintf("Logged out from %s", tui.Bold(user.Email)))
	} else {
		tui.PrintSuccess("Logged out successfully")
	}

	return nil
}

