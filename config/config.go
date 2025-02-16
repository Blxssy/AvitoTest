package config

import (
	"log"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Logger Logger
	Server ServerConfig
	PG     PostgresConfig
	Token  TokenConfig
}

type Logger struct {
	Level string `env:"LOG_LEVEL,required"`
}

type ServerConfig struct {
	Addr string `env:"SERVER_ADDR,required"`
}

type PostgresConfig struct {
	DataSource       string `env:"DB_DATA_SOURCE,required"`
	PathToMigrations string `env:"DB_PATH_TO_MIGRATIONS,required"`
}

type TokenConfig struct {
	TokenKey string        `env:"TOKEN_KEY"`
	TokenTTL time.Duration `env:"TOKEN_TTL"`
}

var (
	config Config
	once   sync.Once
)

func Get() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading env")
		}
		if err := env.Parse(&config); err != nil {
			log.Fatal(err)
		}
	})
	return &config
}
