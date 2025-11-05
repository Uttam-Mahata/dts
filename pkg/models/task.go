package models

import (
	"time"
)

// TaskStatus represents the current status of a task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusScheduled TaskStatus = "scheduled"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

// Task represents a computational task in the distributed system
type Task struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Priority    int                    `json:"priority"`
	CPURequired float64                `json:"cpu_required"`
	MemRequired float64                `json:"mem_required"`
	Duration    time.Duration          `json:"duration"`
	Deadline    time.Time              `json:"deadline"`
	Status      TaskStatus             `json:"status"`
	AssignedTo  string                 `json:"assigned_to,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// TaskRequest represents a task submission request
type TaskRequest struct {
	Name        string                 `json:"name"`
	Priority    int                    `json:"priority"`
	CPURequired float64                `json:"cpu_required"`
	MemRequired float64                `json:"mem_required"`
	Duration    int64                  `json:"duration"` // in seconds
	Deadline    string                 `json:"deadline"` // RFC3339 format
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
