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

// LookupEnvOrBool looks up if key is set as a system environment variable.
// If yes, true is returned.
// If key doesn't exist, defaultVal will be used.
func LookupEnvOrBool(key string, defaultVal bool) bool {
	if _, ok := os.LookupEnv(key); ok {
		return ok
	}
	return defaultVal
}
