package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/auth"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/database"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/models"
)

func CreateDevice(ctx context.Context, userID uuid.UUID) (*models.Device, error) {
	device := &models.Device{
		ID:         uuid.New(),
		DeviceName: "New Device",
		DeviceType: "esp32_scd4x",
		IsActive:   true,
		PairedAt:   ptrTime(time.Now()),
	}

	userIDPtr := &userID

	err := database.DB.QueryRowContext(ctx,
		`INSERT INTO devices (id, user_id, device_name, device_type, is_active, paired_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING created_at, updated_at`,
		device.ID, userIDPtr, device.DeviceName, device.DeviceType, device.IsActive, device.PairedAt,
	).Scan(&device.CreatedAt, &device.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	device.UserID = userIDPtr
	return device, nil
}

func StoreDeviceTokens(ctx context.Context, deviceID, userID uuid.UUID, accessToken, refreshToken string, expiresIn int) error {
	tokenHash := auth.HashToken(accessToken)
	refreshTokenHash := auth.HashToken(refreshToken)
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	_, err := database.DB.ExecContext(ctx,
		`INSERT INTO device_tokens (device_id, token_hash, refresh_token_hash, expires_at)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (device_id) DO UPDATE
		 SET token_hash = EXCLUDED.token_hash,
		     refresh_token_hash = EXCLUDED.refresh_token_hash,
		     expires_at = EXCLUDED.expires_at,
		     created_at = NOW()`,
		deviceID, tokenHash, refreshTokenHash, expiresAt,
	)

	if err != nil {
		return fmt.Errorf("failed to store device tokens: %w", err)
	}

	return nil
}

func ValidateDeviceRefreshToken(ctx context.Context, deviceID uuid.UUID, refreshToken string) (bool, error) {
	refreshTokenHash := auth.HashToken(refreshToken)

	var count int
	err := database.DB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM device_tokens
		 WHERE device_id = $1 AND refresh_token_hash = $2`,
		deviceID, refreshTokenHash,
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	return count > 0, nil
}

func GetDeviceByID(ctx context.Context, deviceID uuid.UUID) (*models.Device, error) {
	var device models.Device

	err := database.DB.QueryRowContext(ctx,
		`SELECT id, user_id, device_name, device_type, last_seen, is_active, paired_at, created_at, updated_at
		 FROM devices WHERE id = $1`,
		deviceID,
	).Scan(
		&device.ID,
		&device.UserID,
		&device.DeviceName,
		&device.DeviceType,
		&device.LastSeen,
		&device.IsActive,
		&device.PairedAt,
		&device.CreatedAt,
		&device.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("device not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	return &device, nil
}

func UpdateDeviceLastSeen(ctx context.Context, deviceID uuid.UUID) error {
	_, err := database.DB.ExecContext(ctx,
		`UPDATE devices SET last_seen = NOW() WHERE id = $1`,
		deviceID,
	)
	if err != nil {
		return fmt.Errorf("failed to update device last_seen: %w", err)
	}
	return nil
}

func GetUserDevices(ctx context.Context, userID uuid.UUID) ([]models.Device, error) {
	rows, err := database.DB.QueryContext(ctx,
		`SELECT id, user_id, device_name, device_type, last_seen, is_active, paired_at, created_at, updated_at,
		 CASE WHEN last_seen > NOW() - INTERVAL '30 seconds' THEN true ELSE false END as online
		 FROM devices
		 WHERE user_id = $1
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user devices: %w", err)
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device
		err := rows.Scan(
			&device.ID,
			&device.UserID,
			&device.DeviceName,
			&device.DeviceType,
			&device.LastSeen,
			&device.IsActive,
			&device.PairedAt,
			&device.CreatedAt,
			&device.UpdatedAt,
			&device.Online,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device: %w", err)
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func UpdateDevice(ctx context.Context, deviceID uuid.UUID, updates *models.DeviceUpdateRequest) error {
	query := `UPDATE devices SET `
	args := []interface{}{}
	argPos := 1

	if updates.DeviceName != nil {
		query += fmt.Sprintf("device_name = $%d, ", argPos)
		args = append(args, *updates.DeviceName)
		argPos++
	}

	if updates.IsActive != nil {
		query += fmt.Sprintf("is_active = $%d, ", argPos)
		args = append(args, *updates.IsActive)
		argPos++
	}

	query += fmt.Sprintf("updated_at = NOW() WHERE id = $%d", argPos)
	args = append(args, deviceID)

	_, err := database.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}

	return nil
}

func DeleteDevice(ctx context.Context, deviceID uuid.UUID) error {
	_, err := database.DB.ExecContext(ctx, `DELETE FROM devices WHERE id = $1`, deviceID)
	if err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}
	return nil
}

func CleanupExpiredDeviceTokens(ctx context.Context) error {
	_, err := database.DB.ExecContext(ctx, "SELECT cleanup_expired_device_tokens()")
	if err != nil {
		return fmt.Errorf("failed to cleanup expired device tokens: %w", err)
	}
	return nil
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
