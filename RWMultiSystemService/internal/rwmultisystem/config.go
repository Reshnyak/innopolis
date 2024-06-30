package rwmultisystem

import (
	"time"

	"github.com/Reshnyak/innopolis/RWMultiSystemService/storage"
)

// General config for rwmultisystem
type Config struct {
	WorkerDuration time.Duration
	LogLevel       string
	Storage        *storage.Config
}

func NewConfig() *Config {
	return &Config{
		WorkerDuration: 1 * time.Second,
		LogLevel:       "debug",
		Storage:        storage.NewConfig(),
	}
}
