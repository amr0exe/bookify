package business

import (
	"context"

	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/google/uuid"
)

type BusinessService interface {
	CreateProfile(ctx context.Context, accountID uuid.UUID, businessInfo *CreateBusiness) error
}

type businessService struct {
	repo BusinessRepository
}

func NewBusinessService(repo BusinessRepository) BusinessService {
	return &businessService{repo: repo}
}

func (s *businessService) CreateProfile(ctx context.Context, accountID uuid.UUID, info *CreateBusiness) error {
	newBusiness := &models.Business{
		AccountId:   accountID,
		Name:        info.Name,
		Description: info.Description,
		Phone:       info.Phone,
		Address:     info.Address,
		Latitude:    info.Latitude,
		Longitude:   info.Longitude,
		Status:      info.Status,
	}

	return s.repo.Create(ctx, newBusiness)
}
