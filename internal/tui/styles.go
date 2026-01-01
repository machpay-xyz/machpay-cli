// ============================================================
// TUI Styles - Terminal UI Styling with Lipgloss
// ============================================================

package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	colorsEnabled = true

	// Colors
	colorPrimary   = lipgloss.Color("#10b981") // Emerald
	colorSecondary = lipgloss.Color("#71717a") // Zinc
	colorError     = lipgloss.Color("#ef4444") // Red
	colorWarning   = lipgloss.Color("#f59e0b") // Amber
	colorSuccess   = lipgloss.Color("#22c55e") // Green
	colorInfo      = lipgloss.Color("#3b82f6") // Blue

	// Styles
	styleBold = lipgloss.NewStyle().Bold(true)

	styleSuccess = lipgloss.NewStyle().
			Foreground(colorSuccess)

	styleError = lipgloss.NewStyle().
			Foreground(colorError)

	styleWarning = lipgloss.NewStyle().
			Foreground(colorWarning)

	styleInfo = lipgloss.NewStyle().
			Foreground(colorInfo)

	styleMuted = lipgloss.NewStyle().
			Foreground(colorSecondary)

	stylePrimary = lipgloss.NewStyle().
			Foreground(colorPrimary)

	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			MarginBottom(1)

	styleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSecondary).
			Padding(1, 2)
)

// DisableColors turns off colored output
func DisableColors() {
	colorsEnabled = false
}

// Bold returns bold text
func Bold(s string) string {
	if !colorsEnabled {
		return s
	}
	return styleBold.Render(s)
}

// Success returns green text
func Success(s string) string {
	if !colorsEnabled {
		return s
	}
	return styleSuccess.Render(s)
}

// Error returns red text
func Error(s string) string {
	if !colorsEnabled {
		return s
	}
	return styleError.Render(s)
}

// Warning returns amber text
func Warning(s string) string {
	if !colorsEnabled {
		return s
	}
	return styleWarning.Render(s)
}

// Info returns blue text
func Info(s string) string {
	if !colorsEnabled {
		return s
	}
	return styleInfo.Render(s)
}

// Muted returns gray text
func Muted(s string) string {
	if !colorsEnabled {
		return s
	}
	return styleMuted.Render(s)
}

// Primary returns emerald text
func Primary(s string) string {
	if !colorsEnabled {
		return s
	}
	return stylePrimary.Render(s)
}

// Header returns a styled header
func Header(s string) string {
	if !colorsEnabled {
		return fmt.Sprintf("=== %s ===", s)
	}
	return styleHeader.Render(s)
}

// Box wraps content in a styled box
func Box(content string) string {
	if !colorsEnabled {
		return content
	}
	return styleBox.Render(content)
}

// SuccessIcon returns a styled checkmark
func SuccessIcon() string {
	return Success("✓")
}

// ErrorIcon returns a styled X
func ErrorIcon() string {
	return Error("✗")
}

// WarningIcon returns a styled warning
func WarningIcon() string {
	return Warning("⚠")
}

// InfoIcon returns a styled info icon
func InfoIcon() string {
	return Info("ℹ")
}

// PrintSuccess prints a success message
func PrintSuccess(msg string) {
	fmt.Printf("%s %s\n", SuccessIcon(), msg)
}

// PrintError prints an error message
func PrintError(msg string) {
	fmt.Printf("%s %s\n", ErrorIcon(), msg)
}

// PrintWarning prints a warning message
func PrintWarning(msg string) {
	fmt.Printf("%s %s\n", WarningIcon(), msg)
}

// PrintInfo prints an info message
func PrintInfo(msg string) {
	fmt.Printf("%s %s\n", InfoIcon(), msg)
}

