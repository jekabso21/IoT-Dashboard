package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Email          string    `json:"email" db:"email"`
	PasswordHash   string    `json:"-" db:"password_hash"`
	FullName       *string   `json:"full_name,omitempty" db:"full_name"`
	Company        *string   `json:"company,omitempty" db:"company"`
	EmailVerified  bool      `json:"email_verified" db:"email_verified"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8"`
	FullName *string `json:"full_name,omitempty"`
	Company  *string `json:"company,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
