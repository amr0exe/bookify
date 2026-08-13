package consumer

import (
	"context"

	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/google/uuid"
)

type ConsumerService interface {
	CreateConsumer(ctx context.Context, accountID uuid.UUID, consumerInfo *CreateConsumer) error
}

type consumerService struct {
	repo ConsumerRepository
}

func NewConsumerService(repo ConsumerRepository) ConsumerService {
	return &consumerService{repo: repo}
}

func (s *consumerService) CreateConsumer(ctx context.Context, accountID uuid.UUID, consumerInfo *CreateConsumer) error {
	info := &models.Consumer{
		AccountId:   accountID,
		Name:        consumerInfo.Name,
		DisplayName: consumerInfo.DisplayName,
		Phone:       consumerInfo.Phone,
	}

	return s.repo.Create(ctx, info)
}
