package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort        int
	ServerHost        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration

	MongoURI    string
	MongoDBName string

	APIToken string
}

func (c *Config) BindAddress() string {
	return c.ServerHost + ":" + strconv.Itoa(c.ServerPort)
}

func LoadConfig() *Config {
	return &Config{
		ServerPort:        getEnvAsInt("SERVER_PORT", 8091),
		ServerHost:        getEnv("SERVER_HOST", "localhost"),
		ReadTimeout:       getEnvAsDuration("READ_TIMEOUT", 10*time.Second),
		WriteTimeout:      getEnvAsDuration("WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:       getEnvAsDuration("IDLE_TIMEOUT", 60*time.Second),
		ReadHeaderTimeout: getEnvAsDuration("READ_HEADER_TIMEOUT", 2*time.Second),
		MongoURI:          getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:       getEnv("MONGO_DATABASE", "products"),
		APIToken:          getEnv("API_TOKEN", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultVal
}
