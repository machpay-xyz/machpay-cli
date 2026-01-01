// ============================================================
// TUI Prompts - Simple stdin prompts with Lipgloss styling
// ============================================================
//
// Provides reusable prompt components:
// - Select: Numbered selection from list
// - TextInput: Text input with validation
// - Confirm: Y/n confirmation
//
// No external TUI dependencies - uses bufio.Reader for input
// and Lipgloss for styling.
//
// ============================================================

package tui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ============================================================
// Prompt Styles
// ============================================================

var (
	questionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff"))

	optionNumberStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#10b981")).
				Bold(true)

	optionLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff"))

	optionDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#71717a")).
			PaddingLeft(4)

	promptArrowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#10b981")).
				Bold(true)

	inputErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ef4444"))

	sectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3f3f46"))

	bannerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#10b981")).
			Padding(0, 2).
			Foreground(lipgloss.Color("#10b981")).
			Bold(true)
)

// ============================================================
// SelectOption
// ============================================================

// SelectOption represents a selectable option
type SelectOption struct {
	Label       string
	Description string
	Value       string
}

// ============================================================
// Select - Numbered selection prompt
// ============================================================

// Select displays a numbered list and returns the selected option
func Select(question string, options []SelectOption) (SelectOption, error) {
	fmt.Println()
	fmt.Println(questionStyle.Render(question))
	fmt.Println()

	for i, opt := range options {
		num := optionNumberStyle.Render(fmt.Sprintf("[%d]", i+1))
		label := optionLabelStyle.Render(opt.Label)
		fmt.Printf("  %s %s\n", num, label)
		if opt.Description != "" {
			fmt.Println(optionDescStyle.Render(opt.Description))
		}
	}

	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	for {
		prompt := promptArrowStyle.Render(fmt.Sprintf("Enter choice [1-%d]: ", len(options)))
		fmt.Print(prompt)

		input, err := reader.ReadString('\n')
		if err != nil {
			return SelectOption{}, fmt.Errorf("read input: %w", err)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(options) {
			fmt.Println(inputErrorStyle.Render(
				fmt.Sprintf("  Please enter a number between 1 and %d", len(options))))
			continue
		}

		return options[choice-1], nil
	}
}

// ============================================================
// TextInput - Text input with validation
// ============================================================

// TextInput prompts for text input with optional validation
func TextInput(question string, placeholder string, validator func(string) error) (string, error) {
	fmt.Println()
	fmt.Println(questionStyle.Render(question))
	if placeholder != "" {
		fmt.Println(optionDescStyle.Render("Example: " + placeholder))
	}
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(promptArrowStyle.Render("> "))

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("read input: %w", err)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println(inputErrorStyle.Render("  Input cannot be empty"))
			continue
		}

		if validator != nil {
			if err := validator(input); err != nil {
				fmt.Println(inputErrorStyle.Render("  " + err.Error()))
				continue
			}
		}

		return input, nil
	}
}

// TextInputOptional allows empty input with a default value
func TextInputOptional(question string, defaultValue string) (string, error) {
	fmt.Println()
	fmt.Println(questionStyle.Render(question))
	if defaultValue != "" {
		fmt.Println(optionDescStyle.Render(
			fmt.Sprintf("Default: %s (press Enter to use)", defaultValue)))
	}
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptArrowStyle.Render("> "))

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue, nil
	}

	return input, nil
}

// ============================================================
// Confirm - Y/n confirmation
// ============================================================

// Confirm prompts for yes/no confirmation
func Confirm(question string, defaultYes bool) (bool, error) {
	fmt.Println()

	var prompt string
	if defaultYes {
		prompt = question + " [Y/n]: "
	} else {
		prompt = question + " [y/N]: "
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(questionStyle.Render(prompt))

		input, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("read input: %w", err)
		}

		input = strings.TrimSpace(strings.ToLower(input))

		if input == "" {
			return defaultYes, nil
		}

		switch input {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			fmt.Println(inputErrorStyle.Render("  Please enter 'y' or 'n'"))
		}
	}
}

// ============================================================
// Decorative Elements
// ============================================================

// PrintBanner displays a styled banner box
func PrintBanner(title string) {
	fmt.Println()
	fmt.Println(bannerStyle.Render(title))
	fmt.Println()
}

// PrintSection prints a section divider line
func PrintSection() {
	fmt.Println()
	line := strings.Repeat("─", 60)
	fmt.Println(sectionStyle.Render(line))
	fmt.Println()
}

// PrintCodeBlock prints code in a styled box
func PrintCodeBlock(code string) {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3f3f46")).
		Padding(0, 1).
		Foreground(lipgloss.Color("#a1a1aa"))

	fmt.Println(style.Render(code))
}

// PrintKeyValue prints a key-value pair with styling
func PrintKeyValue(key, value string) {
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#71717a")).
		Width(12)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10b981"))

	fmt.Printf("  %s %s\n", keyStyle.Render(key+":"), valueStyle.Render(value))
}

