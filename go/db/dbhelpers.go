package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"timesheet/go/handler"
	"timesheet/go/model"
	pkgutil "timesheet/go/util"
)

// ParseNullableTime safely parses a nullable time string from the database
func ParseNullableTime(nullTime sql.NullString) (time.Time, error) {
	if !nullTime.Valid {
		return time.Time{}, nil
	}

	parsedTime, err := time.Parse(time.RFC3339, nullTime.String)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing time '%s': %w", nullTime.String, err)
	}

	return parsedTime, nil
}

// ValidateCategoryExists checks if a category exists in the database
func ValidateCategoryExists(db *sql.DB, categoryName string) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE name = ?)", categoryName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("database error while validating category: %w", err)
	}
	if !exists {
		return fmt.Errorf("invalid category '%s': category does not exist in the system", categoryName)
	}
	return nil
}

// CreateTimeEntryInDB creates a new time entry in the database
func CreateTimeEntryInDB(db *sql.DB, req model.TimeEntryRequest) (*model.TimeEntry, error) {
	// Validate and parse the request
	startTime, endTime, duration, err := handler.ParseAndValidateTimeEntry(req)
	if err != nil {
		return nil, err
	}

	// Validate category exists
	if err := ValidateCategoryExists(db, req.Category); err != nil {
		return nil, err
	}

	// Get current date for compatibility with existing database schema
	currentDate := pkgutil.GetCurrentDateForDB()

	// Insert into database
	result, err := db.Exec(`
		INSERT INTO time_entries (task, description, category, start_time, end_time, duration, date)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, req.Task, req.Description, req.Category, pkgutil.FormatTimeForDB(startTime),
		pkgutil.FormatTimeForDB(endTime), duration, currentDate)

	if err != nil {
		log.Printf("ERROR: Failed to insert time entry - Task: %s, Category: %s, Start: %s, End: %s - Error: %v",
			req.Task, req.Category, req.StartTime, req.EndTime, err)
		return nil, fmt.Errorf("failed to create time entry: %w", err)
	}

	id, _ := result.LastInsertId()
	log.Printf("INSERT: Created time entry ID %d - Task: %s, Category: %s, Duration: %d min, Start: %s, End: %s",
		id, req.Task, req.Category, duration, startTime.Format("2006-01-02 15:04"), endTime.Format("2006-01-02 15:04"))

	return &model.TimeEntry{
		ID:          int(id),
		Task:        req.Task,
		Description: req.Description,
		Category:    req.Category,
		StartTime:   startTime,
		EndTime:     endTime,
	}, nil
}

// UpdateTimeEntryInDB updates an existing time entry in the database
func UpdateTimeEntryInDB(db *sql.DB, id int, req model.TimeEntryRequest) (*model.TimeEntry, error) {
	// Validate and parse the request
	startTime, endTime, duration, err := handler.ParseAndValidateTimeEntry(req)
	if err != nil {
		return nil, err
	}

	// Validate category exists
	if err := ValidateCategoryExists(db, req.Category); err != nil {
		return nil, err
	}

	// Get current date for compatibility with existing database schema
	currentDate := pkgutil.GetCurrentDateForDB()

	// Update in database
	_, err = db.Exec(`
		UPDATE time_entries 
		SET task = ?, description = ?, category = ?, start_time = ?, end_time = ?, duration = ?, date = ?
		WHERE id = ?
	`, req.Task, req.Description, req.Category, pkgutil.FormatTimeForDB(startTime),
		pkgutil.FormatTimeForDB(endTime), duration, currentDate, id)

	if err != nil {
		log.Printf("ERROR: Failed to update time entry ID %d - Task: %s, Category: %s, Start: %s, End: %s - Error: %v",
			id, req.Task, req.Category, req.StartTime, req.EndTime, err)
		return nil, fmt.Errorf("failed to update time entry: %w", err)
	}

	log.Printf("UPDATE: Modified time entry ID %d - Task: %s, Category: %s, Duration: %d min, Start: %s, End: %s",
		id, req.Task, req.Category, duration, startTime.Format("2006-01-02 15:04"), endTime.Format("2006-01-02 15:04"))

	return &model.TimeEntry{
		ID:          id,
		Task:        req.Task,
		Description: req.Description,
		Category:    req.Category,
		StartTime:   startTime,
		EndTime:     endTime,
	}, nil
}
