package loadbalancer

import (
	"sort"
	"sync"

	"github.com/Uttam-Mahata/dts/pkg/edgenode"
	"github.com/Uttam-Mahata/dts/pkg/models"
)

// LoadBalancer manages load balancing across edge nodes
type LoadBalancer struct {
	nodeManager *edgenode.Manager
	mu          sync.RWMutex
}

// NewLoadBalancer creates a new load balancer instance
func NewLoadBalancer(nodeManager *edgenode.Manager) *LoadBalancer {
	return &LoadBalancer{
		nodeManager: nodeManager,
	}
}

// SelectNode selects the best node for a task based on current load
func (lb *LoadBalancer) SelectNode(task *models.Task) (*models.EdgeNode, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	nodes := lb.nodeManager.GetActiveNodes()
	if len(nodes) == 0 {
		return nil, nil
	}

	// Filter nodes that can handle the task
	candidateNodes := make([]*models.EdgeNode, 0)
	for _, node := range nodes {
		if node.AvailableCPU >= task.CPURequired && node.AvailableMemory >= task.MemRequired {
			candidateNodes = append(candidateNodes, node)
		}
	}

	if len(candidateNodes) == 0 {
		return nil, nil
	}

	// Sort by load (ascending)
	sort.Slice(candidateNodes, func(i, j int) bool {
		return candidateNodes[i].GetLoad() < candidateNodes[j].GetLoad()
	})

	// Return the least loaded node
	return candidateNodes[0], nil
}

// BalanceLoad performs load balancing by redistributing tasks
func (lb *LoadBalancer) BalanceLoad() error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	nodes := lb.nodeManager.GetActiveNodes()
	if len(nodes) <= 1 {
		return nil // No balancing needed
	}

	// Calculate average load
	totalLoad := 0.0
	for _, node := range nodes {
		totalLoad += node.GetLoad()
	}
	avgLoad := totalLoad / float64(len(nodes))

	// Identify overloaded and underloaded nodes
	var overloaded, underloaded []*models.EdgeNode
	for _, node := range nodes {
		load := node.GetLoad()
		if load > avgLoad*1.2 { // 20% threshold
			overloaded = append(overloaded, node)
		} else if load < avgLoad*0.8 {
			underloaded = append(underloaded, node)
		}
	}

	// In a real implementation, we would migrate tasks from overloaded to underloaded nodes
	// For now, we just identify the imbalance
	_ = overloaded
	_ = underloaded

	return nil
}

// GetSystemLoad returns the overall system load
func (lb *LoadBalancer) GetSystemLoad() float64 {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	nodes := lb.nodeManager.GetActiveNodes()
	if len(nodes) == 0 {
		return 0.0
	}

	totalLoad := 0.0
	for _, node := range nodes {
		totalLoad += node.GetLoad()
	}

	return totalLoad / float64(len(nodes))
}

// GetNodeLoads returns load information for all nodes
func (lb *LoadBalancer) GetNodeLoads() map[string]float64 {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	loads := make(map[string]float64)
	nodes := lb.nodeManager.GetAllNodes()

	for _, node := range nodes {
		loads[node.ID] = node.GetLoad()
	}

	return loads
}

// GetLeastLoadedNode returns the node with the lowest load
func (lb *LoadBalancer) GetLeastLoadedNode() *models.EdgeNode {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	nodes := lb.nodeManager.GetActiveNodes()
	if len(nodes) == 0 {
		return nil
	}

	minLoad := nodes[0].GetLoad()
	minNode := nodes[0]

	for _, node := range nodes[1:] {
		load := node.GetLoad()
		if load < minLoad {
			minLoad = load
			minNode = node
		}
	}

	return minNode
}

// GetMostLoadedNode returns the node with the highest load
func (lb *LoadBalancer) GetMostLoadedNode() *models.EdgeNode {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	nodes := lb.nodeManager.GetActiveNodes()
	if len(nodes) == 0 {
		return nil
	}

	maxLoad := nodes[0].GetLoad()
	maxNode := nodes[0]

	for _, node := range nodes[1:] {
		load := node.GetLoad()
		if load > maxLoad {
			maxLoad = load
			maxNode = node
		}
	}

	return maxNode
}

// IsBalanced checks if the system is balanced
func (lb *LoadBalancer) IsBalanced(threshold float64) bool {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	nodes := lb.nodeManager.GetActiveNodes()
	if len(nodes) <= 1 {
		return true
	}

	loads := make([]float64, len(nodes))
	avgLoad := 0.0

	for i, node := range nodes {
		loads[i] = node.GetLoad()
		avgLoad += loads[i]
	}
	avgLoad /= float64(len(nodes))

	// Check if all nodes are within threshold of average
	for _, load := range loads {
		deviation := load - avgLoad
		if deviation < 0 {
			deviation = -deviation
		}
		if deviation > threshold*avgLoad {
			return false
		}
	}

	return true
}
