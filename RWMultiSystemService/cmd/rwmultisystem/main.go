package main

import (
	"flag"
	"github.com/Reshnyak/innopolis/RWMultiSystemService/internal/rwmultisystem"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/.env", "path to config file (.env file)")
}

func main() {
	flag.Parse()
	config := rwmultisystem.NewConfig()
	if err := godotenv.Load(configPath); err != nil {
		log.Fatal("could not find .env file:", err)
	} else {
		wd, err := time.ParseDuration(os.Getenv("worker_duration"))
		if err != nil {
			log.Println("could not parse worker_duration. using default duration value ", err)
		} else {
			config.WorkerDuration = wd
		}
		config.LogLevel = os.Getenv("logger_level")
		config.Storage = os.Getenv("logger_level")
		log.Printf("BindAddr:%s LoggerLevel:%s", config.BindAddr, config.LoggerLevel)
	}
}
