package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
)

//go:embed static/*
var staticFiles embed.FS

type TimeEntry struct {
	ID          int       `json:"id"`
	Task        string    `json:"task"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Duration    int       `json:"duration"` // in minutes
	Date        string    `json:"date"`
}

type TimeEntryRequest struct {
	Task        string `json:"task"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Duration    int    `json:"duration"` // in minutes
	Date        string `json:"date"`
	StartTime   string `json:"start_time,omitempty"`
	EndTime     string `json:"end_time,omitempty"`
}

type Category struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CategoryID  int    `json:"category_id"`
	Description string `json:"description"`
}

type CategoryRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TaskRequest struct {
	Name        string `json:"name"`
	CategoryID  int    `json:"category_id"`
	Description string `json:"description"`
}

var db *sql.DB

func main() {
	// Initialize database
	initDB()
	defer db.Close()

	// Setup routes
	r := mux.NewRouter()

	// Serve embedded static files
	staticFS, _ := fs.Sub(staticFiles, "static")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// API routes
	r.HandleFunc("/api/entries", getTimeEntries).Methods("GET")
	r.HandleFunc("/api/entries", createTimeEntry).Methods("POST")
	r.HandleFunc("/api/entries/{id}", updateTimeEntry).Methods("PUT")
	r.HandleFunc("/api/entries/{id}", deleteTimeEntry).Methods("DELETE")

	// Configuration API routes
	r.HandleFunc("/api/categories", getCategories).Methods("GET")
	r.HandleFunc("/api/categories", createCategory).Methods("POST")
	r.HandleFunc("/api/categories/{id}", updateCategory).Methods("PUT")
	r.HandleFunc("/api/categories/{id}", deleteCategory).Methods("DELETE")

	r.HandleFunc("/api/tasks", getTasks).Methods("GET")
	r.HandleFunc("/api/tasks", createTask).Methods("POST")
	r.HandleFunc("/api/tasks/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/api/tasks/{id}", deleteTask).Methods("DELETE")

	// Serve HTML pages
	r.HandleFunc("/", serveIndex).Methods("GET")
	r.HandleFunc("/entries", serveEntries).Methods("GET")
	r.HandleFunc("/config", serveConfig).Methods("GET")

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./timesheet.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS time_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task TEXT NOT NULL,
		description TEXT,
		category TEXT NOT NULL DEFAULT 'other',
		start_time DATETIME,
		end_time DATETIME,
		duration INTEGER NOT NULL,
		date TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Create categories table
	createCategoriesQuery := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		color TEXT NOT NULL DEFAULT '#718096',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createCategoriesQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Create tasks table
	createTasksQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		category_id INTEGER,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
	);`

	_, err = db.Exec(createTasksQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Migration: Add category column to existing entries if it doesn't exist
	// This will be ignored if the column already exists
	db.Exec("ALTER TABLE time_entries ADD COLUMN category TEXT DEFAULT 'other'")

	// Update any existing entries that have NULL category
	db.Exec("UPDATE time_entries SET category = 'other' WHERE category IS NULL")

	// Insert default categories if they don't exist
	db.Exec(`INSERT OR IGNORE INTO categories (name, color) VALUES 
		('project work', '#48bb78'),
		('project support', '#ed8936'),
		('other', '#718096')`)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func getTimeEntries(w http.ResponseWriter, r *http.Request) {
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

func createTimeEntry(w http.ResponseWriter, r *http.Request) {
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

	// Validate category
	validCategories := map[string]bool{
		"project work":    true,
		"project support": true,
		"other":           true,
	}
	if !validCategories[req.Category] {
		http.Error(w, "Invalid category. Must be 'project work', 'project support', or 'other'", http.StatusBadRequest)
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

func updateTimeEntry(w http.ResponseWriter, r *http.Request) {
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

	// Validate category
	validCategories := map[string]bool{
		"project work":    true,
		"project support": true,
		"other":           true,
	}
	if !validCategories[req.Category] {
		http.Error(w, "Invalid category. Must be 'project work', 'project support', or 'other'", http.StatusBadRequest)
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

func deleteTimeEntry(w http.ResponseWriter, r *http.Request) {
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

func serveConfig(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/config.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func serveEntries(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/entries.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

// Category handlers
func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, name, color FROM categories ORDER BY name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name, &category.Color)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	json.NewEncoder(w).Encode(categories)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	if req.Color == "" {
		req.Color = "#718096" // Default color
	}

	result, err := db.Exec("INSERT INTO categories (name, color) VALUES (?, ?)", req.Name, req.Color)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	category := Category{
		ID:    int(id),
		Name:  req.Name,
		Color: req.Color,
	}

	json.NewEncoder(w).Encode(category)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	if req.Color == "" {
		req.Color = "#718096" // Default color
	}

	_, err = db.Exec("UPDATE categories SET name = ?, color = ? WHERE id = ?", req.Name, req.Color, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	category := Category{
		ID:    id,
		Name:  req.Name,
		Color: req.Color,
	}

	json.NewEncoder(w).Encode(category)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Task handlers
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, name, category_id, description FROM tasks ORDER BY name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var categoryID sql.NullInt64
		err := rows.Scan(&task.ID, &task.Name, &categoryID, &task.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if categoryID.Valid {
			task.CategoryID = int(categoryID.Int64)
		}
		tasks = append(tasks, task)
	}

	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Task name is required", http.StatusBadRequest)
		return
	}

	var categoryID interface{} = nil
	if req.CategoryID > 0 {
		categoryID = req.CategoryID
	}

	result, err := db.Exec("INSERT INTO tasks (name, category_id, description) VALUES (?, ?, ?)",
		req.Name, categoryID, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	task := Task{
		ID:          int(id),
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		Description: req.Description,
	}

	json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Task name is required", http.StatusBadRequest)
		return
	}

	var categoryID interface{} = nil
	if req.CategoryID > 0 {
		categoryID = req.CategoryID
	}

	_, err = db.Exec("UPDATE tasks SET name = ?, category_id = ?, description = ? WHERE id = ?",
		req.Name, categoryID, req.Description, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task := Task{
		ID:          id,
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		Description: req.Description,
	}

	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
