package model

// CronTask represents a scheduled cron task.
type CronTask struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Prompt        string  `json:"prompt"`
	ScheduleType  string  `json:"schedule_type"`
	ScheduleExpr  string  `json:"schedule_expr"`
	Status        string  `json:"status"`
	ThreadID      string  `json:"thread_id"`
	UserID        string  `json:"user_id"`
	Channel       string  `json:"channel"`
	LastRun       float64 `json:"last_run"`
	LastRunError  string  `json:"last_run_error"`
	LastRunOutput string  `json:"last_run_output"`
	NextRun       float64 `json:"next_run"`
	CreatedAt     float64 `json:"created_at"`
	UpdatedAt     float64 `json:"updated_at"`
}

// CronCreateRequest is the request body for creating a cron task.
type CronCreateRequest struct {
	Name         string `json:"name"`
	Prompt       string `json:"prompt"`
	ScheduleType string `json:"schedule_type,omitempty"` // default "cron"
	ScheduleExpr string `json:"schedule_expr"`
	Channel      string `json:"channel,omitempty"`
}

// CronUpdateRequest is the request body for updating a cron task.
type CronUpdateRequest struct {
	Name         string `json:"name,omitempty"`
	Prompt       string `json:"prompt,omitempty"`
	ScheduleType string `json:"schedule_type,omitempty"`
	ScheduleExpr string `json:"schedule_expr,omitempty"`
	Channel      string `json:"channel,omitempty"`
}
