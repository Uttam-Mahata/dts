package models

// Schedule represents a task assignment to an edge node
type Schedule struct {
	TaskID     string  `json:"task_id"`
	NodeID     string  `json:"node_id"`
	StartTime  int64   `json:"start_time"`
	Priority   int     `json:"priority"`
	Cost       float64 `json:"cost"` // Cost metric for optimization
	Confidence float64 `json:"confidence"`
}

// ScheduleResult represents the result of scheduling operation
type ScheduleResult struct {
	Success   bool       `json:"success"`
	Schedules []Schedule `json:"schedules"`
	Message   string     `json:"message,omitempty"`
}
