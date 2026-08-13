package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Consumer struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountId   uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"account_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	DisplayName string    `gorm:"type:varchar(255)" json:"display_name,omitempty"`
	Phone       string    `gorm:"type:varchar(60)" json:"phone,omitempty"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
