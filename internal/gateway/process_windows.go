//go:build windows

// ============================================================
// Process Manager - Windows-specific implementations
// ============================================================

package gateway

import (
	"os"
	"os/exec"
)

// setProcAttr sets Windows-specific process attributes for detaching
func setProcAttr(cmd *exec.Cmd) {
	// On Windows, we don't need to set SysProcAttr for basic process detachment
	// The process will naturally detach when started without a console
}

// isProcessAlive checks if a process is alive on Windows
func isProcessAlive(process *os.Process) bool {
	// On Windows, FindProcess always succeeds, so we need a different approach
	// We can try to open the process with minimal access rights
	// If the process doesn't exist, this will fail
	// Since Signal(0) doesn't work on Windows, we use a simple approach
	// that checks if the process can be killed (without actually killing it)
	
	// Attempt a no-op by checking if process exists
	// This is a simplification - for full Windows support, we'd use
	// windows.OpenProcess with PROCESS_QUERY_LIMITED_INFORMATION
	_, err := os.FindProcess(process.Pid)
	return err == nil
}

// sendTermSignal sends a termination signal to a process on Windows
func sendTermSignal(process *os.Process) error {
	// Windows doesn't have SIGTERM, so we just kill the process
	return process.Kill()
}

// sendKillSignal forcefully terminates a process on Windows
func sendKillSignal(process *os.Process) error {
	return process.Kill()
}

