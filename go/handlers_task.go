package timesheet

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Task handlers
func GetTasks(w http.ResponseWriter, r *http.Request) {
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

func CreateTask(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("ERROR: Failed to insert task - Name: %s, CategoryID: %v, Description: %s - Error: %v",
			req.Name, req.CategoryID, req.Description, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	log.Printf("INSERT: Created task ID %d - Name: %s, CategoryID: %d, Description: %s",
		id, req.Name, req.CategoryID, req.Description)
	task := Task{
		ID:          int(id),
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		Description: req.Description,
	}

	json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("ERROR: Failed to update task ID %d - Name: %s, CategoryID: %v, Description: %s - Error: %v",
			id, req.Name, req.CategoryID, req.Description, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("UPDATE: Modified task ID %d - Name: %s, CategoryID: %d, Description: %s",
		id, req.Name, req.CategoryID, req.Description)

	task := Task{
		ID:          id,
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		Description: req.Description,
	}

	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Get task details before deletion for logging
	var name, description string
	var categoryID sql.NullInt64
	err = db.QueryRow("SELECT name, category_id, description FROM tasks WHERE id = ?", id).
		Scan(&name, &categoryID, &description)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("WARNING: Attempted to delete non-existent task ID %d", id)
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		log.Printf("ERROR: Failed to fetch task ID %d for deletion - Error: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		log.Printf("ERROR: Failed to delete task ID %d - Error: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categoryIDVal := 0
	if categoryID.Valid {
		categoryIDVal = int(categoryID.Int64)
	}
	log.Printf("DELETE: Removed task ID %d - Name: %s, CategoryID: %d, Description: %s",
		id, name, categoryIDVal, description)

	w.WriteHeader(http.StatusNoContent)
}
