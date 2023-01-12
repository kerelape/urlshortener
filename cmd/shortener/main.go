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
	var urlShortener = model.NewURLShortener(shortener, conf.BaseURL, conf.ShortenerPath)
	var service = chi.NewRouter()
	service.Mount(conf.ShortenerPath, ui.NewApp(urlShortener).Route())
	service.Mount(conf.APIPath, ui.NewAPI(urlShortener).Route())
	http.ListenAndServe(conf.ServerAddress, service)
}

type config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ShortenerPath string `env:"SHORTENER_PATH" envDefault:"/"`
	APIPath       string `env:"API_PATH" envDefault:"/api"`
}
