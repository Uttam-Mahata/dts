package scheduler

import (
	"testing"
	"time"

	"github.com/Uttam-Mahata/dts/pkg/edgenode"
	"github.com/Uttam-Mahata/dts/pkg/models"
)

func TestSubmitTask(t *testing.T) {
	nodeManager := edgenode.NewManager()
	scheduler := NewScheduler(nodeManager)

	task := &models.Task{
		ID:          "task-1",
		Name:        "Test Task",
		Priority:    1,
		CPURequired: 2.0,
		MemRequired: 1024.0,
		Duration:    10 * time.Second,
		Deadline:    time.Now().Add(1 * time.Hour),
	}

	err := scheduler.SubmitTask(task)
	if err != nil {
		t.Fatalf("Failed to submit task: %v", err)
	}

	if scheduler.GetTaskCount() != 1 {
		t.Errorf("Expected 1 task, got %d", scheduler.GetTaskCount())
	}

	// Try to submit the same task again
	err = scheduler.SubmitTask(task)
	if err == nil {
		t.Error("Expected error when submitting duplicate task")
	}
}

func TestGetTask(t *testing.T) {
	nodeManager := edgenode.NewManager()
	scheduler := NewScheduler(nodeManager)

	task := &models.Task{
		ID:          "task-2",
		Name:        "Test Task 2",
		Priority:    2,
		CPURequired: 1.0,
		MemRequired: 512.0,
		Duration:    5 * time.Second,
		Deadline:    time.Now().Add(30 * time.Minute),
	}

	scheduler.SubmitTask(task)

	retrieved, err := scheduler.GetTask("task-2")
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if retrieved.ID != "task-2" {
		t.Errorf("Expected task ID task-2, got %s", retrieved.ID)
	}

	_, err = scheduler.GetTask("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func TestSchedule(t *testing.T) {
	nodeManager := edgenode.NewManager()
	scheduler := NewScheduler(nodeManager)

	// Register a node
	node := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Test Node",
		TotalCPU:        8.0,
		TotalMemory:     16384.0,
		AvailableCPU:    8.0,
		AvailableMemory: 16384.0,
	}
	nodeManager.RegisterNode(node)

	// Submit a task
	task := &models.Task{
		ID:          "task-3",
		Name:        "Schedulable Task",
		Priority:    3,
		CPURequired: 2.0,
		MemRequired: 1024.0,
		Duration:    10 * time.Second,
		Deadline:    time.Now().Add(1 * time.Hour),
	}
	scheduler.SubmitTask(task)

	// Schedule tasks
	result, err := scheduler.Schedule()
	if err != nil {
		t.Fatalf("Scheduling failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected successful scheduling")
	}

	if len(result.Schedules) == 0 {
		t.Error("Expected at least one schedule")
	}

	// Verify task status changed
	retrieved, _ := scheduler.GetTask("task-3")
	if retrieved.Status != models.TaskStatusScheduled {
		t.Errorf("Expected task status to be scheduled, got %s", retrieved.Status)
	}
}

func TestScheduleNoNodes(t *testing.T) {
	nodeManager := edgenode.NewManager()
	scheduler := NewScheduler(nodeManager)

	// Submit a task but no nodes
	task := &models.Task{
		ID:          "task-4",
		Name:        "Unschedulable Task",
		Priority:    1,
		CPURequired: 2.0,
		MemRequired: 1024.0,
		Duration:    10 * time.Second,
		Deadline:    time.Now().Add(1 * time.Hour),
	}
	scheduler.SubmitTask(task)

	// Try to schedule
	result, err := scheduler.Schedule()
	if err == nil {
		t.Error("Expected error when scheduling with no nodes")
	}

	if result == nil || result.Success {
		t.Error("Expected unsuccessful scheduling result")
	}
}
