package timesheet

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Category handlers
func GetCategories(w http.ResponseWriter, r *http.Request) {
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

func CreateCategory(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("ERROR: Failed to insert category - Name: %s, Color: %s - Error: %v", req.Name, req.Color, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	log.Printf("INSERT: Created category ID %d - Name: %s, Color: %s", id, req.Name, req.Color)
	category := Category{
		ID:    int(id),
		Name:  req.Name,
		Color: req.Color,
	}

	json.NewEncoder(w).Encode(category)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("ERROR: Failed to update category ID %d - Name: %s, Color: %s - Error: %v", id, req.Name, req.Color, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("UPDATE: Modified category ID %d - Name: %s, Color: %s", id, req.Name, req.Color)

	category := Category{
		ID:    id,
		Name:  req.Name,
		Color: req.Color,
	}

	json.NewEncoder(w).Encode(category)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Get category details before deletion for logging
	var name, color string
	err = db.QueryRow("SELECT name, color FROM categories WHERE id = ?", id).Scan(&name, &color)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Printf("WARNING: Attempted to delete non-existent category ID %d", id)
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		log.Printf("ERROR: Failed to fetch category ID %d for deletion - Error: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		log.Printf("ERROR: Failed to delete category ID %d - Error: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("DELETE: Removed category ID %d - Name: %s, Color: %s", id, name, color)

	w.WriteHeader(http.StatusNoContent)
}
