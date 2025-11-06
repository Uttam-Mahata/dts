package config

import (
	"os"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg == nil {
		t.Fatal("Expected non-nil config")
	}

	// Test server defaults
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected server host '0.0.0.0', got '%s'", cfg.Server.Host)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Expected server port 8080, got %d", cfg.Server.Port)
	}

	// Test scheduler defaults
	if cfg.Scheduler.ScheduleInterval != 10*time.Second {
		t.Errorf("Expected schedule interval 10s, got %v", cfg.Scheduler.ScheduleInterval)
	}

	if cfg.Scheduler.MaxRetries != 3 {
		t.Errorf("Expected max retries 3, got %d", cfg.Scheduler.MaxRetries)
	}

	// Test system defaults
	if cfg.System.HeartbeatTimeout != 60*time.Second {
		t.Errorf("Expected heartbeat timeout 60s, got %v", cfg.System.HeartbeatTimeout)
	}

	if cfg.System.HealthCheckInterval != 30*time.Second {
		t.Errorf("Expected health check interval 30s, got %v", cfg.System.HealthCheckInterval)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "config-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Create a custom config
	cfg := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 9090,
		},
		Scheduler: SchedulerConfig{
			ScheduleInterval: 5 * time.Second,
			MaxRetries:       5,
		},
		System: SystemConfig{
			HeartbeatTimeout:    30 * time.Second,
			HealthCheckInterval: 15 * time.Second,
		},
	}

	// Save config
	err = cfg.SaveToFile(tmpPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loadedCfg, err := LoadFromFile(tmpPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify loaded config matches original
	if loadedCfg.Server.Host != cfg.Server.Host {
		t.Errorf("Expected server host '%s', got '%s'", cfg.Server.Host, loadedCfg.Server.Host)
	}

	if loadedCfg.Server.Port != cfg.Server.Port {
		t.Errorf("Expected server port %d, got %d", cfg.Server.Port, loadedCfg.Server.Port)
	}

	if loadedCfg.Scheduler.ScheduleInterval != cfg.Scheduler.ScheduleInterval {
		t.Errorf("Expected schedule interval %v, got %v", cfg.Scheduler.ScheduleInterval, loadedCfg.Scheduler.ScheduleInterval)
	}

	if loadedCfg.Scheduler.MaxRetries != cfg.Scheduler.MaxRetries {
		t.Errorf("Expected max retries %d, got %d", cfg.Scheduler.MaxRetries, loadedCfg.Scheduler.MaxRetries)
	}
}

func TestLoadFromFileNotExists(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/config.json")

	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}
}

func TestLoadFromFileInvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON
	tmpFile, err := os.CreateTemp("", "config-test-invalid-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Write invalid JSON
	_, err = tmpFile.WriteString("{invalid json}")
	if err != nil {
		t.Fatalf("Failed to write invalid JSON: %v", err)
	}
	tmpFile.Close()

	// Try to load
	_, err = LoadFromFile(tmpPath)

	if err == nil {
		t.Error("Expected error when loading invalid JSON")
	}
}

func TestSaveToFileInvalidPath(t *testing.T) {
	cfg := Default()

	err := cfg.SaveToFile("/nonexistent/directory/config.json")

	if err == nil {
		t.Error("Expected error when saving to invalid path")
	}
}

func TestServerConfig(t *testing.T) {
	cfg := ServerConfig{
		Host: "127.0.0.1",
		Port: 3000,
	}

	if cfg.Host != "127.0.0.1" {
		t.Errorf("Expected host '127.0.0.1', got '%s'", cfg.Host)
	}

	if cfg.Port != 3000 {
		t.Errorf("Expected port 3000, got %d", cfg.Port)
	}
}

func TestSchedulerConfig(t *testing.T) {
	cfg := SchedulerConfig{
		ScheduleInterval: 20 * time.Second,
		MaxRetries:       10,
	}

	if cfg.ScheduleInterval != 20*time.Second {
		t.Errorf("Expected schedule interval 20s, got %v", cfg.ScheduleInterval)
	}

	if cfg.MaxRetries != 10 {
		t.Errorf("Expected max retries 10, got %d", cfg.MaxRetries)
	}
}

func TestSystemConfig(t *testing.T) {
	cfg := SystemConfig{
		HeartbeatTimeout:    120 * time.Second,
		HealthCheckInterval: 60 * time.Second,
	}

	if cfg.HeartbeatTimeout != 120*time.Second {
		t.Errorf("Expected heartbeat timeout 120s, got %v", cfg.HeartbeatTimeout)
	}

	if cfg.HealthCheckInterval != 60*time.Second {
		t.Errorf("Expected health check interval 60s, got %v", cfg.HealthCheckInterval)
	}
}
