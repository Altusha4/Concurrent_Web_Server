package worker

import (
	"assignment2/internal/service"
	"context"
	"fmt"
	"time"
)

func StartBackgroundWorker(ctx context.Context, service *service.DataService) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fmt.Println("Background worker started")

	for {
		select {
		case <-ticker.C:
			totalRequests, dbSize := service.GetCurrentStats()
			fmt.Printf("[Worker] Status - Requests: %d, Database size: %d\n",
				totalRequests, dbSize)

		case <-ctx.Done():
			fmt.Println("Background worker stopped")
			return
		}
	}
}
