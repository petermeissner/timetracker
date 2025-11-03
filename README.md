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

4. **Open your web browser and navigate to:**
   ```
   http://localhost:8080
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

## Building for Production

To build a standalone executable:

```powershell
go build -o timesheet.exe main.go
```

Then run:
```powershell
./timesheet.exe
```

## Troubleshooting

1. **Port already in use:** The application runs on port 8080. Make sure no other application is using this port.

2. **Database permissions:** Ensure the application has write permissions in the directory to create `timesheet.db`.

3. **Go dependencies:** If you encounter module errors, run `go mod tidy` to resolve dependencies.

## Future Enhancements

Potential improvements that could be added:

- User authentication and multiple users
- Time tracking with start/stop functionality
- Export data to CSV/PDF
- Time reporting and analytics
- Categories and projects
- Time goals and targets
- Dark mode theme
- Mobile app version