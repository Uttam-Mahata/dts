#!/bin/bash

# DTS Demo Script
# This script demonstrates the basic functionality of the DTS system

set -e

BASE_URL="http://localhost:8080"

echo "=== DTS System Demo ==="
echo ""

# Check if server is running
echo "1. Checking server health..."
curl -s "${BASE_URL}/health" | jq .
echo ""

# Register edge nodes
echo "2. Registering edge nodes..."
curl -s -X POST "${BASE_URL}/api/nodes/register" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "node-1",
    "name": "edge-node-east-1",
    "location": "datacenter-east",
    "total_cpu": 8.0,
    "total_memory": 16384.0
  }' | jq .

curl -s -X POST "${BASE_URL}/api/nodes/register" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "node-2",
    "name": "edge-node-west-1",
    "location": "datacenter-west",
    "total_cpu": 16.0,
    "total_memory": 32768.0
  }' | jq .

curl -s -X POST "${BASE_URL}/api/nodes/register" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "node-3",
    "name": "edge-node-central-1",
    "location": "datacenter-central",
    "total_cpu": 12.0,
    "total_memory": 24576.0
  }' | jq .
echo ""

# List nodes
echo "3. Listing all registered nodes..."
curl -s "${BASE_URL}/api/nodes" | jq .
echo ""

# Submit tasks
echo "4. Submitting tasks..."
for i in {1..5}; do
  curl -s -X POST "${BASE_URL}/api/tasks/submit" \
    -H "Content-Type: application/json" \
    -d "{
      \"name\": \"task-${i}\",
      \"priority\": $((RANDOM % 10 + 1)),
      \"cpu_required\": $((RANDOM % 4 + 1)).0,
      \"mem_required\": $((RANDOM % 2048 + 512)).0,
      \"duration\": $((RANDOM % 300 + 60)),
      \"deadline\": \"2025-11-06T12:00:00Z\",
      \"metadata\": {
        \"batch\": \"demo\",
        \"task_num\": ${i}
      }
    }" | jq .
done
echo ""

# List tasks
echo "5. Listing all tasks..."
curl -s "${BASE_URL}/api/tasks" | jq .
echo ""

# Trigger scheduling
echo "6. Triggering task scheduling..."
curl -s -X POST "${BASE_URL}/api/tasks/schedule" | jq .
echo ""

# Check system status
echo "7. Checking system status..."
curl -s "${BASE_URL}/api/system/status" | jq .
echo ""

# Check system load
echo "8. Checking system load..."
curl -s "${BASE_URL}/api/system/load" | jq .
echo ""

echo "=== Demo Complete ==="
