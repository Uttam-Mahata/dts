"""Edge node representation for the distributed computing system."""

from dataclasses import dataclass, field
from typing import List, Dict, Optional
from enum import Enum
import time


class NodeStatus(Enum):
    """Status of an edge node."""
    ACTIVE = "active"
    INACTIVE = "inactive"
    OVERLOADED = "overloaded"
    MAINTENANCE = "maintenance"


@dataclass
class EdgeNode:
    """Represents an edge computing node in the distributed system."""
    
    node_id: str
    cpu_capacity: float  # Total CPU cores
    memory_capacity: float  # Total memory in GB
    location: str = "unknown"
    status: NodeStatus = NodeStatus.ACTIVE
    cpu_available: Optional[float] = None
    memory_available: Optional[float] = None
    assigned_tasks: List[str] = field(default_factory=list)
    network_latency: float = 0.0  # Network latency in ms
    last_heartbeat: float = field(default_factory=time.time)
    metadata: Dict = field(default_factory=dict)
    
    def __post_init__(self):
        """Initialize available resources if not provided."""
        if self.cpu_available is None:
            self.cpu_available = self.cpu_capacity
        if self.memory_available is None:
            self.memory_available = self.memory_capacity
            
        if self.cpu_capacity <= 0:
            raise ValueError("CPU capacity must be positive")
        if self.memory_capacity <= 0:
            raise ValueError("Memory capacity must be positive")
    
    def can_accommodate(self, cpu_req: float, memory_req: float) -> bool:
        """Check if the node can accommodate the given resource requirements."""
        if self.status != NodeStatus.ACTIVE:
            return False
        return (self.cpu_available >= cpu_req and 
                self.memory_available >= memory_req)
    
    def allocate_resources(self, task_id: str, cpu_req: float, memory_req: float) -> bool:
        """Allocate resources for a task."""
        if not self.can_accommodate(cpu_req, memory_req):
            return False
        
        self.cpu_available -= cpu_req
        self.memory_available -= memory_req
        self.assigned_tasks.append(task_id)
        
        # Update status if overloaded
        if self.get_cpu_utilization() > 0.9 or self.get_memory_utilization() > 0.9:
            self.status = NodeStatus.OVERLOADED
        
        return True
    
    def release_resources(self, task_id: str, cpu_req: float, memory_req: float):
        """Release resources after task completion."""
        if task_id in self.assigned_tasks:
            self.assigned_tasks.remove(task_id)
        
        self.cpu_available = min(self.cpu_available + cpu_req, self.cpu_capacity)
        self.memory_available = min(self.memory_available + memory_req, self.memory_capacity)
        
        # Update status if no longer overloaded
        if self.status == NodeStatus.OVERLOADED:
            if self.get_cpu_utilization() <= 0.9 and self.get_memory_utilization() <= 0.9:
                self.status = NodeStatus.ACTIVE
    
    def get_cpu_utilization(self) -> float:
        """Calculate current CPU utilization as a percentage."""
        return (self.cpu_capacity - self.cpu_available) / self.cpu_capacity
    
    def get_memory_utilization(self) -> float:
        """Calculate current memory utilization as a percentage."""
        return (self.memory_capacity - self.memory_available) / self.memory_capacity
    
    def get_load_score(self) -> float:
        """Calculate a load score for the node (0.0 = no load, 1.0 = full capacity)."""
        cpu_util = self.get_cpu_utilization()
        memory_util = self.get_memory_utilization()
        # Weighted average with CPU having slightly more weight
        return 0.6 * cpu_util + 0.4 * memory_util
    
    def update_heartbeat(self):
        """Update the last heartbeat timestamp."""
        self.last_heartbeat = time.time()
    
    def is_healthy(self, timeout: float = 60.0) -> bool:
        """Check if the node is healthy based on heartbeat."""
        return (time.time() - self.last_heartbeat) < timeout
    
    def __hash__(self):
        """Make EdgeNode hashable for use in sets."""
        return hash(self.node_id)
    
    def __eq__(self, other):
        """Compare nodes by their ID."""
        if not isinstance(other, EdgeNode):
            return False
        return self.node_id == other.node_id
