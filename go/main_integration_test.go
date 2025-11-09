package timesheet

import (
	"os"
	"testing"
	"timesheet/go/serverconfig"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationIntegration(t *testing.T) {
	t.Run("go package getEnvOrDefault function", func(t *testing.T) {
		tests := []struct {
			name       string
			envKey     string
			envValue   string
			setEnv     bool
			defaultVal string
			expected   string
		}{
			{
				name:       "environment variable exists",
				envKey:     "TEST_INTEGRATION_VAR_EXISTS",
				envValue:   "integration_custom_value",
				setEnv:     true,
				defaultVal: "integration_default_value",
				expected:   "integration_custom_value",
			},
			{
				name:       "environment variable does not exist",
				envKey:     "TEST_INTEGRATION_VAR_NOT_EXISTS",
				setEnv:     false,
				defaultVal: "integration_default_value",
				expected:   "integration_default_value",
			},
			{
				name:       "environment variable is empty string",
				envKey:     "TEST_INTEGRATION_VAR_EMPTY",
				envValue:   "",
				setEnv:     true,
				defaultVal: "integration_default_value",
				expected:   "integration_default_value",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Clean up environment variable after test
				defer os.Unsetenv(tt.envKey)

				if tt.setEnv {
					os.Setenv(tt.envKey, tt.envValue)
				}

				// Test the go package's GetEnvOrDefault function
				result := serverconfig.GetEnvOrDefault(tt.envKey, tt.defaultVal)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("configuration parsing matches expected behavior", func(t *testing.T) {
		// Test that go package configuration parsing works as expected
		config := serverconfig.ParseConfig()

		// Verify structure
		assert.NotNil(t, config)
		assert.NotEmpty(t, config.DBPath)
		assert.NotEmpty(t, config.Port)

		// Default values should be reasonable
		assert.Equal(t, "./timesheet.db", config.DBPath)
		assert.Equal(t, "8080", config.Port)
	})
}
