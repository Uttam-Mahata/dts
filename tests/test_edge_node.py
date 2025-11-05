"""Tests for EdgeNode class."""

import unittest
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from dts.core.edge_node import EdgeNode, NodeStatus


class TestEdgeNode(unittest.TestCase):
    """Test cases for EdgeNode class."""
    
    def test_node_creation(self):
        """Test basic node creation."""
        node = EdgeNode(
            node_id="node-1",
            cpu_capacity=8.0,
            memory_capacity=16.0,
            location="datacenter-1"
        )
        
        self.assertEqual(node.node_id, "node-1")
        self.assertEqual(node.cpu_capacity, 8.0)
        self.assertEqual(node.memory_capacity, 16.0)
        self.assertEqual(node.cpu_available, 8.0)
        self.assertEqual(node.memory_available, 16.0)
        self.assertEqual(node.status, NodeStatus.ACTIVE)
    
    def test_node_validation(self):
        """Test node parameter validation."""
        with self.assertRaises(ValueError):
            EdgeNode(
                node_id="node-1",
                cpu_capacity=-1.0,
                memory_capacity=16.0
            )
        
        with self.assertRaises(ValueError):
            EdgeNode(
                node_id="node-1",
                cpu_capacity=8.0,
                memory_capacity=-1.0
            )
    
    def test_can_accommodate(self):
        """Test resource accommodation check."""
        node = EdgeNode(
            node_id="node-1",
            cpu_capacity=8.0,
            memory_capacity=16.0
        )
        
        # Should be able to accommodate
        self.assertTrue(node.can_accommodate(2.0, 4.0))
        self.assertTrue(node.can_accommodate(8.0, 16.0))
        
        # Should not be able to accommodate
        self.assertFalse(node.can_accommodate(10.0, 4.0))
        self.assertFalse(node.can_accommodate(2.0, 20.0))
    
    def test_allocate_resources(self):
        """Test resource allocation."""
        node = EdgeNode(
            node_id="node-1",
            cpu_capacity=8.0,
            memory_capacity=16.0
        )
        
        # Allocate resources
        success = node.allocate_resources("task-1", 2.0, 4.0)
        self.assertTrue(success)
        self.assertEqual(node.cpu_available, 6.0)
        self.assertEqual(node.memory_available, 12.0)
        self.assertIn("task-1", node.assigned_tasks)
        
        # Cannot over-allocate
        success = node.allocate_resources("task-2", 10.0, 4.0)
        self.assertFalse(success)
    
    def test_release_resources(self):
        """Test resource release."""
        node = EdgeNode(
            node_id="node-1",
            cpu_capacity=8.0,
            memory_capacity=16.0
        )
        
        # Allocate and release
        node.allocate_resources("task-1", 2.0, 4.0)
        node.release_resources("task-1", 2.0, 4.0)
        
        self.assertEqual(node.cpu_available, 8.0)
        self.assertEqual(node.memory_available, 16.0)
        self.assertNotIn("task-1", node.assigned_tasks)
    
    def test_utilization_calculations(self):
        """Test utilization calculations."""
        node = EdgeNode(
            node_id="node-1",
            cpu_capacity=8.0,
            memory_capacity=16.0
        )
        
        # Initial utilization
        self.assertEqual(node.get_cpu_utilization(), 0.0)
        self.assertEqual(node.get_memory_utilization(), 0.0)
        
        # After allocation
        node.allocate_resources("task-1", 4.0, 8.0)
        self.assertEqual(node.get_cpu_utilization(), 0.5)
        self.assertEqual(node.get_memory_utilization(), 0.5)
    
    def test_load_score(self):
        """Test load score calculation."""
        node = EdgeNode(
            node_id="node-1",
            cpu_capacity=8.0,
            memory_capacity=16.0
        )
        
        # No load
        self.assertEqual(node.get_load_score(), 0.0)
        
        # Half load
        node.allocate_resources("task-1", 4.0, 8.0)
        load = node.get_load_score()
        self.assertGreater(load, 0.0)
        self.assertLess(load, 1.0)
    
    def test_overloaded_status(self):
        """Test node overloaded status."""
        node = EdgeNode(
            node_id="node-1",
            cpu_capacity=8.0,
            memory_capacity=16.0
        )
        
        # Allocate heavily
        node.allocate_resources("task-1", 7.5, 15.0)
        self.assertEqual(node.status, NodeStatus.OVERLOADED)
        
        # Release to become active again
        node.release_resources("task-1", 7.5, 15.0)
        self.assertEqual(node.status, NodeStatus.ACTIVE)
    
    def test_node_hashable(self):
        """Test that nodes are hashable."""
        node1 = EdgeNode(
            node_id="node-1",
            cpu_capacity=8.0,
            memory_capacity=16.0
        )
        node2 = EdgeNode(
            node_id="node-1",
            cpu_capacity=16.0,
            memory_capacity=32.0
        )
        
        # Same ID should be equal
        self.assertEqual(node1, node2)
        
        # Should be usable in sets
        node_set = {node1, node2}
        self.assertEqual(len(node_set), 1)


if __name__ == '__main__':
    unittest.main()
