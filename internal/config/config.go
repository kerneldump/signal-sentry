package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds all application configuration settings.
type Config struct {
	RouterURL       string `json:"router_url"`
	PingTarget      string `json:"ping_target"`
	RefreshInterval int    `json:"refresh_interval"`
	Format          string `json:"format"`
	Output          string `json:"output"`
	LiveMode        bool   `json:"live_mode"`        // Future proofing for Track 1
	DisableAutoLog  bool   `json:"disable_auto_log"` // Disables the always-on stats.log
	WebEnabled      bool   `json:"web_enabled"`      // Unified Run Mode
	WebPort         int    `json:"web_port"`         // Unified Run Mode
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		RouterURL:       "http://192.168.12.1/TMI/v1/gateway?get=all",
		PingTarget:      "8.8.8.8",
		RefreshInterval: 5,
		WebPort:         8080,
	}
}

// Load reads a JSON configuration file and merges it with defaults.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		return cfg, nil
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // Not an error if file just doesn't exist
		}
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return cfg, nil
}
