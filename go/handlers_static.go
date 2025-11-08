package timesheet

import (
	"net/http"
)

// Static file handlers
func ServeIndexHtml(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func ServeConfigHtml(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/config.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func ServeEntriesHtml(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/entries.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/favicon.ico")
	if err != nil {
		http.Error(w, "Favicon not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/x-icon")
	w.Write(data)
}
