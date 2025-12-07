package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/database"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/models"
)

func CreateUser(ctx context.Context, email, passwordHash, fullName string) (*models.User, error) {
	user := &models.User{
		ID:    uuid.New(),
		Email: email,
	}

	var company sql.NullString
	err := database.DB.QueryRowContext(ctx,
		`INSERT INTO users (id, email, password_hash, full_name, company)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING created_at, updated_at, full_name, company`,
		user.ID, email, passwordHash, fullName, nil,
	).Scan(&user.CreatedAt, &user.UpdatedAt, &user.FullName, &company)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if company.Valid {
		user.Company = &company.String
	}

	return user, nil
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, string, error) {
	var user models.User
	var passwordHash string
	var fullName, company sql.NullString

	err := database.DB.QueryRowContext(ctx,
		`SELECT id, email, password_hash, full_name, company, email_verified, created_at, updated_at
		 FROM users WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&passwordHash,
		&fullName,
		&company,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, "", fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	if fullName.Valid {
		user.FullName = &fullName.String
	}
	if company.Valid {
		user.Company = &company.String
	}

	return &user, passwordHash, nil
}

func GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	var fullName, company sql.NullString

	err := database.DB.QueryRowContext(ctx,
		`SELECT id, email, full_name, company, email_verified, created_at, updated_at
		 FROM users WHERE id = $1`,
		userID,
	).Scan(
		&user.ID,
		&user.Email,
		&fullName,
		&company,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if fullName.Valid {
		user.FullName = &fullName.String
	}
	if company.Valid {
		user.Company = &company.String
	}

	return &user, nil
}

func EmailExists(ctx context.Context, email string) (bool, error) {
	var count int
	err := database.DB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM users WHERE email = $1`,
		email,
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}
