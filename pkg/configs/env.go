package configs

import (
	"fmt"
	"os"

	"github.com/golobby/dotenv"
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

type email struct {
	SMTPHost      string `env:"EMAIL_SMTP_HOST"`
	SMTPPort      int    `env:"EMAIL_SMTP_PORT"`
	Username      string `env:"EMAIL_SMTP_USER"`
	Password      string `env:"EMAIL_SMTP_PASS"`
	SenderAddress string `env:"EMAIL_SENDER_ADDRESS"`
}

type env struct {
	Frontend    string `env:"FRONTEND_URL"`
	Debug       bool   `env:"DEBUG"`
	ServerPort  string `env:"SERVER_PORT"`
	MaxFileSize int64  `env:"MAX_FILE_SIZE"`
	Database    databaseEnv
	Session     session
	Redis       redis
	Email       email
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

	Env.MaxFileSize = Env.MaxFileSize * 1_000_000

	return nil
}
