package log

import (
	"context"
	"fmt"

	"github.com/kerelape/urlshortener/internal/app"
	"github.com/kerelape/urlshortener/internal/app/model"
)

// VerboseShortener is a shortener that writes its actions to the log.
type VerboseShortener struct {
	origin model.Shortener
	log    Log
}

// NewverboseShortener returns a new VerboseShortener.
func NewVerboseShortener(origin model.Shortener, log Log) *VerboseShortener {
	return &VerboseShortener{
		origin: origin,
		log:    log,
	}
}

// Shorten shortens the given origin string.
func (shortener *VerboseShortener) Shorten(ctx context.Context, user app.Token, origin string) (string, error) {
	shortened, shortenError := shortener.origin.Shorten(ctx, user, origin)
	if shortenError != nil {
		go shortener.log.WriteFailure("Failed to shorten: " + shortenError.Error())
	} else {
		go shortener.log.WriteInfo(
			fmt.Sprintf("Shorted \"%s\" to \"%s\"", origin, shortened),
		)
	}
	return shortened, shortenError
}

// Reveal returns the original string by the shortened.
func (shortener *VerboseShortener) Reveal(ctx context.Context, short string) (string, error) {
	origin, err := shortener.origin.Reveal(ctx, short)
	if err != nil {
		go shortener.log.WriteFailure(
			fmt.Sprintf("Failed to reveal \"%s\"", short),
		)
	} else {
		go shortener.log.WriteInfo(
			fmt.Sprintf("Revealed \"%s\" from \"%s\"", origin, short),
		)
	}
	return origin, err
}

// ShortenAll shortens a slice of strings and returns
// a slice of short string in the same order.
func (shortener *VerboseShortener) ShortenAll(ctx context.Context, user app.Token, origins []string) ([]string, error) {
	shortened, shortenError := shortener.origin.ShortenAll(ctx, user, origins)
	if shortenError != nil {
		go shortener.log.WriteFailure("Failed to shorten: " + shortenError.Error())
	} else {
		go shortener.log.WriteInfo(
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

// RevealAll returns a slice of original strings in the same order
// as in shortened.
func (shortener *VerboseShortener) RevealAll(ctx context.Context, shorts []string) ([]string, error) {
	origins, revealError := shortener.origin.RevealAll(ctx, shorts)
	if revealError != nil {
		go shortener.log.WriteFailure(
			fmt.Sprintf("Failed to reveal: " + revealError.Error()),
		)
	} else {
		go shortener.log.WriteInfo(
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

// Delete deletes a string from the shortener.
func (shortener *VerboseShortener) Delete(ctx context.Context, user app.Token, shorts []string) error {
	err := shortener.origin.Delete(ctx, user, shorts)
	if err != nil {
		go shortener.log.WriteFailure("Failed to delete: " + err.Error())
	} else {
		go func() {
			for _, s := range shorts {
				shortener.log.WriteInfo("Deleted url: " + s)
			}
		}()
	}
	return err
}
