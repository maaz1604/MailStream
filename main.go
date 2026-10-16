package main

import "sync"

type Recepient struct {
	Name  string
	Email string
}

func main() {
	recepientChannel := make(chan Recepient)
	go func() {
		loadRecepient("./emails.csv", recepientChannel)
	}()

	var wg sync.WaitGroup
	workerCount := 5

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go emailWorker(i, recepientChannel, &wg)
	}

	wg.Wait()
}
