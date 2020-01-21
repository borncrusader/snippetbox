package models

import (
	"errors"
	"time"
)

// ErrNoRecord is returned when no records are found
var ErrNoRecord = errors.New("models: no matching record found")

// Snippet is the base model for our application
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
