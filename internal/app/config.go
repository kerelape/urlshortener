package app

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

// Config is the application configuration.
type Config struct {
	// ServerAddress is the address that the app runs on.
	ServerAddress string `env:"SERVER_ADDRESS"`

	// BaseURL is the base url for the shortened urls.
	BaseURL string `env:"BASE_URL"`

	// ShortenerPath is the root of the app.
	ShortenerPath string `env:"SHORTENER_PATH"`

	// APIPath is the root of the app's API.
	APIPath string `env:"API_PATH"`

	// FileStoragePath is path to the file database.
	FileStoragePath string `env:"FILE_STORAGE_PATH"`

	// DatabaseDSN is the DSN to connect to.
	DatabaseDSN string `env:"DATABASE_DSN"`
}

// InitConfig initializes Config and returns it.
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
