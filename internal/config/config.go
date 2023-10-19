package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Store   string        `yaml:"store"`
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("cannot load .env file. %s", err)
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist. %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config. %s", err)
	}
	return &cfg
}
