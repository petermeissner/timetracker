package db

import (
	"database/sql"
	"testing"
	"timesheet/go/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	// Create categories table
	_, err = db.Exec(createTableCategories)
	require.NoError(t, err)

	// Create time_entries table
	_, err = db.Exec(createTableTimeEntries)
	require.NoError(t, err)

	// Insert test categories
	_, err = db.Exec("INSERT INTO categories (name) VALUES ('project work'), ('project support'), ('maintenance')")
	require.NoError(t, err)

	return db
}

func TestParseNullableTimeValidRFC3339(t *testing.T) {
	input := sql.NullString{String: "2025-11-09T14:30:00Z", Valid: true}
	result, err := ParseNullableTime(input)

	require.NoError(t, err)
	assert.False(t, result.IsZero())
}

func TestParseNullableTimeNullTime(t *testing.T) {
	input := sql.NullString{Valid: false}
	result, err := ParseNullableTime(input)

	// Note: Based on the actual implementation, null time returns time.Time{}, nil
	require.NoError(t, err)
	assert.True(t, result.IsZero())
}

func TestParseNullableTimeInvalidTimeFormat(t *testing.T) {
	input := sql.NullString{String: "invalid-time", Valid: true}
	result, err := ParseNullableTime(input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing time")
	assert.True(t, result.IsZero())
}

func TestParseNullableTimeEmptyString(t *testing.T) {
	input := sql.NullString{String: "", Valid: true}
	result, err := ParseNullableTime(input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing time")
	assert.True(t, result.IsZero())
}

func TestValidateCategoryExistsValid(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := ValidateCategoryExists(db, "project work")
	require.NoError(t, err)
}

func TestValidateCategoryExistsAnotherValid(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := ValidateCategoryExists(db, "project support")
	require.NoError(t, err)
}

func TestValidateCategoryExistsNonExistent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := ValidateCategoryExists(db, "non existent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid category 'non existent': category does not exist in the system")
}

func TestValidateCategoryExistsEmpty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := ValidateCategoryExists(db, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid category '': category does not exist in the system")
}

func TestValidateCategoryExistsCaseSensitive(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := ValidateCategoryExists(db, "Project Work")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid category 'Project Work': category does not exist in the system")
}

func TestCreateTimeEntryInDBValidEntry(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := model.TimeEntryRequest{
		Task:      "Development",
		Category:  "project work",
		StartTime: "2025-11-09T09:00:00Z",
		EndTime:   "2025-11-09T10:00:00Z",
	}

	entry, err := CreateTimeEntryInDB(db, req)
	require.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Greater(t, entry.ID, 0)
	assert.Equal(t, req.Task, entry.Task)
	assert.Equal(t, req.Category, entry.Category)
}

func TestCreateTimeEntryInDBInvalidCategory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := model.TimeEntryRequest{
		Task:      "Development",
		Category:  "invalid category",
		StartTime: "2025-11-09T09:00:00Z",
		EndTime:   "2025-11-09T10:00:00Z",
	}

	entry, err := CreateTimeEntryInDB(db, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid category 'invalid category': category does not exist in the system")
	assert.Nil(t, entry)
}

func TestCreateTimeEntryInDBInvalidTimeFormat(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := model.TimeEntryRequest{
		Task:      "Development",
		Category:  "project work",
		StartTime: "invalid-time",
		EndTime:   "2025-11-09T10:00:00Z",
	}

	entry, err := CreateTimeEntryInDB(db, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid start time format")
	assert.Nil(t, entry)
}

func TestCreateTimeEntryInDBMissingRequiredFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := model.TimeEntryRequest{
		Category:  "project work",
		StartTime: "2025-11-09T09:00:00Z",
		EndTime:   "2025-11-09T10:00:00Z",
	}

	entry, err := CreateTimeEntryInDB(db, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "task is required")
	assert.Nil(t, entry)
}

func TestUpdateTimeEntryInDBValidUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a time entry to update
	createReq := model.TimeEntryRequest{
		Task:      "Original Task",
		Category:  "project work",
		StartTime: "2025-11-09T09:00:00Z",
		EndTime:   "2025-11-09T10:00:00Z",
	}

	originalEntry, err := CreateTimeEntryInDB(db, createReq)
	require.NoError(t, err)
	require.NotNil(t, originalEntry)

	// Update the entry
	updateReq := model.TimeEntryRequest{
		Task:      "Updated Task",
		Category:  "project support",
		StartTime: "2025-11-09T14:00:00Z",
		EndTime:   "2025-11-09T15:30:00Z",
	}

	entry, err := UpdateTimeEntryInDB(db, originalEntry.ID, updateReq)
	require.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, originalEntry.ID, entry.ID)
	assert.Equal(t, updateReq.Task, entry.Task)
	assert.Equal(t, updateReq.Category, entry.Category)
}

func TestUpdateTimeEntryInDBInvalidCategory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a time entry to update
	createReq := model.TimeEntryRequest{
		Task:      "Original Task",
		Category:  "project work",
		StartTime: "2025-11-09T09:00:00Z",
		EndTime:   "2025-11-09T10:00:00Z",
	}

	originalEntry, err := CreateTimeEntryInDB(db, createReq)
	require.NoError(t, err)
	require.NotNil(t, originalEntry)

	// Try to update with invalid category
	updateReq := model.TimeEntryRequest{
		Task:      "Updated Task",
		Category:  "invalid category",
		StartTime: "2025-11-09T14:00:00Z",
		EndTime:   "2025-11-09T15:30:00Z",
	}

	entry, err := UpdateTimeEntryInDB(db, originalEntry.ID, updateReq)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid category 'invalid category': category does not exist in the system")
	assert.Nil(t, entry)
}
