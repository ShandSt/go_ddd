package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
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
				ServerPort:        8091,
				ServerHost:        "localhost",
				ReadTimeout:       10 * time.Second,
				WriteTimeout:      10 * time.Second,
				IdleTimeout:       60 * time.Second,
				ReadHeaderTimeout: 2 * time.Second,
				MongoURI:          "mongodb://localhost:27017",
				MongoDBName:       "products",
				APIToken:          "",
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
				ServerPort:        9090,
				ServerHost:        "0.0.0.0",
				ReadTimeout:       20 * time.Second,
				WriteTimeout:      20 * time.Second,
				IdleTimeout:       120 * time.Second,
				ReadHeaderTimeout: 5 * time.Second,
				MongoURI:          "mongodb://custom:27017",
				MongoDBName:       "custom_db",
				APIToken:          "test_token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				if v == "" {
					os.Unsetenv(k)
				} else {
					os.Setenv(k, v)
				}
			}

			// Load config
			config := LoadConfig()

			// Check values
			if config.ServerPort != tt.expectedConfig.ServerPort {
				t.Errorf("Expected ServerPort %d, got %d", tt.expectedConfig.ServerPort, config.ServerPort)
			}
			if config.ServerHost != tt.expectedConfig.ServerHost {
				t.Errorf("Expected ServerHost %s, got %s", tt.expectedConfig.ServerHost, config.ServerHost)
			}
			if config.ReadTimeout != tt.expectedConfig.ReadTimeout {
				t.Errorf("Expected ReadTimeout %v, got %v", tt.expectedConfig.ReadTimeout, config.ReadTimeout)
			}
			if config.WriteTimeout != tt.expectedConfig.WriteTimeout {
				t.Errorf("Expected WriteTimeout %v, got %v", tt.expectedConfig.WriteTimeout, config.WriteTimeout)
			}
			if config.IdleTimeout != tt.expectedConfig.IdleTimeout {
				t.Errorf("Expected IdleTimeout %v, got %v", tt.expectedConfig.IdleTimeout, config.IdleTimeout)
			}
			if config.ReadHeaderTimeout != tt.expectedConfig.ReadHeaderTimeout {
				t.Errorf("Expected ReadHeaderTimeout %v, got %v", tt.expectedConfig.ReadHeaderTimeout, config.ReadHeaderTimeout)
			}
			if config.MongoURI != tt.expectedConfig.MongoURI {
				t.Errorf("Expected MongoURI %s, got %s", tt.expectedConfig.MongoURI, config.MongoURI)
			}
			if config.MongoDBName != tt.expectedConfig.MongoDBName {
				t.Errorf("Expected MongoDBName %s, got %s", tt.expectedConfig.MongoDBName, config.MongoDBName)
			}
			if config.APIToken != tt.expectedConfig.APIToken {
				t.Errorf("Expected APIToken %s, got %s", tt.expectedConfig.APIToken, config.APIToken)
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

func TestGetEnvAsDuration(t *testing.T) {
	key := "TEST_DURATION"
	originalValue := os.Getenv(key)

	t.Run("valid duration", func(t *testing.T) {
		os.Setenv(key, "5s")
		duration := getEnvAsDuration(key, 10*time.Second)
		if duration != 5*time.Second {
			t.Errorf("Expected 5s, got %s", duration)
		}
	})

	t.Run("invalid duration", func(t *testing.T) {
		os.Setenv(key, "invalid")
		duration := getEnvAsDuration(key, 10*time.Second)
		if duration != 10*time.Second {
			t.Errorf("Expected 10s, got %s", duration)
		}
	})

	t.Run("no value", func(t *testing.T) {
		os.Unsetenv(key)
		duration := getEnvAsDuration(key, 10*time.Second)
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

func TestGetEnvAsInt(t *testing.T) {
	key := "TEST_INT"
	originalValue := os.Getenv(key)

	t.Run("valid integer", func(t *testing.T) {
		os.Setenv(key, "42")
		value := getEnvAsInt(key, 10)
		if value != 42 {
			t.Errorf("Expected 42, got %d", value)
		}
	})

	t.Run("invalid integer", func(t *testing.T) {
		os.Setenv(key, "invalid")
		value := getEnvAsInt(key, 10)
		if value != 10 {
			t.Errorf("Expected 10, got %d", value)
		}
	})

	t.Run("no value", func(t *testing.T) {
		os.Unsetenv(key)
		value := getEnvAsInt(key, 10)
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

func TestGetEnvAsBool(t *testing.T) {
	key := "TEST_BOOL"
	originalValue := os.Getenv(key)

	t.Run("valid boolean", func(t *testing.T) {
		os.Setenv(key, "true")
		value := getEnvAsBool(key, false)
		if !value {
			t.Error("Expected true, got false")
		}
	})

	t.Run("invalid boolean", func(t *testing.T) {
		os.Setenv(key, "invalid")
		value := getEnvAsBool(key, true)
		if !value {
			t.Error("Expected true, got false")
		}
	})

	t.Run("no value", func(t *testing.T) {
		os.Unsetenv(key)
		value := getEnvAsBool(key, true)
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
		ServerHost: "localhost",
		ServerPort: 8080,
	}

	expected := "localhost:8080"
	if addr := config.BindAddress(); addr != expected {
		t.Errorf("Expected %s, got %s", expected, addr)
	}
}
