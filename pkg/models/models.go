package models

import (
	"errors"
	"time"
)

// ErrNoRecord is returned when no records are found
var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

// Snippet is the base model for our application
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID             string
	Name           string
	Email          string
	SaltedPassword string
	Salt           string
	Created        time.Time
}
