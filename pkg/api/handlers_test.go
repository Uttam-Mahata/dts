package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Uttam-Mahata/dts/pkg/edgenode"
	"github.com/Uttam-Mahata/dts/pkg/loadbalancer"
	"github.com/Uttam-Mahata/dts/pkg/models"
	"github.com/Uttam-Mahata/dts/pkg/scheduler"
)

func setupTestServer() *Server {
	nodeManager := edgenode.NewManager()
	sched := scheduler.NewScheduler(nodeManager)
	lb := loadbalancer.NewLoadBalancer(nodeManager)

	return NewServer(sched, nodeManager, lb)
}

func TestNewServer(t *testing.T) {
	server := setupTestServer()

	if server == nil {
		t.Fatal("Expected non-nil server")
	}

	if server.scheduler == nil {
		t.Error("Expected scheduler to be set")
	}

	if server.nodeManager == nil {
		t.Error("Expected node manager to be set")
	}

	if server.loadBalancer == nil {
		t.Error("Expected load balancer to be set")
	}
}

func TestSetupRoutes(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	if mux == nil {
		t.Fatal("Expected non-nil mux")
	}
}

func TestHandleHealth(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%v'", response["status"])
	}
}

func TestHandleRegisterNode(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	node := models.EdgeNode{
		ID:          "node-1",
		Name:        "Test Node",
		Location:    "datacenter-1",
		TotalCPU:    8.0,
		TotalMemory: 4096.0,
	}

	body, _ := json.Marshal(node)
	req := httptest.NewRequest(http.MethodPost, "/api/nodes/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["success"] != true {
		t.Error("Expected success to be true")
	}
}

func TestHandleRegisterNodeInvalidRequest(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	// Missing required fields
	node := models.EdgeNode{
		Location: "datacenter-1",
	}

	body, _ := json.Marshal(node)
	req := httptest.NewRequest(http.MethodPost, "/api/nodes/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleRegisterNodeWrongMethod(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/nodes/register", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandleGetNodes(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	// Register a node first
	node := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Test Node",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    8.0,
		AvailableMemory: 4096.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}
	server.nodeManager.RegisterNode(node)

	req := httptest.NewRequest(http.MethodGet, "/api/nodes", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	count := int(response["count"].(float64))
	if count != 1 {
		t.Errorf("Expected 1 node, got %d", count)
	}
}

func TestHandleSubmitTask(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	taskReq := models.TaskRequest{
		Name:        "Test Task",
		Priority:    5,
		CPURequired: 2.0,
		MemRequired: 1024.0,
		Duration:    300,
		Deadline:    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		Metadata: map[string]interface{}{
			"user": "test",
		},
	}

	body, _ := json.Marshal(taskReq)
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/submit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["success"] != true {
		t.Error("Expected success to be true")
	}
}

func TestHandleSubmitTaskInvalidDeadline(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	taskReq := models.TaskRequest{
		Name:        "Test Task",
		Priority:    5,
		CPURequired: 2.0,
		MemRequired: 1024.0,
		Duration:    300,
		Deadline:    "invalid-date",
	}

	body, _ := json.Marshal(taskReq)
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/submit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleSubmitTaskMissingName(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	taskReq := models.TaskRequest{
		Priority:    5,
		CPURequired: 2.0,
		MemRequired: 1024.0,
		Duration:    300,
		Deadline:    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
	}

	body, _ := json.Marshal(taskReq)
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/submit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleGetTasks(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	// Submit a task first
	task := &models.Task{
		ID:          "task-1",
		Name:        "Test Task",
		Priority:    5,
		CPURequired: 2.0,
		MemRequired: 1024.0,
		Duration:    300 * time.Second,
		Deadline:    time.Now().Add(1 * time.Hour),
	}
	server.scheduler.SubmitTask(task)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	count := int(response["count"].(float64))
	if count != 1 {
		t.Errorf("Expected 1 task, got %d", count)
	}
}

func TestHandleScheduleTasks(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	// Register a node
	node := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Test Node",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    8.0,
		AvailableMemory: 4096.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}
	server.nodeManager.RegisterNode(node)

	// Submit a task
	task := &models.Task{
		ID:          "task-1",
		Name:        "Test Task",
		Priority:    5,
		CPURequired: 2.0,
		MemRequired: 1024.0,
		Duration:    300 * time.Second,
		Deadline:    time.Now().Add(1 * time.Hour),
	}
	server.scheduler.SubmitTask(task)

	req := httptest.NewRequest(http.MethodPost, "/api/tasks/schedule", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.ScheduleResult
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected scheduling to succeed")
	}
}

func TestHandleScheduleTasksNoNodes(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	// Submit a task without any nodes
	task := &models.Task{
		ID:          "task-1",
		Name:        "Test Task",
		Priority:    5,
		CPURequired: 2.0,
		MemRequired: 1024.0,
		Duration:    300 * time.Second,
		Deadline:    time.Now().Add(1 * time.Hour),
	}
	server.scheduler.SubmitTask(task)

	req := httptest.NewRequest(http.MethodPost, "/api/tasks/schedule", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestHandleSystemStatus(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/system/status", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check that response contains expected fields
	expectedFields := []string{
		"total_nodes", "active_nodes", "total_tasks", "pending_tasks",
		"scheduled_tasks", "running_tasks", "completed_tasks", "failed_tasks",
		"system_load", "is_balanced", "timestamp",
	}

	for _, field := range expectedFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Expected field '%s' in response", field)
		}
	}
}

func TestHandleSystemLoad(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	// Register a node
	node := &models.EdgeNode{
		ID:              "node-1",
		Name:            "Test Node",
		TotalCPU:        8.0,
		TotalMemory:     4096.0,
		AvailableCPU:    4.0,
		AvailableMemory: 2048.0,
		Status:          models.EdgeNodeStatusActive,
		LastHeartbeat:   time.Now(),
		RunningTasks:    []string{},
	}
	server.nodeManager.RegisterNode(node)

	req := httptest.NewRequest(http.MethodGet, "/api/system/load", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, exists := response["node_loads"]; !exists {
		t.Error("Expected 'node_loads' field in response")
	}

	if _, exists := response["system_load"]; !exists {
		t.Error("Expected 'system_load' field in response")
	}

	if _, exists := response["timestamp"]; !exists {
		t.Error("Expected 'timestamp' field in response")
	}
}

func TestHandleTasksWrongMethod(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	req := httptest.NewRequest(http.MethodPost, "/api/tasks", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandleNodesWrongMethod(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	req := httptest.NewRequest(http.MethodPost, "/api/nodes", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandleSystemStatusWrongMethod(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	req := httptest.NewRequest(http.MethodPost, "/api/system/status", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandleSystemLoadWrongMethod(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	req := httptest.NewRequest(http.MethodPost, "/api/system/load", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandleScheduleTasksWrongMethod(t *testing.T) {
	server := setupTestServer()
	mux := server.SetupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/tasks/schedule", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}
