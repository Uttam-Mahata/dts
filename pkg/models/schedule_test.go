package models

import (
	"testing"
)

func TestScheduleCreation(t *testing.T) {
	schedule := Schedule{
		TaskID:     "task-1",
		NodeID:     "node-1",
		StartTime:  1234567890,
		Priority:   5,
		Cost:       0.75,
		Confidence: 0.95,
	}

	if schedule.TaskID != "task-1" {
		t.Errorf("Expected task ID 'task-1', got '%s'", schedule.TaskID)
	}

	if schedule.NodeID != "node-1" {
		t.Errorf("Expected node ID 'node-1', got '%s'", schedule.NodeID)
	}

	if schedule.Priority != 5 {
		t.Errorf("Expected priority 5, got %d", schedule.Priority)
	}

	if schedule.Cost != 0.75 {
		t.Errorf("Expected cost 0.75, got %f", schedule.Cost)
	}

	if schedule.Confidence != 0.95 {
		t.Errorf("Expected confidence 0.95, got %f", schedule.Confidence)
	}
}

func TestScheduleResult(t *testing.T) {
	schedules := []Schedule{
		{
			TaskID:     "task-1",
			NodeID:     "node-1",
			Priority:   5,
			Cost:       0.75,
			Confidence: 0.95,
		},
		{
			TaskID:     "task-2",
			NodeID:     "node-2",
			Priority:   3,
			Cost:       0.80,
			Confidence: 0.90,
		},
	}

	result := ScheduleResult{
		Success:   true,
		Schedules: schedules,
		Message:   "Successfully scheduled 2 tasks",
	}

	if !result.Success {
		t.Error("Expected result to be successful")
	}

	if len(result.Schedules) != 2 {
		t.Errorf("Expected 2 schedules, got %d", len(result.Schedules))
	}

	if result.Message != "Successfully scheduled 2 tasks" {
		t.Errorf("Expected message 'Successfully scheduled 2 tasks', got '%s'", result.Message)
	}
}

func TestScheduleResultFailure(t *testing.T) {
	result := ScheduleResult{
		Success:   false,
		Schedules: []Schedule{},
		Message:   "No active nodes available",
	}

	if result.Success {
		t.Error("Expected result to be unsuccessful")
	}

	if len(result.Schedules) != 0 {
		t.Errorf("Expected 0 schedules, got %d", len(result.Schedules))
	}

	if result.Message != "No active nodes available" {
		t.Errorf("Expected message 'No active nodes available', got '%s'", result.Message)
	}
}

func TestScheduleResultEmpty(t *testing.T) {
	result := ScheduleResult{
		Success:   true,
		Schedules: []Schedule{},
		Message:   "No pending tasks to schedule",
	}

	if !result.Success {
		t.Error("Expected result to be successful")
	}

	if len(result.Schedules) != 0 {
		t.Errorf("Expected 0 schedules, got %d", len(result.Schedules))
	}
}
