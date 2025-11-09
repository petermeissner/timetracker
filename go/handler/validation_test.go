package handler

import (
	"testing"
	"time"
	"timesheet/go/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkgmodel "timesheet/go/model"
	pkgutil "timesheet/go/util"
)

func TestValidateTimeEntryRequest(t *testing.T) {
	tests := []struct {
		name        string
		req         model.TimeEntryRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid request",
			req: pkgmodel.TimeEntryRequest{
				Task:      "Development",
				Category:  "project work",
				StartTime: "2025-11-09T09:00:00Z",
				EndTime:   "2025-11-09T10:00:00Z",
			},
			expectError: false,
		},
		{
			name: "missing task",
			req: pkgmodel.TimeEntryRequest{
				Category:  "project work",
				StartTime: "2025-11-09T09:00:00Z",
				EndTime:   "2025-11-09T10:00:00Z",
			},
			expectError: true,
			errorMsg:    "task is required",
		},
		{
			name: "missing category",
			req: pkgmodel.TimeEntryRequest{
				Task:      "Development",
				StartTime: "2025-11-09T09:00:00Z",
				EndTime:   "2025-11-09T10:00:00Z",
			},
			expectError: true,
			errorMsg:    "category is required",
		},
		{
			name: "missing start time",
			req: pkgmodel.TimeEntryRequest{
				Task:     "Development",
				Category: "project work",
				EndTime:  "2025-11-09T10:00:00Z",
			},
			expectError: true,
			errorMsg:    "start_time is required",
		},
		{
			name: "missing end time",
			req: pkgmodel.TimeEntryRequest{
				Task:      "Development",
				Category:  "project work",
				StartTime: "2025-11-09T09:00:00Z",
			},
			expectError: true,
			errorMsg:    "end_time is required",
		},
		{
			name: "empty strings should fail",
			req: pkgmodel.TimeEntryRequest{
				Task:      "",
				Category:  "",
				StartTime: "",
				EndTime:   "",
			},
			expectError: true,
			errorMsg:    "task is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTimeEntryRequest(tt.req)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParseAndValidateTimeEntry(t *testing.T) {
	tests := []struct {
		name             string
		req              pkgmodel.TimeEntryRequest
		expectError      bool
		errorMsg         string
		expectedDuration int
	}{
		{
			name: "valid 1 hour entry",
			req: pkgmodel.TimeEntryRequest{
				Task:      "Development",
				Category:  "project work",
				StartTime: "2025-11-09T09:00:00Z",
				EndTime:   "2025-11-09T10:00:00Z",
			},
			expectError:      false,
			expectedDuration: 60,
		},
		{
			name: "valid 30 minute entry",
			req: pkgmodel.TimeEntryRequest{
				Task:      "Meeting",
				Category:  "project support",
				StartTime: "2025-11-09T14:00:00Z",
				EndTime:   "2025-11-09T14:30:00Z",
			},
			expectError:      false,
			expectedDuration: 30,
		},
		{
			name: "valid 2.5 hour entry",
			req: pkgmodel.TimeEntryRequest{
				Task:      "Research",
				Category:  "project work",
				StartTime: "2025-11-09T09:00:00Z",
				EndTime:   "2025-11-09T11:30:00Z",
			},
			expectError:      false,
			expectedDuration: 150,
		},
		{
			name: "invalid start time format",
			req: pkgmodel.TimeEntryRequest{
				Task:      "Development",
				Category:  "project work",
				StartTime: "invalid-time",
				EndTime:   "2025-11-09T10:00:00Z",
			},
			expectError: true,
			errorMsg:    "invalid start_time format",
		},
		{
			name: "invalid end time format",
			req: pkgmodel.TimeEntryRequest{
				Task:      "Development",
				Category:  "project work",
				StartTime: "2025-11-09T09:00:00Z",
				EndTime:   "invalid-time",
			},
			expectError: true,
			errorMsg:    "invalid end_time format",
		},
		{
			name: "end time before start time",
			req: pkgmodel.TimeEntryRequest{
				Task:      "Development",
				Category:  "project work",
				StartTime: "2025-11-09T10:00:00Z",
				EndTime:   "2025-11-09T09:00:00Z",
			},
			expectError: true,
			errorMsg:    "end_time must be after start_time",
		},
		{
			name: "equal start and end time",
			req: pkgmodel.TimeEntryRequest{
				Task:      "Development",
				Category:  "project work",
				StartTime: "2025-11-09T09:00:00Z",
				EndTime:   "2025-11-09T09:00:00Z",
			},
			expectError: true,
			errorMsg:    "end_time must be after start_time",
		},
		{
			name: "missing task",
			req: pkgmodel.TimeEntryRequest{
				Category:  "project work",
				StartTime: "2025-11-09T09:00:00Z",
				EndTime:   "2025-11-09T10:00:00Z",
			},
			expectError: true,
			errorMsg:    "task is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime, endTime, duration, err := ParseAndValidateTimeEntry(tt.req)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Zero(t, duration)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedDuration, duration)
				assert.True(t, endTime.After(startTime))
			}
		})
	}
}

func TestValidateTimeSequence(t *testing.T) {
	tests := []struct {
		name        string
		startTime   time.Time
		endTime     time.Time
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid sequence",
			startTime:   time.Date(2025, 11, 9, 9, 0, 0, 0, time.UTC),
			endTime:     time.Date(2025, 11, 9, 10, 0, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "end before start",
			startTime:   time.Date(2025, 11, 9, 10, 0, 0, 0, time.UTC),
			endTime:     time.Date(2025, 11, 9, 9, 0, 0, 0, time.UTC),
			expectError: true,
			errorMsg:    "end_time must be after start_time",
		},
		{
			name:        "equal times",
			startTime:   time.Date(2025, 11, 9, 9, 0, 0, 0, time.UTC),
			endTime:     time.Date(2025, 11, 9, 9, 0, 0, 0, time.UTC),
			expectError: true,
			errorMsg:    "end_time must be after start_time",
		},
		{
			name:        "1 second difference",
			startTime:   time.Date(2025, 11, 9, 9, 0, 0, 0, time.UTC),
			endTime:     time.Date(2025, 11, 9, 9, 0, 1, 0, time.UTC),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTimeSequence(tt.startTime, tt.endTime)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCalculateDurationMinutes(t *testing.T) {
	tests := []struct {
		name             string
		startTime        time.Time
		endTime          time.Time
		expectedDuration int
	}{
		{
			name:             "1 hour",
			startTime:        time.Date(2025, 11, 9, 9, 0, 0, 0, time.UTC),
			endTime:          time.Date(2025, 11, 9, 10, 0, 0, 0, time.UTC),
			expectedDuration: 60,
		},
		{
			name:             "30 minutes",
			startTime:        time.Date(2025, 11, 9, 14, 0, 0, 0, time.UTC),
			endTime:          time.Date(2025, 11, 9, 14, 30, 0, 0, time.UTC),
			expectedDuration: 30,
		},
		{
			name:             "2.5 hours",
			startTime:        time.Date(2025, 11, 9, 9, 0, 0, 0, time.UTC),
			endTime:          time.Date(2025, 11, 9, 11, 30, 0, 0, time.UTC),
			expectedDuration: 150,
		},
		{
			name:             "1 minute",
			startTime:        time.Date(2025, 11, 9, 9, 0, 0, 0, time.UTC),
			endTime:          time.Date(2025, 11, 9, 9, 1, 0, 0, time.UTC),
			expectedDuration: 1,
		},
		{
			name:             "15 seconds (rounds down)",
			startTime:        time.Date(2025, 11, 9, 9, 0, 0, 0, time.UTC),
			endTime:          time.Date(2025, 11, 9, 9, 0, 15, 0, time.UTC),
			expectedDuration: 0,
		},
		{
			name:             "45 seconds (rounds down)",
			startTime:        time.Date(2025, 11, 9, 9, 0, 0, 0, time.UTC),
			endTime:          time.Date(2025, 11, 9, 9, 0, 45, 0, time.UTC),
			expectedDuration: 0,
		},
		{
			name:             "1 minute 30 seconds",
			startTime:        time.Date(2025, 11, 9, 9, 0, 0, 0, time.UTC),
			endTime:          time.Date(2025, 11, 9, 9, 1, 30, 0, time.UTC),
			expectedDuration: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := pkgutil.CalculateDurationMinutes(tt.startTime, tt.endTime)
			assert.Equal(t, tt.expectedDuration, duration)
		})
	}
}

func TestFormatTimeForDB(t *testing.T) {
	testTime := time.Date(2025, 11, 9, 14, 30, 45, 0, time.UTC)
	expected := "2025-11-09 14:30:45"

	result := pkgutil.FormatTimeForDB(testTime)
	assert.Equal(t, expected, result)
}

func TestGetCurrentDateForDB(t *testing.T) {
	result := pkgutil.GetCurrentDateForDB()

	// Should be in format "YYYY-MM-DD"
	assert.Len(t, result, 10)
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}$`, result)

	// Should be today's date
	expectedPrefix := time.Now().Format("2006-01-02")
	assert.Equal(t, expectedPrefix, result)
}
