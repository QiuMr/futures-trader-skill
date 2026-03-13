package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configFileName = "config.json"
)

type Config struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".futures_trader"), nil
}

func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, configFileName), nil
}

func LoadConfig() (*Config, error) {
	configFilePath, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(apiKey, apiSecret string) error {
	configDir, err := GetConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	config := &Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
	}

	configFilePath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFilePath, data, 0600)
}

func ClearConfig() error {
	configFilePath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	return os.Remove(configFilePath)
}
