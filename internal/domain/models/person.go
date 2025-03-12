package models

import "github.com/google/uuid"

type Person struct {
	Id           uuid.UUID `json:"id" db:"id"`
	Login        string    `json:"login" db:"login"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	AboutMe      string    `json:"about_me" db:"about_me"`
	Image        []byte    `json:"image" db:"image"`
}
