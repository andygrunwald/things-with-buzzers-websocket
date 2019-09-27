package main

import (
	"os"
)

// LookupEnvOrString looks up key in the system environment variables.
// If key doesn't exist, defaultVal will be used.
func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
