package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Role string

const (
	RoleBusiness Role = "BUSINESS"
	RoleConsumer Role = "CONSUMER"
)

type Account struct {
	Id       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email    string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PassHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role     Role      `gorm:"type:varchar(20);not null;check:role IN ('BUSINESS', 'CONSUMER')" json:"role"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Business *Business `gorm:"foreignKey:AccountId;constraint:OnDelete:CASCADE;" json:"business,omitempty"`
	Consumer *Consumer `gorm:"foreignKey:AccountId;constraint:OnDelete:CASCADE;" json:"consumer,omitempty"`
}

type RefreshToken struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountId uuid.UUID `gorm:"type:uuid;not null;index" json:"account_id"`
	TokenHash string    `gorm:"not null;uniqueIndex" json:"-"`
	IsRevoked bool      `gorm:"default:false;not null" json:"is_revoked"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
