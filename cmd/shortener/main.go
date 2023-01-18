package main

import (
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/kerelape/urlshortener/internal/app/model"
	"github.com/kerelape/urlshortener/internal/app/ui"
)

func main() {
	var conf config
	var parseError = env.Parse(&conf)
	if parseError != nil {
		panic(parseError)
	}

	var log = model.NewFormattedLog(
		model.NewWriterLog(os.Stdout, os.Stderr),
		time.UnixDate,
	)

	var database model.Database
	if conf.FileStoragePath == "" {
		database = model.NewFakeDatabase()
	} else {
		var file, openFileError = os.OpenFile(
			conf.FileStoragePath,
			os.O_RDWR|os.O_CREATE,
			0644,
		)
		if openFileError != nil {
			panic(openFileError)
		}
		defer (func() {
			var closeError = file.Close()
			if closeError != nil {
				panic(closeError)
			}
		})()
		database = model.NewFileDatabase(file)
	}

	var shortener = model.NewVerboseShortener(
		model.NewAlphabetShortener(
			database,
			model.NewBase62Alphabet(),
		),
		log,
	)

	var urlShortener = model.NewURLShortener(shortener, conf.BaseURL, conf.ShortenerPath)
	var service = chi.NewRouter()
	service.Mount(conf.ShortenerPath, ui.NewApp(urlShortener).Route())
	service.Mount(conf.APIPath, ui.NewAPI(urlShortener).Route())
	http.ListenAndServe(conf.ServerAddress, service)
}

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ShortenerPath   string `env:"SHORTENER_PATH" envDefault:"/"`
	APIPath         string `env:"API_PATH" envDefault:"/api"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}
