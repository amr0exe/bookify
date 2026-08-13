package business

import (
	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/google/uuid"
)

type CreateBusiness struct {
	AccountId   uuid.UUID             `json:"account_id" binding:"omitempty"`
	Name        string                `json:"name" binding:"required,min=3,max=100"`
	Description string                `json:"description" binding:"omitempty,min=10,max=100"`
	Phone       string                `json:"phone" binding:"omitempty,e164"`
	Address     string                `json:"address" binding:"required,min=5,max=100"`
	Latitude    *float64              `json:"latitude" binding:"omitempty,latitude"`
	Longitude   *float64              `json:"longitude" binding:"omitempty,longitude"`
	Status      models.BusinessStatus `json:"status" binding:"required,oneof=OPEN CLOSED"`
}
