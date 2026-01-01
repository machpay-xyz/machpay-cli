// ============================================================
// Root Command - machpay
// ============================================================

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/machpay-xyz/machpay-cli/internal/config"
	"github.com/machpay-xyz/machpay-cli/internal/tui"
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
	Long: `
███╗   ███╗ █████╗  ██████╗██╗  ██╗██████╗  █████╗ ██╗   ██╗
████╗ ████║██╔══██╗██╔════╝██║  ██║██╔══██╗██╔══██╗╚██╗ ██╔╝
██╔████╔██║███████║██║     ███████║██████╔╝███████║ ╚████╔╝ 
██║╚██╔╝██║██╔══██║██║     ██╔══██║██╔═══╝ ██╔══██║  ╚██╔╝  
██║ ╚═╝ ██║██║  ██║╚██████╗██║  ██║██║     ██║  ██║   ██║   
╚═╝     ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝     ╚═╝  ╚═╝   ╚═╝   

MachPay CLI is the unified entry point for the MachPay Network.

ROLES:
  Agent   - Wallet and payment tool for AI services
  Vendor  - Orchestrator that downloads and runs the payment gateway

QUICK START:
  machpay login     Authenticate with MachPay
  machpay setup     Interactive setup wizard
  machpay status    Show current status

VENDOR COMMANDS:
  machpay serve     Start the payment gateway
  machpay stop      Stop the gateway
  machpay logs      View gateway logs

DOCUMENTATION:
  https://docs.machpay.xyz/cli`,
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
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
}

// versionCmd shows version information
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Display the current version of MachPay CLI.

With --debug flag, also shows commit hash and build date.`,
	Example: `  machpay version
  machpay version --debug`,
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

