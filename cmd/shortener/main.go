package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/ui"
)

var (
	Service http.Handler
	Address string
)

func init() {
	var config, configError = initConfig()
	if configError != nil {
		panic(configError)
	}
	var database, databaseError = initDatabase(&config)
	if databaseError != nil {
		panic(databaseError)
	}
	var log = initLog()
	var shortener = initShortener(database, log, &config)
	var service = initService(shortener, &config)
	Address = config.ServerAddress
	Service = service
}

func main() {
	http.ListenAndServe(Address, Service)
}

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	ShortenerPath   string `env:"SHORTENER_PATH"`
	APIPath         string `env:"API_PATH"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func initConfig() (Config, error) {
	var environment Config
	var parseError = env.Parse(&environment)
	if parseError != nil {
		return environment, parseError
	}
	var flags = Config{}
	flag.StringVar(&flags.ServerAddress, "a", "localhost:8080", "Server address")
	flag.StringVar(&flags.BaseURL, "b", "http://localhost:8080", "Base URL")
	flag.StringVar(&flags.FileStoragePath, "f", "/var/cache/urlshortener.db", "Path file DB")
	flag.StringVar(&flags.APIPath, "api-path", "/api", "API root")
	flag.StringVar(&flags.ShortenerPath, "app-path", "/", "Shortener root")
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
	return environment, parseError
}

func initShortener(database model.Database, log model.Log, config *Config) model.Shortener {
	return model.NewURLShortener(
		model.NewVerboseShortener(
			model.NewAlphabetShortener(
				database,
				model.NewBase62Alphabet(),
			),
			log,
		),
		config.BaseURL,
		config.ShortenerPath,
	)
}

func initLog() model.Log {
	return model.NewFormattedLog(
		model.NewWriterLog(os.Stdout, os.Stderr),
		time.UnixDate,
	)
}

func initDatabase(config *Config) (model.Database, error) {
	var database model.Database
	if config.FileStoragePath == "" {
		database = model.NewFakeDatabase()
	} else {
		var file, openFileError = os.OpenFile(
			config.FileStoragePath,
			os.O_RDWR|os.O_CREATE,
			0644,
		)
		if openFileError != nil {
			return nil, openFileError
		}
		database = model.NewFileDatabase(file)
	}
	return database, nil
}

func initService(model model.Shortener, config *Config) http.Handler {
	var router = chi.NewRouter()
	router.Mount(config.ShortenerPath, ui.NewApp(model).Route())
	router.Mount(config.APIPath, ui.NewAPI(model).Route())
	return router
}
