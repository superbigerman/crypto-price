package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ExternalAPIBaseURL string
	ExternalAPITimeout time.Duration
	ExternalAPIRelaxed bool
	ExternalAPITSyms   string
}

func Load() *Config {
	return &Config{
		ExternalAPIBaseURL: getEnv("EXTERNAL_API_URL", "https://min-api.cryptocompare.com"),
		ExternalAPITimeout: time.Duration(getEnvInt("EXTERNAL_API_TIMEOUT_SEC", 10)) * time.Second,
		ExternalAPIRelaxed: getEnvBool("EXTERNAL_API_RELAXED", true),
		ExternalAPITSyms:   getEnv("EXTERNAL_API_TSYMS", "USD"),
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
