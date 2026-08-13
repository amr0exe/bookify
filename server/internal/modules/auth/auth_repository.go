package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ---

type AccountRepository interface {
	GetByEmail(ctx context.Context, email string) (models.Account, error)
	Delete(ctx context.Context, id uuid.UUID) (int, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) GetByEmail(ctx context.Context, email string) (models.Account, error) {
	return gorm.G[models.Account](r.db).Where("email = ?", email).First(ctx)
}

func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID) (int, error) {
	return gorm.G[models.Account](r.db).Select("Consumer", "Business").Where("id = ?", id).Delete(ctx)
}

// ---

type RefreshTknRepository interface {
	// stores refresh token
	Store(ctx context.Context, refrTkn models.RefreshToken) error
	// fetches refreshToken that's valid which hasn't been either revoked or expired for given accountId and tokenHash
	FetchValid(ctx context.Context, account_id uuid.UUID, token_hash string) (models.RefreshToken, error)
	// creates newAccount and stores initial refresh token on register
	InitialStore(ctx context.Context, account *models.Account, token_info *models.RefreshToken) error
	// revokes old/used refreshToken and stores new one
	Rotate(ctx context.Context, oldTknId uuid.UUID, newTkn *models.RefreshToken) error
}

type refreshTknRepository struct {
	db *gorm.DB
}

func NewRefreshTknRepository(db *gorm.DB) RefreshTknRepository {
	return &refreshTknRepository{db: db}
}

// stores refresh token
func (r *refreshTknRepository) Store(ctx context.Context, refrTkn models.RefreshToken) error {
	return gorm.G[models.RefreshToken](r.db).Create(ctx, &refrTkn)
}

// fetches refreshToken that's valid which hasn't been either revoked or expired for given accountId and tokenHash
func (r *refreshTknRepository) FetchValid(ctx context.Context, account_id uuid.UUID, token_hash string) (models.RefreshToken, error) {
	return gorm.G[models.RefreshToken](r.db).
		Where("account_id = ? AND token_hash = ? AND is_revoked = ? AND expires_at > ?", account_id, token_hash, false, time.Now()).
		First(ctx)
}

// creates newAccount and stores initial refresh token on register
func (r *refreshTknRepository) InitialStore(ctx context.Context, account *models.Account, token_info *models.RefreshToken) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[models.Account](tx).Create(ctx, account); err != nil {
			return err
		}

		token_info.AccountId = account.Id
		if err := gorm.G[models.RefreshToken](tx).Create(ctx, token_info); err != nil {
			return err
		}

		return nil
	})
}

// revokes old/used refreshToken and stores new one
func (r *refreshTknRepository) Rotate(ctx context.Context, oldTknId uuid.UUID, newTkn *models.RefreshToken) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// revoke old token
		affectedRows, err := gorm.G[models.RefreshToken](tx).Where("id = ?", oldTknId).Update(ctx, "is_revoked", true)
		if err != nil {
			return err
		}

		if affectedRows == 0 {
			return fmt.Errorf("revoking old token failed, token has already been revoked or used")
		}

		// store new token
		if err := gorm.G[models.RefreshToken](tx).Create(ctx, newTkn); err != nil {
			return err
		}

		return nil
	})
}
