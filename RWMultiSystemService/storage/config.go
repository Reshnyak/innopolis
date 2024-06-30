package storage

type Config struct {
	StorageType string
	Path        string
}

func NewConfig() *Config {
	return &Config{}
}
