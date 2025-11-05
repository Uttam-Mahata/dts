package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Uttam-Mahata/dts/pkg/edgenode"
	"github.com/Uttam-Mahata/dts/pkg/loadbalancer"
	"github.com/Uttam-Mahata/dts/pkg/models"
	"github.com/Uttam-Mahata/dts/pkg/scheduler"
)

// Server represents the API server
type Server struct {
	scheduler    *scheduler.Scheduler
	nodeManager  *edgenode.Manager
	loadBalancer *loadbalancer.LoadBalancer
}

// NewServer creates a new API server
func NewServer(sched *scheduler.Scheduler, nodeMgr *edgenode.Manager, lb *loadbalancer.LoadBalancer) *Server {
	return &Server{
		scheduler:    sched,
		nodeManager:  nodeMgr,
		loadBalancer: lb,
	}
}

// SetupRoutes sets up the HTTP routes
func (s *Server) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Task endpoints
	mux.HandleFunc("/api/tasks", s.handleTasks)
	mux.HandleFunc("/api/tasks/submit", s.handleSubmitTask)
	mux.HandleFunc("/api/tasks/schedule", s.handleScheduleTasks)

	// Node endpoints
	mux.HandleFunc("/api/nodes", s.handleNodes)
	mux.HandleFunc("/api/nodes/register", s.handleRegisterNode)

	// System endpoints
	mux.HandleFunc("/api/system/status", s.handleSystemStatus)
	mux.HandleFunc("/api/system/load", s.handleSystemLoad)

	// Health check
	mux.HandleFunc("/health", s.handleHealth)

	return mux
}

// handleTasks handles GET requests for tasks
func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks := s.scheduler.GetAllTasks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks": tasks,
		"count": len(tasks),
	})
}

// handleSubmitTask handles POST requests for task submission
func (s *Server) handleSubmitTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Name == "" {
		http.Error(w, "Task name is required", http.StatusBadRequest)
		return
	}

	// Parse deadline
	deadline, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid deadline format: %v", err), http.StatusBadRequest)
		return
	}

	// Create task
	task := &models.Task{
		ID:          fmt.Sprintf("task-%d", time.Now().UnixNano()),
		Name:        req.Name,
		Priority:    req.Priority,
		CPURequired: req.CPURequired,
		MemRequired: req.MemRequired,
		Duration:    time.Duration(req.Duration) * time.Second,
		Deadline:    deadline,
		Metadata:    req.Metadata,
	}

	// Submit task
	if err := s.scheduler.SubmitTask(task); err != nil {
		http.Error(w, fmt.Sprintf("Failed to submit task: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"task":    task,
		"message": "Task submitted successfully",
	})
}

// handleScheduleTasks handles POST requests to trigger scheduling
func (s *Server) handleScheduleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result, err := s.scheduler.Schedule()
	if err != nil {
		http.Error(w, fmt.Sprintf("Scheduling failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleNodes handles GET requests for nodes
func (s *Server) handleNodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nodes := s.nodeManager.GetAllNodes()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"nodes": nodes,
		"count": len(nodes),
	})
}

// handleRegisterNode handles POST requests for node registration
func (s *Server) handleRegisterNode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var node models.EdgeNode
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate node
	if node.ID == "" || node.Name == "" {
		http.Error(w, "Node ID and name are required", http.StatusBadRequest)
		return
	}

	// Set available resources equal to total initially
	node.AvailableCPU = node.TotalCPU
	node.AvailableMemory = node.TotalMemory

	// Register node
	if err := s.nodeManager.RegisterNode(&node); err != nil {
		http.Error(w, fmt.Sprintf("Failed to register node: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"node":    node,
		"message": "Node registered successfully",
	})
}

// handleSystemStatus handles GET requests for system status
func (s *Server) handleSystemStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := map[string]interface{}{
		"total_nodes":       s.nodeManager.GetNodeCount(),
		"active_nodes":      s.nodeManager.GetActiveNodeCount(),
		"total_tasks":       s.scheduler.GetTaskCount(),
		"pending_tasks":     len(s.scheduler.GetPendingTasks()),
		"scheduled_tasks":   len(s.scheduler.GetTasksByStatus(models.TaskStatusScheduled)),
		"running_tasks":     len(s.scheduler.GetTasksByStatus(models.TaskStatusRunning)),
		"completed_tasks":   len(s.scheduler.GetTasksByStatus(models.TaskStatusCompleted)),
		"failed_tasks":      len(s.scheduler.GetTasksByStatus(models.TaskStatusFailed)),
		"system_load":       s.loadBalancer.GetSystemLoad(),
		"is_balanced":       s.loadBalancer.IsBalanced(0.2),
		"timestamp":         time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// handleSystemLoad handles GET requests for system load
func (s *Server) handleSystemLoad(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	loads := s.loadBalancer.GetNodeLoads()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"node_loads":  loads,
		"system_load": s.loadBalancer.GetSystemLoad(),
		"timestamp":   time.Now().UTC(),
	})
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"time":   time.Now().UTC(),
	})
}
