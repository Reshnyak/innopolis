package rwmultisystem

import (
	"github.com/Reshnyak/innopolis/RWMultiSystemService/storage"
	"github.com/sirupsen/logrus"
)

type RWMultiSystem struct {
	config  *Config
	logger  *logrus.Logger
	storage storage.Storage
}

// RWMultiSystem constructor
func New(config *Config) *RWMultiSystem {
	return &RWMultiSystem{
		config: config,
		logger: logrus.New(),
	}
}

func (rwms *RWMultiSystem) Process(input <-chan []byte) error {
	if err := rwms.configureLogger(); err != nil {

	}
	if err := rwms.configureStore(); err != nil {
		return err
	}

	return nil
}

// configureStore method
func (rwms *RWMultiSystem) configureStore() error {
	st, err := storage.New(rwms.config.Storage)
	if err != nil {
		return err
	}
	rwms.storage = st
	return nil
}

// func for configureate logger
func (rwms *RWMultiSystem) configureLogger() error {
	level, err := logrus.ParseLevel(rwms.config.LogLevel)
	if err != nil {
		return nil
	}
	rwms.logger.SetLevel(level)

	return nil
}
