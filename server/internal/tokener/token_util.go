package tokener

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("Invalid or expired token.")
	ErrExpiredToken = errors.New("Token has expired.")
)

type Claims struct {
	AccountId uuid.UUID   `json:"account_id"`
	Role      models.Role `json:"role"`
	TokenType string      `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GenerateTokenPair, creates a 15_mins access token and 7_days refresh token
func GenerateTokenPair(acc *models.Account, secret string) (*TokenPair, error) {
	// Generate 15_mins expiry access token
	accessClaims := Claims{
		AccountId: acc.Id,
		Role:      acc.Role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	// CREATE token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	// SIGN token with JWT_SECRET
	SignedAccessToken, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	// Generate 7_days expiry refresh token
	refreshClaims := Claims{
		AccountId: acc.Id,
		Role:      acc.Role,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	// CREATE token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	// Sign token
	refreshSigned, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  SignedAccessToken,
		RefreshToken: refreshSigned,
	}, nil
}

// ValidateToken, verifies the signature and returns the token claims/token-payload
func ValidateToken(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// hash refreshToken with sha256
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
