"""
Matroid-based optimization for task scheduling in edge computing.

This module implements matroid-based greedy algorithms for optimal task allocation
to edge nodes, considering resource constraints and load balancing.
"""

from typing import List, Dict, Set, Tuple, Optional
import numpy as np
from ..core.task import Task
from ..core.edge_node import EdgeNode


class MatroidOptimizer:
    """
    Implements matroid-based optimization for task scheduling.
    
    Uses the greedy algorithm on a matroid structure to find optimal or near-optimal
    task-to-node assignments while respecting resource constraints.
    """
    
    def __init__(self, alpha: float = 0.5, beta: float = 0.3, gamma: float = 0.2):
        """
        Initialize the matroid optimizer.
        
        Args:
            alpha: Weight for load balancing (0-1)
            beta: Weight for resource efficiency (0-1)
            gamma: Weight for network latency (0-1)
            
        Note: alpha + beta + gamma should equal 1.0
        """
        total = alpha + beta + gamma
        self.alpha = alpha / total
        self.beta = beta / total
        self.gamma = gamma / total
    
    def optimize_task_allocation(
        self,
        tasks: List[Task],
        nodes: List[EdgeNode]
    ) -> Dict[str, str]:
        """
        Optimize task allocation using matroid-based greedy algorithm.
        
        Args:
            tasks: List of tasks to allocate
            nodes: List of available edge nodes
            
        Returns:
            Dictionary mapping task_id to node_id
        """
        # Sort tasks by priority (descending) and resource requirements
        sorted_tasks = sorted(
            tasks,
            key=lambda t: (-t.priority, -t.cpu_requirement, -t.memory_requirement)
        )
        
        allocation: Dict[str, str] = {}
        
        for task in sorted_tasks:
            best_node = self._find_best_node(task, nodes, allocation)
            if best_node:
                allocation[task.task_id] = best_node.node_id
        
        return allocation
    
    def _find_best_node(
        self,
        task: Task,
        nodes: List[EdgeNode],
        current_allocation: Dict[str, str]
    ) -> Optional[EdgeNode]:
        """
        Find the best node for a task using matroid independence property.
        
        Args:
            task: Task to allocate
            nodes: Available edge nodes
            current_allocation: Current task-to-node mapping
            
        Returns:
            Best edge node or None if no suitable node found
        """
        best_node = None
        best_score = float('-inf')
        
        for node in nodes:
            # Check matroid independence: can this assignment be made?
            if not self._is_independent(task, node, current_allocation):
                continue
            
            # Calculate score for this node
            score = self._calculate_node_score(task, node)
            
            if score > best_score:
                best_score = score
                best_node = node
        
        return best_node
    
    def _is_independent(
        self,
        task: Task,
        node: EdgeNode,
        current_allocation: Dict[str, str]
    ) -> bool:
        """
        Check if adding this task to the node maintains matroid independence.
        
        In our context, independence means the node can accommodate the task
        without violating resource constraints.
        
        Args:
            task: Task to check
            node: Node to check
            current_allocation: Current allocations
            
        Returns:
            True if the assignment maintains independence
        """
        return node.can_accommodate(task.cpu_requirement, task.memory_requirement)
    
    def _calculate_node_score(self, task: Task, node: EdgeNode) -> float:
        """
        Calculate a score for assigning a task to a node.
        
        Higher score indicates better assignment. Considers:
        - Load balancing (prefer less loaded nodes)
        - Resource efficiency (good fit for resources)
        - Network latency (prefer lower latency)
        
        Args:
            task: Task to assign
            node: Node to evaluate
            
        Returns:
            Score for this assignment (higher is better)
        """
        # Load balancing score (prefer less loaded nodes)
        load_score = 1.0 - node.get_load_score()
        
        # Resource efficiency score (prefer nodes where resources are well-utilized)
        cpu_fit = min(task.cpu_requirement / node.cpu_available, 1.0)
        memory_fit = min(task.memory_requirement / node.memory_available, 1.0)
        efficiency_score = (cpu_fit + memory_fit) / 2.0
        
        # Network latency score (prefer lower latency)
        # Normalize latency to 0-1 range (assuming max latency of 1000ms)
        max_latency = 1000.0
        latency_score = 1.0 - min(node.network_latency / max_latency, 1.0)
        
        # Combine scores with weights
        total_score = (
            self.alpha * load_score +
            self.beta * efficiency_score +
            self.gamma * latency_score
        )
        
        return total_score
    
    def balance_load(
        self,
        nodes: List[EdgeNode],
        threshold: float = 0.8
    ) -> List[Tuple[str, str]]:
        """
        Identify tasks that should be migrated for better load balancing.
        
        Args:
            nodes: List of edge nodes
            threshold: Load threshold above which rebalancing is needed
            
        Returns:
            List of (task_id, target_node_id) tuples for migration
        """
        migrations = []
        
        # Separate overloaded and underloaded nodes
        overloaded = [n for n in nodes if n.get_load_score() > threshold]
        underloaded = [n for n in nodes if n.get_load_score() < threshold * 0.5]
        
        if not overloaded or not underloaded:
            return migrations
        
        # Try to migrate tasks from overloaded to underloaded nodes
        for overloaded_node in overloaded:
            for task_id in overloaded_node.assigned_tasks:
                # Find a suitable underloaded node
                for target_node in underloaded:
                    # Simple check - in practice, would need actual task info
                    if target_node.get_load_score() < threshold * 0.5:
                        migrations.append((task_id, target_node.node_id))
                        break
                
                # Stop if load is balanced enough
                if overloaded_node.get_load_score() <= threshold:
                    break
        
        return migrations
    
    def compute_makespan(
        self,
        tasks: List[Task],
        allocation: Dict[str, str]
    ) -> float:
        """
        Compute the makespan (total completion time) for the allocation.
        
        Args:
            tasks: List of tasks
            allocation: Task to node mapping
            
        Returns:
            Makespan in seconds
        """
        node_times: Dict[str, float] = {}
        
        for task in tasks:
            if task.task_id in allocation:
                node_id = allocation[task.task_id]
                node_times[node_id] = node_times.get(node_id, 0) + task.execution_time
        
        return max(node_times.values()) if node_times else 0.0
