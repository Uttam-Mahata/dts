"""
Simple example demonstrating the distributed task scheduling system.

This example shows how to:
1. Create edge nodes
2. Submit tasks
3. Schedule tasks using matroid-based optimization
4. Monitor system status
"""

import sys
import os

# Add parent directory to path for imports
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from dts import Task, EdgeNode, TaskScheduler, MatroidOptimizer


def main():
    """Run a simple scheduling example."""
    
    print("=" * 70)
    print("Distributed Task Scheduling System - Example")
    print("=" * 70)
    
    # Create edge nodes
    print("\n1. Creating edge nodes...")
    nodes = [
        EdgeNode(
            node_id="edge-1",
            cpu_capacity=8.0,
            memory_capacity=16.0,
            location="datacenter-east",
            network_latency=10.0
        ),
        EdgeNode(
            node_id="edge-2",
            cpu_capacity=16.0,
            memory_capacity=32.0,
            location="datacenter-west",
            network_latency=20.0
        ),
        EdgeNode(
            node_id="edge-3",
            cpu_capacity=4.0,
            memory_capacity=8.0,
            location="datacenter-central",
            network_latency=5.0
        ),
    ]
    
    for node in nodes:
        print(f"   - Created node {node.node_id} with {node.cpu_capacity} CPU cores "
              f"and {node.memory_capacity} GB memory")
    
    # Create task scheduler with matroid optimizer
    print("\n2. Initializing scheduler with matroid-based optimizer...")
    optimizer = MatroidOptimizer(alpha=0.5, beta=0.3, gamma=0.2)
    scheduler = TaskScheduler(optimizer=optimizer)
    
    # Register nodes
    print("   - Registering nodes with scheduler...")
    for node in nodes:
        scheduler.register_node(node)
    
    # Create tasks
    print("\n3. Creating tasks...")
    tasks = [
        Task(
            task_id=f"task-{i}",
            cpu_requirement=2.0,
            memory_requirement=4.0,
            execution_time=10.0 + i,
            priority=3 if i % 3 == 0 else 1
        )
        for i in range(1, 11)
    ]
    
    print(f"   - Created {len(tasks)} tasks with varying priorities")
    
    # Submit tasks
    print("\n4. Submitting tasks to scheduler...")
    scheduler.submit_tasks(tasks)
    print(f"   - Submitted {len(tasks)} tasks")
    
    # Schedule tasks
    print("\n5. Scheduling tasks using matroid-based optimization...")
    allocation = scheduler.schedule_tasks()
    
    print(f"   - Scheduled {len(allocation)} tasks")
    print("\n   Task Allocation:")
    for task_id, node_id in sorted(allocation.items()):
        task = scheduler.tasks[task_id]
        print(f"      {task_id} -> {node_id} (Priority: {task.priority})")
    
    # Show node status
    print("\n6. Node Status After Scheduling:")
    for node_id in sorted([n.node_id for n in nodes]):
        status = scheduler.get_node_status(node_id)
        if status:
            print(f"   - {node_id}:")
            print(f"      CPU Utilization: {status['cpu_utilization']:.2%}")
            print(f"      Memory Utilization: {status['memory_utilization']:.2%}")
            print(f"      Load Score: {status['load_score']:.2f}")
            print(f"      Assigned Tasks: {status['assigned_tasks']}")
    
    # Show system overview
    print("\n7. System Overview:")
    overview = scheduler.get_system_overview()
    print(f"   - Total Nodes: {overview['total_nodes']}")
    print(f"   - Active Nodes: {overview['active_nodes']}")
    print(f"   - Total Tasks: {overview['total_tasks']}")
    print(f"   - Pending Tasks: {overview['pending_tasks']}")
    print(f"   - Average Load: {overview['average_load']:.2f}")
    
    # Calculate makespan
    print("\n8. Performance Metrics:")
    makespan = optimizer.compute_makespan(tasks, allocation)
    print(f"   - Makespan: {makespan:.2f} seconds")
    print(f"   - Tasks Scheduled: {scheduler.metrics['total_tasks_scheduled']}")
    
    print("\n" + "=" * 70)
    print("Example completed successfully!")
    print("=" * 70)


if __name__ == "__main__":
    main()
