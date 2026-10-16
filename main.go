package main

import (
	"bytes"
	"html/template"
	"sync"
)

type Recepient struct {
	Name  string
	Email string
}

func main() {
	recepientChannel := make(chan Recepient)
	deadLetterQueue := make(chan DeadLetter)
	go func() {
		loadRecepient("./emails.csv", recepientChannel)
	}()

	var wg sync.WaitGroup
	var deadLetterWg sync.WaitGroup
	workerCount := 5

	deadLetterWg.Add(1)
	go deadLetterWorker(deadLetterQueue, &deadLetterWg)

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go emailWorker(i, recepientChannel, deadLetterQueue, &wg)
	}

	wg.Wait()
	close(deadLetterQueue)
	deadLetterWg.Wait()
}

func executeTemplate(r Recepient) (string, error) {
	t, err := template.ParseFiles("email.tmpl")
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	err1 := t.Execute(&tpl, r)
	if err1 != nil {
		return "", err1
	}

	return tpl.String(), nil
}
