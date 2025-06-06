package syncd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Config struct {
	Repository struct {
		URL      string `json:"url" yaml:"url"`
		Branch   string `json:"branch" yaml:"branch"`
		AuthType string `json:"auth_type" yaml:"auth_type"` // token, ssh
		Token    string `json:"token" yaml:"token"`
		SSHKey   string `json:"ssh_key" yaml:"ssh_key"`
	} `json:"repository" yaml:"repository"`

	Sync struct {
		Frequency  string `json:"frequency" yaml:"frequency"` // e.g., "5m", "1h"
		LocalPath  string `json:"local_path" yaml:"local_path"`
		RemotePath string `json:"remote_path" yaml:"remote_path"`
	} `json:"sync" yaml:"sync"`
}

const (
	defaultConfigFile = "config.yaml"
)

// getDefaultConfigDir returns the default config directory using XDG Base Directory spec
func getDefaultConfigDir() string {
	// First check XDG_CONFIG_HOME
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		// If not set, default to ~/.config
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Fallback to /etc/hotpot if we can't get home dir
			return "/etc/hotpot/syncd"
		}
		xdgConfigHome = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(xdgConfigHome, "hotpot", "syncd")
}

// LoadConfig loads the syncd configuration from the default location
func LoadConfig() (*Config, error) {
	configPath := filepath.Join(getDefaultConfigDir(), defaultConfigFile)

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to the default location
func SaveConfig(config *Config) error {
	configDir := getDefaultConfigDir()
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	configPath := filepath.Join(configDir, defaultConfigFile)

	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	if err := os.WriteFile(configPath, configJSON, 0600); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}
