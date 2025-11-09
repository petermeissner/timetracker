package util

import "time"

// CalculateDurationMinutes calculates the duration between two times in minutes
func CalculateDurationMinutes(startTime time.Time, endTime time.Time) int {
	return int(endTime.Sub(startTime).Minutes())
}

// FormatTimeForDB formats a time for database storage
func FormatTimeForDB(t time.Time) string {
	return t.Format(time.RFC3339)
}

// GetCurrentDateForDB returns the current date in database format
func GetCurrentDateForDB() string {
	return time.Now().Format("2006-01-02")
}
