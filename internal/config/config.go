package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config represents the application configuration
type Config struct {
	Server    ServerConfig    `json:"server"`
	Scheduler SchedulerConfig `json:"scheduler"`
	System    SystemConfig    `json:"system"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// SchedulerConfig represents the scheduler configuration
type SchedulerConfig struct {
	ScheduleInterval time.Duration `json:"schedule_interval"`
	MaxRetries       int           `json:"max_retries"`
}

// SystemConfig represents the system configuration
type SystemConfig struct {
	HeartbeatTimeout time.Duration `json:"heartbeat_timeout"`
	HealthCheckInterval time.Duration `json:"health_check_interval"`
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Scheduler: SchedulerConfig{
			ScheduleInterval: 10 * time.Second,
			MaxRetries:       3,
		},
		System: SystemConfig{
			HeartbeatTimeout:    60 * time.Second,
			HealthCheckInterval: 30 * time.Second,
		},
	}
}

// LoadFromFile loads configuration from a JSON file
func LoadFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := Default()
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveToFile saves configuration to a JSON file
func (c *Config) SaveToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}
