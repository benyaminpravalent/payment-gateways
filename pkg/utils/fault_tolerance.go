package utils

import (
	"fmt"
	"log"
	"time"
)

// Retry operation with exponential backoff
func RetryOperation(operation func() error, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		if err := operation(); err == nil {
			return nil
		} else {
			log.Printf("Retry %d/%d failed: %v", i+1, maxRetries, err)
		}
		time.Sleep(time.Duration(1<<i) * time.Second) // Exponential backoff
	}
	return fmt.Errorf("operation failed after %d attempts", maxRetries)
}
