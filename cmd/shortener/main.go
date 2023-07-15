package main

import (
	"compress/gzip"
	"context"
	"fmt"
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

	"net/http/pprof"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	fmt.Printf("Build version: %s\n", buildVersion)

	if buildDate == "" {
		buildDate = "N/A"
	}
	fmt.Printf("Build date: %s\n", buildDate)

	if buildCommit == "" {
		buildCommit = "N/A"
	}
	fmt.Printf("Build commit: %s\n", buildCommit)

	config, configError := app.InitConfig()
	if configError != nil {
		panic(configError)
	}
	database, databaseError := initDatabase(&config)
	if databaseError != nil {
		panic(databaseError)
	}
	history, historyError := initHistory()
	if historyError != nil {
		panic(historyError)
	}
	log := initLog()
	address := config.ServerAddress
	shortener := initShortener(database, log, &config)
	service := initService(shortener, &config, log, history, database)
	http.ListenAndServe(address, service)
}

func initShortener(database storage.Database, log logging.Log, config *app.Config) model.Shortener {
	return logging.NewVerboseShortener(
		model.NewURLShortener(
			model.NewAlphabetShortener(
				database,
				model.NewBase62Alphabet(),
			),
			config.BaseURL,
			config.ShortenerPath,
		),
		log,
	)
}

func initLog() logging.Log {
	return logging.NewFormattedLog(
		logging.NewWriterLog(os.Stdout, os.Stderr),
		time.UnixDate,
	)
}

func initDatabase(config *app.Config) (storage.Database, error) {
	if config.DatabaseDSN != "" {
		return storage.DialPostgreSQLDatabase(context.Background(), config.DatabaseDSN)
	}
	var database storage.Database
	if config.FileStoragePath == "" {
		database = storage.NewFakeDatabase()
	} else {
		fileDatabase, openFileDatabaseError := storage.OpenFileDatabase(config.FileStoragePath, true, 0o644, 1024)
		if openFileDatabaseError != nil {
			return nil, openFileDatabaseError
		}
		database = fileDatabase
	}
	return database, nil
}

func initHistory() (storage.History, error) {
	return storage.NewVirtualHistory(), nil
}

func initService(
	model model.Shortener,
	config *app.Config,
	log logging.Log,
	history storage.History,
	database storage.Database,
) http.Handler {
	webUI := ui.NewWebUI(
		map[string]ui.Entry{
			config.ShortenerPath: ui.NewApp(model, history),
			config.APIPath:       api.NewAPI(model, history),
			"/ping":              ui.NewSQLPing(database),
		},
	)
	router := chi.NewRouter()
	router.Use(ui.Decompress())
	router.Use(ui.Tokenize())
	router.Use(middleware.Compress(gzip.BestCompression))
	router.Mount("/", webUI.Route())
	router.Mount("/debug/pprof", http.HandlerFunc(pprof.Index))
	return router
}
