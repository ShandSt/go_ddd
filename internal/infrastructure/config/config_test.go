package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	originalEnvVars := make(map[string]string)
	envVars := []string{
		"SERVER_PORT",
		"SERVER_HOST",
		"READ_TIMEOUT",
		"WRITE_TIMEOUT",
		"IDLE_TIMEOUT",
		"READ_HEADER_TIMEOUT",
		"MONGO_URI",
		"MONGO_DATABASE",
		"API_TOKEN",
	}

	for _, env := range envVars {
		if value, exists := os.LookupEnv(env); exists {
			originalEnvVars[env] = value
			os.Unsetenv(env)
		}
	}

	defer func() {
		for key, value := range originalEnvVars {
			os.Setenv(key, value)
		}
	}()

	t.Run("default values", func(t *testing.T) {
		config := LoadConfig()

		if config.ServerPort != 8080 {
			t.Errorf("Expected ServerPort 8080, got %d", config.ServerPort)
		}

		if config.ServerHost != "localhost" {
			t.Errorf("Expected ServerHost localhost, got %s", config.ServerHost)
		}

		if config.ReadTimeout != 10*time.Second {
			t.Errorf("Expected ReadTimeout 10s, got %s", config.ReadTimeout)
		}

		if config.WriteTimeout != 10*time.Second {
			t.Errorf("Expected WriteTimeout 10s, got %s", config.WriteTimeout)
		}

		if config.IdleTimeout != 60*time.Second {
			t.Errorf("Expected IdleTimeout 60s, got %s", config.IdleTimeout)
		}

		if config.ReadHeaderTimeout != 2*time.Second {
			t.Errorf("Expected ReadHeaderTimeout 2s, got %s", config.ReadHeaderTimeout)
		}

		if config.MongoURI != "mongodb://localhost:27017" {
			t.Errorf("Expected MongoURI mongodb://localhost:27017, got %s", config.MongoURI)
		}

		if config.MongoDBName != "products" {
			t.Errorf("Expected MongoDBName products, got %s", config.MongoDBName)
		}

		if config.APIToken != "" {
			t.Errorf("Expected empty APIToken, got %s", config.APIToken)
		}
	})

	t.Run("custom values", func(t *testing.T) {
		os.Setenv("SERVER_PORT", "3000")
		os.Setenv("SERVER_HOST", "0.0.0.0")
		os.Setenv("READ_TIMEOUT", "20s")
		os.Setenv("WRITE_TIMEOUT", "20s")
		os.Setenv("IDLE_TIMEOUT", "120s")
		os.Setenv("READ_HEADER_TIMEOUT", "5s")
		os.Setenv("MONGO_URI", "mongodb://localhost:27018")
		os.Setenv("MONGO_DATABASE", "test_db")
		os.Setenv("API_TOKEN", "test_token")

		config := LoadConfig()

		if config.ServerPort != 3000 {
			t.Errorf("Expected ServerPort 3000, got %d", config.ServerPort)
		}

		if config.ServerHost != "0.0.0.0" {
			t.Errorf("Expected ServerHost 0.0.0.0, got %s", config.ServerHost)
		}

		if config.ReadTimeout != 20*time.Second {
			t.Errorf("Expected ReadTimeout 20s, got %s", config.ReadTimeout)
		}

		if config.WriteTimeout != 20*time.Second {
			t.Errorf("Expected WriteTimeout 20s, got %s", config.WriteTimeout)
		}

		if config.IdleTimeout != 120*time.Second {
			t.Errorf("Expected IdleTimeout 120s, got %s", config.IdleTimeout)
		}

		if config.ReadHeaderTimeout != 5*time.Second {
			t.Errorf("Expected ReadHeaderTimeout 5s, got %s", config.ReadHeaderTimeout)
		}

		if config.MongoURI != "mongodb://localhost:27018" {
			t.Errorf("Expected MongoURI mongodb://localhost:27018, got %s", config.MongoURI)
		}

		if config.MongoDBName != "test_db" {
			t.Errorf("Expected MongoDBName test_db, got %s", config.MongoDBName)
		}

		if config.APIToken != "test_token" {
			t.Errorf("Expected APIToken test_token, got %s", config.APIToken)
		}
	})
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
