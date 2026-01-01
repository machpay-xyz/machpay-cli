package tui

import (
	"testing"
)

func TestSelectOption(t *testing.T) {
	options := []SelectOption{
		{Label: "Option 1", Description: "First option", Value: "one"},
		{Label: "Option 2", Description: "Second option", Value: "two"},
	}

	if len(options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(options))
	}

	if options[0].Value != "one" {
		t.Errorf("Expected value 'one', got '%s'", options[0].Value)
	}
}

// Note: Interactive prompts (Select, TextInput, Confirm) require
// stdin mocking for proper testing. These would be tested via
// integration tests or with a testing harness that can inject
// input.

func TestPrintBanner(t *testing.T) {
	// Just ensure it doesn't panic
	PrintBanner("Test Banner")
}

func TestPrintSection(t *testing.T) {
	// Just ensure it doesn't panic
	PrintSection()
}

func TestPrintCodeBlock(t *testing.T) {
	// Just ensure it doesn't panic
	PrintCodeBlock("echo 'hello world'")
}

func TestPrintKeyValue(t *testing.T) {
	// Just ensure it doesn't panic
	PrintKeyValue("Key", "Value")
}

