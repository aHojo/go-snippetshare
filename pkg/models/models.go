package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecords = errors.New("models: no records in result set")

	// Invalid email or password
	ErrInvalidCredentials = errors.New("models: invalid email/password")
	ErrDuplicateEmail      = errors.New("models: duplicate email")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}



type User struct {
	ID int
	Name string
	Email string
	HashedPassword []byte
	Created time.Time
}
