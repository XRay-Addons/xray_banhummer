package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"goipban/models"
	"log"
	"os"
)

// Config struct with fields for application configuration
type Config struct {
	BanDuration            models.Duration `json:"ban_duration"`              // ban duration
	UpdateInterval         models.Duration `json:"update_interval"`           // auto-update ban list interval
	XRayLogsFile           string          `json:"xray_logs_file"`            // XRay access.log file path
	XRayBlacklistOutbound  string          `json:"xray_blacklist_outbound"`   // XRay outbound tag for IPs be banned
	BannedDatabaseFile     string          `json:"banned_database_file"`      // banned database file
	ProtectedUFWRulesCount int             `json:"protected_ufw_rules_count"` // first rules of ufw are protected, add bans after them
	DebugMode              bool            `json:"debug_mode"`                // debug mode, don't ban anybody, just log
}

//go:embed default_config.json
var defaultConfigFile embed.FS

// LoadConfig loads the configuration from a file or falls back to the embedded default
func LoadConfig(filePath string) (*Config, error) {
	var data []byte
	var err error

	if filePath != "" {
		// Try loading from the specified file
		data, err = os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to load config from file: %v, falling back to default", err)
		}
	}

	if data == nil {
		// Load the embedded default config
		data, err = defaultConfigFile.ReadFile("default_config.json")
		if err != nil {
			return nil, fmt.Errorf("failed to load embedded default config: %w", err)
		}
	}

	// Parse the YAML data. SYKA BLYAT wtf
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// parse
