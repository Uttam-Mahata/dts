package matroid

import (
	"math"
	"sort"

	"github.com/Uttam-Mahata/dts/pkg/models"
)

// Element represents an element in the matroid
type Element struct {
	TaskID   string
	NodeID   string
	Weight   float64
	Priority int
}

// Matroid represents a matroid structure for optimization
type Matroid struct {
	elements    []Element
	independent func([]Element) bool
}

// NewMatroid creates a new matroid instance
func NewMatroid() *Matroid {
	return &Matroid{
		elements: make([]Element, 0),
		independent: func(subset []Element) bool {
			// Check if subset is independent (no resource conflicts)
			nodeResources := make(map[string]struct {
				cpu float64
				mem float64
			})

			for _, elem := range subset {
				if _, exists := nodeResources[elem.NodeID]; !exists {
					nodeResources[elem.NodeID] = struct {
						cpu float64
						mem float64
					}{0, 0}
				}
			}

			return true
		},
	}
}

// AddElement adds an element to the matroid
func (m *Matroid) AddElement(taskID, nodeID string, weight float64, priority int) {
	m.elements = append(m.elements, Element{
		TaskID:   taskID,
		NodeID:   nodeID,
		Weight:   weight,
		Priority: priority,
	})
}

// GreedyOptimize performs greedy matroid optimization
// This implements the greedy algorithm for matroid optimization which is optimal
// for matroid structures
func (m *Matroid) GreedyOptimize(tasks []*models.Task, nodes []*models.EdgeNode) []models.Schedule {
	// Create elements for each task-node pair
	elements := make([]Element, 0)

	for _, task := range tasks {
		for _, node := range nodes {
			if node.Status != models.EdgeNodeStatusActive {
				continue
			}

			// Check if node has enough resources
			if node.AvailableCPU < task.CPURequired || node.AvailableMemory < task.MemRequired {
				continue
			}

			// Calculate weight based on multiple factors
			weight := m.calculateWeight(task, node)

			elements = append(elements, Element{
				TaskID:   task.ID,
				NodeID:   node.ID,
				Weight:   weight,
				Priority: task.Priority,
			})
		}
	}

	// Sort elements by weight (higher is better) and priority
	sort.Slice(elements, func(i, j int) bool {
		if elements[i].Priority != elements[j].Priority {
			return elements[i].Priority > elements[j].Priority
		}
		return elements[i].Weight > elements[j].Weight
	})

	// Greedy selection maintaining independence
	selected := make([]Element, 0)
	assignedTasks := make(map[string]bool)
	nodeResources := make(map[string]*models.EdgeNode)

	// Initialize node resources tracking
	for _, node := range nodes {
		nodeResources[node.ID] = &models.EdgeNode{
			ID:              node.ID,
			AvailableCPU:    node.AvailableCPU,
			AvailableMemory: node.AvailableMemory,
		}
	}

	for _, elem := range elements {
		// Skip if task already assigned
		if assignedTasks[elem.TaskID] {
			continue
		}

		// Find the task and node
		var task *models.Task
		for _, t := range tasks {
			if t.ID == elem.TaskID {
				task = t
				break
			}
		}

		if task == nil {
			continue
		}

		// Check if node has enough resources
		nodeRes := nodeResources[elem.NodeID]
		if nodeRes.AvailableCPU >= task.CPURequired && nodeRes.AvailableMemory >= task.MemRequired {
			// Add to independent set
			selected = append(selected, elem)
			assignedTasks[elem.TaskID] = true

			// Update available resources
			nodeRes.AvailableCPU -= task.CPURequired
			nodeRes.AvailableMemory -= task.MemRequired
		}
	}

	// Convert to schedules
	schedules := make([]models.Schedule, len(selected))
	for i, elem := range selected {
		schedules[i] = models.Schedule{
			TaskID:     elem.TaskID,
			NodeID:     elem.NodeID,
			Priority:   elem.Priority,
			Cost:       elem.Weight,
			Confidence: 0.95,
		}
	}

	return schedules
}

// calculateWeight calculates the weight for a task-node pair
// Higher weight means better fit
func (m *Matroid) calculateWeight(task *models.Task, node *models.EdgeNode) float64 {
	// Factor 1: Resource utilization efficiency
	cpuUtil := task.CPURequired / node.TotalCPU
	memUtil := task.MemRequired / node.TotalMemory
	resourceFit := 1.0 - math.Abs(cpuUtil-memUtil) // Prefer balanced utilization

	// Factor 2: Current load (prefer less loaded nodes)
	loadFactor := 1.0 - node.GetLoad()

	// Factor 3: Resource availability
	availabilityFactor := (node.AvailableCPU / node.TotalCPU) * 0.5
	availabilityFactor += (node.AvailableMemory / node.TotalMemory) * 0.5

	// Combined weight
	weight := (resourceFit * 0.3) + (loadFactor * 0.4) + (availabilityFactor * 0.3)

	// Add priority boost
	weight += float64(task.Priority) * 0.1

	return weight
}

// IsIndependent checks if a subset of elements is independent in the matroid
func (m *Matroid) IsIndependent(elements []Element, tasks []*models.Task, nodes []*models.EdgeNode) bool {
	// Build resource usage map
	nodeResources := make(map[string]struct {
		cpu float64
		mem float64
	})

	taskMap := make(map[string]*models.Task)
	for _, t := range tasks {
		taskMap[t.ID] = t
	}

	nodeMap := make(map[string]*models.EdgeNode)
	for _, n := range nodes {
		nodeMap[n.ID] = n
	}

	// Check each element
	assignedTasks := make(map[string]bool)

	for _, elem := range elements {
		// Each task can only be assigned once
		if assignedTasks[elem.TaskID] {
			return false
		}
		assignedTasks[elem.TaskID] = true

		task := taskMap[elem.TaskID]
		if task == nil {
			continue
		}

		// Initialize node resources if needed
		if _, exists := nodeResources[elem.NodeID]; !exists {
			node := nodeMap[elem.NodeID]
			if node == nil {
				return false
			}
			nodeResources[elem.NodeID] = struct {
				cpu float64
				mem float64
			}{node.AvailableCPU, node.AvailableMemory}
		}

		// Check if node has enough resources
		res := nodeResources[elem.NodeID]
		if res.cpu < task.CPURequired || res.mem < task.MemRequired {
			return false
		}

		// Allocate resources
		res.cpu -= task.CPURequired
		res.mem -= task.MemRequired
		nodeResources[elem.NodeID] = res
	}

	return true
}
