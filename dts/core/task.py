"""Task representation for the distributed scheduling system."""

from dataclasses import dataclass, field
from typing import Optional, Dict, Any
from enum import Enum
import time


class TaskStatus(Enum):
    """Status of a task in the system."""
    PENDING = "pending"
    SCHEDULED = "scheduled"
    RUNNING = "running"
    COMPLETED = "completed"
    FAILED = "failed"


@dataclass
class Task:
    """Represents a computational task in the edge computing system."""
    
    task_id: str
    cpu_requirement: float  # CPU cores required
    memory_requirement: float  # Memory in GB
    execution_time: float  # Expected execution time in seconds
    priority: int = 1  # Higher value = higher priority
    status: TaskStatus = TaskStatus.PENDING
    assigned_node: Optional[str] = None
    created_at: float = field(default_factory=time.time)
    started_at: Optional[float] = None
    completed_at: Optional[float] = None
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    def __post_init__(self):
        """Validate task parameters."""
        if self.cpu_requirement <= 0:
            raise ValueError("CPU requirement must be positive")
        if self.memory_requirement <= 0:
            raise ValueError("Memory requirement must be positive")
        if self.execution_time <= 0:
            raise ValueError("Execution time must be positive")
    
    def assign_to_node(self, node_id: str):
        """Assign the task to a specific node."""
        self.assigned_node = node_id
        self.status = TaskStatus.SCHEDULED
    
    def start_execution(self):
        """Mark the task as running."""
        self.status = TaskStatus.RUNNING
        self.started_at = time.time()
    
    def complete_execution(self, success: bool = True):
        """Mark the task as completed or failed."""
        self.completed_at = time.time()
        self.status = TaskStatus.COMPLETED if success else TaskStatus.FAILED
    
    def get_resource_requirements(self) -> Dict[str, float]:
        """Get the resource requirements as a dictionary."""
        return {
            'cpu': self.cpu_requirement,
            'memory': self.memory_requirement,
            'time': self.execution_time
        }
    
    def __hash__(self):
        """Make Task hashable for use in sets."""
        return hash(self.task_id)
    
    def __eq__(self, other):
        """Compare tasks by their ID."""
        if not isinstance(other, Task):
            return False
        return self.task_id == other.task_id
