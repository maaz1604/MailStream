package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"
)

type DeadLetter struct {
	Recipient Recepient
	Reason    string
	Error     string
}

func deadLetterWorker(queue <-chan DeadLetter, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.OpenFile("dead-letter.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Unable to open dead-letter.csv: %v\n", err)
		for range queue {
		}
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if fileInfo, statErr := file.Stat(); statErr == nil && fileInfo.Size() == 0 {
		if err := writer.Write([]string{"timestamp", "name", "email", "reason", "error"}); err != nil {
			fmt.Printf("Unable to write dead-letter.csv header: %v\n", err)
			for range queue {
			}
			return
		}
		writer.Flush()
	}

	for deadLetter := range queue {
		if err := writer.Write([]string{
			time.Now().UTC().Format(time.RFC3339),
			deadLetter.Recipient.Name,
			deadLetter.Recipient.Email,
			deadLetter.Reason,
			deadLetter.Error,
		}); err != nil {
			fmt.Printf("Unable to write dead letter for %s: %v\n", deadLetter.Recipient.Email, err)
			continue
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			fmt.Printf("Unable to flush dead letter for %s: %v\n", deadLetter.Recipient.Email, err)
		}
	}
}
