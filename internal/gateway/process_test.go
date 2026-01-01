package gateway

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestNewProcessManager(t *testing.T) {
	pm := NewProcessManager("/usr/bin/test", 8402, "http://localhost:11434")

	if pm.binaryPath != "/usr/bin/test" {
		t.Errorf("binaryPath = %s, want /usr/bin/test", pm.binaryPath)
	}
	if pm.port != 8402 {
		t.Errorf("port = %d, want 8402", pm.port)
	}
	if pm.upstream != "http://localhost:11434" {
		t.Errorf("upstream = %s, want http://localhost:11434", pm.upstream)
	}
}

func TestProcessManager_PIDFile(t *testing.T) {
	pm := NewProcessManager("/usr/bin/test", 8402, "")
	pidFile := pm.PIDFile()

	if filepath.Base(pidFile) != "gateway.pid" {
		t.Errorf("PID file = %s, want gateway.pid", filepath.Base(pidFile))
	}
}

func TestProcessManager_LogFile(t *testing.T) {
	pm := NewProcessManager("/usr/bin/test", 8402, "")
	logFile := pm.LogFile()

	if filepath.Base(logFile) != "gateway.log" {
		t.Errorf("Log file = %s, want gateway.log", filepath.Base(logFile))
	}
}

func TestProcessManager_IsRunning_NoPIDFile(t *testing.T) {
	tmpDir := t.TempDir()
	pm := &ProcessManager{
		binaryPath: "/usr/bin/test",
		configDir:  tmpDir,
		port:       8402,
	}

	if pm.IsRunning() {
		t.Error("Expected IsRunning to return false when no PID file exists")
	}
}

func TestProcessManager_IsRunning_InvalidPID(t *testing.T) {
	tmpDir := t.TempDir()
	pm := &ProcessManager{
		binaryPath: "/usr/bin/test",
		configDir:  tmpDir,
		port:       8402,
	}

	// Write invalid PID
	if err := os.WriteFile(pm.PIDFile(), []byte("not-a-number"), 0644); err != nil {
		t.Fatalf("Write PID file: %v", err)
	}

	if pm.IsRunning() {
		t.Error("Expected IsRunning to return false for invalid PID")
	}
}

func TestProcessManager_IsRunning_DeadProcess(t *testing.T) {
	tmpDir := t.TempDir()
	pm := &ProcessManager{
		binaryPath: "/usr/bin/test",
		configDir:  tmpDir,
		port:       8402,
	}

	// Write a PID that doesn't exist (very high number)
	if err := os.WriteFile(pm.PIDFile(), []byte("999999999"), 0644); err != nil {
		t.Fatalf("Write PID file: %v", err)
	}

	if pm.IsRunning() {
		t.Error("Expected IsRunning to return false for dead process")
	}
}

func TestProcessManager_BuildArgs(t *testing.T) {
	tests := []struct {
		name     string
		port     int
		upstream string
		debug    bool
		wantLen  int
	}{
		{
			name:    "empty",
			wantLen: 0,
		},
		{
			name:    "port only",
			port:    8402,
			wantLen: 2, // --port 8402
		},
		{
			name:     "upstream only",
			upstream: "http://localhost:11434",
			wantLen:  2, // --upstream http://localhost:11434
		},
		{
			name:  "debug only",
			debug: true,
			wantLen: 1, // --debug
		},
		{
			name:     "all",
			port:     8402,
			upstream: "http://localhost:11434",
			debug:    true,
			wantLen:  5, // --port 8402 --upstream http://localhost:11434 --debug
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &ProcessManager{
				port:     tt.port,
				upstream: tt.upstream,
				debug:    tt.debug,
			}

			args := pm.buildArgs()
			if len(args) != tt.wantLen {
				t.Errorf("buildArgs() len = %d, want %d, args = %v", len(args), tt.wantLen, args)
			}
		})
	}
}

func TestProcessManager_SaveLoadPID(t *testing.T) {
	tmpDir := t.TempDir()
	pm := &ProcessManager{
		configDir: tmpDir,
	}

	// Save PID
	expectedPID := 12345
	if err := pm.savePID(expectedPID); err != nil {
		t.Fatalf("savePID: %v", err)
	}

	// Load PID
	pid, err := pm.loadPID()
	if err != nil {
		t.Fatalf("loadPID: %v", err)
	}

	if pid != expectedPID {
		t.Errorf("loadPID = %d, want %d", pid, expectedPID)
	}
}

func TestProcessManager_GetPID_NotRunning(t *testing.T) {
	tmpDir := t.TempDir()
	pm := &ProcessManager{
		configDir: tmpDir,
	}

	_, err := pm.GetPID()
	if err == nil {
		t.Error("Expected error when PID file doesn't exist")
	}
}

func TestProcessManager_Stop_NotRunning(t *testing.T) {
	tmpDir := t.TempDir()
	pm := &ProcessManager{
		configDir: tmpDir,
		port:      8402,
	}

	err := pm.Stop()
	if err != ErrNotRunning {
		t.Errorf("Stop() error = %v, want ErrNotRunning", err)
	}
}

func TestProcessManager_SetDebug(t *testing.T) {
	pm := NewProcessManager("/usr/bin/test", 8402, "")

	if pm.debug {
		t.Error("Expected debug to be false initially")
	}

	pm.SetDebug(true)

	if !pm.debug {
		t.Error("Expected debug to be true after SetDebug(true)")
	}
}

func TestProcessManager_SetPort(t *testing.T) {
	pm := NewProcessManager("/usr/bin/test", 8402, "")

	pm.SetPort(9000)

	if pm.port != 9000 {
		t.Errorf("port = %d, want 9000", pm.port)
	}
}

func TestProcessManager_SetUpstream(t *testing.T) {
	pm := NewProcessManager("/usr/bin/test", 8402, "")

	pm.SetUpstream("http://newhost:8080")

	if pm.upstream != "http://newhost:8080" {
		t.Errorf("upstream = %s, want http://newhost:8080", pm.upstream)
	}
}

func TestProcessManager_ClearLogs(t *testing.T) {
	tmpDir := t.TempDir()
	pm := &ProcessManager{
		configDir: tmpDir,
	}

	// Create log file with content
	logPath := pm.LogFile()
	if err := os.WriteFile(logPath, []byte("test log content"), 0644); err != nil {
		t.Fatalf("Write log file: %v", err)
	}

	// Clear logs
	if err := pm.ClearLogs(); err != nil {
		t.Fatalf("ClearLogs: %v", err)
	}

	// Check file is empty
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Stat log file: %v", err)
	}

	if info.Size() != 0 {
		t.Errorf("Log file size = %d, want 0", info.Size())
	}
}

func TestProcessManager_Start_AlreadyRunning(t *testing.T) {
	tmpDir := t.TempDir()
	pm := &ProcessManager{
		binaryPath: "/usr/bin/test",
		configDir:  tmpDir,
		port:       8402,
	}

	// Write current process PID (which is running)
	if err := os.WriteFile(pm.PIDFile(), []byte(strconv.Itoa(os.Getpid())), 0644); err != nil {
		t.Fatalf("Write PID file: %v", err)
	}

	err := pm.Start()
	if err != ErrAlreadyRunning {
		t.Errorf("Start() error = %v, want ErrAlreadyRunning", err)
	}
}

// Note: Full integration tests for Start/Stop require a real binary
// and are covered by integration tests

