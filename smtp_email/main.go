package main

import (
	"context"
	"sync"
	"time"
)

func main() {
	db := initializeDB()

	// Define email addresses
	from := SMTP_CONFIG.Sender
	to := EMAIL_LIST

	var wg sync.WaitGroup

	// Create a new context with timeout
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wg.Add(1)
	go TaskEmail(db, timeoutCtx, from, to[0], &wg)
	// Wait for all goroutines to finish
	wg.Wait()
}
