package business

import (
	"context"
	"errors"

	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessRepository interface {
	Create(ctx context.Context, businessInfo *models.Business) error
	Fetch(ctx context.Context, accountID uuid.UUID) (models.Business, error)
}

type businessRepository struct {
	db *gorm.DB
}

func NewBusinessRepository(db *gorm.DB) BusinessRepository {
	return &businessRepository{db: db}
}

func (r *businessRepository) Create(ctx context.Context, businessInfo *models.Business) error {
	err := gorm.G[models.Business](r.db).Create(ctx, businessInfo)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrBusinessAlreadyExists
	}
	return err
}

func (r *businessRepository) Fetch(ctx context.Context, accountID uuid.UUID) (models.Business, error) {
	return gorm.G[models.Business](r.db).Where("account_id = ?", accountID).First(ctx)
}
