package log

import (
	"fmt"

	"github.com/kerelape/urlshortener/internal/app/model"
)

type VerboseShortener struct {
	Origin model.Shortener
	Log    Log
}

func NewVerboseShortener(origin model.Shortener, log Log) *VerboseShortener {
	return &VerboseShortener{
		Origin: origin,
		Log:    log,
	}
}

func (shortener *VerboseShortener) Shorten(origin string) (string, error) {
	shortened, shortenError := shortener.Origin.Shorten(origin)
	if shortenError != nil {
		shortener.Log.WriteFailure("Failed to shorten: " + shortenError.Error())
	} else {
		shortener.Log.WriteInfo(
			fmt.Sprintf("Shorted \"%s\" to \"%s\"", origin, shortened),
		)
	}
	return shortened, shortenError
}

func (shortener *VerboseShortener) Reveal(short string) (string, error) {
	origin, err := shortener.Origin.Reveal(short)
	if err != nil {
		shortener.Log.WriteFailure(
			fmt.Sprintf("Failed to reveal \"%s\"", short),
		)
	} else {
		shortener.Log.WriteInfo(
			fmt.Sprintf("Revealed \"%s\" from \"%s\"", origin, short),
		)
	}
	return origin, err
}

func (shortener *VerboseShortener) ShortenAll(origins []string) ([]string, error) {
	shortened, shortenError := shortener.Origin.ShortenAll(origins)
	if shortenError != nil {
		shortener.Log.WriteFailure("Failed to shorten: " + shortenError.Error())
	} else {
		shortener.Log.WriteInfo(
			fmt.Sprintf(
				"Shortened %d from:\n\t\t%v\n\tto:\n\t\t%v",
				len(origins),
				origins,
				shortened,
			),
		)
	}
	return shortened, shortenError
}

func (shortener *VerboseShortener) RevealAll(shorts []string) ([]string, error) {
	origins, revealError := shortener.Origin.RevealAll(shorts)
	if revealError != nil {
		shortener.Log.WriteFailure(
			fmt.Sprintf("Failed to reveal: " + revealError.Error()),
		)
	} else {
		shortener.Log.WriteInfo(
			fmt.Sprintf(
				"Revealed %d from:\n\t\t%v\n\tto:\n\t\t%v",
				len(shorts),
				shorts,
				origins,
			),
		)
	}
	return origins, revealError
}
