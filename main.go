package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"

	timesheet "timesheet/go"

	_ "modernc.org/sqlite"
)

//go:embed static/*
var staticFiles embed.FS

var db *sql.DB

func main() {
	var err error

	dbPath := "./timesheet.db"

	// Check database version and create backup if needed
	if err := timesheet.CheckAndBackupDatabase(dbPath); err != nil {
		log.Fatalf("Database backup failed: %v", err)
	}

	// connect to the database
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Set shared resources for the timesheet package
	timesheet.SetStaticFiles(staticFiles)
	timesheet.SetDB(db)

	// Initialize database
	timesheet.InitDB()
	defer db.Close()

	// Setup routes
	router := timesheet.SetUpRouter()

	// spin up server
	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", router))
}
