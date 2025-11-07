package timesheet

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Time entry handlers
func GetTimeEntries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query(`
		SELECT id, task, description, category, start_time, end_time 
		FROM time_entries 
		ORDER BY start_time DESC, id DESC
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []TimeEntry
	for rows.Next() {
		var entry TimeEntry
		var startTime, endTime sql.NullString

		err := rows.Scan(&entry.ID, &entry.Task, &entry.Description, &entry.Category,
			&startTime, &endTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if startTime.Valid {
			parsedStartTime, err := time.Parse(time.RFC3339, startTime.String)
			if err != nil {
				log.Printf("Error parsing start_time '%s': %v", startTime.String, err)
			} else {
				entry.StartTime = parsedStartTime
			}
		}
		if endTime.Valid {
			parsedEndTime, err := time.Parse(time.RFC3339, endTime.String)
			if err != nil {
				log.Printf("Error parsing end_time '%s': %v", endTime.String, err)
			} else {
				entry.EndTime = parsedEndTime
			}
		}

		entries = append(entries, entry)
	}

	json.NewEncoder(w).Encode(entries)
}

func CreateTimeEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req TimeEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Task == "" || req.Category == "" || req.StartTime == "" || req.EndTime == "" {
		http.Error(w, "Task, category, start_time, and end_time are required", http.StatusBadRequest)
		return
	}

	// Validate category exists in database
	var categoryExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE name = ?)", req.Category).Scan(&categoryExists)
	if err != nil {
		http.Error(w, "Database error while validating category", http.StatusInternalServerError)
		return
	}
	if !categoryExists {
		http.Error(w, "Invalid category. Category does not exist in the system", http.StatusBadRequest)
		return
	}

	// Parse start and end times (required fields)
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format. Expected ISO timestamp", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format. Expected ISO timestamp", http.StatusBadRequest)
		return
	}

	// Validate that end time is after start time
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	// Calculate duration from the time difference
	duration := int(endTime.Sub(startTime).Minutes())

	// Get current date for compatibility with existing database schema
	currentDate := time.Now().Format("2006-01-02")

	result, err := db.Exec(`
		INSERT INTO time_entries (task, description, category, start_time, end_time, duration, date)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, req.Task, req.Description, req.Category, startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339), duration, currentDate)

	if err != nil {
		log.Printf("ERROR: Failed to insert time entry - Task: %s, Category: %s, Start: %s, End: %s - Error: %v",
			req.Task, req.Category, req.StartTime, req.EndTime, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	log.Printf("INSERT: Created time entry ID %d - Task: %s, Category: %s, Duration: %d min, Start: %s, End: %s",
		id, req.Task, req.Category, duration, startTime.Format("2006-01-02 15:04"), endTime.Format("2006-01-02 15:04"))

	entry := TimeEntry{
		ID:          int(id),
		Task:        req.Task,
		Description: req.Description,
		Category:    req.Category,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	json.NewEncoder(w).Encode(entry)
}

func UpdateTimeEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req TimeEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Task == "" || req.Category == "" || req.StartTime == "" || req.EndTime == "" {
		http.Error(w, "Task, category, start_time, and end_time are required", http.StatusBadRequest)
		return
	}

	// Validate category exists in database
	var categoryExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE name = ?)", req.Category).Scan(&categoryExists)
	if err != nil {
		http.Error(w, "Database error while validating category", http.StatusInternalServerError)
		return
	}
	if !categoryExists {
		http.Error(w, "Invalid category. Category does not exist in the system", http.StatusBadRequest)
		return
	}

	// Parse start and end times (required fields)
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format. Expected ISO timestamp", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format. Expected ISO timestamp", http.StatusBadRequest)
		return
	}

	// Validate that end time is after start time
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	// Calculate duration from the time difference
	duration := int(endTime.Sub(startTime).Minutes())

	// Get current date for compatibility with existing database schema
	currentDate := time.Now().Format("2006-01-02")

	_, err = db.Exec(`
		UPDATE time_entries 
		SET task = ?, description = ?, category = ?, start_time = ?, end_time = ?, duration = ?, date = ?
		WHERE id = ?
	`, req.Task, req.Description, req.Category, startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339), duration, currentDate, id)

	if err != nil {
		log.Printf("ERROR: Failed to update time entry ID %d - Task: %s, Category: %s, Start: %s, End: %s - Error: %v",
			id, req.Task, req.Category, req.StartTime, req.EndTime, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("UPDATE: Modified time entry ID %d - Task: %s, Category: %s, Duration: %d min, Start: %s, End: %s",
		id, req.Task, req.Category, duration, startTime.Format("2006-01-02 15:04"), endTime.Format("2006-01-02 15:04"))

	entry := TimeEntry{
		ID:          id,
		Task:        req.Task,
		Description: req.Description,
		Category:    req.Category,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	json.NewEncoder(w).Encode(entry)
}

func DeleteTimeEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Get entry details before deletion for logging
	var task, category string
	var startTime, endTime sql.NullString
	err = db.QueryRow("SELECT task, category, start_time, end_time FROM time_entries WHERE id = ?", id).
		Scan(&task, &category, &startTime, &endTime)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("WARNING: Attempted to delete non-existent time entry ID %d", id)
			http.Error(w, "Time entry not found", http.StatusNotFound)
			return
		}
		log.Printf("ERROR: Failed to fetch time entry ID %d for deletion - Error: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM time_entries WHERE id = ?", id)
	if err != nil {
		log.Printf("ERROR: Failed to delete time entry ID %d - Error: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Format times for logging
	startTimeStr := "N/A"
	endTimeStr := "N/A"
	if startTime.Valid {
		if parsed, parseErr := time.Parse(time.RFC3339, startTime.String); parseErr == nil {
			startTimeStr = parsed.Format("2006-01-02 15:04")
		}
	}
	if endTime.Valid {
		if parsed, parseErr := time.Parse(time.RFC3339, endTime.String); parseErr == nil {
			endTimeStr = parsed.Format("2006-01-02 15:04")
		}
	}

	log.Printf("DELETE: Removed time entry ID %d - Task: %s, Category: %s, Time: %s to %s",
		id, task, category, startTimeStr, endTimeStr)

	w.WriteHeader(http.StatusNoContent)
}
