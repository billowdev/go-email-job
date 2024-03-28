package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var EMAIL_LIST = []string{
	"a1@gmail.com",
	"a2@gmail.com",
	"a3@gmail.com",
	"a4@gmail.com",
	"a5@gmail.com",
	"a6@gmail.com",
	"a7@gmail.com",
	"a8@gmail.com",
	"a9@gmail.com",
}

type Email struct {
	gorm.Model
	From      string
	To        string
	Timestamp time.Time
	Success   bool
}

func TaskEmail(db *gorm.DB, ctx context.Context, from string, to []string, wg *sync.WaitGroup) {
	defer wg.Done()

	successCh := make(chan bool)
	errorCh := make(chan error)

	go func() {
		for {
			select {
			// Check if the context is canceled
			case <-ctx.Done():
				fmt.Println("\n---\nsent")
				return // Exit the goroutine
			default:
				// Simulate sending email
				for _, e := range to {
					err := sendEmail(db, from, e)
					if err != nil {
						errorCh <- err // Send error to error channel
					} else {
						successCh <- true // Send success to success channel
					}
					// time.Sleep(500 * time.Millisecond) // Simulate some work
					randomSleep := time.Duration(rand.Intn(901)+100) * time.Millisecond
					time.Sleep(randomSleep)
				}
			}
		}
	}()

	for {
		select {
		case <-successCh:
			fmt.Println("Email sent successfully!")
		case err := <-errorCh:
			fmt.Printf("Error sending email: %v\n", err)
		case <-ctx.Done():
			// Context is canceled, exit the loop
			fmt.Printf("Error sending email: %v\n", "timeout")
			return
		}
	}
}
func initializeDB() *gorm.DB {
	// Read configuration from config.yaml
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// Read database configuration
	dbConfig := viper.GetStringMapString("database")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbConfig["host"],
		dbConfig["user"],
		dbConfig["password"],
		dbConfig["dbname"],
		dbConfig["port"],
		dbConfig["sslmode"],
		dbConfig["timezone"],
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto Migrate the Email model
	db.AutoMigrate(&Email{})

	return db
}

func main() {
	db := initializeDB()
	// Define email addresses
	from := "abc@gmail.com"
	to := EMAIL_LIST

	var wg sync.WaitGroup

	// Create a new context with timeout
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wg.Add(1)
	go TaskEmail(db, timeoutCtx, from, to, &wg)

	// Wait for all goroutines to finish
	wg.Wait()
}

// Dummy function to simulate sending an email
func sendEmail(db *gorm.DB, from, to string) error {
	randomIndex := rand.Intn(len(EMAIL_LIST))
	randomTo := EMAIL_LIST[randomIndex]
	if to == randomTo {
		// Simulate sending email
		time.Sleep(100 * time.Millisecond)

		// Save email to database
		db.Create(&Email{From: from, To: to, Timestamp: time.Now(), Success: false})

		return errors.New("failed to send email") // Return error to simulate failure
	} else {
		// Simulate sending email
		time.Sleep(100 * time.Millisecond)
		db.Create(&Email{From: from, To: to, Timestamp: time.Now(), Success: true})
		return nil // Return nil to indicate success
	}
}
