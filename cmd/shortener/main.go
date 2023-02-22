package main

import (
	"compress/gzip"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kerelape/urlshortener/internal/app"
	logging "github.com/kerelape/urlshortener/internal/app/log"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
	"github.com/kerelape/urlshortener/internal/app/ui"
	"github.com/kerelape/urlshortener/internal/app/ui/api"
)

func main() {
	var config, configError = app.InitConfig()
	if configError != nil {
		panic(configError)
	}
	var database, databaseError = initDatabase(&config)
	if databaseError != nil {
		panic(databaseError)
	}
	var log = initLog()
	var address = config.ServerAddress
	var shortener = initShortener(database, log, &config)
	var service = initService(shortener, &config, log)
	http.ListenAndServe(address, service)
}

func initShortener(database storage.Database, log logging.Log, config *app.Config) model.Shortener {
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

func initDatabase(config *app.Config) (storage.Database, error) {
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

func initService(model model.Shortener, config *app.Config, log logging.Log) http.Handler {
	var router = chi.NewRouter()
	router.Use(middleware.Compress(gzip.BestCompression))
	router.Use(ui.Decompress())
	router.Mount(config.ShortenerPath, ui.NewApp(model).Route())
	var api = api.NewAPI(
		api.NewShortenAPI(model),
	)
	router.Mount(config.APIPath, api.Route())
	return router
}
