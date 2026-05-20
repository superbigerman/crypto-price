package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ExternalAPIBaseURL string
	ExternalAPITimeout time.Duration
	ExternalAPIRelaxed bool
	ExternalAPITSyms   string

	// PostgresSQL
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		ExternalAPIBaseURL: getEnv("EXTERNAL_API_URL", "https://min-api.cryptocompare.com"),
		ExternalAPITimeout: time.Duration(getEnvInt("EXTERNAL_API_TIMEOUT_SEC", 10)) * time.Second,
		ExternalAPIRelaxed: getEnvBool("EXTERNAL_API_RELAXED", true),
		ExternalAPITSyms:   getEnv("EXTERNAL_API_TSYMS", "USD"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "crypto"),
		DBSSLMode:          getEnv("DB_SSL_MODE", "disbale"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
