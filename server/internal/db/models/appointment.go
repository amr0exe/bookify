package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppointmentStatus string

const (
	AppointmentCreated   AppointmentStatus = "created"
	AppointmentRejected  AppointmentStatus = "rejected"
	AppointmentAccepted  AppointmentStatus = "accepted"
	AppointmentStarted   AppointmentStatus = "started"
	AppointmentCompleted AppointmentStatus = "completed"
)

type Appointment struct {
	Id          uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ConsumerId  uuid.UUID         `gorm:"type:uuid;index;not null" json:"consumer_id"`
	BusinessId  uuid.UUID         `gorm:"type:uuid;index;not null" json:"business_id"`
	ServiceId   uuid.UUID         `gorm:"type:uuid;index;not null" json:"service_id"`
	Status      AppointmentStatus `gorm:"type:varchar(20);not null;default:'created'" json:"status"`
	Duration    int               `gorm:"type:integer;not null;check:duration > 0" json:"duration"`
	Remarks     *string           `gorm:"type:text" json:"remarks,omitempty"`
	ScheduledAt *time.Time        `gorm:"type:timestamptz" json:"scheduled_at,omitempty"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
