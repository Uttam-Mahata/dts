# Getting Started with DTS

This guide will help you get started with the Distributed Task Scheduling (DTS) system.

## Prerequisites

- Go 1.21 or higher
- Unix-like operating system (Linux, macOS, or WSL on Windows)
- curl and jq (for testing API endpoints)

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/Uttam-Mahata/dts.git
cd dts
```

### 2. Install Dependencies

```bash
make install
```

### 3. Build the Server

```bash
make build
```

This will create a `dts-server` executable in the current directory.

## Running the Server

### Basic Usage

```bash
./dts-server
```

The server will start on `http://localhost:8080` by default.

### With Custom Configuration

Create a configuration file:

```bash
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
```

Run with configuration:

```bash
./dts-server -config config.json
```

## Quick Example

Once the server is running, open a new terminal and follow these steps:

### 1. Check Health

```bash
curl http://localhost:8080/health
```

### 2. Register an Edge Node

```bash
curl -X POST http://localhost:8080/api/nodes/register \
  -H "Content-Type: application/json" \
  -d '{
    "id": "node-1",
    "name": "edge-node-1",
    "location": "datacenter-east",
    "total_cpu": 8.0,
    "total_memory": 16384.0
  }'
```

### 3. Submit a Task

```bash
curl -X POST http://localhost:8080/api/tasks/submit \
  -H "Content-Type: application/json" \
  -d '{
    "name": "data-processing",
    "priority": 5,
    "cpu_required": 2.0,
    "mem_required": 1024.0,
    "duration": 300,
    "deadline": "2025-11-06T12:00:00Z"
  }'
```

### 4. Check Task Status

```bash
curl http://localhost:8080/api/tasks
```

### 5. View System Status

```bash
curl http://localhost:8080/api/system/status
```

## Running the Demo Script

We provide a demo script that demonstrates all features:

```bash
# Make the script executable
chmod +x examples/demo.sh

# Run the demo (server must be running)
./examples/demo.sh
```

## Development

### Running Tests

```bash
make test
```

### Running Tests with Coverage

```bash
make test-coverage
```

This generates a `coverage.html` file that you can open in a browser.

### Code Formatting

```bash
make fmt
```

### Building for Multiple Platforms

```bash
make build-all
```

This creates binaries for:
- Linux (amd64)
- macOS (amd64 and arm64)
- Windows (amd64)

## Understanding the System

### Key Concepts

1. **Tasks**: Computational jobs that need to be executed
   - Have resource requirements (CPU, memory)
   - Have priorities and deadlines
   - Go through lifecycle: pending → scheduled → running → completed/failed

2. **Edge Nodes**: Computing resources that execute tasks
   - Have finite CPU and memory capacity
   - Can be active or inactive
   - Send heartbeats to indicate health

3. **Scheduling**: Process of assigning tasks to nodes
   - Uses matroid-based optimization
   - Considers resource availability and load
   - Runs automatically every 10 seconds

4. **Load Balancing**: Distribution of work across nodes
   - Monitors CPU and memory utilization
   - Calculates system-wide metrics
   - Identifies imbalanced nodes

### Matroid Optimization

The system uses a matroid-based greedy algorithm for optimal task scheduling:

1. **Element Generation**: Create all valid task-node pairs
2. **Weight Calculation**: Score each pair based on:
   - Resource utilization efficiency
   - Current node load
   - Resource availability
   - Task priority
3. **Greedy Selection**: Select highest-weight pairs while maintaining independence
4. **Resource Allocation**: Assign tasks to nodes and allocate resources

This approach guarantees optimal scheduling in O(n log n) time.

## Configuration Options

### Server Configuration

- `host`: Server bind address (default: "0.0.0.0")
- `port`: Server port (default: 8080)

### Scheduler Configuration

- `schedule_interval`: How often to run scheduling (in nanoseconds)
- `max_retries`: Maximum retries for failed tasks

### System Configuration

- `heartbeat_timeout`: How long before a node is considered inactive (in nanoseconds)
- `health_check_interval`: How often to check node health (in nanoseconds)

## Troubleshooting

### Server Won't Start

- Check if port 8080 is already in use: `lsof -i :8080`
- Try a different port in the configuration

### Tasks Not Being Scheduled

- Ensure at least one node is registered and active
- Check that nodes have sufficient resources
- Verify task resource requirements are realistic

### Node Shows as Inactive

- Check network connectivity
- Verify heartbeat timeout is appropriate
- Ensure node registration was successful

## Next Steps

- Read the [API Documentation](API.md) for detailed endpoint information
- Check the [Architecture Documentation](ARCHITECTURE.md) to understand the system design
- Explore the [examples](../examples/) directory for more code samples
- Contribute improvements or report issues on GitHub

## Support

For issues and questions:
- Open an issue on GitHub
- Check existing documentation
- Review the code examples
