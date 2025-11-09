package global

import (
	"database/sql"
	"embed"
)

// Package-level variables that will be set by main
var Db *sql.DB
var StaticFiles embed.FS

// SetDB sets the database connection for the handlers to use
func SetDB(database *sql.DB) {
	Db = database
}

// SetStaticFiles sets the embedded static files for the handlers to use
func SetStaticFiles(files embed.FS) {
	StaticFiles = files
}
