package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {
	lineChan := make(chan string)
	errorChan := make(chan error)
	sigchan := make(chan os.Signal, 1) //signal.Notify не настойчиво просит буфферезированный канал
	signal.Notify(sigchan, os.Interrupt)
	go func() {
		<-sigchan
		log.Println("Closed!")
		close(lineChan)
		os.Exit(0)
	}()
	ReadConsole(os.Stdin, lineChan, errorChan)
	WriteFile("file.txt", lineChan, errorChan)
	for err := range errorChan {
		log.Println(err)
	}
}
