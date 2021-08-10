package models

import (
	"errors"
	"time"
)

var ErrNoRecords = errors.New("models: no records in result set")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}


