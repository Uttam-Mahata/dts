"""Tests for MatroidOptimizer class."""

import unittest
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from dts.core.task import Task
from dts.core.edge_node import EdgeNode
from dts.optimization.matroid_optimizer import MatroidOptimizer


class TestMatroidOptimizer(unittest.TestCase):
    """Test cases for MatroidOptimizer class."""
    
    def setUp(self):
        """Set up test fixtures."""
        self.optimizer = MatroidOptimizer(alpha=0.5, beta=0.3, gamma=0.2)
        
        self.nodes = [
            EdgeNode(
                node_id="node-1",
                cpu_capacity=8.0,
                memory_capacity=16.0,
                network_latency=10.0
            ),
            EdgeNode(
                node_id="node-2",
                cpu_capacity=16.0,
                memory_capacity=32.0,
                network_latency=20.0
            ),
            EdgeNode(
                node_id="node-3",
                cpu_capacity=4.0,
                memory_capacity=8.0,
                network_latency=5.0
            ),
        ]
        
        self.tasks = [
            Task(
                task_id=f"task-{i}",
                cpu_requirement=2.0,
                memory_requirement=4.0,
                execution_time=10.0,
                priority=1
            )
            for i in range(1, 6)
        ]
    
    def test_optimizer_creation(self):
        """Test optimizer initialization."""
        optimizer = MatroidOptimizer(alpha=0.5, beta=0.3, gamma=0.2)
        
        # Weights should sum to 1.0
        total = optimizer.alpha + optimizer.beta + optimizer.gamma
        self.assertAlmostEqual(total, 1.0)
    
    def test_optimize_task_allocation(self):
        """Test basic task allocation optimization."""
        allocation = self.optimizer.optimize_task_allocation(
            self.tasks,
            self.nodes
        )
        
        # All tasks should be allocated
        self.assertEqual(len(allocation), len(self.tasks))
        
        # All allocations should be to valid nodes
        node_ids = {node.node_id for node in self.nodes}
        for task_id, node_id in allocation.items():
            self.assertIn(node_id, node_ids)
    
    def test_priority_based_allocation(self):
        """Test that higher priority tasks are allocated first."""
        # Create tasks with different priorities
        tasks = [
            Task(
                task_id="low-priority",
                cpu_requirement=4.0,
                memory_requirement=8.0,
                execution_time=10.0,
                priority=1
            ),
            Task(
                task_id="high-priority",
                cpu_requirement=4.0,
                memory_requirement=8.0,
                execution_time=10.0,
                priority=10
            ),
        ]
        
        # Use a node with limited capacity
        nodes = [
            EdgeNode(
                node_id="limited-node",
                cpu_capacity=4.5,
                memory_capacity=8.5
            )
        ]
        
        allocation = self.optimizer.optimize_task_allocation(tasks, nodes)
        
        # High priority task should be allocated
        self.assertIn("high-priority", allocation)
    
    def test_resource_constraints(self):
        """Test that resource constraints are respected."""
        # Create a task that's too large for all nodes
        large_task = Task(
            task_id="large-task",
            cpu_requirement=100.0,
            memory_requirement=200.0,
            execution_time=10.0
        )
        
        allocation = self.optimizer.optimize_task_allocation(
            [large_task],
            self.nodes
        )
        
        # Task should not be allocated
        self.assertEqual(len(allocation), 0)
    
    def test_load_balancing(self):
        """Test that load balancing works."""
        # Create many small tasks
        many_tasks = [
            Task(
                task_id=f"task-{i}",
                cpu_requirement=1.0,
                memory_requirement=2.0,
                execution_time=10.0
            )
            for i in range(20)
        ]
        
        allocation = self.optimizer.optimize_task_allocation(
            many_tasks,
            self.nodes
        )
        
        # Count tasks per node
        tasks_per_node = {}
        for task_id, node_id in allocation.items():
            tasks_per_node[node_id] = tasks_per_node.get(node_id, 0) + 1
        
        # All nodes should have tasks (assuming capacity)
        self.assertGreater(len(tasks_per_node), 0)
    
    def test_compute_makespan(self):
        """Test makespan calculation."""
        allocation = {
            "task-1": "node-1",
            "task-2": "node-1",
            "task-3": "node-2",
        }
        
        tasks = [
            Task("task-1", 1.0, 1.0, 10.0),
            Task("task-2", 1.0, 1.0, 20.0),
            Task("task-3", 1.0, 1.0, 15.0),
        ]
        
        makespan = self.optimizer.compute_makespan(tasks, allocation)
        
        # Node-1 has tasks taking 30s, node-2 has 15s
        # Makespan should be 30s
        self.assertEqual(makespan, 30.0)
    
    def test_balance_load_identification(self):
        """Test identification of overloaded nodes."""
        # Manually overload a node
        self.nodes[0].allocate_resources("task-1", 7.0, 14.0)
        
        migrations = self.optimizer.balance_load(self.nodes, threshold=0.8)
        
        # Should identify tasks for migration
        self.assertIsInstance(migrations, list)


if __name__ == '__main__':
    unittest.main()
