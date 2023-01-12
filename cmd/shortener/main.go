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

const (
	DefaultHost = "localhost:8080"
	Path        = "/"
)

type config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func main() {
	var log = model.NewFormattedLog(
		model.NewWriterLog(os.Stdout, os.Stderr),
		time.UnixDate,
	)
	var shortener = model.NewVerboseShortener(
		model.NewAlphabetShortener(
			model.NewFakeDatabase(),
			model.NewBase62Alphabet(),
		),
		log,
	)
	var conf config
	var parseError = env.Parse(&conf)
	if parseError != nil {
		panic(parseError)
	}
	log.WriteInfo("BASE_URL " + conf.BaseURL)
	log.WriteInfo("SERVER_ADDRESS" + conf.ServerAddress)
	var urlShortener = model.NewURLShortener(shortener, "http://"+conf.BaseURL+Path)
	var service = chi.NewRouter()
	service.Mount("/", ui.NewApp(urlShortener).Route())
	service.Mount("/api", ui.NewAPI(urlShortener).Route())
	http.ListenAndServe(DefaultHost, service)
}
