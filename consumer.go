package main

import (
	"fmt"
	"log"
	"net/smtp"
	"sync"
	"time"
)

func emailWorker(id int, ch <-chan Recepient, deadLetterQueue chan<- DeadLetter, wg *sync.WaitGroup) {
	defer wg.Done()
	for recipient := range ch {
		smtpHost := "localhost"
		smtpPort := "1025"

		// formattedMsg := fmt.Sprintf("To: %s\r\nSubject: Test Email\r\n\r\n%s\r\n", recipient.Email, "Just testing our email campaign.")
		// msg := []byte(formattedMsg)

		msg, err := executeTemplate(recipient)
		if err != nil {
			log.Printf("Worker %d: error parsing template for %s: %v", id, recipient.Email, err)
			deadLetterQueue <- DeadLetter{
				Recipient: recipient,
				Reason:    "template",
				Error:     err.Error(),
			}
			continue
		}

		fmt.Printf("Worker %d: Sending email to %s \n", id, recipient.Email)

		err2 := smtp.SendMail(smtpHost+":"+smtpPort, nil, "rakesh@codersgyan.com", []string{recipient.Email}, []byte(msg))
		if err2 != nil {
			log.Printf("Worker %d: failed to send email to %s: %v", id, recipient.Email, err2)
			deadLetterQueue <- DeadLetter{
				Recipient: recipient,
				Reason:    "smtp",
				Error:     err2.Error(),
			}
			continue
		}

		time.Sleep(50 * time.Millisecond)

		fmt.Printf("Worker %d: Sent email to %s \n", id, recipient.Email)
	}
}
