package model

import "time"

// User represents a row in the users table.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // never serialised to JSON
	CreatedAt    time.Time `json:"created_at"`
}
