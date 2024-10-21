package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env             string        `yaml:"env" env-required:"true"`
	AliasLen        int           `yaml:"alias_len" env-default:"8"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"5s"`
	TokenTTL        time.Duration `yaml:"token_ttl" env-default:"1h"`
	PostgresServer  `yaml:"postgres_server"`
	HTTPServer      `yaml:"http_server"`
}

type PostgresServer struct {
	Address  string `yaml:"address" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DbName   string `yaml:"db_name" env-required:"true"`
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

	if config.Env != "prod" && config.Env != "local" {
		log.Fatal("CONFIG_ENV must be either 'prod' or 'local'")
	}

	return &config
}
