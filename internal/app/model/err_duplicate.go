package model

import "errors"

// ErrDuplicate is returned when trying to shorten a URL that has
// already been shortened.
var ErrDuplicate = errors.New("duplicate")
