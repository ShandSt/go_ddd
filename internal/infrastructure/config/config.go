package config

import (
	"os"
	"time"
)

type Config struct {
	Server  ServerConfig
	MongoDB MongoDBConfig
	API     APIConfig
}

type ServerConfig struct {
	Host              string
	Port              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type APIConfig struct {
	Token string
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Host:              getEnv("SERVER_HOST", "localhost"),
			Port:              getEnv("SERVER_PORT", "8091"),
			ReadTimeout:       getDurationEnv("READ_TIMEOUT", 10*time.Second),
			WriteTimeout:      getDurationEnv("WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:       getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
			ReadHeaderTimeout: getDurationEnv("READ_HEADER_TIMEOUT", 2*time.Second),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGO_DATABASE", "products"),
		},
		API: APIConfig{
			Token: getEnv("API_TOKEN", ""),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
