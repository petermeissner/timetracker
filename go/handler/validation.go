package handler

import (
	"errors"
	"fmt"
	"time"
	pkgmodel "timesheet/go/model"
	pkgutil "timesheet/go/util"
)

// ValidateTimeEntryRequest validates the required fields of a time entry request
func ValidateTimeEntryRequest(req pkgmodel.TimeEntryRequest) error {
	if req.Task == "" {
		return errors.New("task is required")
	}
	if req.Category == "" {
		return errors.New("category is required")
	}
	if req.StartTime == "" {
		return errors.New("start_time is required")
	}
	if req.EndTime == "" {
		return errors.New("end_time is required")
	}
	return nil
}

// ParseAndValidateTimeEntry parses and validates time entry times, returning parsed times and calculated duration
func ParseAndValidateTimeEntry(req pkgmodel.TimeEntryRequest) (startTime, endTime time.Time, duration int, err error) {
	// First validate required fields
	if err = ValidateTimeEntryRequest(req); err != nil {
		return time.Time{}, time.Time{}, 0, err
	}

	// Parse start time
	startTime, err = time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return time.Time{}, time.Time{}, 0, fmt.Errorf("invalid start time format. Expected ISO timestamp: %w", err)
	}

	// Parse end time
	endTime, err = time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return time.Time{}, time.Time{}, 0, fmt.Errorf("invalid end time format. Expected ISO timestamp: %w", err)
	}

	// Validate time sequence
	if err = ValidateTimeSequence(startTime, endTime); err != nil {
		return time.Time{}, time.Time{}, 0, err
	}

	// Calculate duration
	duration = pkgutil.CalculateDurationMinutes(startTime, endTime)

	return startTime, endTime, duration, nil
}

// ValidateTimeSequence ensures end time is after start time
func ValidateTimeSequence(startTime time.Time, endTime time.Time) error {
	if endTime.Before(startTime) {
		return errors.New("end time must be after start time")
	}
	if endTime.Equal(startTime) {
		return errors.New("end time must be after start time")
	}
	return nil
}
