// ============================================================
// Progress Bar - Terminal progress indicators
// ============================================================
//
// Provides:
// - ProgressReader: Wrap io.Reader to show download progress
// - Spinner: Indeterminate progress indicator
//
// No external dependencies - uses only stdlib.
//
// ============================================================

package gateway

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

// ============================================================
// ProgressReader - Progress bar for known-size downloads
// ============================================================

// ProgressReader wraps an io.Reader to report download progress
type ProgressReader struct {
	Reader  io.Reader
	Total   int64
	Current int64
	Writer  io.Writer

	mu      sync.Mutex
	lastPct int
	started time.Time
}

// NewProgressReader creates a new progress reader
func NewProgressReader(reader io.Reader, total int64, writer io.Writer) *ProgressReader {
	return &ProgressReader{
		Reader:  reader,
		Total:   total,
		Writer:  writer,
		started: time.Now(),
	}
}

// Read implements io.Reader and updates progress
func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)

	pr.mu.Lock()
	pr.Current += int64(n)
	pct := int(float64(pr.Current) / float64(pr.Total) * 100)

	// Only update on percentage change (avoid flickering)
	if pct != pr.lastPct {
		pr.lastPct = pct
		pr.render()
	}
	pr.mu.Unlock()

	return n, err
}

// render draws the progress bar
func (pr *ProgressReader) render() {
	width := 40
	filled := int(float64(width) * float64(pr.Current) / float64(pr.Total))
	if filled > width {
		filled = width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	// Calculate speed
	elapsed := time.Since(pr.started).Seconds()
	speed := float64(pr.Current) / elapsed / 1024 / 1024 // MB/s

	// Format sizes
	currentMB := float64(pr.Current) / 1024 / 1024
	totalMB := float64(pr.Total) / 1024 / 1024

	fmt.Fprintf(pr.Writer, "\r  %s %3d%%  %.1f/%.1f MB  (%.1f MB/s)",
		bar, pr.lastPct, currentMB, totalMB, speed)

	// Newline on completion
	if pr.Current >= pr.Total {
		fmt.Fprintln(pr.Writer)
	}
}

// Finish ensures the progress bar is complete
func (pr *ProgressReader) Finish() {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	if pr.lastPct < 100 {
		pr.lastPct = 100
		pr.Current = pr.Total
		pr.render()
	}
}

// ============================================================
// Spinner - Indeterminate progress indicator
// ============================================================

// Spinner shows a spinning indicator for operations without known size
type Spinner struct {
	Writer  io.Writer
	Message string
	done    chan struct{}
	wg      sync.WaitGroup
	mu      sync.Mutex
	running bool
}

// NewSpinner creates a new spinner
func NewSpinner(writer io.Writer, message string) *Spinner {
	return &Spinner{
		Writer:  writer,
		Message: message,
		done:    make(chan struct{}),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.done = make(chan struct{})
	// Copy values to avoid race with Stop()
	writer := s.Writer
	message := s.Message
	done := s.done
	s.wg.Add(1)
	s.mu.Unlock()

	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	go func() {
		defer s.wg.Done()
		i := 0
		for {
			select {
			case <-done:
				return
			default:
				fmt.Fprintf(writer, "\r  %s %s", frames[i%len(frames)], message)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner with a status
func (s *Spinner) Stop(success bool) {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}

	close(s.done)
	s.running = false
	s.mu.Unlock()

	// Wait for goroutine to exit before writing final message
	s.wg.Wait()

	if success {
		fmt.Fprintf(s.Writer, "\r  ✓ %s\n", s.Message)
	} else {
		fmt.Fprintf(s.Writer, "\r  ✗ %s\n", s.Message)
	}
}

// StopWithMessage stops the spinner with a custom message
func (s *Spinner) StopWithMessage(success bool, message string) {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}

	close(s.done)
	s.running = false
	s.mu.Unlock()

	// Wait for goroutine to exit before writing final message
	s.wg.Wait()

	if success {
		fmt.Fprintf(s.Writer, "\r  ✓ %s\n", message)
	} else {
		fmt.Fprintf(s.Writer, "\r  ✗ %s\n", message)
	}
}

// ============================================================
// Step Progress - For multi-step operations
// ============================================================

// StepProgress tracks progress through multiple steps
type StepProgress struct {
	Writer   io.Writer
	Total    int
	Current  int
	StepName string
}

// NewStepProgress creates a new step progress tracker
func NewStepProgress(writer io.Writer, total int) *StepProgress {
	return &StepProgress{
		Writer: writer,
		Total:  total,
	}
}

// Step advances to the next step
func (sp *StepProgress) Step(name string) {
	sp.Current++
	sp.StepName = name
	fmt.Fprintf(sp.Writer, "  [%d/%d] %s\n", sp.Current, sp.Total, name)
}

// Complete marks all steps as done
func (sp *StepProgress) Complete() {
	fmt.Fprintf(sp.Writer, "  ✓ All %d steps complete\n", sp.Total)
}

