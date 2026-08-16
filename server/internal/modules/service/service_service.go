package service

import (
	"context"

	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/google/uuid"
)

type ServiceService interface {
	CreateService(ctx context.Context, accId uuid.UUID, info *CreateService) (*models.Service, error)
	GetAllServices(ctx context.Context, accId uuid.UUID) ([]models.Service, error)
	UpdateService(ctx context.Context, accId uuid.UUID, serviceId uuid.UUID, info *UpdateService) (*models.Service, error)
	DeleteService(ctx context.Context, accId uuid.UUID, serviceId uuid.UUID) error
}

type serviceService struct {
	repo ServiceRepository
}

func NewServiceService(repo ServiceRepository) ServiceService {
	return &serviceService{repo: repo}
}

func (s *serviceService) CreateService(ctx context.Context, accId uuid.UUID, info *CreateService) (*models.Service, error) {
	businessId, err := s.repo.FetchBusinessId(ctx, accId)
	if err != nil {
		return nil, err
	}

	newService := &models.Service{
		BusinessId: businessId,
		Name:       info.Name,
		Desc:       info.Desc,
		Duration:   int(info.Duration),
		Charge:     int(info.Charge),
	}
	if err := s.repo.Create(ctx, businessId, newService); err != nil {
		return nil, err
	}
	return newService, nil
}

func (s *serviceService) GetAllServices(ctx context.Context, accId uuid.UUID) ([]models.Service, error) {
	businessId, err := s.repo.FetchBusinessId(ctx, accId)
	if err != nil {
		return nil, err
	}

	return s.repo.FindAll(ctx, businessId)
}

func (s *serviceService) UpdateService(ctx context.Context, accId uuid.UUID, serviceId uuid.UUID, info *UpdateService) (*models.Service, error) {
	businessId, err := s.repo.FetchBusinessId(ctx, accId)
	if err != nil {
		return nil, err
	}

	update := &models.Service{
		Name:     info.Name,
		Desc:     info.Desc,
		Duration: int(info.Duration),
		Charge:   int(info.Charge),
	}

	if err := s.repo.Update(ctx, businessId, serviceId, update); err != nil {
		return nil, err
	}

	service, err := s.repo.FindById(ctx, businessId, serviceId)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (s *serviceService) DeleteService(ctx context.Context, accId uuid.UUID, serviceId uuid.UUID) error {
	businessId, err := s.repo.FetchBusinessId(ctx, accId)
	if err != nil {
		return err
	}

	if err := s.repo.Remove(ctx, businessId, serviceId); err != nil {
		return err
	}
	return nil
}
