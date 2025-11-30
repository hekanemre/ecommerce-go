package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DB DatabaseConfig `json:"db"`
}

type DatabaseConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	// Open the config file
	file, err := os.Open("config/config.json")
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	// Parse the JSON config
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("could not parse config: %w", err)
	}

	return &config, nil
}
