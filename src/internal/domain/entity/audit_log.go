package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	UserID        *uuid.UUID     `gorm:"type:uuid;index"` // Nullable for system actions
	Action        string         `gorm:"size:100;not null;index"`
	ResourceType  string         `gorm:"size:100;not null;index"`
	ResourceID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	PayloadBefore datatypes.JSON `gorm:"type:jsonb"`
	PayloadAfter  datatypes.JSON `gorm:"type:jsonb"`
	Timestamp     time.Time      `gorm:"not null;index"`
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}
	return nil
}

func (a *AuditLog) TableName() string {
	return "audit_logs"
}
