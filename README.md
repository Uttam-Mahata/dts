# DTS - Distributed Task Scheduling System

Distributed Task Scheduling and Load Balancing System for Edge Computing using Matroid-Based Optimization

## Overview

DTS is a sophisticated task scheduling system designed for edge computing environments. It uses matroid-based optimization algorithms to efficiently allocate computational tasks to edge nodes while maintaining optimal load balance and resource utilization.

### Key Features

- **Matroid-Based Optimization**: Uses mathematical matroid theory for optimal task allocation
- **Load Balancing**: Automatically balances workload across edge nodes
- **Resource-Aware Scheduling**: Considers CPU, memory, and network latency
- **Priority-Based Allocation**: Supports task priorities for critical workloads
- **Real-Time Monitoring**: Track system metrics and node status
- **Thread-Safe**: Safe for concurrent operations

## Installation

### Requirements

- Python 3.7 or higher
- numpy >= 1.21.0
- scipy >= 1.7.0
- networkx >= 2.6.0

### Install Dependencies

```bash
pip install -r requirements.txt
```

## Quick Start

```python
from dts import Task, EdgeNode, TaskScheduler, MatroidOptimizer

# Create edge nodes
nodes = [
    EdgeNode(
        node_id="edge-1",
        cpu_capacity=8.0,
        memory_capacity=16.0,
        location="datacenter-east"
    ),
    EdgeNode(
        node_id="edge-2",
        cpu_capacity=16.0,
        memory_capacity=32.0,
        location="datacenter-west"
    ),
]

# Initialize scheduler with matroid optimizer
optimizer = MatroidOptimizer(alpha=0.5, beta=0.3, gamma=0.2)
scheduler = TaskScheduler(optimizer=optimizer)

# Register nodes
for node in nodes:
    scheduler.register_node(node)

# Create and submit tasks
tasks = [
    Task(
        task_id=f"task-{i}",
        cpu_requirement=2.0,
        memory_requirement=4.0,
        execution_time=10.0,
        priority=1
    )
    for i in range(10)
]

scheduler.submit_tasks(tasks)

# Schedule tasks using matroid optimization
allocation = scheduler.schedule_tasks()

# View results
print(f"Scheduled {len(allocation)} tasks")
overview = scheduler.get_system_overview()
print(f"System Overview: {overview}")
```

## Architecture

### Core Components

#### 1. Task (`dts.core.task.Task`)

Represents a computational task with resource requirements:
- CPU requirement (cores)
- Memory requirement (GB)
- Execution time (seconds)
- Priority level

#### 2. EdgeNode (`dts.core.edge_node.EdgeNode`)

Represents an edge computing node with:
- CPU and memory capacity
- Resource availability tracking
- Load monitoring
- Health status

#### 3. MatroidOptimizer (`dts.optimization.matroid_optimizer.MatroidOptimizer`)

Implements matroid-based greedy algorithm for optimal task allocation:
- Respects resource constraints (matroid independence property)
- Optimizes for load balancing, resource efficiency, and network latency
- Computes makespan (total completion time)

#### 4. TaskScheduler (`dts.scheduler.task_scheduler.TaskScheduler`)

Main scheduling engine that:
- Manages task queues and node registry
- Coordinates task allocation using matroid optimizer
- Tracks metrics and system state
- Supports dynamic load rebalancing

### Matroid-Based Optimization

The system uses matroid theory to model resource constraints:

1. **Ground Set**: All possible task-to-node assignments
2. **Independent Sets**: Valid assignments respecting resource constraints
3. **Greedy Algorithm**: Achieves optimal or near-optimal solutions

The optimizer evaluates assignments based on:
- **Load Balancing** (α): Prefers less-loaded nodes
- **Resource Efficiency** (β): Prefers good resource fit
- **Network Latency** (γ): Prefers lower latency nodes

Weights (α, β, γ) are configurable and sum to 1.0.

## Running the Example

```bash
python examples/simple_example.py
```

This demonstrates:
1. Creating edge nodes with different capacities
2. Submitting multiple tasks with priorities
3. Scheduling using matroid optimization
4. Viewing node status and system metrics

## Testing

Run the test suite:

```bash
python -m unittest discover tests
```

Run individual test modules:

```bash
python -m unittest tests.test_task
python -m unittest tests.test_edge_node
python -m unittest tests.test_matroid_optimizer
python -m unittest tests.test_scheduler
```

## API Reference

### Task Class

```python
Task(
    task_id: str,
    cpu_requirement: float,
    memory_requirement: float,
    execution_time: float,
    priority: int = 1,
    status: TaskStatus = TaskStatus.PENDING,
    metadata: Dict[str, Any] = None
)
```

### EdgeNode Class

```python
EdgeNode(
    node_id: str,
    cpu_capacity: float,
    memory_capacity: float,
    location: str = "unknown",
    status: NodeStatus = NodeStatus.ACTIVE,
    network_latency: float = 0.0,
    metadata: Dict = None
)
```

### MatroidOptimizer Class

```python
MatroidOptimizer(
    alpha: float = 0.5,  # Load balancing weight
    beta: float = 0.3,   # Resource efficiency weight
    gamma: float = 0.2   # Network latency weight
)
```

Methods:
- `optimize_task_allocation(tasks, nodes)`: Find optimal allocation
- `compute_makespan(tasks, allocation)`: Calculate completion time
- `balance_load(nodes, threshold)`: Identify rebalancing opportunities

### TaskScheduler Class

```python
TaskScheduler(optimizer: Optional[MatroidOptimizer] = None)
```

Methods:
- `register_node(node)`: Register an edge node
- `submit_task(task)`: Submit a task
- `submit_tasks(tasks)`: Submit multiple tasks
- `schedule_tasks()`: Schedule pending tasks
- `complete_task(task_id, success)`: Mark task as completed
- `balance_load(threshold)`: Rebalance system load
- `get_node_status(node_id)`: Get node information
- `get_task_status(task_id)`: Get task information
- `get_system_overview()`: Get system statistics

## Configuration

The system supports configuration via the `Config` class:

```python
from dts.utils import Config

config = Config({
    'optimizer': {
        'alpha': 0.5,
        'beta': 0.3,
        'gamma': 0.2,
    },
    'scheduler': {
        'rebalance_threshold': 0.8,
        'heartbeat_timeout': 60.0,
    }
})
```

## Performance Considerations

- **Scheduling Complexity**: O(n log n + nm) where n=tasks, m=nodes
- **Thread Safety**: All scheduler operations are thread-safe
- **Memory Usage**: Efficient for large-scale deployments
- **Makespan**: Near-optimal due to matroid greedy algorithm

## Contributing

Contributions are welcome! Please ensure:
1. All tests pass
2. Code follows existing style
3. New features include tests
4. Documentation is updated

## License

See LICENSE file for details.

## References

- Matroid Theory: Oxley, J. G. (2011). "Matroid Theory" (2nd ed.)
- Edge Computing: Shi, W., et al. (2016). "Edge Computing: Vision and Challenges"
- Task Scheduling: Topcuoglu, H., et al. (2002). "Performance-effective and low-complexity task scheduling"
