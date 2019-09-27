package main

import (
	"os"
	"testing"
)

func TestLookupEnvOrString_EnvVar(t *testing.T) {
	envKey := "TWB_ENV_TEST_KEY_ENVVALUE"
	envValue := "env-value"
	defaultVal := "default"

	err := os.Setenv(envKey, envValue)
	if err != nil {
		t.Errorf("not able to set environment variable %q: %q", envKey, err)
	}
	defer os.Unsetenv(envKey)

	if v := LookupEnvOrString(envKey, defaultVal); v != envValue {
		t.Errorf("got %q, want %q", v, envValue)
	}
}

func TestLookupEnvOrString_DefaultVal(t *testing.T) {
	envKey := "TWB_ENV_TEST_KEY_DEFAULTVALUE"
	defaultVal := "default"

	if v := LookupEnvOrString(envKey, defaultVal); v != defaultVal {
		t.Errorf("got %q, want %q", v, defaultVal)
	}
}
