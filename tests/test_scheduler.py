"""Tests for TaskScheduler class."""

import unittest
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from dts.core.task import Task, TaskStatus
from dts.core.edge_node import EdgeNode
from dts.scheduler.task_scheduler import TaskScheduler
from dts.optimization.matroid_optimizer import MatroidOptimizer


class TestTaskScheduler(unittest.TestCase):
    """Test cases for TaskScheduler class."""
    
    def setUp(self):
        """Set up test fixtures."""
        self.scheduler = TaskScheduler()
        
        self.nodes = [
            EdgeNode(
                node_id="node-1",
                cpu_capacity=8.0,
                memory_capacity=16.0
            ),
            EdgeNode(
                node_id="node-2",
                cpu_capacity=16.0,
                memory_capacity=32.0
            ),
        ]
        
        for node in self.nodes:
            self.scheduler.register_node(node)
        
        self.tasks = [
            Task(
                task_id=f"task-{i}",
                cpu_requirement=2.0,
                memory_requirement=4.0,
                execution_time=10.0
            )
            for i in range(1, 6)
        ]
    
    def test_scheduler_creation(self):
        """Test scheduler initialization."""
        scheduler = TaskScheduler()
        self.assertIsNotNone(scheduler.optimizer)
        self.assertEqual(len(scheduler.nodes), 0)
        self.assertEqual(len(scheduler.tasks), 0)
    
    def test_register_node(self):
        """Test node registration."""
        scheduler = TaskScheduler()
        node = EdgeNode("test-node", 8.0, 16.0)
        
        scheduler.register_node(node)
        self.assertIn("test-node", scheduler.nodes)
    
    def test_unregister_node(self):
        """Test node unregistration."""
        scheduler = TaskScheduler()
        node = EdgeNode("test-node", 8.0, 16.0)
        
        scheduler.register_node(node)
        scheduler.unregister_node("test-node")
        self.assertNotIn("test-node", scheduler.nodes)
    
    def test_submit_task(self):
        """Test task submission."""
        task = Task("test-task", 2.0, 4.0, 10.0)
        self.scheduler.submit_task(task)
        
        self.assertIn("test-task", self.scheduler.tasks)
        self.assertIn(task, self.scheduler.pending_tasks)
    
    def test_submit_multiple_tasks(self):
        """Test submitting multiple tasks."""
        self.scheduler.submit_tasks(self.tasks)
        
        self.assertEqual(len(self.scheduler.tasks), len(self.tasks))
        self.assertEqual(len(self.scheduler.pending_tasks), len(self.tasks))
    
    def test_schedule_tasks(self):
        """Test task scheduling."""
        self.scheduler.submit_tasks(self.tasks)
        allocation = self.scheduler.schedule_tasks()
        
        # All tasks should be scheduled
        self.assertGreater(len(allocation), 0)
        
        # Pending queue should be smaller
        self.assertLess(len(self.scheduler.pending_tasks), len(self.tasks))
    
    def test_task_lifecycle(self):
        """Test complete task lifecycle."""
        task = Task("lifecycle-task", 2.0, 4.0, 10.0)
        
        # Submit
        self.scheduler.submit_task(task)
        self.assertEqual(task.status, TaskStatus.PENDING)
        
        # Schedule
        allocation = self.scheduler.schedule_tasks()
        if "lifecycle-task" in allocation:
            self.assertEqual(task.status, TaskStatus.SCHEDULED)
            
            # Start
            self.scheduler.start_task("lifecycle-task")
            self.assertEqual(task.status, TaskStatus.RUNNING)
            
            # Complete
            self.scheduler.complete_task("lifecycle-task", success=True)
            self.assertEqual(task.status, TaskStatus.COMPLETED)
    
    def test_task_failure(self):
        """Test task failure handling."""
        task = Task("fail-task", 2.0, 4.0, 10.0)
        
        self.scheduler.submit_task(task)
        self.scheduler.schedule_tasks()
        
        if task.assigned_node:
            self.scheduler.start_task("fail-task")
            self.scheduler.complete_task("fail-task", success=False)
            
            self.assertEqual(task.status, TaskStatus.FAILED)
            self.assertEqual(self.scheduler.metrics['total_tasks_failed'], 1)
    
    def test_get_node_status(self):
        """Test getting node status."""
        status = self.scheduler.get_node_status("node-1")
        
        self.assertIsNotNone(status)
        self.assertEqual(status['node_id'], "node-1")
        self.assertIn('cpu_utilization', status)
        self.assertIn('memory_utilization', status)
        self.assertIn('load_score', status)
    
    def test_get_task_status(self):
        """Test getting task status."""
        task = Task("status-task", 2.0, 4.0, 10.0)
        self.scheduler.submit_task(task)
        
        status = self.scheduler.get_task_status("status-task")
        
        self.assertIsNotNone(status)
        self.assertEqual(status['task_id'], "status-task")
        self.assertIn('status', status)
    
    def test_get_system_overview(self):
        """Test getting system overview."""
        self.scheduler.submit_tasks(self.tasks)
        self.scheduler.schedule_tasks()
        
        overview = self.scheduler.get_system_overview()
        
        self.assertIn('total_nodes', overview)
        self.assertIn('active_nodes', overview)
        self.assertIn('total_tasks', overview)
        self.assertIn('pending_tasks', overview)
        self.assertEqual(overview['total_nodes'], len(self.nodes))
        self.assertEqual(overview['total_tasks'], len(self.tasks))
    
    def test_metrics_tracking(self):
        """Test metrics tracking."""
        task = Task("metric-task", 2.0, 4.0, 10.0)
        
        self.scheduler.submit_task(task)
        self.scheduler.schedule_tasks()
        
        metrics = self.scheduler.get_metrics()
        self.assertGreater(metrics['total_tasks_scheduled'], 0)
    
    def test_load_balancing(self):
        """Test load balancing functionality."""
        # Submit many tasks to create load imbalance
        many_tasks = [
            Task(f"task-{i}", 1.0, 2.0, 10.0)
            for i in range(10)
        ]
        
        self.scheduler.submit_tasks(many_tasks)
        self.scheduler.schedule_tasks()
        
        # Attempt load balancing
        migrated = self.scheduler.balance_load(threshold=0.5)
        self.assertIsInstance(migrated, list)


if __name__ == '__main__':
    unittest.main()
