package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *Config
	}{
		{
			name: "default values",
			envVars: map[string]string{
				"SERVER_PORT":         "",
				"SERVER_HOST":         "",
				"READ_TIMEOUT":        "",
				"WRITE_TIMEOUT":       "",
				"IDLE_TIMEOUT":        "",
				"READ_HEADER_TIMEOUT": "",
				"MONGO_URI":           "",
				"MONGO_DATABASE":      "",
				"API_TOKEN":           "",
			},
			expectedConfig: &Config{
				Server: ServerConfig{
					Host:              "localhost",
					Port:              "8091",
					ReadTimeout:       10 * time.Second,
					WriteTimeout:      10 * time.Second,
					IdleTimeout:       60 * time.Second,
					ReadHeaderTimeout: 2 * time.Second,
				},
				MongoDB: MongoDBConfig{
					URI:      "mongodb://localhost:27017",
					Database: "products",
				},
				API: APIConfig{
					Token: "",
				},
			},
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"SERVER_PORT":         "9090",
				"SERVER_HOST":         "0.0.0.0",
				"READ_TIMEOUT":        "20s",
				"WRITE_TIMEOUT":       "20s",
				"IDLE_TIMEOUT":        "120s",
				"READ_HEADER_TIMEOUT": "5s",
				"MONGO_URI":           "mongodb://custom:27017",
				"MONGO_DATABASE":      "custom_db",
				"API_TOKEN":           "test_token",
			},
			expectedConfig: &Config{
				Server: ServerConfig{
					Host:              "0.0.0.0",
					Port:              "9090",
					ReadTimeout:       20 * time.Second,
					WriteTimeout:      20 * time.Second,
					IdleTimeout:       120 * time.Second,
					ReadHeaderTimeout: 5 * time.Second,
				},
				MongoDB: MongoDBConfig{
					URI:      "mongodb://custom:27017",
					Database: "custom_db",
				},
				API: APIConfig{
					Token: "test_token",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				if v == "" {
					os.Unsetenv(k)
				} else {
					os.Setenv(k, v)
				}
			}

			config, err := Load()
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}

			if config.Server.Host != tt.expectedConfig.Server.Host {
				t.Errorf("Expected Server.Host %s, got %s", tt.expectedConfig.Server.Host, config.Server.Host)
			}
			if config.Server.Port != tt.expectedConfig.Server.Port {
				t.Errorf("Expected Server.Port %s, got %s", tt.expectedConfig.Server.Port, config.Server.Port)
			}
			if config.Server.ReadTimeout != tt.expectedConfig.Server.ReadTimeout {
				t.Errorf("Expected Server.ReadTimeout %v, got %v", tt.expectedConfig.Server.ReadTimeout, config.Server.ReadTimeout)
			}
			if config.Server.WriteTimeout != tt.expectedConfig.Server.WriteTimeout {
				t.Errorf("Expected Server.WriteTimeout %v, got %v", tt.expectedConfig.Server.WriteTimeout, config.Server.WriteTimeout)
			}
			if config.Server.IdleTimeout != tt.expectedConfig.Server.IdleTimeout {
				t.Errorf("Expected Server.IdleTimeout %v, got %v", tt.expectedConfig.Server.IdleTimeout, config.Server.IdleTimeout)
			}
			if config.Server.ReadHeaderTimeout != tt.expectedConfig.Server.ReadHeaderTimeout {
				t.Errorf("Expected Server.ReadHeaderTimeout %v, got %v", tt.expectedConfig.Server.ReadHeaderTimeout, config.Server.ReadHeaderTimeout)
			}
			if config.MongoDB.URI != tt.expectedConfig.MongoDB.URI {
				t.Errorf("Expected MongoDB.URI %s, got %s", tt.expectedConfig.MongoDB.URI, config.MongoDB.URI)
			}
			if config.MongoDB.Database != tt.expectedConfig.MongoDB.Database {
				t.Errorf("Expected MongoDB.Database %s, got %s", tt.expectedConfig.MongoDB.Database, config.MongoDB.Database)
			}
			if config.API.Token != tt.expectedConfig.API.Token {
				t.Errorf("Expected API.Token %s, got %s", tt.expectedConfig.API.Token, config.API.Token)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	key := "TEST_ENV_VAR"
	originalValue := os.Getenv(key)

	t.Run("with value", func(t *testing.T) {
		os.Setenv(key, "test_value")
		value := getEnv(key, "default")
		if value != "test_value" {
			t.Errorf("Expected test_value, got %s", value)
		}
	})

	t.Run("without value", func(t *testing.T) {
		os.Unsetenv(key)
		value := getEnv(key, "default")
		if value != "default" {
			t.Errorf("Expected default, got %s", value)
		}
	})

	if originalValue != "" {
		os.Setenv(key, originalValue)
	} else {
		os.Unsetenv(key)
	}
}

func TestGetDurationEnv(t *testing.T) {
	key := "TEST_DURATION"
	originalValue := os.Getenv(key)

	t.Run("valid duration", func(t *testing.T) {
		os.Setenv(key, "5s")
		duration := getDurationEnv(key, 10*time.Second)
		if duration != 5*time.Second {
			t.Errorf("Expected 5s, got %s", duration)
		}
	})

	t.Run("invalid duration", func(t *testing.T) {
		os.Setenv(key, "invalid")
		duration := getDurationEnv(key, 10*time.Second)
		if duration != 10*time.Second {
			t.Errorf("Expected 10s, got %s", duration)
		}
	})

	t.Run("no value", func(t *testing.T) {
		os.Unsetenv(key)
		duration := getDurationEnv(key, 10*time.Second)
		if duration != 10*time.Second {
			t.Errorf("Expected 10s, got %s", duration)
		}
	})

	if originalValue != "" {
		os.Setenv(key, originalValue)
	} else {
		os.Unsetenv(key)
	}
}

func TestGetIntEnv(t *testing.T) {
	key := "TEST_INT"
	originalValue := os.Getenv(key)

	t.Run("valid integer", func(t *testing.T) {
		os.Setenv(key, "42")
		value := getIntEnv(key, 10)
		if value != 42 {
			t.Errorf("Expected 42, got %d", value)
		}
	})

	t.Run("invalid integer", func(t *testing.T) {
		os.Setenv(key, "invalid")
		value := getIntEnv(key, 10)
		if value != 10 {
			t.Errorf("Expected 10, got %d", value)
		}
	})

	t.Run("no value", func(t *testing.T) {
		os.Unsetenv(key)
		value := getIntEnv(key, 10)
		if value != 10 {
			t.Errorf("Expected 10, got %d", value)
		}
	})

	if originalValue != "" {
		os.Setenv(key, originalValue)
	} else {
		os.Unsetenv(key)
	}
}

func TestGetBoolEnv(t *testing.T) {
	key := "TEST_BOOL"
	originalValue := os.Getenv(key)

	t.Run("valid boolean", func(t *testing.T) {
		os.Setenv(key, "true")
		value := getBoolEnv(key, false)
		if !value {
			t.Error("Expected true, got false")
		}
	})

	t.Run("invalid boolean", func(t *testing.T) {
		os.Setenv(key, "invalid")
		value := getBoolEnv(key, true)
		if !value {
			t.Error("Expected true, got false")
		}
	})

	t.Run("no value", func(t *testing.T) {
		os.Unsetenv(key)
		value := getBoolEnv(key, true)
		if !value {
			t.Error("Expected true, got false")
		}
	})

	if originalValue != "" {
		os.Setenv(key, originalValue)
	} else {
		os.Unsetenv(key)
	}
}

func TestBindAddress(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
	}

	expected := "localhost:8080"
	if addr := config.BindAddress(); addr != expected {
		t.Errorf("Expected %s, got %s", expected, addr)
	}
}
