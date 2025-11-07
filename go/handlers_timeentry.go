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
		SELECT id, task, description, category, start_time, end_time, duration, date 
		FROM time_entries 
		ORDER BY date DESC, id DESC
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
			&startTime, &endTime, &entry.Duration, &entry.Date)
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
	if req.Task == "" || req.Duration <= 0 || req.Date == "" || req.Category == "" {
		http.Error(w, "Task, category, duration, and date are required", http.StatusBadRequest)
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

	// Calculate start and end times
	var startTime, endTime time.Time

	if req.StartTime != "" && req.EndTime != "" {
		// Use provided start and end times from time slot selection (ISO format)
		var err error
		startTime, err = time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			http.Error(w, "Invalid start time format. Expected ISO timestamp", http.StatusBadRequest)
			return
		}

		endTime, err = time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			http.Error(w, "Invalid end time format. Expected ISO timestamp", http.StatusBadRequest)
			return
		}

		// Calculate duration from the time difference
		req.Duration = int(endTime.Sub(startTime).Minutes())
	} else {
		// Fall back to calculating from duration (legacy behavior)
		now := time.Now()
		startTime = now.Add(-time.Duration(req.Duration) * time.Minute)
		endTime = now
	}

	result, err := db.Exec(`
		INSERT INTO time_entries (task, description, category, start_time, end_time, duration, date)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, req.Task, req.Description, req.Category, startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339), req.Duration, req.Date)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	entry := TimeEntry{
		ID:          int(id),
		Task:        req.Task,
		Description: req.Description,
		Category:    req.Category,
		StartTime:   startTime,
		EndTime:     endTime,
		Duration:    req.Duration,
		Date:        req.Date,
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
	if req.Task == "" || req.Duration <= 0 || req.Date == "" || req.Category == "" {
		http.Error(w, "Task, category, duration, and date are required", http.StatusBadRequest)
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

	// Calculate start and end times based on duration
	now := time.Now()
	startTime := now.Add(-time.Duration(req.Duration) * time.Minute)
	endTime := now

	_, err = db.Exec(`
		UPDATE time_entries 
		SET task = ?, description = ?, category = ?, start_time = ?, end_time = ?, duration = ?, date = ?
		WHERE id = ?
	`, req.Task, req.Description, req.Category, startTime.Format("2006-01-02 15:04:05"),
		endTime.Format("2006-01-02 15:04:05"), req.Duration, req.Date, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	entry := TimeEntry{
		ID:          id,
		Task:        req.Task,
		Description: req.Description,
		Category:    req.Category,
		StartTime:   startTime,
		EndTime:     endTime,
		Duration:    req.Duration,
		Date:        req.Date,
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

	_, err = db.Exec("DELETE FROM time_entries WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
