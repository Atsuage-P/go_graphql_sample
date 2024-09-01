package env

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DB DBConfig `envPrefix:"DB_"`
}

type DBConfig struct {
	MS               string `env:"MS,notEmpty"`
	Host             string `env:"HOST,notEmpty"`
	Name             string `env:"NAME,notEmpty"`
	User             string `env:"USER,notEmpty"`
	Password         string `env:"PASSWORD,notEmpty"`
	Port             string `env:"PORT,notEmpty"`
	MaxLifeTimeMin   int    `env:"MAX_LIFE_TIME_MIN,notEmpty"`
	MaxOpenConns     int    `env:"MAX_OPEN_CONNS,notEmpty"`
	MaxOpenIdleConns int    `env:"MAX_IDLE_CONNS,notEmpty"`
}

func LoadEnv() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Load Env Error: %v", err)
	}
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Env Parse Error: %v", err)
	}

	return &cfg
}
