package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	lineChan := make(chan string)
	sigchan := make(chan os.Signal, 1) //signal.Notify не настойчиво просит буфферезированный канал
	signal.Notify(sigchan, os.Interrupt)

	go func() {
		<-sigchan
		log.Println("Closed!")
		cancel()
		os.Exit(0)
	}()

	wg.Add(2)
	{
		ReadConsole(ctx, wg, os.Stdin, lineChan)
		WriteFile(ctx, wg, "file.txt", lineChan)
	}
	wg.Wait()

}
