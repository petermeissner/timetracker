package handler

import (
	"net/http"
	pkgglobal "timesheet/go/global"
)

// Static file handlers
func ServeIndexHtml(w http.ResponseWriter, r *http.Request) {
	data, err := pkgglobal.StaticFiles.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func ServeConfigHtml(w http.ResponseWriter, r *http.Request) {
	data, err := pkgglobal.StaticFiles.ReadFile("static/config.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func ServeEntriesHtml(w http.ResponseWriter, r *http.Request) {
	data, err := pkgglobal.StaticFiles.ReadFile("static/entries.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	data, err := pkgglobal.StaticFiles.ReadFile("static/favicon.ico")
	if err != nil {
		http.Error(w, "Favicon not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/x-icon")
	w.Write(data)
}
