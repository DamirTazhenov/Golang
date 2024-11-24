package config

import (
	"errors"
	"log"
	"os"
)

type Config struct {
	BindAddr     string
	LogLevel     string
	DatabaseURL  string
	JWTSecretKey string
	CSRFSecret   string
}

func NewConfig() *Config {

	config := &Config{
		BindAddr:     os.Getenv("BIND_ADDR"),
		LogLevel:     os.Getenv("LOG_LEVEL"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
		CSRFSecret:   os.Getenv("CSRF_SECRET"),
	}

	if err := config.Validate(); err != nil {
		log.Fatalf("Configuration validation error: %v", err)
	}

	return config
}

func (c *Config) Validate() error {
	if c.BindAddr == "" {
		return errors.New("BindAddr must not be empty")
	}
	if c.LogLevel == "" {
		return errors.New("LogLevel must not be empty")
	}
	if c.DatabaseURL == "" {
		return errors.New("DatabaseURL must not be empty")
	}
	if c.JWTSecretKey == "" {
		return errors.New("JWTSecretKey must not be empty")
	}
	if c.CSRFSecret == "" {
		return errors.New("JWTSecretKey must not be empty")
	}

	return nil
}
