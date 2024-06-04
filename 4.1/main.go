package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {
	lineChan := make(chan string)

	go ReadConsole(os.Stdin, lineChan)
	go WriteFile("file.txt", lineChan)

	sigchan := make(chan os.Signal, 1) //signal.Notify не настойчиво просит буфферезированный канал
	signal.Notify(sigchan, os.Interrupt)

	<-sigchan
	log.Println("Closed!")
	close(lineChan)
	os.Exit(0)
}
