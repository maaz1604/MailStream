package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func loadRecepient(filepath string, ch chan Recepient) error {

	defer close(ch)
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil
	}

	for _, record := range records[1:] {
		fmt.Println(record)
		ch <- Recepient{
			Name:  record[0],
			Email: record[1],
		}
	}
	// send -> consumer (to implement queue like str) -> channels

	return nil
}
