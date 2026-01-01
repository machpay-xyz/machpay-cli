// ============================================================
// Logs Command - View gateway logs
// ============================================================
//
// Usage: machpay logs [-f] [--clear]
//
// Shows gateway log output. Use -f to follow in real-time.
//
// ============================================================

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/machpay-xyz/machpay-cli/internal/gateway"
	"github.com/machpay-xyz/machpay-cli/internal/tui"
)

var (
	logsFollow bool
	logsClear  bool
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View gateway logs",
	Long: `View the gateway log output.

Examples:
  machpay logs           # Show recent logs
  machpay logs -f        # Follow logs in real-time
  machpay logs --clear   # Clear the log file`,
	RunE: runLogs,
}

func init() {
	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Follow log output")
	logsCmd.Flags().BoolVar(&logsClear, "clear", false, "Clear log file")

	// Add logs command to root
	rootCmd.AddCommand(logsCmd)
}

func runLogs(cmd *cobra.Command, args []string) error {
	dl := gateway.NewDownloader()
	pm := gateway.NewProcessManager(dl.BinaryPath(), 0, "")

	// Handle clear flag
	if logsClear {
		if err := pm.ClearLogs(); err != nil {
			if os.IsNotExist(err) {
				fmt.Println(tui.Muted("No logs to clear"))
				return nil
			}
			return fmt.Errorf("clear logs: %w", err)
		}
		tui.PrintSuccess("Logs cleared")
		return nil
	}

	// Set up context for follow mode
	ctx := context.Background()

	if logsFollow {
		// Show status before following
		if pm.IsRunning() {
			pid, _ := pm.GetPID()
			fmt.Printf("Following logs (PID %d)... Press Ctrl+C to stop\n", pid)
		} else {
			fmt.Println("Following logs... Press Ctrl+C to stop")
			fmt.Println(tui.Muted("(Gateway is not currently running)"))
		}
		fmt.Println()

		// Set up signal handling
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigChan
			fmt.Println()
			cancel()
		}()
	}

	// Tail logs
	if err := pm.TailLogs(ctx, logsFollow, os.Stdout); err != nil {
		if os.IsNotExist(err) {
			fmt.Println(tui.Muted("No logs found"))
			fmt.Println(tui.Muted("  Run 'machpay serve' to start the gateway"))
			return nil
		}
		return err
	}

	return nil
}

