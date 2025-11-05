# Architecture Documentation

## Distributed Task Scheduling System for Edge Computing

### Overview

This system implements a distributed task scheduling and load balancing solution for edge computing environments using matroid-based optimization. The design prioritizes efficiency, scalability, and optimal resource utilization.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Task Scheduler                          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │          Matroid-Based Optimizer                    │   │
│  │  • Greedy Algorithm                                 │   │
│  │  • Independence Property                            │   │
│  │  • Multi-objective Optimization                     │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  Task Queue  │  │ Node Manager │  │   Metrics    │     │
│  │  (Pending)   │  │  (Registry)  │  │  Collector   │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
        ┌──────────────────────────────────────────┐
        │         Edge Computing Nodes              │
        ├──────────────┬──────────────┬────────────┤
        │   Node 1     │   Node 2     │   Node 3   │
        │  ┌────────┐  │  ┌────────┐  │ ┌────────┐ │
        │  │Task A  │  │  │Task B  │  │ │Task C  │ │
        │  │Task D  │  │  │Task E  │  │ │        │ │
        │  └────────┘  │  └────────┘  │ └────────┘ │
        └──────────────┴──────────────┴────────────┘
```

## Core Components

### 1. Task (`dts.core.task.Task`)

**Purpose**: Represents a computational task in the system.

**Key Attributes**:
- `task_id`: Unique identifier
- `cpu_requirement`: CPU cores needed
- `memory_requirement`: Memory in GB
- `execution_time`: Expected duration
- `priority`: Task priority (higher = more important)
- `status`: Current state (PENDING, SCHEDULED, RUNNING, COMPLETED, FAILED)

**Lifecycle**:
```
PENDING → SCHEDULED → RUNNING → COMPLETED/FAILED
```

### 2. EdgeNode (`dts.core.edge_node.EdgeNode`)

**Purpose**: Represents an edge computing node with resources.

**Key Attributes**:
- `node_id`: Unique identifier
- `cpu_capacity`: Total CPU cores
- `memory_capacity`: Total memory
- `cpu_available`: Available CPU
- `memory_available`: Available memory
- `status`: Node state (ACTIVE, INACTIVE, OVERLOADED, MAINTENANCE)
- `network_latency`: Network delay in ms

**Resource Management**:
- Allocates resources to tasks
- Tracks utilization
- Monitors health status
- Manages overload conditions

### 3. MatroidOptimizer (`dts.optimization.matroid_optimizer.MatroidOptimizer`)

**Purpose**: Implements matroid-based optimization for task allocation.

**Mathematical Foundation**:

A matroid M = (E, I) where:
- E: Ground set (all possible task-to-node assignments)
- I: Independent sets (valid assignments respecting constraints)

**Independence Property**:
An assignment is independent if:
1. Node has sufficient CPU: `cpu_available >= task.cpu_requirement`
2. Node has sufficient memory: `memory_available >= task.memory_requirement`

**Greedy Algorithm**:
```
1. Sort tasks by priority (descending)
2. For each task:
   a. Evaluate all nodes with scoring function
   b. Select node with highest score that maintains independence
   c. Assign task to selected node
```

**Scoring Function**:
```
score = α × load_score + β × efficiency_score + γ × latency_score

where:
- load_score = 1 - node.load (prefer less loaded)
- efficiency_score = avg(cpu_fit, memory_fit) (prefer good fit)
- latency_score = 1 - normalized_latency (prefer lower latency)
- α + β + γ = 1.0 (configurable weights)
```

**Optimality**: The greedy algorithm on a matroid achieves a (1 - 1/e) approximation for submodular functions, providing near-optimal solutions efficiently.

### 4. TaskScheduler (`dts.scheduler.task_scheduler.TaskScheduler`)

**Purpose**: Main scheduling engine coordinating all operations.

**Key Operations**:

1. **Task Submission**:
   - Adds tasks to pending queue
   - Thread-safe operation

2. **Scheduling**:
   - Gets pending tasks and active nodes
   - Uses MatroidOptimizer for allocation
   - Updates task and node states

3. **Load Balancing**:
   - Monitors node utilization
   - Identifies overloaded nodes
   - Migrates tasks when needed

4. **Monitoring**:
   - Tracks system metrics
   - Reports node status
   - Provides system overview

## Data Structures

### Task Queue
- **Type**: OrderedDict
- **Complexity**: O(1) insertion, O(1) deletion
- **Purpose**: Maintains pending tasks efficiently

### Node Registry
- **Type**: Dictionary
- **Complexity**: O(1) lookup
- **Purpose**: Fast node access

### Task Registry
- **Type**: Dictionary
- **Complexity**: O(1) lookup
- **Purpose**: Fast task access

## Algorithms

### Task Allocation Algorithm

```python
def allocate_tasks(tasks, nodes):
    # Sort by priority (descending)
    sorted_tasks = sort(tasks, key=priority, reverse=True)
    
    allocation = {}
    for task in sorted_tasks:
        best_node = None
        best_score = -∞
        
        for node in nodes:
            if is_independent(task, node):
                score = calculate_score(task, node)
                if score > best_score:
                    best_score = score
                    best_node = node
        
        if best_node:
            allocation[task.id] = best_node.id
            allocate_resources(task, best_node)
    
    return allocation
```

**Time Complexity**: O(n log n + nm) where n = tasks, m = nodes

### Load Balancing Algorithm

```python
def balance_load(nodes, threshold):
    overloaded = [n for n in nodes if n.load > threshold]
    underloaded = [n for n in nodes if n.load < threshold/2]
    
    migrations = []
    for overloaded_node in overloaded:
        for task_id in overloaded_node.tasks:
            for target_node in underloaded:
                if can_migrate(task_id, target_node):
                    migrations.append((task_id, target_node))
                    break
            
            if overloaded_node.load <= threshold:
                break
    
    return migrations
```

## Thread Safety

All scheduler operations are protected by locks:

```python
with self.lock:
    # Critical section
    # Modify shared state
```

This ensures:
- Consistent state
- No race conditions
- Safe concurrent access

## Performance Characteristics

### Scheduling Performance

| Operation | Complexity | Notes |
|-----------|------------|-------|
| Submit Task | O(1) | Add to queue |
| Schedule Tasks | O(n log n + nm) | Sort + allocation |
| Complete Task | O(1) | Update state |
| Node Status | O(1) | Dictionary lookup |
| Load Balance | O(mk) | m nodes, k tasks/node |

### Space Complexity

| Component | Space | Notes |
|-----------|-------|-------|
| Tasks | O(n) | n tasks |
| Nodes | O(m) | m nodes |
| Allocation | O(n) | Task-to-node mapping |
| Total | O(n + m) | Linear in inputs |

## Configuration

### Optimizer Weights

Adjust based on priorities:

```python
# Load balancing focused
MatroidOptimizer(alpha=0.7, beta=0.2, gamma=0.1)

# Network latency focused  
MatroidOptimizer(alpha=0.3, beta=0.2, gamma=0.5)

# Resource efficiency focused
MatroidOptimizer(alpha=0.2, beta=0.7, gamma=0.1)
```

### Scheduler Parameters

```python
config = {
    'scheduler': {
        'rebalance_threshold': 0.8,  # Load threshold
        'heartbeat_timeout': 60.0,   # Node health timeout
    },
    'system': {
        'max_tasks_per_node': 100,
        'scheduling_interval': 5.0,
    }
}
```

## Extension Points

### 1. Custom Scoring Functions

Implement custom node scoring:

```python
class CustomOptimizer(MatroidOptimizer):
    def _calculate_node_score(self, task, node):
        # Custom scoring logic
        return score
```

### 2. Task Priorities

Implement custom priority schemes:

```python
class UrgentTask(Task):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, priority=100, **kwargs)
```

### 3. Node Selection

Implement custom node selection:

```python
class CustomScheduler(TaskScheduler):
    def _find_best_node(self, task, nodes):
        # Custom selection logic
        return best_node
```

## Monitoring and Metrics

### System Metrics

- `total_tasks_scheduled`: Total tasks allocated
- `total_tasks_completed`: Successfully completed tasks
- `total_tasks_failed`: Failed tasks
- `average_waiting_time`: Average time in queue

### Node Metrics

- `cpu_utilization`: CPU usage percentage
- `memory_utilization`: Memory usage percentage
- `load_score`: Combined load metric
- `assigned_tasks`: Number of tasks
- `is_healthy`: Health status

### Performance Metrics

- `makespan`: Total completion time
- `throughput`: Tasks per second
- `utilization`: Resource usage efficiency

## Best Practices

### 1. Task Sizing

- Keep tasks reasonably sized
- Avoid extremely long-running tasks
- Use appropriate resource requirements

### 2. Priority Assignment

- Use priorities judiciously
- Avoid priority inversion
- Balance urgent and regular tasks

### 3. Node Configuration

- Configure realistic capacities
- Monitor network latency
- Update node status regularly

### 4. Load Management

- Run periodic load balancing
- Monitor overloaded nodes
- Scale nodes as needed

## Future Enhancements

1. **Dynamic Scaling**: Auto-scale nodes based on load
2. **Predictive Scheduling**: ML-based task time prediction
3. **Fault Tolerance**: Automatic task retry and failover
4. **Multi-tenancy**: Support for multiple users/applications
5. **Advanced Metrics**: More detailed performance analytics
6. **Resource Pools**: Group nodes by capabilities
7. **Task Dependencies**: DAG-based task workflows
8. **Heterogeneous Tasks**: GPU, storage, network tasks

## References

1. Oxley, J. G. (2011). "Matroid Theory" (2nd ed.). Oxford University Press.
2. Korte, B., & Vygen, J. (2018). "Combinatorial Optimization: Theory and Algorithms".
3. Shi, W., et al. (2016). "Edge Computing: Vision and Challenges". IEEE IoT Journal.
4. Topcuoglu, H., et al. (2002). "Performance-effective and low-complexity task scheduling".
