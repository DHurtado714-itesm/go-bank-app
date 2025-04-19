package auth

import "time"

type User struct {
	ID string `json:"id"`
	Email string `json:"email"`
	HashedPassword string `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}