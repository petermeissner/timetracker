package timesheet

import (
	"database/sql"
	"embed"
)

// Package-level variables that will be set by main
var db *sql.DB
var staticFiles embed.FS

// SetDB sets the database connection for the handlers to use
func SetDB(database *sql.DB) {
	db = database
}

// SetStaticFiles sets the embedded static files for the handlers to use
func SetStaticFiles(files embed.FS) {
	staticFiles = files
}
