package loadbalancer

import (
	"testing"
	"time"

	"github.com/Uttam-Mahata/dts/pkg/edgenode"
	"github.com/Uttam-Mahata/dts/pkg/models"
)

func TestNewLoadBalancer(t *testing.T) {
	manager := edgenode.NewManager()
	lb := NewLoadBalancer(manager)

	if lb == nil {
		t.Fatal("Expected non-nil load balancer")
	}

	if lb.nodeManager == nil {
		t.Error("Expected node manager to be set")
	}
}

func TestSelectNode(t *testing.T) {
	manager := edgenode.NewManager()
	lb := NewLoadBalancer(manager)

	// Register nodes
	node1 := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Node 1",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    8.0,
		AvailableMemory: 4096.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	node2 := &models.EdgeNode{
		ID:              "node-2",
		Name:            "Node 2",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    4.0, // More loaded
		AvailableMemory: 2048.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	manager.RegisterNode(node1)
	manager.RegisterNode(node2)

	task := &models.Task{
		ID:          "task-1",
		CPURequired: 1.0,
		MemRequired: 512.0,
	}

	// Should select least loaded node (node-1)
	selectedNode, err := lb.SelectNode(task)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if selectedNode == nil {
		t.Fatal("Expected a node to be selected")
	}

	if selectedNode.ID != "node-1" {
		t.Errorf("Expected node-1 to be selected (least loaded), got %s", selectedNode.ID)
	}
}

func TestSelectNodeInsufficientResources(t *testing.T) {
	manager := edgenode.NewManager()
	lb := NewLoadBalancer(manager)

	// Register node with limited resources
	node := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Node 1",
		TotalCPU:        2.0,
		TotalMemory:     1024.0,
		AvailableCPU:    2.0,
		AvailableMemory: 1024.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	manager.RegisterNode(node)

	// Task requires more resources than available
	task := &models.Task{
		ID:          "task-1",
		CPURequired: 4.0,
		MemRequired: 2048.0,
	}

	selectedNode, err := lb.SelectNode(task)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if selectedNode != nil {
		t.Error("Expected no node to be selected due to insufficient resources")
	}
}

func TestSelectNodeNoNodes(t *testing.T) {
	manager := edgenode.NewManager()
	lb := NewLoadBalancer(manager)

	task := &models.Task{
		ID:          "task-1",
		CPURequired: 1.0,
		MemRequired: 512.0,
	}

	selectedNode, err := lb.SelectNode(task)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if selectedNode != nil {
		t.Error("Expected no node to be selected when no nodes are available")
	}
}

func TestGetSystemLoad(t *testing.T) {
	manager := edgenode.NewManager()
	lb := NewLoadBalancer(manager)

	// Register nodes with different loads
	node1 := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Node 1",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    4.0, // 50% load
		AvailableMemory: 2048.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	node2 := &models.EdgeNode{
		ID:              "node-2",
		Name:            "Node 2",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    6.0, // 25% load
		AvailableMemory: 3072.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	manager.RegisterNode(node1)
	manager.RegisterNode(node2)

	systemLoad := lb.GetSystemLoad()

	// Average load should be (0.5 + 0.25) / 2 = 0.375
	expectedLoad := 0.375

	if systemLoad != expectedLoad {
		t.Errorf("Expected system load %f, got %f", expectedLoad, systemLoad)
	}
}

func TestGetSystemLoadNoNodes(t *testing.T) {
	manager := edgenode.NewManager()
	lb := NewLoadBalancer(manager)

	systemLoad := lb.GetSystemLoad()

	if systemLoad != 0.0 {
		t.Errorf("Expected system load 0.0 with no nodes, got %f", systemLoad)
	}
}

func TestGetNodeLoads(t *testing.T) {
	manager := edgenode.NewManager()
	lb := NewLoadBalancer(manager)

	// Register nodes
	node1 := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Node 1",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    4.0,
		AvailableMemory: 2048.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	node2 := &models.EdgeNode{
		ID:              "node-2",
		Name:            "Node 2",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    8.0,
		AvailableMemory: 4096.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	manager.RegisterNode(node1)
	manager.RegisterNode(node2)

	loads := lb.GetNodeLoads()

	if len(loads) != 2 {
		t.Errorf("Expected 2 node loads, got %d", len(loads))
	}

	if _, exists := loads["node-1"]; !exists {
		t.Error("Expected load for node-1")
	}

	if _, exists := loads["node-2"]; !exists {
		t.Error("Expected load for node-2")
	}

	// Verify loads are correct
	if loads["node-1"] != 0.5 {
		t.Errorf("Expected load 0.5 for node-1, got %f", loads["node-1"])
	}

	if loads["node-2"] != 0.0 {
		t.Errorf("Expected load 0.0 for node-2, got %f", loads["node-2"])
	}
}

func TestGetLeastLoadedNode(t *testing.T) {
	manager := edgenode.NewManager()
	lb := NewLoadBalancer(manager)

	// Register nodes with different loads
	node1 := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Node 1",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    2.0, // 75% load
		AvailableMemory: 1024.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	node2 := &models.EdgeNode{
		ID:              "node-2",
		Name:            "Node 2",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    6.0, // 25% load
		AvailableMemory: 3072.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	manager.RegisterNode(node1)
	manager.RegisterNode(node2)

	leastLoaded := lb.GetLeastLoadedNode()

	if leastLoaded == nil {
		t.Fatal("Expected a node to be returned")
	}

	if leastLoaded.ID != "node-2" {
		t.Errorf("Expected node-2 to be least loaded, got %s", leastLoaded.ID)
	}
}

func TestGetMostLoadedNode(t *testing.T) {
	manager := edgenode.NewManager()
	lb := NewLoadBalancer(manager)

	// Register nodes with different loads
	node1 := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Node 1",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    2.0, // 75% load
		AvailableMemory: 1024.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	node2 := &models.EdgeNode{
		ID:              "node-2",
		Name:            "Node 2",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    6.0, // 25% load
		AvailableMemory: 3072.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	manager.RegisterNode(node1)
	manager.RegisterNode(node2)

	mostLoaded := lb.GetMostLoadedNode()

	if mostLoaded == nil {
		t.Fatal("Expected a node to be returned")
	}

	if mostLoaded.ID != "node-1" {
		t.Errorf("Expected node-1 to be most loaded, got %s", mostLoaded.ID)
	}
}

func TestIsBalanced(t *testing.T) {
	manager := edgenode.NewManager()
	lb := NewLoadBalancer(manager)

	// Register balanced nodes
	node1 := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Node 1",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    4.0, // 50% load
		AvailableMemory: 2048.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	node2 := &models.EdgeNode{
		ID:              "node-2",
		Name:            "Node 2",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    4.0, // 50% load
		AvailableMemory: 2048.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	manager.RegisterNode(node1)
	manager.RegisterNode(node2)

	isBalanced := lb.IsBalanced(0.2)

	if !isBalanced {
		t.Error("Expected system to be balanced")
	}
}

func TestIsBalancedUnbalanced(t *testing.T) {
	manager := edgenode.NewManager()
	lb := NewLoadBalancer(manager)

	// Register unbalanced nodes
	node1 := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Node 1",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    1.0, // 87.5% load
		AvailableMemory: 512.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	node2 := &models.EdgeNode{
		ID:              "node-2",
		Name:            "Node 2",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    8.0, // 0% load
		AvailableMemory: 4096.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	manager.RegisterNode(node1)
	manager.RegisterNode(node2)

	isBalanced := lb.IsBalanced(0.2)

	if isBalanced {
		t.Error("Expected system to be unbalanced")
	}
}

func TestIsBalancedSingleNode(t *testing.T) {
	manager := edgenode.NewManager()
	lb := NewLoadBalancer(manager)

	node := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Node 1",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    4.0,
		AvailableMemory: 2048.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	manager.RegisterNode(node)

	isBalanced := lb.IsBalanced(0.2)

	if !isBalanced {
		t.Error("Expected single node system to be balanced")
	}
}

func TestBalanceLoad(t *testing.T) {
	manager := edgenode.NewManager()
	lb := NewLoadBalancer(manager)

	// Register nodes
	node1 := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Node 1",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    2.0,
		AvailableMemory: 1024.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	node2 := &models.EdgeNode{
		ID:              "node-2",
		Name:            "Node 2",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    8.0,
		AvailableMemory: 4096.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}

	manager.RegisterNode(node1)
	manager.RegisterNode(node2)

	// BalanceLoad currently doesn't implement task migration
	// Just verify it doesn't error
	err := lb.BalanceLoad()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
