package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Config struct {
	App struct {
		Env     string
		Port    string
		BaseURL string
	}

	Database DatabaseConfig
}

func Load() (*Config, error) {

	_ = godotenv.Load()

	config := &Config{}

	config.App.Env = os.Getenv("APP_ENV")
	config.App.Port = os.Getenv("APP_PORT")
	config.App.BaseURL = os.Getenv("APP_BASE_URL")

	config.Database.Host = os.Getenv("DATABASE_HOST")
	config.Database.Port = os.Getenv("DATABASE_PORT")
	config.Database.User = os.Getenv("DATABASE_USER")
	config.Database.Password = os.Getenv("DATABASE_PASSWORD")
	config.Database.Name = os.Getenv("DATABASE_NAME")

	if config.App.Port == "" {
		config.App.Port = "8080"
	}

	if config.App.BaseURL == "" {
		return nil, fmt.Errorf(
			"APP_BASE_URL is required",
		)
	}

	if config.Database.Host == "" {
		return nil, fmt.Errorf(
			"DATABASE_HOST is required",
		)
	}

	if config.Database.Port == "" {
		return nil, fmt.Errorf(
			"DATABASE_PORT is required",
		)
	}

	if config.Database.User == "" {
		return nil, fmt.Errorf(
			"DATABASE_USER is required",
		)
	}

	if config.Database.Name == "" {
		return nil, fmt.Errorf(
			"DATABASE_NAME is required",
		)
	}

	return config, nil
}
