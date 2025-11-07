package timesheet

import "time"

type TimeEntry struct {
	ID          int       `json:"id"`
	Task        string    `json:"task"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

type TimeEntryRequest struct {
	Task        string `json:"task"`
	Description string `json:"description"`
	Category    string `json:"category"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type Category struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CategoryID  int    `json:"category_id"`
	Description string `json:"description"`
}

type CategoryRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TaskRequest struct {
	Name        string `json:"name"`
	CategoryID  int    `json:"category_id"`
	Description string `json:"description"`
}
