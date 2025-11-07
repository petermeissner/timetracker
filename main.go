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

	// connect to the database
	db, err = sql.Open("sqlite", "./timesheet.db")
	if err != nil {
		log.Fatal(err)
	}

	// make the db accessible to the timesheet package
	timesheet.SetDB(db)

	// Initialize database
	timesheet.InitDB()
	defer db.Close()

	// Set shared resources for the timesheet package
	timesheet.SetStaticFiles(staticFiles)

	// Setup routes
	r := timesheet.SetUpRouter(staticFiles)

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
