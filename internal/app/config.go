package app

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	ShortenerPath   string `env:"SHORTENER_PATH"`
	APIPath         string `env:"API_PATH"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func InitConfig() (Config, error) {
	var environment Config
	parseError := env.Parse(&environment)
	if parseError != nil {
		return environment, parseError
	}
	flags := Config{}
	flag.StringVar(&flags.ServerAddress, "a", "localhost:8080", "Server address")
	flag.StringVar(&flags.BaseURL, "b", "http://localhost:8080", "Base URL")
	flag.StringVar(&flags.FileStoragePath, "f", "urlshortener.db", "Path file DB")
	flag.StringVar(&flags.APIPath, "api-path", "/api", "API root")
	flag.StringVar(&flags.ShortenerPath, "app-path", "/", "Shortener root")
	flag.StringVar(&flags.DatabaseDSN, "d", "", "")
	flag.Parse()
	if environment.ServerAddress == "" {
		environment.ServerAddress = flags.ServerAddress
	}
	if environment.BaseURL == "" {
		environment.BaseURL = flags.BaseURL
	}
	if environment.FileStoragePath == "" {
		environment.FileStoragePath = flags.FileStoragePath
	}
	if environment.APIPath == "" {
		environment.APIPath = flags.APIPath
	}
	if environment.ShortenerPath == "" {
		environment.ShortenerPath = flags.ShortenerPath
	}
	if environment.DatabaseDSN == "" {
		environment.DatabaseDSN = flags.DatabaseDSN
	}
	return environment, parseError
}
