"""Tests for Task class."""

import unittest
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from dts.core.task import Task, TaskStatus


class TestTask(unittest.TestCase):
    """Test cases for Task class."""
    
    def test_task_creation(self):
        """Test basic task creation."""
        task = Task(
            task_id="test-1",
            cpu_requirement=2.0,
            memory_requirement=4.0,
            execution_time=10.0
        )
        
        self.assertEqual(task.task_id, "test-1")
        self.assertEqual(task.cpu_requirement, 2.0)
        self.assertEqual(task.memory_requirement, 4.0)
        self.assertEqual(task.execution_time, 10.0)
        self.assertEqual(task.status, TaskStatus.PENDING)
    
    def test_task_validation(self):
        """Test task parameter validation."""
        with self.assertRaises(ValueError):
            Task(
                task_id="test-1",
                cpu_requirement=-1.0,
                memory_requirement=4.0,
                execution_time=10.0
            )
        
        with self.assertRaises(ValueError):
            Task(
                task_id="test-1",
                cpu_requirement=2.0,
                memory_requirement=-1.0,
                execution_time=10.0
            )
    
    def test_task_assignment(self):
        """Test task assignment to node."""
        task = Task(
            task_id="test-1",
            cpu_requirement=2.0,
            memory_requirement=4.0,
            execution_time=10.0
        )
        
        task.assign_to_node("node-1")
        self.assertEqual(task.assigned_node, "node-1")
        self.assertEqual(task.status, TaskStatus.SCHEDULED)
    
    def test_task_execution_lifecycle(self):
        """Test task execution lifecycle."""
        task = Task(
            task_id="test-1",
            cpu_requirement=2.0,
            memory_requirement=4.0,
            execution_time=10.0
        )
        
        # Start execution
        task.start_execution()
        self.assertEqual(task.status, TaskStatus.RUNNING)
        self.assertIsNotNone(task.started_at)
        
        # Complete execution
        task.complete_execution(success=True)
        self.assertEqual(task.status, TaskStatus.COMPLETED)
        self.assertIsNotNone(task.completed_at)
    
    def test_task_failure(self):
        """Test task failure handling."""
        task = Task(
            task_id="test-1",
            cpu_requirement=2.0,
            memory_requirement=4.0,
            execution_time=10.0
        )
        
        task.start_execution()
        task.complete_execution(success=False)
        self.assertEqual(task.status, TaskStatus.FAILED)
    
    def test_task_resource_requirements(self):
        """Test getting resource requirements."""
        task = Task(
            task_id="test-1",
            cpu_requirement=2.0,
            memory_requirement=4.0,
            execution_time=10.0
        )
        
        requirements = task.get_resource_requirements()
        self.assertEqual(requirements['cpu'], 2.0)
        self.assertEqual(requirements['memory'], 4.0)
        self.assertEqual(requirements['time'], 10.0)
    
    def test_task_hashable(self):
        """Test that tasks are hashable."""
        task1 = Task(
            task_id="test-1",
            cpu_requirement=2.0,
            memory_requirement=4.0,
            execution_time=10.0
        )
        task2 = Task(
            task_id="test-1",
            cpu_requirement=3.0,
            memory_requirement=5.0,
            execution_time=15.0
        )
        
        # Same ID should be equal
        self.assertEqual(task1, task2)
        
        # Should be usable in sets
        task_set = {task1, task2}
        self.assertEqual(len(task_set), 1)


if __name__ == '__main__':
    unittest.main()
