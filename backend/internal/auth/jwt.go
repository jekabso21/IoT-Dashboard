package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/config"
	"golang.org/x/crypto/bcrypt"
)

type UserClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type DeviceClaims struct {
	DeviceID uuid.UUID `json:"device_id"`
	UserID   uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateUserToken(userID uuid.UUID, email string) (string, error) {
	claims := UserClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.AppConfig.JWT.ExpiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.SecretKey))
}

func GenerateDeviceToken(deviceID, userID uuid.UUID) (accessToken, refreshToken string, expiresIn int, err error) {
	expiresIn = 24 * 60 * 60

	accessClaims := DeviceClaims{
		DeviceID: deviceID,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = token.SignedString([]byte(config.AppConfig.JWT.SecretKey))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err = GenerateRandomToken(32)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, expiresIn, nil
}

func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func ValidateUserToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func ValidateDeviceToken(tokenString string) (*DeviceClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &DeviceClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*DeviceClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
