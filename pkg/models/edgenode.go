package models

import (
	"sync"
	"time"
)

// EdgeNodeStatus represents the current status of an edge node
type EdgeNodeStatus string

const (
	EdgeNodeStatusActive   EdgeNodeStatus = "active"
	EdgeNodeStatusInactive EdgeNodeStatus = "inactive"
	EdgeNodeStatusBusy     EdgeNodeStatus = "busy"
)

// EdgeNode represents a computing node in the edge network
type EdgeNode struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Location        string         `json:"location"`
	TotalCPU        float64        `json:"total_cpu"`
	TotalMemory     float64        `json:"total_memory"`
	AvailableCPU    float64        `json:"available_cpu"`
	AvailableMemory float64        `json:"available_memory"`
	Status          EdgeNodeStatus `json:"status"`
	LastHeartbeat   time.Time      `json:"last_heartbeat"`
	RunningTasks    []string       `json:"running_tasks"`
	mu              sync.RWMutex
}

// AllocateResources allocates resources for a task
func (n *EdgeNode) AllocateResources(cpu, memory float64) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.AvailableCPU >= cpu && n.AvailableMemory >= memory {
		n.AvailableCPU -= cpu
		n.AvailableMemory -= memory
		return true
	}
	return false
}

// ReleaseResources releases resources after task completion
func (n *EdgeNode) ReleaseResources(cpu, memory float64) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.AvailableCPU += cpu
	n.AvailableMemory += memory

	// Ensure resources don't exceed total capacity
	if n.AvailableCPU > n.TotalCPU {
		n.AvailableCPU = n.TotalCPU
	}
	if n.AvailableMemory > n.TotalMemory {
		n.AvailableMemory = n.TotalMemory
	}
}

// GetLoad returns the current load percentage of the node
func (n *EdgeNode) GetLoad() float64 {
	n.mu.RLock()
	defer n.mu.RUnlock()

	cpuLoad := (n.TotalCPU - n.AvailableCPU) / n.TotalCPU
	memLoad := (n.TotalMemory - n.AvailableMemory) / n.TotalMemory

	// Return average load
	return (cpuLoad + memLoad) / 2.0
}

// AddTask adds a task to the running tasks list
func (n *EdgeNode) AddTask(taskID string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.RunningTasks = append(n.RunningTasks, taskID)
}

// RemoveTask removes a task from the running tasks list
func (n *EdgeNode) RemoveTask(taskID string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	for i, id := range n.RunningTasks {
		if id == taskID {
			n.RunningTasks = append(n.RunningTasks[:i], n.RunningTasks[i+1:]...)
			break
		}
	}
}
