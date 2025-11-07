package timesheet

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
)

func SetUpRouter(staticFiles embed.FS) *mux.Router {
	r := mux.NewRouter()

	// Serve embedded static files
	staticFS, _ := fs.Sub(staticFiles, "static")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// API routes
	r.HandleFunc("/api/entries", GetTimeEntries).Methods("GET")
	r.HandleFunc("/api/entries", CreateTimeEntry).Methods("POST")
	r.HandleFunc("/api/entries/{id}", UpdateTimeEntry).Methods("PUT")
	r.HandleFunc("/api/entries/{id}", DeleteTimeEntry).Methods("DELETE")

	// Configuration API routes
	r.HandleFunc("/api/categories", GetCategories).Methods("GET")
	r.HandleFunc("/api/categories", CreateCategory).Methods("POST")
	r.HandleFunc("/api/categories/{id}", UpdateCategory).Methods("PUT")
	r.HandleFunc("/api/categories/{id}", DeleteCategory).Methods("DELETE")

	r.HandleFunc("/api/tasks", GetTasks).Methods("GET")
	r.HandleFunc("/api/tasks", CreateTask).Methods("POST")
	r.HandleFunc("/api/tasks/{id}", UpdateTask).Methods("PUT")
	r.HandleFunc("/api/tasks/{id}", DeleteTask).Methods("DELETE")

	// Serve HTML pages
	r.HandleFunc("/", ServeIndexHtml).Methods("GET")
	r.HandleFunc("/entries", ServeEntriesHtml).Methods("GET")
	r.HandleFunc("/config", ServeConfigHtml).Methods("GET")
	return r
}
