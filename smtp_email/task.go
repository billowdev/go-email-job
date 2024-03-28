package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

// Dummy function to simulate sending an email
// func sendEmailWithTemplateHelper(db *gorm.DB, from, to string) error {
// 	randomIndex := rand.Intn(len(EMAIL_LIST))
// 	randomTo := EMAIL_LIST[randomIndex]
// 	if to == randomTo {
// 		// Simulate sending email
// 		time.Sleep(100 * time.Millisecond)

// 		// Save email to database
// 		db.Create(&Email{From: from, To: to, Timestamp: time.Now(), Success: false})

// 		return errors.New("failed to send email") // Return error to simulate failure
// 	} else {
// 		// Simulate sending email
// 		time.Sleep(100 * time.Millisecond)
// 		db.Create(&Email{From: from, To: to, Timestamp: time.Now(), Success: true})
// 		return nil // Return nil to indicate success
// 	}
// }

func TaskEmail(db *gorm.DB, ctx context.Context, from string, to string, wg *sync.WaitGroup) {
	defer wg.Done()

	successCh := make(chan bool)
	errorCh := make(chan error)
	emailLog := EmailLog{From: from, To: to, Timestamp: time.Now(), Success: true}
	err := db.Create(&emailLog).Error
	if err != nil {
		return
	}

	go func() {
		subject := "Email Testing"
		htmlTemplate := HTML_TEMPLATE
		type SHTMLArgs struct {
			CustomerName string `json:"customer_name"`
		}
		htmlArgs := SHTMLArgs{
			CustomerName: "mr testing",
		}
		err := SendEmailWithTemplate(to, subject, htmlTemplate, htmlArgs)
		if err != nil {
			errorCh <- err // Send error to error channel
		} else {
			successCh <- true // Send success to success channel
		}
	}()

	select {
	case <-successCh:
		emailLog.Success = true
		db.Save(&emailLog)
		fmt.Println("Email sent successfully!")
	case err := <-errorCh:
		emailLog.Error = err.Error()
		emailLog.Success = false
		db.Save(&emailLog)
		fmt.Printf("Error sending email: %v\n", err)
	case <-ctx.Done():
		// Context is canceled, exit the block
		emailLog.Error = errors.New("timeout").Error()
		emailLog.Success = false
		db.Save(&emailLog)
		fmt.Printf("Error sending email: %v\n", "timeout")
	}

	// go func() {
	// 	for {
	// 		select {
	// 		// Check if the context is canceled
	// 		case <-ctx.Done():
	// 			fmt.Println("\n---\nsent")
	// 			return // Exit the goroutine
	// 		default:
	// 			// Simulate sending email
	// 			// err := sendEmailWithTemplateHelper(db, from, to)
	// 			subject := "Email Testing"
	// 			htmlTemplate := HTML_TEMPLATE
	// 			type SHTMLArgs struct {
	// 				CustomerName string `json:"customer_name"`
	// 			}
	// 			htmlArgs := SHTMLArgs{
	// 				CustomerName: "mr testing",
	// 			}
	// 			err := SendEmailWithTemplate(to, subject, htmlTemplate, htmlArgs)
	// 			if err != nil {
	// 				errorCh <- err // Send error to error channel
	// 			} else {
	// 				successCh <- true // Send success to success channel
	// 			}
	// 		}
	// 	}
	// }()

	// for {
	// 	select {
	// 	case <-successCh:
	// 		emailLog.Success = true
	// 		db.Save(&emailLog)
	// 		fmt.Println("Email sent successfully!")
	// 	case err := <-errorCh:
	// 		emailLog.Error = err.Error()
	// 		emailLog.Success = false
	// 		db.Save(&emailLog)
	// 		fmt.Printf("Error sending email: %v\n", err)
	// 	case <-ctx.Done():
	// 		// Context is canceled, exit the loop
	// 		emailLog.Error = errors.New("timeout").Error()
	// 		emailLog.Success = false
	// 		db.Save(&emailLog)
	// 		fmt.Printf("Error sending email: %v\n", "timeout")
	// 		return
	// 	}
	// }
}
