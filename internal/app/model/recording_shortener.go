package model

import "github.com/kerelape/urlshortener/internal/app/model/storage"

type RecordingShortener struct {
	origin  Shortener
	history storage.History
}
