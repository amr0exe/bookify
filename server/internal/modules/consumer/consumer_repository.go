package consumer

import (
	"context"
	"errors"

	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConsumerRepository interface {
	Create(ctx context.Context, consumerInfo *models.Consumer) error
	Fetch(ctx context.Context, accountID uuid.UUID) (*models.Consumer, error)
}

type consumerRepository struct {
	db *gorm.DB
}

func NewConsumerRepository(db *gorm.DB) ConsumerRepository {
	return &consumerRepository{db: db}
}

func (r *consumerRepository) Create(ctx context.Context, consumerInfo *models.Consumer) error {
	err := gorm.G[models.Consumer](r.db).Create(ctx, consumerInfo)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrConsumerAlreadyExists
	}
	return err
}

func (r *consumerRepository) Fetch(ctx context.Context, accountID uuid.UUID) (*models.Consumer, error) {
	var consumer models.Consumer

	consumer, err := gorm.G[models.Consumer](r.db).Where("account_id = ?", accountID).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrConsumerNotFound
	}

	if err != nil {
		return nil, err
	}

	return &consumer, nil
}
