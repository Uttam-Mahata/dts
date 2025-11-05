package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Example client demonstrating DTS API usage

const baseURL = "http://localhost:8080"

type TaskRequest struct {
	Name        string                 `json:"name"`
	Priority    int                    `json:"priority"`
	CPURequired float64                `json:"cpu_required"`
	MemRequired float64                `json:"mem_required"`
	Duration    int64                  `json:"duration"`
	Deadline    string                 `json:"deadline"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type NodeRegistration struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Location    string  `json:"location"`
	TotalCPU    float64 `json:"total_cpu"`
	TotalMemory float64 `json:"total_memory"`
}

func main() {
	// Register a node
	node := NodeRegistration{
		ID:          "example-node-1",
		Name:        "Example Edge Node",
		Location:    "example-location",
		TotalCPU:    8.0,
		TotalMemory: 16384.0,
	}

	if err := registerNode(node); err != nil {
		fmt.Printf("Failed to register node: %v\n", err)
		return
	}
	fmt.Println("✓ Node registered successfully")

	// Submit a task
	deadline := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	task := TaskRequest{
		Name:        "example-task",
		Priority:    5,
		CPURequired: 2.0,
		MemRequired: 1024.0,
		Duration:    300,
		Deadline:    deadline,
		Metadata: map[string]interface{}{
			"user": "example-user",
			"type": "demo",
		},
	}

	if err := submitTask(task); err != nil {
		fmt.Printf("Failed to submit task: %v\n", err)
		return
	}
	fmt.Println("✓ Task submitted successfully")

	// Trigger scheduling
	if err := triggerScheduling(); err != nil {
		fmt.Printf("Failed to trigger scheduling: %v\n", err)
		return
	}
	fmt.Println("✓ Scheduling triggered successfully")

	// Get system status
	if err := getSystemStatus(); err != nil {
		fmt.Printf("Failed to get system status: %v\n", err)
		return
	}
	fmt.Println("✓ System status retrieved successfully")
}

func registerNode(node NodeRegistration) error {
	data, err := json.Marshal(node)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		baseURL+"/api/nodes/register",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	return nil
}

func submitTask(task TaskRequest) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		baseURL+"/api/tasks/submit",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	return nil
}

func triggerScheduling() error {
	resp, err := http.Post(baseURL+"/api/tasks/schedule", "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	return nil
}

func getSystemStatus() error {
	resp, err := http.Get(baseURL + "/api/system/status")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	var status map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return err
	}

	fmt.Println("\nSystem Status:")
	data, _ := json.MarshalIndent(status, "", "  ")
	fmt.Println(string(data))

	return nil
}
