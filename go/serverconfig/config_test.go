package serverconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault(t *testing.T) {
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
			envKey:     "TEST_VAR_EXISTS",
			envValue:   "custom_value",
			setEnv:     true,
			defaultVal: "default_value",
			expected:   "custom_value",
		},
		{
			name:       "environment variable does not exist",
			envKey:     "TEST_VAR_NOT_EXISTS",
			setEnv:     false,
			defaultVal: "default_value",
			expected:   "default_value",
		},
		{
			name:       "environment variable is empty string",
			envKey:     "TEST_VAR_EMPTY",
			envValue:   "",
			setEnv:     true,
			defaultVal: "default_value",
			expected:   "default_value",
		},
		{
			name:       "environment variable with spaces",
			envKey:     "TEST_VAR_SPACES",
			envValue:   "  value_with_spaces  ",
			setEnv:     true,
			defaultVal: "default_value",
			expected:   "value_with_spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment variable after test
			defer os.Unsetenv(tt.envKey)

			if tt.setEnv {
				os.Setenv(tt.envKey, tt.envValue)
			}

			result := GetEnvOrDefault(tt.envKey, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name         string
		envVars      map[string]string
		expectedDB   string
		expectedPort string
	}{
		{
			name:         "default values",
			envVars:      map[string]string{},
			expectedDB:   "./timesheet.db",
			expectedPort: "8080",
		},
		{
			name: "environment variables set",
			envVars: map[string]string{
				"DB_PATH": "custom.db",
				"PORT":    "9000",
			},
			expectedDB:   "custom.db",
			expectedPort: "9000",
		},
		{
			name: "partial environment variables",
			envVars: map[string]string{
				"PORT": "3000",
			},
			expectedDB:   "./timesheet.db",
			expectedPort: "3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment variables after test
			defer func() {
				os.Unsetenv("DB_PATH")
				os.Unsetenv("PORT")
			}()

			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			config := ParseConfig()
			assert.Equal(t, tt.expectedDB, config.DBPath)
			assert.Equal(t, tt.expectedPort, config.Port)
		})
	}
}
