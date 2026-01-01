// ============================================================
// Setup Command - Interactive Setup Wizard
// ============================================================
//
// Usage: machpay setup [--non-interactive]
//
// Guides users through initial configuration:
// - Role selection (Agent or Vendor)
// - Network selection (Devnet or Mainnet)
// - Wallet generation or import
// - Role-specific configuration
//
// ============================================================

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/machpay-xyz/machpay-cli/internal/auth"
	"github.com/machpay-xyz/machpay-cli/internal/config"
	"github.com/machpay-xyz/machpay-cli/internal/tui"
	"github.com/machpay-xyz/machpay-cli/internal/wallet"
)

var setupNonInteractive bool

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup wizard",
	Long: `Configure your MachPay CLI as an agent or vendor.

This wizard will guide you through:
  - Choosing your role (Agent or Vendor)
  - Selecting network (Devnet or Mainnet)
  - Setting up your wallet
  - Generating API keys or configuring your service

Non-interactive mode for CI/CD:
  MACHPAY_ROLE=agent MACHPAY_NETWORK=devnet machpay setup --non-interactive`,
	RunE: runSetup,
}

func init() {
	setupCmd.Flags().BoolVar(&setupNonInteractive, "non-interactive", false,
		"Use environment variables instead of prompts (for CI/CD)")
}

func runSetup(cmd *cobra.Command, args []string) error {
	// 1. Check if logged in
	if !auth.IsLoggedIn() {
		tui.PrintError("Not logged in")
		fmt.Println(tui.Muted("  Run 'machpay login' first to authenticate."))
		return fmt.Errorf("authentication required")
	}

	// 2. Show banner
	printSetupBanner()

	// 3. Check if already configured
	cfg := config.Get()
	if cfg.Role != "" {
		reconfigure, err := tui.Confirm(
			fmt.Sprintf("Already configured as %s. Reconfigure?", tui.Primary(cfg.Role)),
			false,
		)
		if err != nil {
			return err
		}
		if !reconfigure {
			fmt.Println(tui.Muted("Setup cancelled."))
			return nil
		}
	}

	// 4. Handle non-interactive mode
	if setupNonInteractive {
		return runNonInteractiveSetup()
	}

	// 5. Run interactive wizard
	return runInteractiveSetup()
}

func printSetupBanner() {
	tui.PrintBanner("MachPay Setup Wizard")

	user := auth.GetUser()
	if user != nil {
		fmt.Printf("  Logged in as: %s\n", tui.Primary(user.Email))
	}
}

func runInteractiveSetup() error {
	// Step 1: Role selection
	tui.PrintSection()
	role, err := tui.Select("What do you want to do?", []tui.SelectOption{
		{
			Label:       "Run an AI Agent",
			Description: "I want to use APIs and pay for services",
			Value:       "agent",
		},
		{
			Label:       "Run a Vendor Node",
			Description: "I want to sell my APIs and earn money",
			Value:       "vendor",
		},
	})
	if err != nil {
		return err
	}

	// Step 2: Network selection
	tui.PrintSection()
	network, err := tui.Select("Select network:", []tui.SelectOption{
		{
			Label:       "Devnet",
			Description: "Testing network - free tokens, no real money",
			Value:       "devnet",
		},
		{
			Label:       "Mainnet",
			Description: "Production network - real USDC transactions",
			Value:       "mainnet",
		},
	})
	if err != nil {
		return err
	}

	// Warning for mainnet
	if network.Value == "mainnet" {
		tui.PrintSection()
		fmt.Println(tui.Warning("⚠️  MAINNET WARNING"))
		fmt.Println()
		fmt.Println(tui.Muted("  You are about to configure for Mainnet."))
		fmt.Println(tui.Muted("  All transactions will use real USDC."))

		proceed, err := tui.Confirm("Continue with Mainnet?", false)
		if err != nil {
			return err
		}
		if !proceed {
			fmt.Println(tui.Muted("  Switching to Devnet..."))
			network.Value = "devnet"
		}
	}

	// Step 3: Branch by role
	tui.PrintSection()
	switch role.Value {
	case "agent":
		return setupAgent(network.Value)
	case "vendor":
		return setupVendor(network.Value)
	}

	return nil
}

// ============================================================
// Agent Setup
// ============================================================

func setupAgent(network string) error {
	fmt.Println(tui.Bold("Agent Setup"))
	fmt.Println()

	// 1. Wallet setup
	kp, err := promptWallet()
	if err != nil {
		return err
	}

	// 2. Save config
	cfg := config.Get()
	cfg.Role = "agent"
	cfg.Network = network
	cfg.Wallet.KeypairPath = filepath.Join(config.GetDir(), "wallet.json")
	cfg.Wallet.PublicKey = kp.PublicKeyBase58()

	if err := config.Save(); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	// 3. Show success
	tui.PrintSection()
	printAgentSuccess(kp.PublicKeyBase58(), network)

	return nil
}

func printAgentSuccess(address, network string) {
	tui.PrintSuccess("Agent setup complete!")
	fmt.Println()

	tui.PrintKeyValue("Wallet", address)
	tui.PrintKeyValue("Network", network)
	tui.PrintKeyValue("Config", config.GetPath())

	fmt.Println()
	fmt.Println(tui.Bold("Quick Start:"))
	fmt.Println()

	tui.PrintCodeBlock(`pip install machpay

from machpay import MachPay
client = MachPay()
response = client.call("weather-api", "/forecast")`)

	fmt.Println()
	fmt.Printf("  Next: Fund your wallet at %s\n",
		tui.Primary(config.GetConsoleURL()+"/agent/finance"))
	fmt.Println()
}

// ============================================================
// Vendor Setup
// ============================================================

func setupVendor(network string) error {
	fmt.Println(tui.Bold("Vendor Setup"))
	fmt.Println()

	// 1. Collect service info
	serviceName, err := tui.TextInput("Service name", "My LLM API", nil)
	if err != nil {
		return err
	}

	category, err := tui.Select("Category:", []tui.SelectOption{
		{Label: "AI/ML", Value: "ai"},
		{Label: "Data", Value: "data"},
		{Label: "Finance", Value: "finance"},
		{Label: "Compute", Value: "compute"},
		{Label: "Other", Value: "other"},
	})
	if err != nil {
		return err
	}

	tui.PrintSection()

	// 2. Service configuration
	upstreamURL, err := tui.TextInput("Upstream API URL", "http://localhost:11434", validateURL)
	if err != nil {
		return err
	}

	priceStr, err := tui.TextInput("Price per request (USDC)", "0.001", validatePrice)
	if err != nil {
		return err
	}

	price, _ := strconv.ParseFloat(priceStr, 64)

	tui.PrintSection()

	// 3. Wallet setup
	kp, err := promptWallet()
	if err != nil {
		return err
	}

	// 4. Save config
	cfg := config.Get()
	cfg.Role = "vendor"
	cfg.Network = network
	cfg.Wallet.KeypairPath = filepath.Join(config.GetDir(), "wallet.json")
	cfg.Wallet.PublicKey = kp.PublicKeyBase58()
	cfg.Vendor.UpstreamURL = upstreamURL
	cfg.Vendor.PricePerRequest = price

	if err := config.Save(); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	// 5. Show success
	tui.PrintSection()
	printVendorSuccess(serviceName, category.Value, upstreamURL, priceStr, network)

	// 6. Offer to start gateway (placeholder for Phase 3)
	fmt.Println()
	startNow, err := tui.Confirm("Start gateway now?", true)
	if err != nil {
		return err
	}
	if startNow {
		fmt.Println()
		fmt.Println(tui.Info("ℹ") + " Gateway will be available in a future update.")
		fmt.Println(tui.Muted("  Run 'machpay serve' when available."))
	}

	return nil
}

func printVendorSuccess(name, category, upstream, price, network string) {
	tui.PrintSuccess("Vendor setup complete!")
	fmt.Println()

	tui.PrintKeyValue("Service", name+" ("+category+")")
	tui.PrintKeyValue("Upstream", upstream)
	tui.PrintKeyValue("Price", price+" USDC/request")
	tui.PrintKeyValue("Network", network)
	tui.PrintKeyValue("Config", config.GetPath())

	fmt.Println()
	fmt.Println(tui.Bold("Next Steps:"))
	fmt.Println("  1. Run 'machpay serve' to start your gateway")
	fmt.Println("  2. Your API will be available for agents to discover")
	fmt.Println()
}

// ============================================================
// Wallet Prompt
// ============================================================

func promptWallet() (*wallet.Keypair, error) {
	choice, err := tui.Select("Wallet setup:", []tui.SelectOption{
		{Label: "Generate new wallet", Description: "Recommended for new users", Value: "generate"},
		{Label: "Import existing keypair", Description: "Use existing Solana wallet", Value: "import"},
	})
	if err != nil {
		return nil, err
	}

	switch choice.Value {
	case "generate":
		return generateNewWallet()
	case "import":
		return importExistingWallet()
	}
	return nil, fmt.Errorf("invalid choice")
}

func generateNewWallet() (*wallet.Keypair, error) {
	fmt.Println()
	fmt.Println(tui.Muted("  Generating new wallet..."))

	kp, err := wallet.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate wallet: %w", err)
	}

	walletPath := filepath.Join(config.GetDir(), "wallet.json")
	if err := kp.SaveToFile(walletPath); err != nil {
		return nil, fmt.Errorf("save wallet: %w", err)
	}

	fmt.Println()
	tui.PrintSuccess("Generated new wallet")
	fmt.Println()
	tui.PrintKeyValue("Address", kp.PublicKeyBase58())
	tui.PrintKeyValue("Saved to", walletPath)
	fmt.Println()
	fmt.Println(tui.Warning("⚠️  BACKUP THIS FILE! It contains your private key."))

	return kp, nil
}

func importExistingWallet() (*wallet.Keypair, error) {
	path, err := tui.TextInput("Path to keypair file", "~/.config/solana/id.json", func(s string) error {
		// Expand ~ to home directory
		expandedPath := s
		if strings.HasPrefix(s, "~") {
			home, _ := os.UserHomeDir()
			expandedPath = filepath.Join(home, s[1:])
		}
		if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", s)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Load keypair
	kp, err := wallet.LoadFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("load keypair: %w", err)
	}

	// Copy to MachPay directory
	destPath := filepath.Join(config.GetDir(), "wallet.json")
	if err := kp.SaveToFile(destPath); err != nil {
		return nil, fmt.Errorf("save keypair: %w", err)
	}

	fmt.Println()
	tui.PrintSuccess("Imported wallet")
	tui.PrintKeyValue("Address", kp.PublicKeyBase58())

	return kp, nil
}

// ============================================================
// Non-Interactive Setup (CI/CD)
// ============================================================

func runNonInteractiveSetup() error {
	role := os.Getenv("MACHPAY_ROLE")
	network := os.Getenv("MACHPAY_NETWORK")
	walletPath := os.Getenv("MACHPAY_WALLET_PATH")

	if role == "" {
		return fmt.Errorf("MACHPAY_ROLE environment variable required (agent or vendor)")
	}
	if role != "agent" && role != "vendor" {
		return fmt.Errorf("MACHPAY_ROLE must be 'agent' or 'vendor'")
	}
	if network == "" {
		network = "devnet"
	}

	fmt.Printf("Configuring as %s on %s...\n", tui.Primary(role), tui.Primary(network))

	cfg := config.Get()
	cfg.Role = role
	cfg.Network = network

	// Handle wallet
	if walletPath != "" {
		kp, err := wallet.LoadFromFile(walletPath)
		if err != nil {
			return fmt.Errorf("load wallet: %w", err)
		}
		cfg.Wallet.KeypairPath = walletPath
		cfg.Wallet.PublicKey = kp.PublicKeyBase58()
		fmt.Printf("Wallet: %s\n", tui.Primary(kp.PublicKeyBase58()))
	} else {
		// Generate new wallet
		kp, err := wallet.Generate()
		if err != nil {
			return fmt.Errorf("generate wallet: %w", err)
		}
		walletPath := filepath.Join(config.GetDir(), "wallet.json")
		if err := kp.SaveToFile(walletPath); err != nil {
			return fmt.Errorf("save wallet: %w", err)
		}
		cfg.Wallet.KeypairPath = walletPath
		cfg.Wallet.PublicKey = kp.PublicKeyBase58()
		fmt.Printf("Generated wallet: %s\n", tui.Primary(kp.PublicKeyBase58()))
	}

	// Vendor-specific config
	if role == "vendor" {
		cfg.Vendor.UpstreamURL = os.Getenv("MACHPAY_UPSTREAM_URL")
		if priceStr := os.Getenv("MACHPAY_PRICE"); priceStr != "" {
			price, err := strconv.ParseFloat(priceStr, 64)
			if err != nil {
				return fmt.Errorf("invalid MACHPAY_PRICE: %w", err)
			}
			cfg.Vendor.PricePerRequest = price
		}
	}

	if err := config.Save(); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	fmt.Println()
	tui.PrintSuccess("Setup complete")
	fmt.Printf("  Config saved to: %s\n", tui.Muted(config.GetPath()))

	return nil
}

// ============================================================
// Validators
// ============================================================

func validateURL(s string) error {
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}
	return nil
}

func validatePrice(s string) error {
	price, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("invalid number")
	}
	if price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if price > 1000 {
		return fmt.Errorf("price seems too high (max 1000 USDC)")
	}
	return nil
}

