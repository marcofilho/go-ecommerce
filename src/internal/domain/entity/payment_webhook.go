package entity

import (
	"time"

	"github.com/google/uuid"
)

// PaymentWebhookRequest represents a simplified payment webhook payload
type PaymentWebhookRequest struct {
	OrderID       string        `json:"order_id"`
	TransactionID string        `json:"transaction_id"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	Timestamp     int64         `json:"timestamp"`
}

// WebhookStatus represents the processing status of a webhook
type WebhookStatus string

const (
	WebhookStatusPending    WebhookStatus = "pending"
	WebhookStatusProcessing WebhookStatus = "processing"
	WebhookStatusCompleted  WebhookStatus = "completed"
	WebhookStatusFailed     WebhookStatus = "failed"
)

// WebhookLog stores webhook events for audit
type WebhookLog struct {
	ID            uuid.UUID     `gorm:"type:uuid;primaryKey"`
	OrderID       uuid.UUID     `gorm:"type:uuid;not null;index"`
	TransactionID string        `gorm:"type:varchar(255);not null;uniqueIndex"`
	PaymentStatus PaymentStatus `gorm:"type:varchar(20);not null"`
	Status        WebhookStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	RetryCount    int           `gorm:"default:0"`
	NextRetryAt   *time.Time
	RawPayload    string `gorm:"type:text"`
	ProcessedAt   *time.Time
	CreatedAt     time.Time
}
