package config

import (
	"database/sql"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	DbConfig   `yaml:"db"`
}
type DbConfig struct {
	Host     string `yaml:"host" env-default:"db"`
	Port     string `yaml:"port" env-default:"5432"`
	Username string `yaml:"username" env-default:"postgres"`
	Password string `yaml:"password" env-default:"qwerty"`
	DBName   string `yaml:"dbname" env-default:"postgres"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does nob exist: %s", configPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: %s", err)
	}
	return &cfg
}

func NewPostgresDB(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DbConfig.Host, cfg.DbConfig.Port, cfg.DbConfig.Username, cfg.DbConfig.DBName, cfg.DbConfig.Password, cfg.DbConfig.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
