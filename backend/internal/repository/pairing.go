package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/database"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/models"
)

func GeneratePairingCode(ctx context.Context, userID uuid.UUID) (*models.PairingCode, error) {
	var code string
	err := database.DB.QueryRowContext(ctx, "SELECT generate_pairing_code()").Scan(&code)
	if err != nil {
		return nil, fmt.Errorf("failed to generate pairing code: %w", err)
	}

	expiresAt := time.Now().Add(10 * time.Minute)

	_, err = database.DB.ExecContext(ctx,
		`INSERT INTO pairing_codes (code, user_id, expires_at) VALUES ($1, $2, $3)`,
		code, userID, expiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert pairing code: %w", err)
	}

	return &models.PairingCode{
		Code:      code,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Used:      false,
		CreatedAt: time.Now(),
	}, nil
}

func ValidatePairingCode(ctx context.Context, code string) (*models.PairingCode, error) {
	var pairingCode models.PairingCode

	err := database.DB.QueryRowContext(ctx,
		`SELECT code, user_id, expires_at, used, used_at, created_at
		 FROM pairing_codes
		 WHERE code = $1 AND used = FALSE AND expires_at > NOW()`,
		code,
	).Scan(
		&pairingCode.Code,
		&pairingCode.UserID,
		&pairingCode.ExpiresAt,
		&pairingCode.Used,
		&pairingCode.UsedAt,
		&pairingCode.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid or expired pairing code")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to validate pairing code: %w", err)
	}

	return &pairingCode, nil
}

func MarkPairingCodeAsUsed(ctx context.Context, code string) error {
	_, err := database.DB.ExecContext(ctx,
		`UPDATE pairing_codes SET used = TRUE, used_at = NOW() WHERE code = $1`,
		code,
	)
	if err != nil {
		return fmt.Errorf("failed to mark pairing code as used: %w", err)
	}
	return nil
}

func CleanupExpiredPairingCodes(ctx context.Context) error {
	_, err := database.DB.ExecContext(ctx, "SELECT cleanup_expired_pairing_codes()")
	if err != nil {
		return fmt.Errorf("failed to cleanup expired pairing codes: %w", err)
	}
	return nil
}
