package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"gorm.io/gorm"
)

type WebhookRepositoryPostgres struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) *WebhookRepositoryPostgres {
	return &WebhookRepositoryPostgres{db: db}
}

func (r *WebhookRepositoryPostgres) Create(ctx context.Context, log *entity.WebhookLog) error {
	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *WebhookRepositoryPostgres) Update(ctx context.Context, log *entity.WebhookLog) error {
	return r.db.WithContext(ctx).Save(log).Error
}

func (r *WebhookRepositoryPostgres) GetByOrderID(ctx context.Context, orderID string) ([]entity.WebhookLog, error) {
	var logs []entity.WebhookLog
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}
