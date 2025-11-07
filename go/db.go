package timesheet

import (
	"log"
)

const CURRENT_DB_VERSION = 1

// GetTargetDBVersion returns the target database version for migration planning
func GetTargetDBVersion() int {
	return CURRENT_DB_VERSION
}

func InitDB() {
	var err error

	// Create version table first
	createVersionQuery := `
	CREATE TABLE IF NOT EXISTS db_version (
		id INTEGER PRIMARY KEY,
		version INTEGER NOT NULL,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createVersionQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Check current database version
	var currentVersion int
	err = db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM db_version").Scan(&currentVersion)
	if err != nil {
		log.Fatal(err)
	}

	// Apply migrations if needed
	if currentVersion < CURRENT_DB_VERSION {
		applyMigrations(currentVersion)
	}
}

func applyMigrations(fromVersion int) {
	log.Printf("Applying database migrations from version %d to %d", fromVersion, CURRENT_DB_VERSION)

	// Migration 1: Initial schema
	if fromVersion < 1 {
		applyMigration1()
		recordMigration(1)
	}

	log.Printf("Database migrations completed. Current version: %d", CURRENT_DB_VERSION)
}

func applyMigration1() {
	log.Println("Applying migration 1: Creating initial tables")

	createTimeEntryQuery := `
	CREATE TABLE IF NOT EXISTS time_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task TEXT NOT NULL,
		description TEXT,
		category TEXT NOT NULL DEFAULT 'other',
		start_time DATETIME,
		end_time DATETIME,
		duration INTEGER NOT NULL,
		date TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createTimeEntryQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Create categories table
	createCategoriesQuery := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		color TEXT NOT NULL DEFAULT '#718096',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createCategoriesQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Create tasks table
	createTasksQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		category_id INTEGER,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
	);`

	_, err = db.Exec(createTasksQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Migration: Add category column to existing entries if it doesn't exist
	// This will be ignored if the column already exists
	db.Exec("ALTER TABLE time_entries ADD COLUMN category TEXT DEFAULT 'other'")

	// Update any existing entries that have NULL category
	db.Exec("UPDATE time_entries SET category = 'other' WHERE category IS NULL")

	// Insert default categories if they don't exist
	db.Exec(`INSERT OR IGNORE INTO categories (name, color) VALUES 
		('project work', '#48bb78'),
		('project support', '#ed8936'),
		('other', '#718096')`)
}

func recordMigration(version int) {
	_, err := db.Exec("INSERT INTO db_version (version) VALUES (?)", version)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Migration %d applied successfully", version)
}
