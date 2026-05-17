package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Server ServerConfig
	DB     DatabaseConfig
	Logger LoggerConfig
}

// ServerConfig holds server-specific settings
type ServerConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
	MaxHeaderMB  int64
}

// DatabaseConfig holds database-specific settings
type DatabaseConfig struct {
	Path string
}

// LoggerConfig holds logger settings
type LoggerConfig struct {
	Level string // debug, info, warn, error
}

// Load reads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvInt("SERVER_READ_TIMEOUT", 15),
			WriteTimeout: getEnvInt("SERVER_WRITE_TIMEOUT", 15),
			MaxHeaderMB:  1,
		},
		DB: DatabaseConfig{
			Path: getEnv("DB_PATH", "./crud.db"),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}
}

// getEnv returns an environment variable or a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt returns an environment variable as int or a default value
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("file:%s?cache=shared&mode=rwc", c.DB.Path)
}
