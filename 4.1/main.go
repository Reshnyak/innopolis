package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {

	sigchan := make(chan os.Signal, 1) //signal.Notify не настойчиво просит буфферезированный канал
	signal.Notify(sigchan, os.Interrupt)
	go func() {
		<-sigchan
		log.Println("Closed!")
		os.Exit(0)
	}()

	lineChan := ReadAndSendString(os.Stdin)
	for err := range WriteFile("file.txt", lineChan) {
		log.Println(err)
	}
}
