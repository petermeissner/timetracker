package serverconfig

import "os"

// Config holds application configuration
type Config struct {
	DBPath string
	Port   string
}

// GetEnvOrDefault returns the value of an environment variable or a default value if not set
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ParseConfig creates configuration from environment variables with defaults
func ParseConfig() *Config {
	return &Config{
		DBPath: GetEnvOrDefault("DB_PATH", "./timesheet.db"),
		Port:   GetEnvOrDefault("PORT", "8080"),
	}
}
