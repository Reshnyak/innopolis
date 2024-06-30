package storage

import "fmt"

type Storage interface {
	WriteData(fName string, data []byte) error
	ReadData(fName string) ([]byte, error)
}

func New(config *Config) (Storage, error) {
	switch config.StorageType {
	case "file":
		return &StorageFile{}, nil
	default:
		return nil, fmt.Errorf("unsupported storage type")
	}

}
