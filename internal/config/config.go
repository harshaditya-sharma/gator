package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error finding User Home Directory: %w", err)
	}
	configDir := homeDir + "/" + configFileName
	return configDir, nil
}

func write(path string, data []byte) error {
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}

func Read() (*Config, error) {

	configDir, err := getConfigFilePath()
	if err != nil {
		return &Config{}, err
	}

	content, err := os.ReadFile(configDir)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, fmt.Errorf("Config file does not exit: %w", err)
		}
		return &Config{}, fmt.Errorf("Error reading file: %w", err)
	}

	var config Config
	err = json.Unmarshal([]byte(content), &config)
	if err != nil {
		return &Config{}, fmt.Errorf("Error parsing config json: %w", err)
	}

	return &config, nil
}

func (c *Config) SetUser(username string) error {
	if username == "" {
		return fmt.Errorf("Empty username given")
	}

	c.Current_user_name = username

	configDir, err := getConfigFilePath()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling config json: %w", err)
	}

	err = write(configDir, b)
	if err != nil {
		c.Current_user_name = ""
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
