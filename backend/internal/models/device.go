package models

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	DeviceName string     `json:"device_name" db:"device_name"`
	DeviceType string     `json:"device_type" db:"device_type"`
	LastSeen   *time.Time `json:"last_seen,omitempty" db:"last_seen"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	PairedAt   *time.Time `json:"paired_at,omitempty" db:"paired_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	Online     bool       `json:"online" db:"online"`
}

type DeviceToken struct {
	ID               uuid.UUID `json:"id" db:"id"`
	DeviceID         uuid.UUID `json:"device_id" db:"device_id"`
	TokenHash        string    `json:"-" db:"token_hash"`
	RefreshTokenHash string    `json:"-" db:"refresh_token_hash"`
	ExpiresAt        time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

type DevicePairRequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

type DevicePairResponse struct {
	DeviceID     uuid.UUID `json:"device_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
}

type DeviceTokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type DeviceTokenRefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type DeviceUpdateRequest struct {
	DeviceName *string `json:"device_name,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}
