package edgenode

import (
	"fmt"
	"sync"
	"time"

	"github.com/Uttam-Mahata/dts/pkg/models"
)

// Manager manages edge nodes in the system
type Manager struct {
	nodes map[string]*models.EdgeNode
	mu    sync.RWMutex
}

// NewManager creates a new edge node manager
func NewManager() *Manager {
	return &Manager{
		nodes: make(map[string]*models.EdgeNode),
	}
}

// RegisterNode registers a new edge node
func (m *Manager) RegisterNode(node *models.EdgeNode) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.nodes[node.ID]; exists {
		return fmt.Errorf("node %s already registered", node.ID)
	}

	node.Status = models.EdgeNodeStatusActive
	node.LastHeartbeat = time.Now()
	node.RunningTasks = make([]string, 0)

	m.nodes[node.ID] = node
	return nil
}

// UnregisterNode removes a node from the system
func (m *Manager) UnregisterNode(nodeID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.nodes[nodeID]; !exists {
		return fmt.Errorf("node %s not found", nodeID)
	}

	delete(m.nodes, nodeID)
	return nil
}

// GetNode retrieves a node by ID
func (m *Manager) GetNode(nodeID string) (*models.EdgeNode, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	node, exists := m.nodes[nodeID]
	if !exists {
		return nil, fmt.Errorf("node %s not found", nodeID)
	}

	return node, nil
}

// GetAllNodes returns all registered nodes
func (m *Manager) GetAllNodes() []*models.EdgeNode {
	m.mu.RLock()
	defer m.mu.RUnlock()

	nodes := make([]*models.EdgeNode, 0, len(m.nodes))
	for _, node := range m.nodes {
		nodes = append(nodes, node)
	}

	return nodes
}

// GetActiveNodes returns all active nodes
func (m *Manager) GetActiveNodes() []*models.EdgeNode {
	m.mu.RLock()
	defer m.mu.RUnlock()

	nodes := make([]*models.EdgeNode, 0)
	for _, node := range m.nodes {
		if node.Status == models.EdgeNodeStatusActive {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

// UpdateHeartbeat updates the last heartbeat time for a node
func (m *Manager) UpdateHeartbeat(nodeID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	node, exists := m.nodes[nodeID]
	if !exists {
		return fmt.Errorf("node %s not found", nodeID)
	}

	node.LastHeartbeat = time.Now()
	return nil
}

// UpdateNodeResources updates the available resources for a node
func (m *Manager) UpdateNodeResources(nodeID string, availableCPU, availableMemory float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	node, exists := m.nodes[nodeID]
	if !exists {
		return fmt.Errorf("node %s not found", nodeID)
	}

	node.AvailableCPU = availableCPU
	node.AvailableMemory = availableMemory
	return nil
}

// CheckNodeHealth checks the health of all nodes based on heartbeat
func (m *Manager) CheckNodeHealth(timeout time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for _, node := range m.nodes {
		if now.Sub(node.LastHeartbeat) > timeout {
			node.Status = models.EdgeNodeStatusInactive
		}
	}
}

// GetNodeCount returns the total number of nodes
func (m *Manager) GetNodeCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.nodes)
}

// GetActiveNodeCount returns the number of active nodes
func (m *Manager) GetActiveNodeCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, node := range m.nodes {
		if node.Status == models.EdgeNodeStatusActive {
			count++
		}
	}

	return count
}
