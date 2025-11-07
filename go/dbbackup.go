package timesheet

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"time"
)

// CheckAndBackupDatabase checks if there's a version difference and creates a backup if needed
func CheckAndBackupDatabase(dbPath string) error {
	// Check if database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("No database file found, nothing to back up.")
		return nil
	}

	// Open database to check current version
	tempDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database for version check: %v", err)
	}
	defer tempDB.Close()

	// Check if version table exists and get current version
	var currentVersion int
	err = tempDB.QueryRow("SELECT COALESCE(MAX(version), 0) FROM db_version").Scan(&currentVersion)
	if err != nil {
		// If version table doesn't exist, assume version 0
		currentVersion = 0
		fmt.Println("No version table found, assuming database version 0")
	}

	// Get target version from timesheet package
	targetVersion := GetTargetDBVersion()
	fmt.Printf("Current database version: %d, Target version: %d\n", currentVersion, targetVersion)

	// If versions differ, create backup
	if currentVersion != targetVersion {
		backupPath := fmt.Sprintf("timesheet_backup_v%d_%s.db", currentVersion, time.Now().Format("20060102_150405"))
		fmt.Printf("Version difference detected. Creating backup: %s\n", backupPath)

		if err := copyFile(dbPath, backupPath); err != nil {
			return fmt.Errorf("failed to create database backup: %v", err)
		}
		fmt.Printf("Database backup created successfully: %s\n", backupPath)
	} else {
		fmt.Println("Database version matches target, no backup needed")
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}
