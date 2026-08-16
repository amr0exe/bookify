package service

import (
	"context"
	"errors"

	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	Create(ctx context.Context, businessId uuid.UUID, serviceInfo *models.Service) error
	FindAll(ctx context.Context, businessId uuid.UUID) ([]models.Service, error)
	FindById(ctx context.Context, businessId uuid.UUID, serviceId uuid.UUID) (models.Service, error)
	Update(ctx context.Context, businessId uuid.UUID, serviceId uuid.UUID, updateInfo *models.Service) error
	Remove(ctx context.Context, businessId uuid.UUID, serviceId uuid.UUID) error
	FetchBusinessId(ctx context.Context, accId uuid.UUID) (uuid.UUID, error)
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(ctx context.Context, businessId uuid.UUID, serviceInfo *models.Service) error {
	return gorm.G[models.Service](r.db).Create(ctx, serviceInfo)
}

func (r *serviceRepository) FetchBusinessId(ctx context.Context, accId uuid.UUID) (uuid.UUID, error) {
	business, err := gorm.G[models.Business](r.db).Where("account_id = ?", accId).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.UUID{}, ErrNoBusinessFound
	}
	return business.Id, nil
}

func (r *serviceRepository) FindAll(ctx context.Context, businessId uuid.UUID) ([]models.Service, error) {
	return gorm.G[models.Service](r.db).Where("business_id = ?", businessId).Find(ctx)
}

func (r *serviceRepository) FindById(ctx context.Context, businessId uuid.UUID, serviceId uuid.UUID) (models.Service, error) {
	service, err := gorm.G[models.Service](r.db).
		Where("business_id = ? AND id = ?", businessId, serviceId).
		First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Service{}, ErrNoServiceFound
	}
	if err != nil {
		return models.Service{}, err
	}
	return service, nil
}

func (r *serviceRepository) Update(ctx context.Context, businessId uuid.UUID, serviceId uuid.UUID, updateInfo *models.Service) error {
	rows, err := gorm.G[models.Service](r.db).
		Where("business_id = ? AND id = ?", businessId, serviceId).
		Updates(ctx, *updateInfo)
	if rows == 0 {
		return ErrNoServiceFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *serviceRepository) Remove(ctx context.Context, businessId uuid.UUID, serviceId uuid.UUID) error {
	rows, err := gorm.G[models.Service](r.db).Where("business_id = ? AND id = ?", businessId, serviceId).Delete(ctx)
	if rows == 0 {
		return ErrNoServiceFound
	}
	if err != nil {
		return err
	}
	return nil
}
