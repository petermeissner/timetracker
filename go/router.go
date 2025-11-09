package timesheet

import (
	"io/fs"
	"net/http"
	pkgglobal "timesheet/go/global"
	pkghandler "timesheet/go/handler"

	"github.com/gorilla/mux"
)

func SetUpRouter() *mux.Router {
	r := mux.NewRouter()

	// Serve embedded static files
	staticFS, _ := fs.Sub(pkgglobal.StaticFiles, "static")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// API routes
	r.HandleFunc("/api/entries", pkghandler.GetTimeEntries).Methods("GET")
	r.HandleFunc("/api/entries", pkghandler.CreateTimeEntry).Methods("POST")
	r.HandleFunc("/api/entries/{id}", pkghandler.UpdateTimeEntry).Methods("PUT")
	r.HandleFunc("/api/entries/{id}", pkghandler.DeleteTimeEntry).Methods("DELETE")

	// Configuration API routes
	r.HandleFunc("/api/categories", pkghandler.GetCategories).Methods("GET")
	r.HandleFunc("/api/categories", pkghandler.CreateCategory).Methods("POST")
	r.HandleFunc("/api/categories/{id}", pkghandler.UpdateCategory).Methods("PUT")
	r.HandleFunc("/api/categories/{id}", pkghandler.DeleteCategory).Methods("DELETE")

	r.HandleFunc("/api/tasks", pkghandler.GetTasks).Methods("GET")
	r.HandleFunc("/api/tasks", pkghandler.CreateTask).Methods("POST")
	r.HandleFunc("/api/tasks/{id}", pkghandler.UpdateTask).Methods("PUT")
	r.HandleFunc("/api/tasks/{id}", pkghandler.DeleteTask).Methods("DELETE")

	// Serve HTML pages
	r.HandleFunc("/", pkghandler.ServeIndexHtml).Methods("GET")
	r.HandleFunc("/entries", pkghandler.ServeEntriesHtml).Methods("GET")
	r.HandleFunc("/config", pkghandler.ServeConfigHtml).Methods("GET")

	// Serve favicon
	r.HandleFunc("/favicon.ico", pkghandler.ServeFavicon).Methods("GET")

	return r
}
