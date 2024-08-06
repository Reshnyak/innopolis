package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime/trace"
	"syscall"

	"github.com/Reshnyak/innopolis/rwmultiservice/configs"
	"github.com/Reshnyak/innopolis/rwmultiservice/setup"
)

func main() {
	cfg, err := configs.ParseFlags()
	if err != nil {
		log.Printf("Config error. Set default values: %s", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// setup tracing
	fileTrace, _ := os.Create(cfg.FilePath + "trace.out")
	_ = trace.Start(fileTrace)
	defer trace.Stop()
	//instace rwmultisystem
	rwms, err := New(cfg)
	if err != nil {
		log.Fatalf("Could not initialize service: %s", err)
	}
	// Будем считать что пользователи отправляют обновления в существующие файлы.
	// Создание будем обрабатывать отдельно... когда-нибудь

	//сгенерируем файлы
	err = rwms.storage.File().SetupFiles()
	if err != nil {
		log.Fatalf("Could not generete files: %s", err)
	}
	//сгенерируем пользователей
	err = rwms.storage.User().SetupUsers(cfg)
	if err != nil {
		log.Fatalf("Could not generete users: %s", err)
	}
	// имитация отправки пользователями сообщений в каналы для обработки и записи
	inputChans, err := setup.SetupProcess(rwms.storage.GetAllUsers())
	if err != nil {
		log.Fatal(err)
	}
	err = rwms.Process(ctx, inputChans)
	if err != nil {
		log.Printf("Could not generete users: %s", err)
	}

}
