package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort            string
	AppEnv             string
	DatabaseHost       string
	DatabasePort       string
	DatabaseUser       string
	DatabasePassword   string
	DatabaseName       string
	DatabaseSSLMode    string
	DatabaseMaxConns   int
	DatabaseMinConns   int
	DatabaseMaxIdle    time.Duration
	DatabaseMaxLife    time.Duration
	CorsAllowedOrigins string
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists, ignore error if file not found (e.g., in production/docker)
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:            getEnv("PORT", "8080"),
		AppEnv:             getEnv("ENV", "development"),
		DatabaseHost:       getEnv("DB_HOST", "localhost"),
		DatabasePort:       getEnv("DB_PORT", "5432"),
		DatabaseUser:       getEnv("DB_USER", "postgres"),
		DatabasePassword:   getEnv("DB_PASSWORD", "postgres"),
		DatabaseName:       getEnv("DB_NAME", "photo_studio_db"),
		DatabaseSSLMode:    getEnv("DB_SSLMODE", "disable"),
		DatabaseMaxConns:   getEnvAsInt("DB_MAX_CONNS", 25),
		DatabaseMinConns:   getEnvAsInt("DB_MIN_CONNS", 5),
		DatabaseMaxIdle:    time.Duration(getEnvAsInt("DB_MAX_IDLE_MINS", 15)) * time.Minute,
		DatabaseMaxLife:    time.Duration(getEnvAsInt("DB_MAX_LIFE_MINS", 60)) * time.Minute,
		CorsAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
	}

	return cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseName,
		c.DatabaseSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return fallback
}
