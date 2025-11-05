package edgenode

import (
	"testing"
	"time"

	"github.com/Uttam-Mahata/dts/pkg/models"
)

func TestRegisterNode(t *testing.T) {
	manager := NewManager()

	node := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Test Node",
		TotalCPU:        8.0,
		TotalMemory:     16384.0,
		AvailableCPU:    8.0,
		AvailableMemory: 16384.0,
	}

	err := manager.RegisterNode(node)
	if err != nil {
		t.Fatalf("Failed to register node: %v", err)
	}

	if manager.GetNodeCount() != 1 {
		t.Errorf("Expected 1 node, got %d", manager.GetNodeCount())
	}

	// Try to register the same node again
	err = manager.RegisterNode(node)
	if err == nil {
		t.Error("Expected error when registering duplicate node")
	}
}

func TestGetNode(t *testing.T) {
	manager := NewManager()

	node := &models.EdgeNode{
		ID:              "node-2",
		Name:            "Test Node 2",
		TotalCPU:        4.0,
		TotalMemory:     8192.0,
		AvailableCPU:    4.0,
		AvailableMemory: 8192.0,
	}

	manager.RegisterNode(node)

	retrieved, err := manager.GetNode("node-2")
	if err != nil {
		t.Fatalf("Failed to get node: %v", err)
	}

	if retrieved.ID != "node-2" {
		t.Errorf("Expected node ID node-2, got %s", retrieved.ID)
	}

	_, err = manager.GetNode("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent node")
	}
}

func TestGetActiveNodes(t *testing.T) {
	manager := NewManager()

	node1 := &models.EdgeNode{
		ID:              "node-3",
		Name:            "Active Node",
		TotalCPU:        8.0,
		TotalMemory:     16384.0,
		AvailableCPU:    8.0,
		AvailableMemory: 16384.0,
	}

	manager.RegisterNode(node1)

	activeNodes := manager.GetActiveNodes()
	if len(activeNodes) != 1 {
		t.Errorf("Expected 1 active node, got %d", len(activeNodes))
	}
}

func TestCheckNodeHealth(t *testing.T) {
	manager := NewManager()

	node := &models.EdgeNode{
		ID:              "node-4",
		Name:            "Health Test Node",
		TotalCPU:        4.0,
		TotalMemory:     8192.0,
		AvailableCPU:    4.0,
		AvailableMemory: 8192.0,
	}

	manager.RegisterNode(node)

	// Set last heartbeat to past
	n, _ := manager.GetNode("node-4")
	n.LastHeartbeat = time.Now().Add(-2 * time.Minute)

	// Check health with 1 minute timeout
	manager.CheckNodeHealth(1 * time.Minute)

	n, _ = manager.GetNode("node-4")
	if n.Status != models.EdgeNodeStatusInactive {
		t.Errorf("Expected node to be inactive, got status %s", n.Status)
	}
}
