package config

import (
	"os"
	"strconv"
)

// GetIntEnv retrieves the value of the environment variable named by the key,
// converts it to an integer, and returns it. If the variable is not present in
// the environment or cannot be converted to an integer, it returns the specified
// defaultValue.
func GetIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); len(value) == 0 {
		return defaultValue
	} else if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	} else {
		return defaultValue
	}
}

// GetStrEnv retrieves the value of the environment variable named by the key.
// If the variable is not present in the environment, it returns the specified
// defaultValue.
func GetStrEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); len(value) == 0 {
		return defaultValue
	} else {
		return value
	}
}
