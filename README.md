# Timesheet Tracker

A simple web application for tracking time spent on different tasks, built with Go backend, HTML/CSS/JavaScript frontend, and SQLite database.

## Features

- ✅ Add time entries with task name, description, duration, and date
- ✅ View all time entries grouped by date
- ✅ Edit existing time entries
- ✅ Delete time entries
- ✅ Calculate total time spent per day
- ✅ Responsive web design
- ✅ SQLite database storage
- ✅ RESTful API

## Project Structure

```
timesheet/
├── main.go              # Go server with API endpoints
├── go.mod              # Go module dependencies
├── static/             # Frontend assets
│   ├── index.html      # Main HTML page
│   ├── styles.css      # CSS styling
│   └── script.js       # JavaScript functionality
├── timesheet.db        # SQLite database (created automatically)
└── README.md           # This file
```

## Prerequisites

- Go 1.21 or higher
- Git (for version control)

## Installation and Setup

1. **Navigate to the project directory:**
   ```powershell
   cd c:\Users\peter\git\timesheet
   ```

2. **Initialize Go modules and download dependencies:**
   ```powershell
   go mod tidy
   ```

3. **Run the application:**
   ```powershell
   go run main.go
   ```

   Or with custom parameters:
   ```powershell
   go run main.go -port 8081 -db ./custom.db
   ```

4. **Open your web browser and navigate to:**
   ```
   http://localhost:8080
   ```
   (or the custom port you specified)

## Command-Line Parameters

The application supports the following command-line parameters and environment variables:

### Command-Line Flags:
- `-port` - Port to run the server on (default: "8080")
- `-db` - Path to the SQLite database file (default: "./timesheet.db")
- `-help` - Show usage information

### Environment Variables:
- `PORT` - Port to run the server on (overridden by -port flag)
- `DB_PATH` - Path to the SQLite database file (overridden by -db flag)

### Examples:

**Using Command-Line Flags:**
```powershell
# Use default settings
go run main.go

# Use custom port
go run main.go -port 8081

# Use custom database file
go run main.go -db ./my-timesheet.db

# Use custom port and database
go run main.go -port 8081 -db ./my-timesheet.db

# Show help
go run main.go -help
```

**Using Environment Variables:**
```powershell
# PowerShell syntax
$env:PORT="8081"; go run main.go
$env:DB_PATH="./my-timesheet.db"; go run main.go
$env:PORT="8081"; $env:DB_PATH="./my-timesheet.db"; go run main.go

# Command Prompt syntax
set PORT=8081 && go run main.go
set DB_PATH=./my-timesheet.db && go run main.go
```

**Precedence:** Command-line flags take priority over environment variables.

### Built Executable Usage:

After building the application (`go build .`), you can use the same parameters:

```powershell
# Default settings
.\timesheet.exe

# Custom port and database
.\timesheet.exe -port 8081 -db ./my-timesheet.db
```

## API Endpoints

The application provides the following REST API endpoints:

- `GET /` - Serve the main HTML page
- `GET /api/entries` - Get all time entries
- `POST /api/entries` - Create a new time entry
- `PUT /api/entries/{id}` - Update an existing time entry
- `DELETE /api/entries/{id}` - Delete a time entry

### API Request/Response Examples

**Create a time entry (POST /api/entries):**
```json
{
  "task": "Development",
  "description": "Working on timesheet application",
  "duration": 120,
  "date": "2025-11-03"
}
```

**Response:**
```json
{
  "id": 1,
  "task": "Development",
  "description": "Working on timesheet application",
  "start_time": "2025-11-03T14:00:00Z",
  "end_time": "2025-11-03T16:00:00Z",
  "duration": 120,
  "date": "2025-11-03"
}
```

## Database Schema

The application uses SQLite with the following table structure:

```sql
CREATE TABLE time_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task TEXT NOT NULL,
    description TEXT,
    start_time DATETIME,
    end_time DATETIME,
    duration INTEGER NOT NULL,
    date TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Usage

1. **Adding a Time Entry:**
   - Fill in the task name (required)
   - Add an optional description
   - Enter the duration in minutes
   - Select the date
   - Click "Add Entry"

2. **Viewing Entries:**
   - All entries are displayed grouped by date
   - Total time for each day is shown
   - Current day's total is highlighted at the top

3. **Editing Entries:**
   - Click the "Edit" button on any entry
   - Modify the details in the modal dialog
   - Click "Update Entry" to save changes

4. **Deleting Entries:**
   - Click the "Delete" button on any entry
   - Confirm the deletion in the dialog

## Development

To modify the application:

1. **Backend changes:** Edit `main.go` and restart the server
2. **Frontend changes:** Edit files in the `static/` directory and refresh the browser
3. **Database changes:** The database is created automatically on first run

### Development History

For a complete record of the project's development process, including all user requests and implementation decisions, see [ai_chat.md](ai_chat.md).

## Building for Production

The application can be built into a single, self-contained executable that includes all static files (HTML, CSS, JavaScript) embedded within the binary. This makes distribution extremely easy.

### Build Single Executable

To build a standalone executable with all assets embedded:

```powershell
go build -o timesheet.exe main.go
```

Then run:
```powershell
./timesheet.exe
```

Or with custom parameters:
```powershell
./timesheet.exe -port 8081 -db ./production.db
```

### Cross-Platform Builds

Use the automated build script to create binaries for all platforms:

```powershell
.\build-cross-platform.ps1
```

This script builds optimized binaries for:
- **Windows** (amd64) - `timesheet-windows-amd64.exe`
- **Linux** (amd64) - `timesheet-linux-amd64` 
- **Linux** (arm64) - `timesheet-linux-arm64` (Raspberry Pi, ARM servers)
- **macOS** (amd64) - `timesheet-macos-amd64` (Intel Macs)
- **macOS** (arm64) - `timesheet-macos-arm64` (Apple Silicon Macs)

All binaries are placed in the `distribution/` folder and include size optimization (`-ldflags "-s -w"`).

### Distribution Benefits

- ✅ **Single file distribution** - No need to bundle static folder
- ✅ **No external dependencies** - All assets embedded in binary
- ✅ **Cross-platform** - Build for Windows, Linux, macOS
- ✅ **Easy deployment** - Just copy the executable
- ✅ **No installation required** - Run directly

The executable will create the SQLite database file (`timesheet.db`) in the same directory when first run.

### How File Embedding Works

The application uses Go 1.16+ `embed` package to bundle all static files directly into the binary:

```go
//go:embed static/*
var staticFiles embed.FS
```

This means:
- All files in the `static/` directory are compiled into the executable
- No separate `static/` folder needed for distribution
- Files are served from memory, improving performance
- Works across all platforms without path issues

## Troubleshooting

1. **Port already in use:** The application runs on port 8080. Make sure no other application is using this port.

2. **Database permissions:** Ensure the application has write permissions in the directory to create `timesheet.db`.

3. **Go dependencies:** If you encounter module errors, run `go mod tidy` to resolve dependencies.

## Future Enhancements

Potential improvements that could be added:

- User authentication and multiple users
- Time reporting and analytics
- Time goals and targets