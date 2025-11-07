package timesheet

import (
	"log"
)

func InitDB() {
	var err error

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

	_, err = db.Exec(createTimeEntryQuery)
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
