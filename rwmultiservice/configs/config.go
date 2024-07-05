package configs

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var configPath = flag.String("config-path", "configs/cfg.yaml", "path to config file (.yaml or .env file)")

type Config struct {
	TokenLen        int           `yaml:"token_len"`
	WorkerDuration  time.Duration `yaml:"worker_duration"`
	WorkersCount    int           `yaml:"worker_count"`
	FilePath        string        `yaml:"file_path"`
	FilesCount      int           `yaml:"files_count"`
	UsersCount      int           `yaml:"users_count"`
	MessageMaxCount int           `yaml:"message_max_count"`
	MessageMaxLen   int           `yaml:"message_max_len"`
	StorageType     string        `yaml:"storage_type"`
}

func DefaultInitialize() *Config {
	return &Config{
		TokenLen:        16,
		WorkerDuration:  1 * time.Second,
		WorkersCount:    4,
		FilePath:        "/data",
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
	case "env":
		return parseEnvCfg(configPath)
	default:
		return DefaultInitialize(), fmt.Errorf("incorrect config path")
	}
}
func ParseFlags() (*Config, error) {
	flag.Parse()

	return GetConfig(*configPath)
}

func parseYamlCfg(configPath string) (*Config, error) {
	config := &Config{}
	yamlConf, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultInitialize(), fmt.Errorf("could not read config file:%s", err)
	}
	err = yaml.Unmarshal(yamlConf, config)
	if err != nil {
		return DefaultInitialize(), fmt.Errorf("could not unmarshal config:%s", err)
	}
	return config, nil
}

func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("environment variable %s not set", key)
}
func getEnvAsInt(name string) (int, error) {
	var value int
	valueStr, err := getEnv(name)
	if err != nil {
		return 0, err
	}
	value, err = strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s not cast to int", name)
	}
	return value, nil
}
func getEnvAsDuration(name string) (time.Duration, error) {
	var value time.Duration
	valueStr, err := getEnv(name)
	if err != nil {
		return 0, err
	}
	value, err = time.ParseDuration(valueStr)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s not cast to time.Duration", name)
	}
	return value, nil
}

func parseEnvCfg(configPath string) (*Config, error) {
	config := &Config{}
	cfg := DefaultInitialize()
	err := godotenv.Load(configPath)
	if err != nil {
		return cfg, fmt.Errorf("could not load config from .env:%s", err)
	}
	config.TokenLen, err = getEnvAsInt("TOKEN_LEN")
	if err != nil {
		return cfg, fmt.Errorf("could not parse config from .env:%s", err)
	}
	config.WorkerDuration, err = getEnvAsDuration("WORKER_DURATION")
	if err != nil {
		return cfg, fmt.Errorf("could not parse config from .env:%s", err)
	}
	config.WorkersCount, err = getEnvAsInt("WORKER_COUNT")
	if err != nil {
		return cfg, fmt.Errorf("could not parse config from .env:%s", err)
	}
	config.FilePath, err = getEnv("FILE_PATH")
	if err != nil {
		return cfg, fmt.Errorf("could not parse config from .env:%s", err)
	}
	config.FilesCount, err = getEnvAsInt("FILE_COUNT")
	if err != nil {
		return cfg, fmt.Errorf("could not parse config from .env:%s", err)
	}
	config.UsersCount, err = getEnvAsInt("USER_COUNT")
	if err != nil {
		return cfg, fmt.Errorf("could not parse config from .env:%s", err)
	}
	config.MessageMaxCount, err = getEnvAsInt("MESSAGE_MAX_COUNT")
	if err != nil {
		return cfg, fmt.Errorf("could not parse config from .env:%s", err)
	}
	config.MessageMaxLen, err = getEnvAsInt("MESSAGE_MAX_LEN")
	if err != nil {
		return cfg, fmt.Errorf("could not parse config from .env:%s", err)
	}
	config.FilePath, err = getEnv("STORAGE_TYPE")
	if err != nil {
		return cfg, fmt.Errorf("could not parse config from .env:%s", err)
	}
	return config, nil
}
