package configs

import (
	"flag"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var configPath = flag.String("config-path", "configs/cfg.yaml", "path to config file (.yaml or .env file)")

type Config struct {
	TokenLen        int           `yaml:"token_len" env:"TOKEN_LEN" env-default:"16"`
	WorkerDuration  time.Duration `yaml:"worker_duration" env:"WORKER_DURATION" env-default:"1s"`
	WorkersCount    int           `yaml:"worker_count" env:"WORKER_COUNT" env-default:"4"`
	FilePath        string        `yaml:"file_path" env:"FILE_PATH" env-default:"data/"`
	FilesCount      int           `yaml:"files_count" env:"FILE_COUNT" env-default:"4"`
	UsersCount      int           `yaml:"users_count" env:"USER_COUNT" env-default:"3"`
	MessageMaxCount int           `yaml:"message_max_count" env:"MESSAGE_MAX_COUNT" env-default:"40"`
	MessageMaxLen   int           `yaml:"message_max_len" env:"MESSAGE_MAX_LEN" env-default:"24"`
	StorageType     string        `yaml:"storage_type" env:"STORAGE_TYPE" env-default:"FS"`
}

func DefaultInitialize() *Config {
	return &Config{
		TokenLen:        16,
		WorkerDuration:  1 * time.Second,
		WorkersCount:    4,
		FilePath:        "data/",
		FilesCount:      4,
		UsersCount:      3,
		MessageMaxCount: 40,
		MessageMaxLen:   24,
		StorageType:     "FS",
	}
}

func getConfigType(configPath string) string {
	return configPath[strings.LastIndex(configPath, ".")+1:]
}
func GetConfig(configPath string) (*Config, error) {
	switch getConfigType(configPath) {
	case "yaml":
		return parseYamlCfg(configPath)
	default:
		return parseEnvCfg(configPath)
	}
}
func ParseFlags() (*Config, error) {
	flag.Parse()

	return GetConfig(*configPath)
}

func parseYamlCfg(configPath string) (*Config, error) {
	config := &Config{}
	err := cleanenv.ReadConfig(configPath, config)
	return config, err
}

func parseEnvCfg(configPath string) (*Config, error) {
	config := &Config{}
	_ = godotenv.Load(configPath)
	err := cleanenv.ReadEnv(config)
	return config, err
}
