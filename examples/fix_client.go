package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/kolobublik/limit-order-book/internal/data/fix"
	"github.com/quickfixgo/quickfix"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Data struct {
		Fix fix.Config `yaml:"fix"`
	} `yaml:"data"`
}

func main() {
	// Load application config
	appCfg, err := os.Open(filepath.Join("configs", "config.yaml"))
	if err != nil {
		log.Fatalf("Error opening config.yaml: %v", err)
	}
	defer appCfg.Close()

	var config Config
	if err := yaml.NewDecoder(appCfg).Decode(&config); err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	// Load FIX settings
	fixCfg, err := os.Open(filepath.Join("configs", "fix.cfg"))
	if err != nil {
		log.Fatalf("Error opening fix.cfg: %v", err)
	}
	defer fixCfg.Close()

	settings, err := quickfix.ParseSettings(fixCfg)
	if err != nil {
		log.Fatalf("Error parsing settings: %v", err)
	}

	// Create and start FIX client
	client, err := fix.NewClient(settings, &config.Data.Fix)
	if err != nil {
		log.Fatalf("Error creating FIX client: %v", err)
	}

	if err := client.Start(); err != nil {
		log.Fatalf("Error starting FIX client: %v", err)
	}
	defer client.Stop()

	// Example: Place a limit order
	err = client.PlaceOrder(
		"BTC-USD",                   // symbol
		"LIMIT",                     // order type
		"BUY",                       // side
		"BASE",                      // quantity type
		"0.1",                       // quantity
		"25000.00",                  // price
		config.Data.Fix.PortfolioID, // portfolio
	)
	if err != nil {
		log.Printf("Error placing order: %v", err)
	}

	// Process messages from the FIX session
	for msg := range client.Messages() {
		log.Printf("Received message: %v", msg)
	}
}
