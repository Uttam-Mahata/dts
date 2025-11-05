# Distributed Task Scheduling System (DTS)

A high-performance distributed task scheduling and load balancing system for edge computing environments, powered by matroid-based optimization algorithms.

## Overview

DTS is a Go-based system that efficiently schedules computational tasks across edge computing nodes using advanced matroid theory optimization. The system provides:

- **Optimal Scheduling**: Matroid-based greedy algorithm for provably optimal task assignments
- **Load Balancing**: Intelligent distribution of tasks across edge nodes
- **Resource Management**: Fine-grained CPU and memory allocation tracking
- **REST API**: Easy-to-use HTTP API for task submission and monitoring
- **Real-time Monitoring**: System status and load metrics
- **Health Checks**: Automatic node health monitoring with heartbeats

## Features

- ✅ Matroid-based optimization for optimal task scheduling
- ✅ Multi-priority task queuing
- ✅ Dynamic resource allocation and tracking
- ✅ Load balancing across edge nodes
- ✅ RESTful API for integration
- ✅ Real-time system monitoring
- ✅ Automatic health checks and failover
- ✅ Thread-safe concurrent operations
- ✅ Comprehensive test coverage

## Architecture

The system consists of several key components:

- **Task Scheduler**: Manages task lifecycle and scheduling using matroid optimization
- **Edge Node Manager**: Tracks and manages edge computing nodes
- **Load Balancer**: Monitors and balances load across nodes
- **Matroid Optimizer**: Implements greedy algorithm for optimal task-node assignments
- **REST API Server**: Provides HTTP endpoints for external integration

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed architecture documentation.

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Unix-like operating system (Linux, macOS)

### Installation

```bash
# Clone the repository
git clone https://github.com/Uttam-Mahata/dts.git
cd dts

# Build the server
go build -o dts-server ./cmd/server

# Run the server
./dts-server
```

The server will start on `http://localhost:8080` by default.

### Using with Configuration File

```bash
# Create a config file (optional)
cat > config.json <<EOF
{
  "server": {
    "host": "0.0.0.0",
    "port": 8080
  },
  "scheduler": {
    "schedule_interval": 10000000000,
    "max_retries": 3
  },
  "system": {
    "heartbeat_timeout": 60000000000,
    "health_check_interval": 30000000000
  }
}
EOF

# Run with config
./dts-server -config config.json
```

## Usage Examples

### 1. Register an Edge Node

```bash
curl -X POST http://localhost:8080/api/nodes/register \
  -H "Content-Type: application/json" \
  -d '{
    "id": "node-1",
    "name": "edge-node-east-1",
    "location": "datacenter-east",
    "total_cpu": 8.0,
    "total_memory": 16384.0
  }'
```

### 2. Submit a Task

```bash
curl -X POST http://localhost:8080/api/tasks/submit \
  -H "Content-Type: application/json" \
  -d '{
    "name": "data-processing-job",
    "priority": 5,
    "cpu_required": 2.0,
    "mem_required": 1024.0,
    "duration": 300,
    "deadline": "2025-11-06T12:00:00Z",
    "metadata": {
      "user": "admin",
      "type": "batch"
    }
  }'
```

### 3. Trigger Scheduling

```bash
curl -X POST http://localhost:8080/api/tasks/schedule
```

### 4. Check System Status

```bash
curl http://localhost:8080/api/system/status
```

### 5. Monitor System Load

```bash
curl http://localhost:8080/api/system/load
```

## API Documentation

Complete API documentation is available in [docs/API.md](docs/API.md).

### Key Endpoints

- `POST /api/tasks/submit` - Submit a new task
- `GET /api/tasks` - List all tasks
- `POST /api/tasks/schedule` - Trigger scheduling
- `POST /api/nodes/register` - Register an edge node
- `GET /api/nodes` - List all nodes
- `GET /api/system/status` - Get system status
- `GET /api/system/load` - Get load metrics
- `GET /health` - Health check

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## Development

### Project Structure

```
dts/
├── cmd/
│   └── server/           # Main server application
├── pkg/
│   ├── api/              # REST API handlers
│   ├── edgenode/         # Edge node management
│   ├── loadbalancer/     # Load balancing logic
│   ├── matroid/          # Matroid-based optimization
│   ├── models/           # Data models
│   └── scheduler/        # Task scheduling
├── internal/
│   └── config/           # Configuration management
├── docs/                 # Documentation
└── README.md
```

### Building

```bash
# Build for current platform
go build -o dts-server ./cmd/server

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o dts-server-linux ./cmd/server

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o dts-server-macos ./cmd/server
```

## How It Works

### Matroid-Based Optimization

The system uses matroid theory to achieve optimal task scheduling:

1. **Element Generation**: Creates task-node pairs as matroid elements
2. **Weight Calculation**: Assigns weights based on:
   - Resource utilization efficiency
   - Current node load
   - Resource availability
   - Task priority
3. **Greedy Selection**: Selects highest-weight elements maintaining independence
4. **Independence Check**: Ensures resource constraints are satisfied

This approach guarantees optimal scheduling in O(n log n) time complexity.

### Load Balancing

The load balancer:
- Monitors CPU and memory utilization across nodes
- Calculates system-wide load metrics
- Identifies overloaded and underloaded nodes
- Provides load information to the scheduler

### Resource Management

Resources are tracked at multiple levels:
- Per-task requirements (CPU, memory)
- Per-node availability (total and available)
- System-wide capacity and utilization

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Matroid theory optimization algorithms
- Go standard library and ecosystem
- Edge computing research community
