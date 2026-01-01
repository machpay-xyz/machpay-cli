// ============================================================
// Root Command - machpay
// ============================================================

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/machpay/machpay-cli/internal/config"
	"github.com/machpay/machpay-cli/internal/tui"
)

var (
	// Version info (set by main.go)
	versionInfo struct {
		Version string
		Commit  string
		Date    string
	}

	// Global flags
	cfgFile string
	debug   bool
	noColor bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "machpay",
	Short: "MachPay CLI - Orchestrator for AI agent payments",
	Long: `MachPay CLI is the unified entry point for the MachPay Network.

For Agents: A wallet and payment tool for AI services.
For Vendors: An orchestrator that downloads and runs the payment gateway.
For Everyone: The onboarding wizard.

Get started:
  machpay login     Authenticate with MachPay
  machpay setup     Interactive setup wizard
  machpay status    Show current status

Documentation: https://docs.machpay.xyz/cli`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip initialization for version and help
		if cmd.Name() == "version" || cmd.Name() == "help" {
			return nil
		}

		// Initialize config
		if err := config.Init(cfgFile); err != nil {
			return fmt.Errorf("config init: %w", err)
		}

		// Apply color settings
		if noColor {
			tui.DisableColors()
		}

		return nil
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// SetVersionInfo sets version information from build flags
func SetVersionInfo(version, commit, date string) {
	versionInfo.Version = version
	versionInfo.Commit = commit
	versionInfo.Date = date
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ~/.machpay/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")

	// Bind flags to viper
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(statusCmd)
	// rootCmd.AddCommand(setupCmd)   // TODO: Phase 2
	// rootCmd.AddCommand(serveCmd)   // TODO: Phase 3
	// rootCmd.AddCommand(openCmd)    // TODO: Phase 2
}

// versionCmd shows version information
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version %s\n", tui.Bold("machpay"), tui.Success(versionInfo.Version))
		if debug {
			fmt.Printf("  Commit: %s\n", versionInfo.Commit)
			fmt.Printf("  Built:  %s\n", versionInfo.Date)
		}
	},
}

// getConfigDir returns the MachPay config directory path
func getConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home + "/.machpay"
}

