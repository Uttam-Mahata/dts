package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/Uttam-Mahata/dts/pkg/edgenode"
	"github.com/Uttam-Mahata/dts/pkg/matroid"
	"github.com/Uttam-Mahata/dts/pkg/models"
)

// Scheduler manages task scheduling using matroid-based optimization
type Scheduler struct {
	tasks       map[string]*models.Task
	nodeManager *edgenode.Manager
	matroid     *matroid.Matroid
	mu          sync.RWMutex
}

// NewScheduler creates a new scheduler instance
func NewScheduler(nodeManager *edgenode.Manager) *Scheduler {
	return &Scheduler{
		tasks:       make(map[string]*models.Task),
		nodeManager: nodeManager,
		matroid:     matroid.NewMatroid(),
	}
}

// SubmitTask submits a new task to the scheduler
func (s *Scheduler) SubmitTask(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.ID]; exists {
		return fmt.Errorf("task %s already exists", task.ID)
	}

	task.Status = models.TaskStatusPending
	task.CreatedAt = time.Now()
	s.tasks[task.ID] = task

	return nil
}

// GetTask retrieves a task by ID
func (s *Scheduler) GetTask(taskID string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task %s not found", taskID)
	}

	return task, nil
}

// GetAllTasks returns all tasks
func (s *Scheduler) GetAllTasks() []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// GetPendingTasks returns all pending tasks
func (s *Scheduler) GetPendingTasks() []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0)
	for _, task := range s.tasks {
		if task.Status == models.TaskStatusPending {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

// Schedule performs task scheduling using matroid-based optimization
func (s *Scheduler) Schedule() (*models.ScheduleResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get pending tasks
	pendingTasks := make([]*models.Task, 0)
	for _, task := range s.tasks {
		if task.Status == models.TaskStatusPending {
			pendingTasks = append(pendingTasks, task)
		}
	}

	if len(pendingTasks) == 0 {
		return &models.ScheduleResult{
			Success:   true,
			Schedules: []models.Schedule{},
			Message:   "No pending tasks to schedule",
		}, nil
	}

	// Get active nodes
	nodes := s.nodeManager.GetActiveNodes()
	if len(nodes) == 0 {
		return &models.ScheduleResult{
			Success: false,
			Message: "No active nodes available",
		}, fmt.Errorf("no active nodes available")
	}

	// Use matroid optimization for scheduling
	schedules := s.matroid.GreedyOptimize(pendingTasks, nodes)

	if len(schedules) == 0 {
		return &models.ScheduleResult{
			Success: false,
			Message: "Could not schedule any tasks - insufficient resources",
		}, nil
	}

	// Apply schedules
	for _, schedule := range schedules {
		task := s.tasks[schedule.TaskID]
		node, err := s.nodeManager.GetNode(schedule.NodeID)
		if err != nil {
			continue
		}

		// Allocate resources
		if node.AllocateResources(task.CPURequired, task.MemRequired) {
			task.Status = models.TaskStatusScheduled
			task.AssignedTo = schedule.NodeID
			node.AddTask(task.ID)
		}
	}

	return &models.ScheduleResult{
		Success:   true,
		Schedules: schedules,
		Message:   fmt.Sprintf("Successfully scheduled %d tasks", len(schedules)),
	}, nil
}

// StartTask marks a task as running
func (s *Scheduler) StartTask(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}

	if task.Status != models.TaskStatusScheduled {
		return fmt.Errorf("task %s is not in scheduled state", taskID)
	}

	now := time.Now()
	task.Status = models.TaskStatusRunning
	task.StartedAt = &now

	return nil
}

// CompleteTask marks a task as completed
func (s *Scheduler) CompleteTask(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}

	if task.Status != models.TaskStatusRunning && task.Status != models.TaskStatusScheduled {
		return fmt.Errorf("task %s is not in running or scheduled state", taskID)
	}

	now := time.Now()
	task.Status = models.TaskStatusCompleted
	task.CompletedAt = &now

	// Release resources
	if task.AssignedTo != "" {
		node, err := s.nodeManager.GetNode(task.AssignedTo)
		if err == nil {
			node.ReleaseResources(task.CPURequired, task.MemRequired)
			node.RemoveTask(task.ID)
		}
	}

	return nil
}

// FailTask marks a task as failed
func (s *Scheduler) FailTask(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}

	task.Status = models.TaskStatusFailed

	// Release resources
	if task.AssignedTo != "" {
		node, err := s.nodeManager.GetNode(task.AssignedTo)
		if err == nil {
			node.ReleaseResources(task.CPURequired, task.MemRequired)
			node.RemoveTask(task.ID)
		}
	}

	return nil
}

// GetTaskCount returns the total number of tasks
func (s *Scheduler) GetTaskCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.tasks)
}

// GetTasksByStatus returns tasks filtered by status
func (s *Scheduler) GetTasksByStatus(status models.TaskStatus) []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0)
	for _, task := range s.tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}

	return tasks
}
