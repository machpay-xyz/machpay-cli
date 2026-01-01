package gateway

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

func TestProgressReader_Read(t *testing.T) {
	data := make([]byte, 1000)
	for i := range data {
		data[i] = byte(i % 256)
	}

	reader := bytes.NewReader(data)
	output := &bytes.Buffer{}

	pr := NewProgressReader(reader, int64(len(data)), output)

	// Read all data
	result := make([]byte, 0, len(data))
	buf := make([]byte, 100)

	for {
		n, err := pr.Read(buf)
		if n > 0 {
			result = append(result, buf[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Read error: %v", err)
		}
	}

	// Verify data integrity
	if !bytes.Equal(result, data) {
		t.Error("Data mismatch after reading through ProgressReader")
	}

	// Verify progress was written
	outputStr := output.String()
	if !strings.Contains(outputStr, "100%") {
		t.Error("Expected 100% in progress output")
	}
}

func TestProgressReader_PartialRead(t *testing.T) {
	data := make([]byte, 1000)
	reader := bytes.NewReader(data)
	output := &bytes.Buffer{}

	pr := NewProgressReader(reader, int64(len(data)), output)

	// Read half
	buf := make([]byte, 500)
	n, err := pr.Read(buf)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	if n != 500 {
		t.Errorf("Expected to read 500 bytes, got %d", n)
	}

	// Should show 50%
	outputStr := output.String()
	if !strings.Contains(outputStr, "50%") {
		t.Errorf("Expected 50%% in output, got: %s", outputStr)
	}
}

func TestSpinner_StartStop(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := NewSpinner(output, "Loading")

	spinner.Start()
	time.Sleep(200 * time.Millisecond) // Let it spin a bit
	spinner.Stop(true)

	outputStr := output.String()
	if !strings.Contains(outputStr, "Loading") {
		t.Error("Expected message in output")
	}
	if !strings.Contains(outputStr, "✓") {
		t.Error("Expected success checkmark")
	}
}

func TestSpinner_StopFailure(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := NewSpinner(output, "Downloading")

	spinner.Start()
	time.Sleep(100 * time.Millisecond)
	spinner.Stop(false)

	outputStr := output.String()
	if !strings.Contains(outputStr, "✗") {
		t.Error("Expected failure X mark")
	}
}

func TestSpinner_DoubleStart(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := NewSpinner(output, "Test")

	spinner.Start()
	spinner.Start() // Should not panic or create duplicate goroutines
	time.Sleep(100 * time.Millisecond)
	spinner.Stop(true)
}

func TestSpinner_DoubleStop(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := NewSpinner(output, "Test")

	spinner.Start()
	time.Sleep(100 * time.Millisecond)
	spinner.Stop(true)
	spinner.Stop(true) // Should not panic
}

func TestStepProgress(t *testing.T) {
	output := &bytes.Buffer{}
	sp := NewStepProgress(output, 3)

	sp.Step("First step")
	sp.Step("Second step")
	sp.Step("Third step")
	sp.Complete()

	outputStr := output.String()

	if !strings.Contains(outputStr, "[1/3]") {
		t.Error("Expected [1/3] in output")
	}
	if !strings.Contains(outputStr, "[2/3]") {
		t.Error("Expected [2/3] in output")
	}
	if !strings.Contains(outputStr, "[3/3]") {
		t.Error("Expected [3/3] in output")
	}
	if !strings.Contains(outputStr, "All 3 steps complete") {
		t.Error("Expected completion message")
	}
}

