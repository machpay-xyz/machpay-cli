// ============================================================
// Process Manager - Gateway process lifecycle
// ============================================================
//
// Manages the gateway process:
// - Start in foreground or background (detached)
// - Stop with graceful shutdown
// - Check if running via PID file
// - Stream logs
// - Health checks
//
// ============================================================

package gateway

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Process errors
var (
	ErrNotRunning     = fmt.Errorf("gateway is not running")
	ErrAlreadyRunning = fmt.Errorf("gateway is already running")
)

// ProcessManager manages the gateway process lifecycle
type ProcessManager struct {
	binaryPath string
	configDir  string
	port       int
	upstream   string
	debug      bool
}

// NewProcessManager creates a new process manager
func NewProcessManager(binaryPath string, port int, upstream string) *ProcessManager {
	home, _ := os.UserHomeDir()
	return &ProcessManager{
		binaryPath: binaryPath,
		configDir:  filepath.Join(home, ".machpay"),
		port:       port,
		upstream:   upstream,
	}
}

// SetDebug enables debug mode
func (pm *ProcessManager) SetDebug(debug bool) {
	pm.debug = debug
}

// SetPort sets the port
func (pm *ProcessManager) SetPort(port int) {
	pm.port = port
}

// SetUpstream sets the upstream URL
func (pm *ProcessManager) SetUpstream(upstream string) {
	pm.upstream = upstream
}

// PIDFile returns the path to the PID file
func (pm *ProcessManager) PIDFile() string {
	return filepath.Join(pm.configDir, "gateway.pid")
}

// LogFile returns the path to the log file
func (pm *ProcessManager) LogFile() string {
	return filepath.Join(pm.configDir, "gateway.log")
}

// ============================================================
// Start Methods
// ============================================================

// Start starts the gateway process in the background (detached)
func (pm *ProcessManager) Start() error {
	if pm.IsRunning() {
		return ErrAlreadyRunning
	}

	// Ensure config dir exists
	if err := os.MkdirAll(pm.configDir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	// Build command arguments
	args := pm.buildArgs()

	cmd := exec.Command(pm.binaryPath, args...)

	// Set up logging to file
	logFile, err := os.OpenFile(pm.LogFile(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// Detach from parent process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Start process
	if err := cmd.Start(); err != nil {
		logFile.Close()
		return fmt.Errorf("start gateway: %w", err)
	}

	// Save PID
	if err := pm.savePID(cmd.Process.Pid); err != nil {
		cmd.Process.Kill()
		logFile.Close()
		return fmt.Errorf("save PID: %w", err)
	}

	// Don't wait for process (it's detached)
	go func() {
		cmd.Wait()
		logFile.Close()
		os.Remove(pm.PIDFile())
	}()

	return nil
}

// StartForeground starts the gateway in the foreground (blocking)
func (pm *ProcessManager) StartForeground(ctx context.Context, stdout, stderr io.Writer) error {
	if pm.IsRunning() {
		return ErrAlreadyRunning
	}

	// Ensure config dir exists
	if err := os.MkdirAll(pm.configDir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	// Build command arguments
	args := pm.buildArgs()

	cmd := exec.CommandContext(ctx, pm.binaryPath, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	// Start process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start gateway: %w", err)
	}

	// Save PID
	pm.savePID(cmd.Process.Pid)
	defer os.Remove(pm.PIDFile())

	// Wait for process to exit
	return cmd.Wait()
}

// buildArgs builds command line arguments
func (pm *ProcessManager) buildArgs() []string {
	var args []string

	if pm.port > 0 {
		args = append(args, "--port", strconv.Itoa(pm.port))
	}

	if pm.upstream != "" {
		args = append(args, "--upstream", pm.upstream)
	}

	if pm.debug {
		args = append(args, "--debug")
	}

	return args
}

// ============================================================
// Stop Methods
// ============================================================

// Stop stops the gateway process gracefully
func (pm *ProcessManager) Stop() error {
	pid, err := pm.loadPID()
	if err != nil {
		return ErrNotRunning
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		os.Remove(pm.PIDFile())
		return ErrNotRunning
	}

	// Check if process is actually running
	if !pm.isProcessAlive(process) {
		os.Remove(pm.PIDFile())
		return ErrNotRunning
	}

	// Send SIGTERM for graceful shutdown
	if err := process.Signal(syscall.SIGTERM); err != nil {
		os.Remove(pm.PIDFile())
		return ErrNotRunning
	}

	// Wait for process to exit (with timeout)
	done := make(chan error, 1)
	go func() {
		_, err := process.Wait()
		done <- err
	}()

	select {
	case <-done:
		os.Remove(pm.PIDFile())
		return nil
	case <-time.After(10 * time.Second):
		// Force kill if graceful shutdown takes too long
		process.Signal(syscall.SIGKILL)
		os.Remove(pm.PIDFile())
		return nil
	}
}

// Kill forcefully kills the gateway process
func (pm *ProcessManager) Kill() error {
	pid, err := pm.loadPID()
	if err != nil {
		return ErrNotRunning
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		os.Remove(pm.PIDFile())
		return ErrNotRunning
	}

	if err := process.Signal(syscall.SIGKILL); err != nil {
		os.Remove(pm.PIDFile())
		return ErrNotRunning
	}

	os.Remove(pm.PIDFile())
	return nil
}

// ============================================================
// Status Methods
// ============================================================

// IsRunning checks if the gateway is running
func (pm *ProcessManager) IsRunning() bool {
	pid, err := pm.loadPID()
	if err != nil {
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	return pm.isProcessAlive(process)
}

// GetPID returns the PID of the running gateway
func (pm *ProcessManager) GetPID() (int, error) {
	return pm.loadPID()
}

// isProcessAlive checks if a process is alive
func (pm *ProcessManager) isProcessAlive(process *os.Process) bool {
	// Sending signal 0 checks if process exists without doing anything
	err := process.Signal(syscall.Signal(0))
	return err == nil
}

// ============================================================
// Health Check
// ============================================================

// HealthCheck checks if the gateway is responding to HTTP requests
func (pm *ProcessManager) HealthCheck() error {
	url := fmt.Sprintf("http://localhost:%d/healthz", pm.port)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned %d", resp.StatusCode)
	}

	return nil
}

// WaitForHealthy waits for the gateway to become healthy
func (pm *ProcessManager) WaitForHealthy(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if err := pm.HealthCheck(); err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("gateway did not become healthy within %s", timeout)
}

// ============================================================
// Log Streaming
// ============================================================

// TailLogs streams logs from the log file
func (pm *ProcessManager) TailLogs(ctx context.Context, follow bool, writer io.Writer) error {
	file, err := os.Open(pm.LogFile())
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("log file not found: %s", pm.LogFile())
		}
		return fmt.Errorf("open log file: %w", err)
	}
	defer file.Close()

	if !follow {
		// Just dump the entire file
		_, err := io.Copy(writer, file)
		return err
	}

	// Follow mode - start from end and watch for new lines
	// First, dump existing content
	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	// Then follow new lines
	reader := bufio.NewReader(file)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				// No new data, wait a bit
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if err != nil {
				return err
			}
			fmt.Fprint(writer, line)
		}
	}
}

// ClearLogs clears the log file
func (pm *ProcessManager) ClearLogs() error {
	return os.Truncate(pm.LogFile(), 0)
}

// ============================================================
// PID File Helpers
// ============================================================

func (pm *ProcessManager) savePID(pid int) error {
	return os.WriteFile(pm.PIDFile(), []byte(strconv.Itoa(pid)), 0644)
}

func (pm *ProcessManager) loadPID() (int, error) {
	data, err := os.ReadFile(pm.PIDFile())
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

