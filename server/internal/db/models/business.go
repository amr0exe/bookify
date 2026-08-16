package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessStatus string

const (
	StatusOpen   BusinessStatus = "OPEN"
	StatusClosed BusinessStatus = "CLOSED"
)

type Business struct {
	Id          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountId   uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"account_id"`
	Name        string         `gorm:"type:varchar(254);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	Phone       string         `gorm:"type:varchar(50);not null" json:"phone"`
	Address     string         `gorm:"type:text;not null" json:"address"`
	Latitude    *float64       `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude   *float64       `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	Status      BusinessStatus `gorm:"type:varchar(20);not null;default:'OPEN'" json:"status"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Service []Service `gorm:"foreignKey:BusinessId;constraint:OnDelete:CASCADE;" json:"services,omitempty"`
}
