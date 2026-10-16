package main

type Recepient struct {
	Name  string
	Email string
}

func main() {
	recepientChannel := make(chan Recepient)
	loadRecepient("./emails.csv", recepientChannel)
}
