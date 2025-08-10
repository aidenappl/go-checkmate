package env

import (
	"fmt"
	"os"
)

// Environment Variables
var (
	PORT = getEnv("PORT", "8080")
)

// Basic Environment Variables with fallback
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

// Get Environment Variable or Panic server
func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s is not set", key))
	}
	return value
}
