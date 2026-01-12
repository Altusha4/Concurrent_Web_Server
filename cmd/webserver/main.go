package main

import (
	"assignment2/internal/handler"
	"assignment2/internal/service"
	"assignment2/internal/storage"
	"assignment2/internal/worker"
	"assignment2/pkg/frontend"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	storage := storage.NewDataStorage()
	service := service.NewDataService(storage)
	handler := handler.NewDataHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/data", handler.PostData)
	mux.HandleFunc("GET /api/data", handler.GetData)
	mux.HandleFunc("DELETE /api/data/{key}", handler.DeleteData)
	mux.HandleFunc("GET /api/stats", handler.GetStats)

	mux.HandleFunc("/", frontend.ServeFrontend)
	mux.HandleFunc("/index.html", frontend.ServeFrontend)
	mux.HandleFunc("/static/", frontend.ServeStatic)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker.StartBackgroundWorker(ctx, service)

	go func() {
		fmt.Println("Server started on http://localhost:8080")
		fmt.Println("Available endpoints:")
		fmt.Println("  Web UI:           http://localhost:8080/")
		fmt.Println("  POST   /api/data  - Store key-value pair")
		fmt.Println("  GET    /api/data  - Get all data")
		fmt.Println("  DELETE /api/data/{key} - Delete by key")
		fmt.Println("  GET    /api/stats - Get server statistics")
		fmt.Println("\nPress Ctrl+C to shutdown server")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\n Shutdown signal received")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Shutdown error: %v\n", err)
	}

	fmt.Println("Server stopped gracefully")
}
