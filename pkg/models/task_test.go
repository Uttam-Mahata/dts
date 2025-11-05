package models

import (
	"testing"
	"time"
)

func TestTaskStatus(t *testing.T) {
	task := &Task{
		ID:          "test-1",
		Name:        "Test Task",
		Status:      TaskStatusPending,
		CPURequired: 2.0,
		MemRequired: 1024.0,
		CreatedAt:   time.Now(),
	}

	if task.Status != TaskStatusPending {
		t.Errorf("Expected status %s, got %s", TaskStatusPending, task.Status)
	}

	task.Status = TaskStatusRunning
	if task.Status != TaskStatusRunning {
		t.Errorf("Expected status %s, got %s", TaskStatusRunning, task.Status)
	}
}

func TestTaskCreation(t *testing.T) {
	deadline := time.Now().Add(1 * time.Hour)
	task := &Task{
		ID:          "test-2",
		Name:        "Test Task 2",
		Priority:    5,
		CPURequired: 4.0,
		MemRequired: 2048.0,
		Duration:    30 * time.Second,
		Deadline:    deadline,
		Status:      TaskStatusPending,
		Metadata:    map[string]interface{}{"key": "value"},
	}

	if task.ID != "test-2" {
		t.Errorf("Expected ID test-2, got %s", task.ID)
	}

	if task.Priority != 5 {
		t.Errorf("Expected priority 5, got %d", task.Priority)
	}

	if task.CPURequired != 4.0 {
		t.Errorf("Expected CPU 4.0, got %f", task.CPURequired)
	}
}
