package matroid

import (
	"testing"
	"time"

	"github.com/Uttam-Mahata/dts/pkg/models"
)

func TestNewMatroid(t *testing.T) {
	m := NewMatroid()

	if m == nil {
		t.Fatal("Expected non-nil matroid")
	}

	if m.elements == nil {
		t.Error("Expected elements slice to be initialized")
	}
}

func TestAddElement(t *testing.T) {
	m := NewMatroid()

	m.AddElement("task-1", "node-1", 0.75, 5)

	if len(m.elements) != 1 {
		t.Errorf("Expected 1 element, got %d", len(m.elements))
	}

	elem := m.elements[0]
	if elem.TaskID != "task-1" {
		t.Errorf("Expected task ID 'task-1', got '%s'", elem.TaskID)
	}

	if elem.NodeID != "node-1" {
		t.Errorf("Expected node ID 'node-1', got '%s'", elem.NodeID)
	}

	if elem.Weight != 0.75 {
		t.Errorf("Expected weight 0.75, got %f", elem.Weight)
	}

	if elem.Priority != 5 {
		t.Errorf("Expected priority 5, got %d", elem.Priority)
	}
}

func TestGreedyOptimize(t *testing.T) {
	m := NewMatroid()

	// Create test tasks
	tasks := []*models.Task{
		{
			ID:          "task-1",
			Name:        "Task 1",
			Priority:    5,
			CPURequired: 2.0,
			MemRequired: 1024.0,
			Deadline:    time.Now().Add(1 * time.Hour),
		},
		{
			ID:          "task-2",
			Name:        "Task 2",
			Priority:    3,
			CPURequired: 1.0,
			MemRequired: 512.0,
			Deadline:    time.Now().Add(2 * time.Hour),
		},
	}

	// Create test nodes
	nodes := []*models.EdgeNode{
		{
			ID:              "node-1",
			Name:            "Node 1",
			TotalCPU:        8.0,
			TotalMemory:     4096.0,
			AvailableCPU:    8.0,
			AvailableMemory: 4096.0,
			Status:          models.EdgeNodeStatusActive,
		},
		{
			ID:              "node-2",
			Name:            "Node 2",
			TotalCPU:        4.0,
			TotalMemory:     2048.0,
			AvailableCPU:    4.0,
			AvailableMemory: 2048.0,
			Status:          models.EdgeNodeStatusActive,
		},
	}

	schedules := m.GreedyOptimize(tasks, nodes)

	if len(schedules) != 2 {
		t.Errorf("Expected 2 schedules, got %d", len(schedules))
	}

	// Verify all tasks are scheduled
	scheduledTasks := make(map[string]bool)
	for _, schedule := range schedules {
		scheduledTasks[schedule.TaskID] = true
	}

	if !scheduledTasks["task-1"] || !scheduledTasks["task-2"] {
		t.Error("Not all tasks were scheduled")
	}
}

func TestGreedyOptimizeInsufficientResources(t *testing.T) {
	m := NewMatroid()

	// Create tasks requiring more resources than available
	tasks := []*models.Task{
		{
			ID:          "task-1",
			Name:        "Large Task",
			Priority:    5,
			CPURequired: 10.0,
			MemRequired: 8192.0,
			Deadline:    time.Now().Add(1 * time.Hour),
		},
	}

	// Create node with insufficient resources
	nodes := []*models.EdgeNode{
		{
			ID:              "node-1",
			Name:            "Small Node",
			TotalCPU:        4.0,
			TotalMemory:     2048.0,
			AvailableCPU:    4.0,
			AvailableMemory: 2048.0,
			Status:          models.EdgeNodeStatusActive,
		},
	}

	schedules := m.GreedyOptimize(tasks, nodes)

	if len(schedules) != 0 {
		t.Errorf("Expected 0 schedules due to insufficient resources, got %d", len(schedules))
	}
}

func TestGreedyOptimizeInactiveNodes(t *testing.T) {
	m := NewMatroid()

	tasks := []*models.Task{
		{
			ID:          "task-1",
			Name:        "Task 1",
			Priority:    5,
			CPURequired: 2.0,
			MemRequired: 1024.0,
			Deadline:    time.Now().Add(1 * time.Hour),
		},
	}

	// All nodes are inactive
	nodes := []*models.EdgeNode{
		{
			ID:              "node-1",
			Name:            "Inactive Node",
			TotalCPU:        8.0,
			TotalMemory:     4096.0,
			AvailableCPU:    8.0,
			AvailableMemory: 4096.0,
			Status:          models.EdgeNodeStatusInactive,
		},
	}

	schedules := m.GreedyOptimize(tasks, nodes)

	if len(schedules) != 0 {
		t.Errorf("Expected 0 schedules due to inactive nodes, got %d", len(schedules))
	}
}

func TestGreedyOptimizePriorityOrdering(t *testing.T) {
	m := NewMatroid()

	// Create tasks with different priorities
	tasks := []*models.Task{
		{
			ID:          "task-low",
			Name:        "Low Priority Task",
			Priority:    1,
			CPURequired: 2.0,
			MemRequired: 1024.0,
			Deadline:    time.Now().Add(1 * time.Hour),
		},
		{
			ID:          "task-high",
			Name:        "High Priority Task",
			Priority:    10,
			CPURequired: 2.0,
			MemRequired: 1024.0,
			Deadline:    time.Now().Add(1 * time.Hour),
		},
	}

	// Create node with limited resources (can only handle one task)
	nodes := []*models.EdgeNode{
		{
			ID:              "node-1",
			Name:            "Limited Node",
			TotalCPU:        2.0,
			TotalMemory:     1024.0,
			AvailableCPU:    2.0,
			AvailableMemory: 1024.0,
			Status:          models.EdgeNodeStatusActive,
		},
	}

	schedules := m.GreedyOptimize(tasks, nodes)

	if len(schedules) != 1 {
		t.Errorf("Expected 1 schedule, got %d", len(schedules))
	}

	// High priority task should be scheduled
	if schedules[0].TaskID != "task-high" {
		t.Errorf("Expected high priority task to be scheduled, got %s", schedules[0].TaskID)
	}
}

func TestCalculateWeight(t *testing.T) {
	m := NewMatroid()

	task := &models.Task{
		ID:          "task-1",
		Name:        "Test Task",
		Priority:    5,
		CPURequired: 2.0,
		MemRequired: 1024.0,
	}

	node := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Test Node",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    8.0,
		AvailableMemory: 4096.0,
		Status:          models.EdgeNodeStatusActive,
	}

	weight := m.calculateWeight(task, node)

	// Weight should be positive
	if weight <= 0 {
		t.Errorf("Expected positive weight, got %f", weight)
	}

	// Weight should be reasonable (between 0 and some upper bound)
	if weight > 10.0 {
		t.Errorf("Weight seems unreasonably high: %f", weight)
	}
}

func TestIsIndependent(t *testing.T) {
	m := NewMatroid()

	tasks := []*models.Task{
		{
			ID:          "task-1",
			Name:        "Task 1",
			CPURequired: 2.0,
			MemRequired: 1024.0,
		},
		{
			ID:          "task-2",
			Name:        "Task 2",
			CPURequired: 1.0,
			MemRequired: 512.0,
		},
	}

	nodes := []*models.EdgeNode{
		{
			ID:              "node-1",
			Name:            "Node 1",
			TotalCPU:        8.0,
			TotalMemory:     4096.0,
			AvailableCPU:    8.0,
			AvailableMemory: 4096.0,
			Status:          models.EdgeNodeStatusActive,
		},
	}

	// Test independent set (both tasks can fit on the node)
	elements := []Element{
		{TaskID: "task-1", NodeID: "node-1"},
		{TaskID: "task-2", NodeID: "node-1"},
	}

	if !m.IsIndependent(elements, tasks, nodes) {
		t.Error("Expected elements to be independent")
	}

	// Test dependent set (duplicate task assignment)
	duplicateElements := []Element{
		{TaskID: "task-1", NodeID: "node-1"},
		{TaskID: "task-1", NodeID: "node-1"},
	}

	if m.IsIndependent(duplicateElements, tasks, nodes) {
		t.Error("Expected duplicate task assignment to be dependent")
	}
}

func TestGreedyOptimizeEmptyInputs(t *testing.T) {
	m := NewMatroid()

	// Test with no tasks
	schedules := m.GreedyOptimize([]*models.Task{}, []*models.EdgeNode{
		{
			ID:              "node-1",
			TotalCPU:        8.0,
			TotalMemory:     4096.0,
			AvailableCPU:    8.0,
			AvailableMemory: 4096.0,
			Status:          models.EdgeNodeStatusActive,
		},
	})

	if len(schedules) != 0 {
		t.Errorf("Expected 0 schedules with no tasks, got %d", len(schedules))
	}

	// Test with no nodes
	schedules = m.GreedyOptimize([]*models.Task{
		{
			ID:          "task-1",
			CPURequired: 2.0,
			MemRequired: 1024.0,
		},
	}, []*models.EdgeNode{})

	if len(schedules) != 0 {
		t.Errorf("Expected 0 schedules with no nodes, got %d", len(schedules))
	}
}

func TestGreedyOptimizeResourceExhaustion(t *testing.T) {
	m := NewMatroid()

	// Create multiple tasks
	tasks := []*models.Task{
		{
			ID:          "task-1",
			Priority:    5,
			CPURequired: 2.0,
			MemRequired: 1024.0,
		},
		{
			ID:          "task-2",
			Priority:    5,
			CPURequired: 2.0,
			MemRequired: 1024.0,
		},
		{
			ID:          "task-3",
			Priority:    5,
			CPURequired: 2.0,
			MemRequired: 1024.0,
		},
	}

	// Create node that can only handle 2 tasks
	nodes := []*models.EdgeNode{
		{
			ID:              "node-1",
			TotalCPU:        4.0,
			TotalMemory:     2048.0,
			AvailableCPU:    4.0,
			AvailableMemory: 2048.0,
			Status:          models.EdgeNodeStatusActive,
		},
	}

	schedules := m.GreedyOptimize(tasks, nodes)

	if len(schedules) != 2 {
		t.Errorf("Expected 2 schedules (node can only handle 2 tasks), got %d", len(schedules))
	}
}
