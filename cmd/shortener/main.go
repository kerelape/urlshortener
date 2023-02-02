package main

import (
	"compress/gzip"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	logging "github.com/kerelape/urlshortener/internal/app/log"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
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
	var service = initService(shortener, &config, log)
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

func initShortener(database storage.Database, log logging.Log, config *Config) model.Shortener {
	return model.NewURLShortener(
		logging.NewVerboseShortener(
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

func initLog() logging.Log {
	return logging.NewFormattedLog(
		logging.NewWriterLog(os.Stdout, os.Stderr),
		time.UnixDate,
	)
}

func initDatabase(config *Config) (storage.Database, error) {
	var database storage.Database
	if config.FileStoragePath == "" {
		database = storage.NewFakeDatabase()
	} else {
		var fileDatabase, openFileDatabaseError = storage.OpenFileDatabase(config.FileStoragePath, true, 0644)
		if openFileDatabaseError != nil {
			return nil, openFileDatabaseError
		}
		database = fileDatabase
	}
	return database, nil
}

func initService(model model.Shortener, config *Config, log logging.Log) http.Handler {
	var router = chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.WriteInfo("content encoding: " + r.Header.Get("Content-Encoding"))
			h.ServeHTTP(w, r)
		})
	})
	router.Use(middleware.Compress(gzip.BestCompression))
	router.Use(ui.Decompress())
	router.Mount(config.ShortenerPath, ui.NewApp(model).Route())
	router.Mount(config.APIPath, ui.NewAPI(model).Route())
	return router
}
