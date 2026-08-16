package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	Id         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BusinessId uuid.UUID      `gorm:"type:uuid;index;not null" json:"business_id"`
	Name       string         `gorm:"type:varchar(255);not null" json:"name"`
	Desc       string         `gorm:"type:varchar(255);not null" json:"desc"`
	Duration   int            `gorm:"type:integer;not null;check:duration > 0" json:"duration"`
	Charge     int            `gorm:"type:integer;not null;check:charge > 0" json:"charge"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
