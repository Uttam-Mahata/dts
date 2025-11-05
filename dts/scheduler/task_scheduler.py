"""
Task scheduler for the distributed edge computing system.

This module implements the main task scheduling logic using matroid-based
optimization for load balancing and efficient resource allocation.
"""

from typing import List, Dict, Optional, Set
import threading
import time
from collections import deque

from ..core.task import Task, TaskStatus
from ..core.edge_node import EdgeNode, NodeStatus
from ..optimization.matroid_optimizer import MatroidOptimizer


class TaskScheduler:
    """
    Main task scheduler for distributed edge computing.
    
    Uses matroid-based optimization to allocate tasks to edge nodes
    while maintaining load balance and efficient resource utilization.
    """
    
    def __init__(self, optimizer: Optional[MatroidOptimizer] = None):
        """
        Initialize the task scheduler.
        
        Args:
            optimizer: Matroid optimizer instance (creates default if None)
        """
        self.optimizer = optimizer or MatroidOptimizer()
        self.nodes: Dict[str, EdgeNode] = {}
        self.tasks: Dict[str, Task] = {}
        self.pending_tasks: deque = deque()
        self.lock = threading.Lock()
        self.metrics = {
            'total_tasks_scheduled': 0,
            'total_tasks_completed': 0,
            'total_tasks_failed': 0,
            'average_waiting_time': 0.0,
        }
    
    def register_node(self, node: EdgeNode):
        """
        Register an edge node with the scheduler.
        
        Args:
            node: EdgeNode to register
        """
        with self.lock:
            self.nodes[node.node_id] = node
    
    def unregister_node(self, node_id: str):
        """
        Unregister an edge node from the scheduler.
        
        Args:
            node_id: ID of the node to unregister
        """
        with self.lock:
            if node_id in self.nodes:
                # Reassign tasks from this node
                node = self.nodes[node_id]
                for task_id in list(node.assigned_tasks):
                    if task_id in self.tasks:
                        task = self.tasks[task_id]
                        if task.status in [TaskStatus.SCHEDULED, TaskStatus.RUNNING]:
                            task.status = TaskStatus.PENDING
                            task.assigned_node = None
                            self.pending_tasks.append(task)
                
                del self.nodes[node_id]
    
    def submit_task(self, task: Task):
        """
        Submit a task to the scheduler.
        
        Args:
            task: Task to schedule
        """
        with self.lock:
            self.tasks[task.task_id] = task
            self.pending_tasks.append(task)
    
    def submit_tasks(self, tasks: List[Task]):
        """
        Submit multiple tasks to the scheduler.
        
        Args:
            tasks: List of tasks to schedule
        """
        with self.lock:
            for task in tasks:
                self.tasks[task.task_id] = task
                self.pending_tasks.append(task)
    
    def schedule_tasks(self) -> Dict[str, str]:
        """
        Schedule pending tasks to available nodes using matroid optimization.
        
        Returns:
            Dictionary mapping task_id to node_id for scheduled tasks
        """
        with self.lock:
            if not self.pending_tasks or not self.nodes:
                return {}
            
            # Get active nodes
            active_nodes = [
                node for node in self.nodes.values()
                if node.status == NodeStatus.ACTIVE and node.is_healthy()
            ]
            
            if not active_nodes:
                return {}
            
            # Convert pending tasks to list
            pending_list = list(self.pending_tasks)
            
            # Get optimal allocation using matroid optimizer
            allocation = self.optimizer.optimize_task_allocation(
                pending_list,
                active_nodes
            )
            
            # Apply the allocation
            scheduled = {}
            for task_id, node_id in allocation.items():
                if self._assign_task_to_node(task_id, node_id):
                    scheduled[task_id] = node_id
                    # Remove from pending queue
                    task = self.tasks[task_id]
                    if task in self.pending_tasks:
                        self.pending_tasks.remove(task)
                    self.metrics['total_tasks_scheduled'] += 1
            
            return scheduled
    
    def _assign_task_to_node(self, task_id: str, node_id: str) -> bool:
        """
        Assign a task to a node and allocate resources.
        
        Args:
            task_id: ID of task to assign
            node_id: ID of node to assign to
            
        Returns:
            True if assignment successful, False otherwise
        """
        if task_id not in self.tasks or node_id not in self.nodes:
            return False
        
        task = self.tasks[task_id]
        node = self.nodes[node_id]
        
        # Allocate resources
        if node.allocate_resources(
            task.task_id,
            task.cpu_requirement,
            task.memory_requirement
        ):
            task.assign_to_node(node_id)
            return True
        
        return False
    
    def start_task(self, task_id: str):
        """
        Mark a task as started.
        
        Args:
            task_id: ID of task to start
        """
        with self.lock:
            if task_id in self.tasks:
                self.tasks[task_id].start_execution()
    
    def complete_task(self, task_id: str, success: bool = True):
        """
        Mark a task as completed and release resources.
        
        Args:
            task_id: ID of task to complete
            success: Whether task completed successfully
        """
        with self.lock:
            if task_id not in self.tasks:
                return
            
            task = self.tasks[task_id]
            task.complete_execution(success)
            
            # Release resources on the node
            if task.assigned_node and task.assigned_node in self.nodes:
                node = self.nodes[task.assigned_node]
                node.release_resources(
                    task_id,
                    task.cpu_requirement,
                    task.memory_requirement
                )
            
            # Update metrics
            if success:
                self.metrics['total_tasks_completed'] += 1
            else:
                self.metrics['total_tasks_failed'] += 1
    
    def balance_load(self, threshold: float = 0.8) -> List[str]:
        """
        Rebalance load across nodes.
        
        Args:
            threshold: Load threshold for rebalancing
            
        Returns:
            List of task IDs that were migrated
        """
        with self.lock:
            active_nodes = [
                node for node in self.nodes.values()
                if node.status in [NodeStatus.ACTIVE, NodeStatus.OVERLOADED]
            ]
            
            migrations = self.optimizer.balance_load(active_nodes, threshold)
            
            migrated_tasks = []
            for task_id, target_node_id in migrations:
                if self._migrate_task(task_id, target_node_id):
                    migrated_tasks.append(task_id)
            
            return migrated_tasks
    
    def _migrate_task(self, task_id: str, target_node_id: str) -> bool:
        """
        Migrate a task from its current node to a target node.
        
        Args:
            task_id: ID of task to migrate
            target_node_id: ID of target node
            
        Returns:
            True if migration successful, False otherwise
        """
        if task_id not in self.tasks or target_node_id not in self.nodes:
            return False
        
        task = self.tasks[task_id]
        
        # Can only migrate scheduled or running tasks
        if task.status not in [TaskStatus.SCHEDULED, TaskStatus.RUNNING]:
            return False
        
        # Release from current node
        if task.assigned_node and task.assigned_node in self.nodes:
            current_node = self.nodes[task.assigned_node]
            current_node.release_resources(
                task_id,
                task.cpu_requirement,
                task.memory_requirement
            )
        
        # Assign to new node
        return self._assign_task_to_node(task_id, target_node_id)
    
    def get_node_status(self, node_id: str) -> Optional[Dict]:
        """
        Get status information for a node.
        
        Args:
            node_id: ID of node
            
        Returns:
            Dictionary with node status information
        """
        with self.lock:
            if node_id not in self.nodes:
                return None
            
            node = self.nodes[node_id]
            return {
                'node_id': node.node_id,
                'status': node.status.value,
                'cpu_utilization': node.get_cpu_utilization(),
                'memory_utilization': node.get_memory_utilization(),
                'load_score': node.get_load_score(),
                'assigned_tasks': len(node.assigned_tasks),
                'is_healthy': node.is_healthy(),
            }
    
    def get_task_status(self, task_id: str) -> Optional[Dict]:
        """
        Get status information for a task.
        
        Args:
            task_id: ID of task
            
        Returns:
            Dictionary with task status information
        """
        with self.lock:
            if task_id not in self.tasks:
                return None
            
            task = self.tasks[task_id]
            return {
                'task_id': task.task_id,
                'status': task.status.value,
                'assigned_node': task.assigned_node,
                'priority': task.priority,
                'cpu_requirement': task.cpu_requirement,
                'memory_requirement': task.memory_requirement,
            }
    
    def get_metrics(self) -> Dict:
        """
        Get scheduler metrics.
        
        Returns:
            Dictionary with scheduler metrics
        """
        with self.lock:
            return self.metrics.copy()
    
    def get_system_overview(self) -> Dict:
        """
        Get an overview of the entire system.
        
        Returns:
            Dictionary with system overview
        """
        with self.lock:
            total_nodes = len(self.nodes)
            active_nodes = sum(
                1 for n in self.nodes.values()
                if n.status == NodeStatus.ACTIVE
            )
            
            total_tasks = len(self.tasks)
            pending_tasks = len(self.pending_tasks)
            running_tasks = sum(
                1 for t in self.tasks.values()
                if t.status == TaskStatus.RUNNING
            )
            
            avg_load = (
                sum(n.get_load_score() for n in self.nodes.values()) / total_nodes
                if total_nodes > 0 else 0.0
            )
            
            return {
                'total_nodes': total_nodes,
                'active_nodes': active_nodes,
                'total_tasks': total_tasks,
                'pending_tasks': pending_tasks,
                'running_tasks': running_tasks,
                'average_load': avg_load,
                'metrics': self.metrics.copy(),
            }
