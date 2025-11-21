package repository

import (
	"context"

	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
)

type WebhookRepository interface {
	Create(ctx context.Context, log *entity.WebhookLog) error
	Update(ctx context.Context, log *entity.WebhookLog) error
	GetByOrderID(ctx context.Context, orderID string) ([]entity.WebhookLog, error)
}
