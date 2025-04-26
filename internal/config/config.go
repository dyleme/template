package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `env:"ENV" env-required:"true"`
	Database Database
	JWT      JWT
	APIKey   APIKey
	Server   Server
}

type Server struct {
	Port                    int           `env:"APP_PORT"               env-required:"true"`
	MaxHeaderBytes          int           `env:"MAX_HEADER"             env-required:"true"`
	ReadTimeout             time.Duration `env:"READ_TIMEOUT"           env-required:"true"`
	WriteTimeout            time.Duration `env:"WRITE_TIMEOUT"          env-required:"true"`
	TimeForGracefulShutdown time.Duration `env:"GRACEFUL_SHUTDOWN_TIME" env-required:"true"`
}

type Database struct {
	Port     int    `env:"DB_PORT"           env-required:"true"`
	Host     string `env:"DB_HOST"           env-required:"true"`
	SSLMode  string `env:"DB_SSL_MODE"       env-required:"true"`
	User     string `env:"POSTGRES_USER"     env-required:"true"`
	Database string `env:"POSTGRES_DB"       env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
}

func (db *Database) ConnectionString() string {
	address := net.JoinHostPort(db.Host, strconv.Itoa(db.Port))

	return fmt.Sprintf("postgres://%s:%s@%s/%s", db.User, db.Password, address, db.Database)
}

type JWT struct {
	TokenTTL  time.Duration `env:"TOKEN_TTL"  env-required:"true"`
	SignedKey string        `env:"SIGNED_KEY" env-required:"true"`
}

type APIKey struct {
	Key string `env:"API_KEY" env-required:"true"`
}

func Load() (Config, error) {
	var cfg Config
	filename := ".env"
	folderPath, err := os.Getwd()
	if err != nil {
		return Config{}, fmt.Errorf("can't get working dir: %w", err)
	}

	var pathToFile string
	for folderPath != "/" {
		pathToFile = folderPath + "/" + filename
		_, err = os.Stat(pathToFile)
		if err == nil {
			break
		}
		idx := strings.LastIndex(folderPath, "/")
		if idx == -1 {
			break
		}
		folderPath = folderPath[:idx]
	}

	err = cleanenv.ReadConfig(pathToFile, &cfg)
	if err != nil {
		err = cleanenv.ReadEnv(&cfg)
		if err != nil {
			return Config{}, fmt.Errorf("can't read env: %w", err)
		}
	}

	return cfg, nil
}
