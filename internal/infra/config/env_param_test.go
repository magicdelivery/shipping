package config

import (
	"os"
	"testing"
)

func TestGetIntEnv(t *testing.T) {
	const KEY = "ENV_INT_KEY"
	const DEFAULT_VALUE = 11111111

	test_data := []struct {
		name     string
		set_env  func()
		expected int
	}{
		{
			name:     "environment variable is set to a valid integer",
			set_env:  func() { os.Setenv(KEY, "22222222") },
			expected: 22222222,
		},
		{
			name:     "environment variable is set to an invalid integer",
			set_env:  func() { os.Setenv(KEY, "invalid") },
			expected: DEFAULT_VALUE,
		},
		{
			name:     "environment variable is not set",
			set_env:  func() { os.Unsetenv(KEY) },
			expected: DEFAULT_VALUE,
		},
	}

	for _, td := range test_data {
		t.Run(td.name, func(t *testing.T) {
			td.set_env()
			if result := GetIntEnv(KEY, DEFAULT_VALUE); result != td.expected {
				t.Errorf("expected %d, got %d", td.expected, result)
			}
		})
	}
}

func TestGetStrEnv(t *testing.T) {
	const KEY = "ENV_STR_KEY"
	const DEFAULT_VALUE = "default"

	test_data := []struct {
		name     string
		set_env  func()
		expected string
	}{
		{
			name:     "environment variable is set",
			set_env:  func() { os.Setenv(KEY, "text") },
			expected: "text",
		},
		{
			name:     "environment variable is not set",
			set_env:  func() { os.Unsetenv(KEY) },
			expected: DEFAULT_VALUE,
		},
	}

	for _, td := range test_data {
		t.Run(td.name, func(t *testing.T) {
			td.set_env()
			if result := GetStrEnv(KEY, DEFAULT_VALUE); result != td.expected {
				t.Errorf("expected %s, got %s", td.expected, result)
			}
		})
	}
}
