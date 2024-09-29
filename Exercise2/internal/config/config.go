package config

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

func NewConfig() *Config {
	// Загрузка переменных из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	// Пример парсинга переменной окружения как времени (если нужно)
	_ = os.Getenv("DB_PORT")

	// Создание структуры конфигурации
	config := &Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
	}

	// Валидация конфигурации
	if err := config.Validate(); err != nil {
		log.Fatalf("Configuration validation error: %v", err)
	}

	return config
}

func (c *Config) Validate() error {
	if c.Host == "" {
		return errors.New("Host must not be empty")
	}
	if c.Port == "" {
		return errors.New("Port must not be empty")
	}
	if c.User == "" {
		return errors.New("User must not be empty")
	}
	if c.Password == "" {
		return errors.New("Password must not be empty")
	}
	if c.DbName == "" {
		return errors.New("DbName must not be empty")
	}

	return nil
}
