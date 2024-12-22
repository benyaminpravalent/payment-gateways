package cmd

import (
	"context"
	"log"
	"payment-gateway/pkg/constants"

	"github.com/robfig/cron/v3"
)

func InitCron() *cron.Cron {
	c := cron.New()
	ctx := context.Background()

	// Schedule a job to check gateway health status every 30 seconds
	_, err := c.AddFunc("@every 30s", func() {
		// log.Println("Cron Job: Checking gateway health status...")

		// 1. fetch all active gatewayID
		gatewayIDs := []int{1, 2, 3}

		// 2. check if external gateway is reachable and ready to accept traffic (e.g via API)

		// 3. update the status accordingly (HEALTHY/UNHEALTHY)
		for _, v := range gatewayIDs {
			GatewayService.UpdateGatewayHealthStatusByID(ctx, v, constants.HEALTHY)
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule cron job: %v", err)
	}

	c.Start()

	log.Println("Cron scheduler initialized successfully (every 30 seconds)")

	return c
}
