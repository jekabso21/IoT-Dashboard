package models

import (
	"time"

	"github.com/google/uuid"
)

type PairingCode struct {
	Code      string     `json:"code" db:"code"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	Used      bool       `json:"used" db:"used"`
	UsedAt    *time.Time `json:"used_at,omitempty" db:"used_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

type GeneratePairingCodeResponse struct {
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	ExpiresIn int       `json:"expires_in"`
}
