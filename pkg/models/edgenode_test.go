package models

import (
	"testing"
	"time"
)

func TestEdgeNodeAllocateResources(t *testing.T) {
	node := &EdgeNode{
		ID:              "node-1",
		Name:            "test-node",
		TotalCPU:        8.0,
		TotalMemory:     16384.0,
		AvailableCPU:    8.0,
		AvailableMemory: 16384.0,
		Status:          EdgeNodeStatusActive,
	}

	// Test successful allocation
	success := node.AllocateResources(2.0, 4096.0)
	if !success {
		t.Error("Expected successful resource allocation")
	}

	if node.AvailableCPU != 6.0 {
		t.Errorf("Expected available CPU to be 6.0, got %f", node.AvailableCPU)
	}

	if node.AvailableMemory != 12288.0 {
		t.Errorf("Expected available memory to be 12288.0, got %f", node.AvailableMemory)
	}

	// Test failed allocation (insufficient resources)
	success = node.AllocateResources(10.0, 1024.0)
	if success {
		t.Error("Expected allocation to fail due to insufficient CPU")
	}

	// Resources should remain unchanged
	if node.AvailableCPU != 6.0 {
		t.Errorf("Available CPU should not change on failed allocation, got %f", node.AvailableCPU)
	}
}

func TestEdgeNodeReleaseResources(t *testing.T) {
	node := &EdgeNode{
		ID:              "node-1",
		Name:            "test-node",
		TotalCPU:        8.0,
		TotalMemory:     16384.0,
		AvailableCPU:    4.0,
		AvailableMemory: 8192.0,
		Status:          EdgeNodeStatusActive,
	}

	// Release resources
	node.ReleaseResources(2.0, 4096.0)

	if node.AvailableCPU != 6.0 {
		t.Errorf("Expected available CPU to be 6.0, got %f", node.AvailableCPU)
	}

	if node.AvailableMemory != 12288.0 {
		t.Errorf("Expected available memory to be 12288.0, got %f", node.AvailableMemory)
	}

	// Test resource cap at total capacity
	node.ReleaseResources(10.0, 10000.0)

	if node.AvailableCPU != 8.0 {
		t.Errorf("Available CPU should be capped at total CPU (8.0), got %f", node.AvailableCPU)
	}

	if node.AvailableMemory != 16384.0 {
		t.Errorf("Available memory should be capped at total memory (16384.0), got %f", node.AvailableMemory)
	}
}

func TestEdgeNodeGetLoad(t *testing.T) {
	node := &EdgeNode{
		ID:              "node-1",
		Name:            "test-node",
		TotalCPU:        8.0,
		TotalMemory:     16384.0,
		AvailableCPU:    4.0,
		AvailableMemory: 8192.0,
		Status:          EdgeNodeStatusActive,
	}

	load := node.GetLoad()
	expectedLoad := 0.5 // 50% average load

	if load != expectedLoad {
		t.Errorf("Expected load to be %f, got %f", expectedLoad, load)
	}

	// Test with no load
	node.AvailableCPU = 8.0
	node.AvailableMemory = 16384.0
	load = node.GetLoad()

	if load != 0.0 {
		t.Errorf("Expected load to be 0.0, got %f", load)
	}

	// Test with full load
	node.AvailableCPU = 0.0
	node.AvailableMemory = 0.0
	load = node.GetLoad()

	if load != 1.0 {
		t.Errorf("Expected load to be 1.0, got %f", load)
	}
}

func TestEdgeNodeAddTask(t *testing.T) {
	node := &EdgeNode{
		ID:           "node-1",
		Name:         "test-node",
		RunningTasks: []string{},
	}

	node.AddTask("task-1")
	if len(node.RunningTasks) != 1 {
		t.Errorf("Expected 1 running task, got %d", len(node.RunningTasks))
	}

	if node.RunningTasks[0] != "task-1" {
		t.Errorf("Expected task-1, got %s", node.RunningTasks[0])
	}

	node.AddTask("task-2")
	if len(node.RunningTasks) != 2 {
		t.Errorf("Expected 2 running tasks, got %d", len(node.RunningTasks))
	}
}

func TestEdgeNodeRemoveTask(t *testing.T) {
	node := &EdgeNode{
		ID:           "node-1",
		Name:         "test-node",
		RunningTasks: []string{"task-1", "task-2", "task-3"},
	}

	node.RemoveTask("task-2")
	if len(node.RunningTasks) != 2 {
		t.Errorf("Expected 2 running tasks, got %d", len(node.RunningTasks))
	}

	// Verify task-2 is removed
	for _, task := range node.RunningTasks {
		if task == "task-2" {
			t.Error("task-2 should have been removed")
		}
	}

	// Verify remaining tasks
	if node.RunningTasks[0] != "task-1" || node.RunningTasks[1] != "task-3" {
		t.Errorf("Expected remaining tasks to be task-1 and task-3, got %v", node.RunningTasks)
	}
}

func TestEdgeNodeStatus(t *testing.T) {
	statuses := []EdgeNodeStatus{
		EdgeNodeStatusActive,
		EdgeNodeStatusInactive,
		EdgeNodeStatusBusy,
	}

	expectedValues := []string{"active", "inactive", "busy"}

	for i, status := range statuses {
		if string(status) != expectedValues[i] {
			t.Errorf("Expected status %s, got %s", expectedValues[i], string(status))
		}
	}
}

func TestEdgeNodeConcurrentOperations(t *testing.T) {
	node := &EdgeNode{
		ID:              "node-1",
		Name:            "test-node",
		TotalCPU:        8.0,
		TotalMemory:     16384.0,
		AvailableCPU:    8.0,
		AvailableMemory: 16384.0,
		Status:          EdgeNodeStatusActive,
		RunningTasks:    []string{},
		LastHeartbeat:   time.Now(),
	}

	// Test concurrent allocations and releases
	done := make(chan bool)

	// Concurrent allocations
	for i := 0; i < 10; i++ {
		go func() {
			node.AllocateResources(0.1, 100.0)
			done <- true
		}()
	}

	// Concurrent releases
	for i := 0; i < 10; i++ {
		go func() {
			node.ReleaseResources(0.1, 100.0)
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}

	// Verify resources are still within valid range
	if node.AvailableCPU < 0 || node.AvailableCPU > node.TotalCPU {
		t.Errorf("Available CPU out of valid range: %f", node.AvailableCPU)
	}

	if node.AvailableMemory < 0 || node.AvailableMemory > node.TotalMemory {
		t.Errorf("Available memory out of valid range: %f", node.AvailableMemory)
	}
}
