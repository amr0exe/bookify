package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/amr0exe/bookify/internal/tokener"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	CreateAccount(ctx context.Context, req CreateAccount) (*CreateAccResp, error)
	LoginToAccount(ctx context.Context, email string) (*models.Account, *tokener.TokenPair, error)
	RefreshTheToken(ctx context.Context, token string, claims *tokener.Claims) (*tokener.TokenPair, error)

	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

type authService struct {
	accRepo   AccountRepository
	rfrRepo   RefreshTknRepository
	jwtSecret string
}

func NewAuthService(accRepo AccountRepository, rfrRepo RefreshTknRepository, jwtSecret string) AuthService {
	return &authService{accRepo: accRepo, rfrRepo: rfrRepo, jwtSecret: jwtSecret}
}

func (a *authService) CreateAccount(ctx context.Context, req CreateAccount) (*CreateAccResp, error) {
	// check for duplicate email
	_, err := a.accRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed checking email: %w", err)
	}

	// hash password
	hashed_pwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	account := models.Account{
		Id:       uuid.New(),
		Email:    req.Email,
		PassHash: string(hashed_pwd),
		Role:     req.Role,
	}

	// generate tokens
	tokens, err := tokener.GenerateTokenPair(&account, a.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("Token creation failed")
	}

	// hash refresh token
	hash := sha256.Sum256([]byte(tokens.RefreshToken))
	refrToken := models.RefreshToken{
		AccountId: account.Id,
		TokenHash: hex.EncodeToString(hash[:]),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	// create account, store refresh token
	err = a.rfrRepo.InitialStore(ctx, &account, &refrToken)
	if err != nil {
		return nil, fmt.Errorf("Failed creating account/storing-token")
	}

	return &CreateAccResp{
		Account: &account,
		Tokens: &tokener.TokenPair{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}, nil
}

func (a *authService) LoginToAccount(ctx context.Context, email string) (*models.Account, *tokener.TokenPair, error) {
	// query for user's email
	usr, err := a.accRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, ErrAccountNotFound
	}

	// generate tokens for user
	tokens, err := tokener.GenerateTokenPair(&usr, a.jwtSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("failed generating tokens: %w", err)
	}

	refrHash := tokener.HashToken(tokens.RefreshToken)

	refrTkn := models.RefreshToken{
		Id:        uuid.New(),
		AccountId: usr.Id,
		TokenHash: refrHash,
		IsRevoked: false,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := a.rfrRepo.Store(ctx, refrTkn); err != nil {
		return nil, nil, fmt.Errorf("failed storing refr tokens to db: %w", err)
	}

	return &usr, tokens, nil
}

// create new access_token, rotate new refresh_token
func (a *authService) RefreshTheToken(ctx context.Context, token string, claims *tokener.Claims) (*tokener.TokenPair, error) {
	// fetch for stored refresh token associated with current_user
	token_hash := tokener.HashToken(token)
	refreshToken, err := a.rfrRepo.FetchValid(ctx, claims.AccountId, token_hash)
	if err != nil {
		return nil, err
	}

	// see if token provider by user matches the stored one
	if refreshToken.TokenHash != token_hash {
		return nil, fmt.Errorf("token's don't match")
	}

	// craft current user
	usr := &models.Account{
		Id:   claims.AccountId,
		Role: models.Role(claims.Role),
	}

	// generate new tokens
	tknPair, err := tokener.GenerateTokenPair(usr, a.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed generating tokens")
	}

	// revoke old and store new
	tokenHash := tokener.HashToken(tknPair.RefreshToken)

	refrTkn := &models.RefreshToken{
		Id:        uuid.New(),
		AccountId: usr.Id,
		TokenHash: tokenHash,
		IsRevoked: false,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := a.rfrRepo.Rotate(ctx, refreshToken.Id, refrTkn); err != nil {
		return nil, err
	}

	return tknPair, nil
}

func (a *authService) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	rows_affected, err := a.accRepo.Delete(ctx, id)
	// Confirm deletions, by checking for affected rows
	if err != nil {
		return err
	}
	if rows_affected == 0 {
		return ErrAccountNotFound
	}
	return nil
}
