package env

import "os"

// Environment Variables
var (
	PORT = getEnv("PORT", "8080")
)

// getEnv returns the environment variable value or a fallback default
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
