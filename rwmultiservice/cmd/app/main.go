package main

import (
	"context"
	"github.com/Reshnyak/innopolis/rwmultiservice/configs"
	"github.com/Reshnyak/innopolis/rwmultiservice/setup"
	"log"
	"os"
	"os/signal"
	"runtime/trace"
	"syscall"
)

func main() {
	cfg, err := configs.ParseFlags()
	if err != nil {
		log.Printf("Config error. Set default values: %s", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// setup tracing
	fileTrace, _ := os.Create("trace.out")
	trace.Start(fileTrace)
	defer trace.Stop()
	//instace rwmultisystem
	rwms, err := New(cfg)
	if err != nil {
		log.Fatalf("Could not initialize service: %s", err)
	}
	// Будем считать что пользователи отправляют обновления в существующие файлы.
	// Создание будем обрабатывать отдельно, и тогда создавать новые файлы
	inputChans, err := setup.SetupProcess(rwms.storage.GetAllUsers())
	if err != nil {
		log.Fatal(err)
	}

	rwms.Process(ctx, inputChans)

}
