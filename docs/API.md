# API Documentation

## Overview

The Distributed Task Scheduling (DTS) system provides a REST API for task submission, node management, and system monitoring.

## Base URL

```
http://localhost:8080/api
```

## Endpoints

### Task Management

#### Submit Task

```http
POST /api/tasks/submit
```

**Request Body:**
```json
{
  "name": "example-task",
  "priority": 5,
  "cpu_required": 2.0,
  "mem_required": 1024.0,
  "duration": 300,
  "deadline": "2025-11-06T12:00:00Z",
  "metadata": {
    "user": "admin",
    "department": "engineering"
  }
}
```

**Response:**
```json
{
  "success": true,
  "task": {
    "id": "task-1730834400123456789",
    "name": "example-task",
    "priority": 5,
    "cpu_required": 2.0,
    "mem_required": 1024.0,
    "duration": 300000000000,
    "deadline": "2025-11-06T12:00:00Z",
    "status": "pending",
    "created_at": "2025-11-05T19:00:00Z"
  },
  "message": "Task submitted successfully"
}
```

#### Get All Tasks

```http
GET /api/tasks
```

**Response:**
```json
{
  "tasks": [
    {
      "id": "task-1",
      "name": "example-task",
      "status": "pending",
      ...
    }
  ],
  "count": 1
}
```

#### Trigger Scheduling

```http
POST /api/tasks/schedule
```

**Response:**
```json
{
  "success": true,
  "schedules": [
    {
      "task_id": "task-1",
      "node_id": "node-1",
      "priority": 5,
      "cost": 0.85,
      "confidence": 0.95
    }
  ],
  "message": "Successfully scheduled 1 tasks"
}
```

### Node Management

#### Register Node

```http
POST /api/nodes/register
```

**Request Body:**
```json
{
  "id": "node-1",
  "name": "edge-node-1",
  "location": "datacenter-east",
  "total_cpu": 8.0,
  "total_memory": 16384.0
}
```

**Response:**
```json
{
  "success": true,
  "node": {
    "id": "node-1",
    "name": "edge-node-1",
    "status": "active",
    "total_cpu": 8.0,
    "total_memory": 16384.0,
    "available_cpu": 8.0,
    "available_memory": 16384.0
  },
  "message": "Node registered successfully"
}
```

#### Get All Nodes

```http
GET /api/nodes
```

**Response:**
```json
{
  "nodes": [
    {
      "id": "node-1",
      "name": "edge-node-1",
      "status": "active",
      ...
    }
  ],
  "count": 1
}
```

### System Monitoring

#### Get System Status

```http
GET /api/system/status
```

**Response:**
```json
{
  "total_nodes": 3,
  "active_nodes": 3,
  "total_tasks": 10,
  "pending_tasks": 2,
  "scheduled_tasks": 3,
  "running_tasks": 4,
  "completed_tasks": 1,
  "failed_tasks": 0,
  "system_load": 0.65,
  "is_balanced": true,
  "timestamp": "2025-11-05T19:10:00Z"
}
```

#### Get System Load

```http
GET /api/system/load
```

**Response:**
```json
{
  "node_loads": {
    "node-1": 0.60,
    "node-2": 0.70,
    "node-3": 0.65
  },
  "system_load": 0.65,
  "timestamp": "2025-11-05T19:10:00Z"
}
```

### Health Check

```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "time": "2025-11-05T19:10:00Z"
}
```

## Error Responses

All endpoints may return error responses in the following format:

**4xx Client Errors:**
```json
{
  "error": "Invalid request body: ...",
  "status": 400
}
```

**5xx Server Errors:**
```json
{
  "error": "Internal server error: ...",
  "status": 500
}
```

## Rate Limiting

Currently, there are no rate limits enforced. This may be added in future versions.

## Authentication

Currently, the API does not require authentication. This should be added before production deployment.
