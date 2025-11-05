package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Uttam-Mahata/dts/internal/config"
	"github.com/Uttam-Mahata/dts/pkg/api"
	"github.com/Uttam-Mahata/dts/pkg/edgenode"
	"github.com/Uttam-Mahata/dts/pkg/loadbalancer"
	"github.com/Uttam-Mahata/dts/pkg/scheduler"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	// Load configuration
	var cfg *config.Config
	var err error

	if *configPath != "" {
		cfg, err = config.LoadFromFile(*configPath)
		if err != nil {
			log.Printf("Failed to load config from file: %v, using defaults", err)
			cfg = config.Default()
		}
	} else {
		cfg = config.Default()
	}

	log.Printf("Starting Distributed Task Scheduling System")
	log.Printf("Server: %s:%d", cfg.Server.Host, cfg.Server.Port)

	// Initialize components
	nodeManager := edgenode.NewManager()
	sched := scheduler.NewScheduler(nodeManager)
	lb := loadbalancer.NewLoadBalancer(nodeManager)

	// Create API server
	apiServer := api.NewServer(sched, nodeManager, lb)
	mux := apiServer.SetupRoutes()

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start background tasks
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start periodic scheduling
	go func() {
		ticker := time.NewTicker(cfg.Scheduler.ScheduleInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				result, err := sched.Schedule()
				if err != nil {
					log.Printf("Scheduling error: %v", err)
				} else if result.Success && len(result.Schedules) > 0 {
					log.Printf("Scheduled %d tasks", len(result.Schedules))
				}
			}
		}
	}()

	// Start periodic health check
	go func() {
		ticker := time.NewTicker(cfg.System.HealthCheckInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				nodeManager.CheckNodeHealth(cfg.System.HeartbeatTimeout)
			}
		}
	}()

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on %s", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Cancel background tasks
	cancel()

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
