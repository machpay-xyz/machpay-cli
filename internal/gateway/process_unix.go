//go:build !windows

// ============================================================
// Process Manager - Unix-specific implementations
// ============================================================

package gateway

import (
	"os"
	"os/exec"
	"syscall"
)

// setProcAttr sets Unix-specific process attributes for detaching
func setProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

// isProcessAlive checks if a process is alive on Unix
func isProcessAlive(process *os.Process) bool {
	// Sending signal 0 checks if process exists without doing anything
	err := process.Signal(syscall.Signal(0))
	return err == nil
}

// sendTermSignal sends SIGTERM to gracefully stop a process
func sendTermSignal(process *os.Process) error {
	return process.Signal(syscall.SIGTERM)
}

// sendKillSignal sends SIGKILL to forcefully stop a process
func sendKillSignal(process *os.Process) error {
	return process.Signal(syscall.SIGKILL)
}

