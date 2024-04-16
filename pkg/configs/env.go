package configs

import (
	"fmt"
	"github.com/golobby/dotenv"
	"os"
)

type session struct {
	LifespanSeconds int `env:"SESSION_LIFESPAN"`
}

type redis struct {
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
	User     string `env:"REDIS_USER"`
	Password string `env:"REDIS_PASS"`
	DB       int    `env:"REDIS_DB"`
}

type databaseEnv struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	Username string `env:"DB_USER"`
	Password string `env:"DB_PASS"`
	Database string `env:"DB_NAME"`
}

type env struct {
	Debug      bool   `env:"DEBUG"`
	ServerPort string `env:"SERVER_PORT"`
	Database   databaseEnv
	Session    session
	Redis      redis
}

var Env env

func LoadEnv() error {
	file, err := os.Open("./.env")
	if err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}

	err = dotenv.NewDecoder(file).Decode(&Env)
	if err != nil {
		return fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return nil
}
