package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	timesheet "timesheet/go"
	pkgdb "timesheet/go/db"
	tserverconfig "timesheet/go/serverconfig"

	pkgglobal "timesheet/go/global"

	_ "modernc.org/sqlite"
)

//go:embed static/*
var mainStaticFiles embed.FS

var mainDb *sql.DB

func main() {
	// Define command-line flags with environment variable fallbacks
	var dbPath = flag.String("db", tserverconfig.GetEnvOrDefault("DB_PATH", "./timesheet.db"), "Path to the SQLite database file")
	var port = flag.String("port", tserverconfig.GetEnvOrDefault("PORT", "8080"), "Port to run the server on")
	var help = flag.Bool("help", false, "Show usage information")

	// Parse command-line flags
	flag.Parse()

	// Show help if requested
	if *help {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  PORT      Port to run the server on (overridden by -port flag)\n")
		fmt.Fprintf(os.Stderr, "  DB_PATH   Path to the SQLite database file (overridden by -db flag)\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                              # Use default database and port\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -port 8081                   # Use port 8081\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -db ./custom.db              # Use custom database file\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -db ./custom.db -port 8081   # Use custom database and port\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n  # Using environment variables:\n")
		fmt.Fprintf(os.Stderr, "  PORT=8081 %s                    # Use port 8081\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  DB_PATH=./custom.db %s          # Use custom database file\n", os.Args[0])
		os.Exit(0)
	}

	var err error

	// Check database version and create backup if needed
	if err := pkgdb.CheckAndBackupDatabase(*dbPath); err != nil {
		log.Fatalf("Database backup failed: %v", err)
	}

	// connect to the database
	mainDb, err = sql.Open("sqlite", *dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Set shared resources for the timesheet package
	pkgglobal.SetStaticFiles(mainStaticFiles)
	pkgglobal.SetDB(mainDb)

	// Initialize database
	pkgdb.InitDB()
	defer mainDb.Close()

	// Setup routes
	router := timesheet.SetUpRouter()

	// spin up server
	serverAddr := "127.0.0.1:" + *port
	fmt.Printf("Server starting on http://localhost:%s\n", *port)
	fmt.Printf("Database file: %s\n", *dbPath)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
