package remote

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type SSHConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	LogPath  string `json:"log_path"`
}

type Config struct {
	Servers []SSHConfig `json:"servers"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// Set default values
	for i := range config.Servers {
		if config.Servers[i].Port == 0 {
			config.Servers[i].Port = 22
		}
		if config.Servers[i].LogPath == "" {
			config.Servers[i].LogPath = "/var/log/nginx/access.log"
		}
	}

	return &config, nil
}

func CreateSampleConfig(filename string) error {
	config := Config{
		Servers: []SSHConfig{
			{
				Host:     "your-server.com",
				Port:     22,
				Username: "root",
				Password: "your-password",
				LogPath:  "/var/log/nginx/access.log",
			},
		},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}