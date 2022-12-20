package model

import "fmt"

type VerboseShortener struct {
	Origin Shortener
	Log    Log
}

func NewVerboseShortener(origin Shortener, log Log) *VerboseShortener {
	return &VerboseShortener{
		Origin: origin,
		Log:    log,
	}
}

func (shortener *VerboseShortener) Shorten(origin string) string {
	var shortened = shortener.Origin.Shorten(origin)
	shortener.Log.WriteInfo(
		fmt.Sprintf("Shorted \"%s\" to \"%s\"", origin, shortened),
	)
	return shortened
}

func (shortener *VerboseShortener) Reveal(short string) (string, error) {
	var origin, err = shortener.Origin.Reveal(short)
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
