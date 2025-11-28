package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server      string `yaml:"server"`
	AccessToken string `yaml:"access_token"`
	Username    string `yaml:"username"`
}

var configDir = filepath.Join(os.Getenv("HOME"), ".cutter")
var configFile = filepath.Join(configDir, "config.yaml")

func Load() (*Config, error) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &Config{}, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config dir: %v", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}

	return nil
}

func (c *Config) IsAuthenticated() bool {
	return c.AccessToken != "" && c.Server != ""
}

func GetConfigPath() string {
	return configFile
}
