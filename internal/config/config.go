package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env           string `yaml:"env" env-required:"true"`
	StorageServer `yaml:"storage_server"`
	HTTPServer    `yaml:"http_server"`
}

type StorageServer struct {
	Address string `yaml:"address" env-required:"true"`
	Port    string `yaml:"port" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Port        string        `yaml:"port" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoadConfig() *Config {
	configPath := "config/config.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("CONFIG_PATH does not exist")
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	return &config
}
