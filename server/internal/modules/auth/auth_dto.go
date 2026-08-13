package auth

import (
	"time"

	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/amr0exe/bookify/internal/tokener"
	"github.com/google/uuid"
)

// request structure used for registration
type CreateAccount struct {
	Email    string      `json:"email" binding:"required,email,max=255"`
	Password string      `json:"password" binding:"required,min=8,max=72"`
	Role     models.Role `json:"role" binding:"required,oneof=BUSINESS CONSUMER"`
}

// response structure used for registration
type AccountResponse struct {
	Id        uuid.UUID          `json:"id"`
	Email     string             `json:"email"`
	Role      string             `json:"role"`
	Tokens    *tokener.TokenPair `json:"tokens"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// Maps models.Account to an AccountResponse dto
func ToAccountResponse(a *models.Account, tokens *tokener.TokenPair) AccountResponse {
	return AccountResponse{
		Id:        a.Id,
		Email:     a.Email,
		Role:      string(a.Role),
		Tokens:    tokens,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

type CreateAccResp struct {
	Account *models.Account
	Tokens  *tokener.TokenPair
}

// Login req/res structure
type LoginAccount struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

// refresh-token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
