"""
Advanced example demonstrating load balancing and monitoring in the distributed task scheduling system.

This example shows:
1. Creating nodes with different capacities and latencies
2. Submitting tasks with priorities
3. Monitoring system metrics
4. Dynamic load balancing
5. Task lifecycle management
"""

import sys
import os
import time

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from dts import Task, EdgeNode, TaskScheduler, MatroidOptimizer
from dts.core.task import TaskStatus
from dts.core.edge_node import NodeStatus


def print_section(title):
    """Print a formatted section header."""
    print(f"\n{'='*70}")
    print(f"  {title}")
    print('='*70)


def main():
    """Run advanced scheduling example."""
    
    print_section("Advanced Distributed Task Scheduling System")
    
    # 1. Create nodes with varying capabilities
    print("\n1. Creating heterogeneous edge nodes...")
    nodes = [
        EdgeNode(
            node_id="edge-powerful",
            cpu_capacity=32.0,
            memory_capacity=64.0,
            location="datacenter-central",
            network_latency=5.0
        ),
        EdgeNode(
            node_id="edge-medium-1",
            cpu_capacity=16.0,
            memory_capacity=32.0,
            location="datacenter-east",
            network_latency=15.0
        ),
        EdgeNode(
            node_id="edge-medium-2",
            cpu_capacity=16.0,
            memory_capacity=32.0,
            location="datacenter-west",
            network_latency=20.0
        ),
        EdgeNode(
            node_id="edge-small",
            cpu_capacity=8.0,
            memory_capacity=16.0,
            location="datacenter-remote",
            network_latency=50.0
        ),
    ]
    
    for node in nodes:
        print(f"   ✓ {node.node_id}: {node.cpu_capacity} cores, "
              f"{node.memory_capacity} GB RAM, {node.network_latency}ms latency")
    
    # 2. Initialize scheduler with custom weights
    print("\n2. Configuring matroid optimizer...")
    print("   - Load balancing weight (α): 0.6")
    print("   - Resource efficiency (β): 0.2")
    print("   - Network latency (γ): 0.2")
    
    optimizer = MatroidOptimizer(alpha=0.6, beta=0.2, gamma=0.2)
    scheduler = TaskScheduler(optimizer=optimizer)
    
    for node in nodes:
        scheduler.register_node(node)
    print(f"   ✓ Registered {len(nodes)} nodes")
    
    # 3. Create diverse task workload
    print("\n3. Creating diverse task workload...")
    tasks = []
    
    # High priority critical tasks
    for i in range(3):
        tasks.append(Task(
            task_id=f"critical-{i}",
            cpu_requirement=4.0,
            memory_requirement=8.0,
            execution_time=20.0,
            priority=10,
            metadata={'type': 'critical'}
        ))
    
    # Medium priority compute tasks
    for i in range(5):
        tasks.append(Task(
            task_id=f"compute-{i}",
            cpu_requirement=2.0,
            memory_requirement=4.0,
            execution_time=15.0,
            priority=5,
            metadata={'type': 'compute'}
        ))
    
    # Low priority batch tasks
    for i in range(7):
        tasks.append(Task(
            task_id=f"batch-{i}",
            cpu_requirement=1.0,
            memory_requirement=2.0,
            execution_time=10.0,
            priority=1,
            metadata={'type': 'batch'}
        ))
    
    print(f"   ✓ Created {len(tasks)} tasks:")
    print(f"      - 3 critical tasks (priority 10)")
    print(f"      - 5 compute tasks (priority 5)")
    print(f"      - 7 batch tasks (priority 1)")
    
    # 4. Submit and schedule tasks
    print("\n4. Submitting and scheduling tasks...")
    scheduler.submit_tasks(tasks)
    allocation = scheduler.schedule_tasks()
    
    print(f"   ✓ Scheduled {len(allocation)}/{len(tasks)} tasks")
    
    # 5. Display allocation by priority
    print("\n5. Task Allocation by Priority:")
    priority_groups = {}
    for task_id, node_id in allocation.items():
        task = scheduler.tasks[task_id]
        priority = task.priority
        if priority not in priority_groups:
            priority_groups[priority] = []
        priority_groups[priority].append((task_id, node_id))
    
    for priority in sorted(priority_groups.keys(), reverse=True):
        print(f"\n   Priority {priority}:")
        for task_id, node_id in priority_groups[priority]:
            print(f"      {task_id} → {node_id}")
    
    # 6. Show detailed node status
    print("\n6. Detailed Node Status:")
    for node_id in sorted([n.node_id for n in nodes]):
        status = scheduler.get_node_status(node_id)
        if status:
            print(f"\n   {node_id}:")
            print(f"      Status: {status['status']}")
            print(f"      CPU Utilization: {status['cpu_utilization']:.1%}")
            print(f"      Memory Utilization: {status['memory_utilization']:.1%}")
            print(f"      Load Score: {status['load_score']:.2f}")
            print(f"      Assigned Tasks: {status['assigned_tasks']}")
            print(f"      Healthy: {'✓' if status['is_healthy'] else '✗'}")
    
    # 7. Simulate task execution
    print("\n7. Simulating task execution...")
    executed = 0
    for task_id in list(allocation.keys())[:5]:
        scheduler.start_task(task_id)
        scheduler.complete_task(task_id, success=True)
        executed += 1
    print(f"   ✓ Completed {executed} tasks")
    
    # 8. Check for load imbalance and rebalance
    print("\n8. Load Balancing Analysis:")
    overview = scheduler.get_system_overview()
    print(f"   Average system load: {overview['average_load']:.2f}")
    
    if overview['average_load'] > 0.6:
        print("   ⚠ System load is high, checking for rebalancing opportunities...")
        migrated = scheduler.balance_load(threshold=0.7)
        if migrated:
            print(f"   ✓ Migrated {len(migrated)} tasks for better balance")
        else:
            print("   ℹ No rebalancing needed")
    else:
        print("   ✓ System load is balanced")
    
    # 9. Performance metrics
    print("\n9. Performance Metrics:")
    makespan = optimizer.compute_makespan(
        [scheduler.tasks[tid] for tid in allocation.keys()],
        allocation
    )
    metrics = scheduler.get_metrics()
    
    print(f"   - Makespan: {makespan:.1f} seconds")
    print(f"   - Total Scheduled: {metrics['total_tasks_scheduled']}")
    print(f"   - Total Completed: {metrics['total_tasks_completed']}")
    print(f"   - Total Failed: {metrics['total_tasks_failed']}")
    print(f"   - Success Rate: {metrics['total_tasks_completed']/(metrics['total_tasks_completed']+metrics['total_tasks_failed'])*100 if metrics['total_tasks_completed']+metrics['total_tasks_failed'] > 0 else 0:.1f}%")
    
    # 10. System summary
    print("\n10. Final System Overview:")
    overview = scheduler.get_system_overview()
    print(f"   - Total Nodes: {overview['total_nodes']}")
    print(f"   - Active Nodes: {overview['active_nodes']}")
    print(f"   - Total Tasks: {overview['total_tasks']}")
    print(f"   - Running Tasks: {overview['running_tasks']}")
    print(f"   - Pending Tasks: {overview['pending_tasks']}")
    print(f"   - Average Load: {overview['average_load']:.2f}")
    
    # Distribution of tasks across nodes
    task_distribution = {}
    for node_id in [n.node_id for n in nodes]:
        node = scheduler.nodes[node_id]
        task_distribution[node_id] = len(node.assigned_tasks)
    
    print("\n   Task Distribution:")
    for node_id, count in sorted(task_distribution.items()):
        bar = '█' * count
        print(f"      {node_id:20s}: {bar} ({count})")
    
    print_section("Advanced Example Completed Successfully!")


if __name__ == "__main__":
    main()
