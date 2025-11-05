# System Architecture

## Overview

The Distributed Task Scheduling (DTS) system is designed to efficiently schedule and load balance computational tasks across edge computing nodes using matroid-based optimization.

## Components

### 1. Core Models (`pkg/models`)

**Task (`task.go`)**
- Represents a computational task with resource requirements
- Tracks task lifecycle: pending → scheduled → running → completed/failed
- Stores metadata, deadlines, and resource allocations

**EdgeNode (`edgenode.go`)**
- Represents a computing node in the edge network
- Manages resource allocation and availability
- Thread-safe operations for concurrent access
- Tracks running tasks and health status

**Schedule (`schedule.go`)**
- Represents task-to-node assignments
- Contains optimization metrics (cost, confidence)

### 2. Matroid-Based Optimization (`pkg/matroid`)

The system uses matroid theory for optimal task scheduling:

**Key Concepts:**
- **Elements**: Task-node pairs with associated weights
- **Independence**: Resource constraint satisfaction
- **Greedy Algorithm**: Optimal for matroid structures

**Weight Calculation:**
- Resource utilization efficiency (30%)
- Current node load (40%)
- Resource availability (30%)
- Priority boost

**Algorithm:**
1. Generate all valid task-node pairs
2. Calculate weight for each pair
3. Sort by priority and weight
4. Greedily select maintaining independence
5. Verify resource constraints

### 3. Edge Node Manager (`pkg/edgenode`)

**Responsibilities:**
- Node registration and deregistration
- Health monitoring via heartbeats
- Resource tracking
- Node status management

**Features:**
- Thread-safe node operations
- Automatic health checks
- Resource allocation/deallocation
- Active node filtering

### 4. Task Scheduler (`pkg/scheduler`)

**Responsibilities:**
- Task submission and lifecycle management
- Periodic scheduling using matroid optimization
- Task state transitions
- Resource coordination with nodes

**Scheduling Process:**
1. Collect pending tasks
2. Get active nodes
3. Apply matroid optimization
4. Allocate resources
5. Update task and node states

### 5. Load Balancer (`pkg/loadbalancer`)

**Responsibilities:**
- Monitor system load
- Select optimal nodes for tasks
- Identify load imbalances
- Provide load metrics

**Metrics:**
- Per-node load (CPU + memory utilization)
- System-wide average load
- Load balance status

### 6. REST API (`pkg/api`)

**Responsibilities:**
- HTTP endpoint handling
- Request/response serialization
- Error handling
- System monitoring endpoints

**Endpoints:**
- Task management
- Node management
- System status and monitoring
- Health checks

### 7. Configuration (`internal/config`)

**Features:**
- JSON-based configuration
- Default values
- File loading/saving
- Configurable intervals and timeouts

## System Flow

### Task Submission Flow
```
Client → API → Scheduler → Task Queue
```

### Scheduling Flow
```
Periodic Trigger → Scheduler → Matroid Optimizer → Resource Allocation → Node Assignment
```

### Health Check Flow
```
Periodic Trigger → Node Manager → Check Heartbeats → Update Status
```

## Thread Safety

All core components use appropriate synchronization:
- `sync.RWMutex` for read-heavy operations
- `sync.Mutex` for critical sections
- Thread-safe models for concurrent access

## Scalability Considerations

1. **Horizontal Scaling**: Multiple scheduler instances can be added
2. **Load Distribution**: Matroid optimization distributes tasks efficiently
3. **Resource Management**: Fine-grained resource tracking prevents overallocation
4. **Health Monitoring**: Automatic detection of failed nodes

## Optimization Strategy

The matroid-based approach provides:
- **Optimality**: Greedy algorithm is optimal for matroid structures
- **Efficiency**: O(n log n) time complexity
- **Flexibility**: Weight function can be tuned for different objectives
- **Constraints**: Naturally handles resource constraints

## Future Enhancements

1. **Task Migration**: Move tasks between nodes for better balance
2. **Priority Queues**: Multiple scheduling queues by priority
3. **Predictive Scaling**: ML-based resource prediction
4. **Fault Tolerance**: Automatic task rescheduling on node failure
5. **Multi-tenancy**: Resource isolation and quotas
6. **Authentication**: API security and access control
