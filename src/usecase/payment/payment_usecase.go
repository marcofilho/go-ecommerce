package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
)

type PaymentService interface {
	ProcessWebhook(ctx context.Context, req *entity.PaymentWebhookRequest) error
	GetWebhookHistory(ctx context.Context, orderID string) ([]entity.WebhookLog, error)
}

type PaymentUseCase struct {
	orderRepo   repository.OrderRepository
	webhookRepo repository.WebhookRepository
}

func NewPaymentUseCase(
	orderRepo repository.OrderRepository,
	webhookRepo repository.WebhookRepository,
) *PaymentUseCase {
	return &PaymentUseCase{
		orderRepo:   orderRepo,
		webhookRepo: webhookRepo,
	}
}

func (uc *PaymentUseCase) ProcessWebhook(ctx context.Context, req *entity.PaymentWebhookRequest) error {
	if req.TransactionID == "" {
		return errors.New("transaction_id is required")
	}

	existingLogs, err := uc.webhookRepo.GetByOrderID(ctx, req.OrderID)
	if err == nil {
		for _, log := range existingLogs {
			if log.TransactionID == req.TransactionID {
				return nil
			}
		}
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return errors.New("invalid order_id format")
	}

	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.Status != entity.Pending {
		return fmt.Errorf("order status must be 'pending' to process payment, current status: %s", order.Status)
	}

	if req.PaymentStatus != entity.Paid && req.PaymentStatus != entity.Failed {
		return errors.New("payment_status must be either 'paid' or 'failed'")
	}

	// Create webhook log first with pending status
	rawPayload, _ := json.Marshal(req)
	now := time.Now()
	webhookLog := &entity.WebhookLog{
		ID:            uuid.New(),
		OrderID:       orderID,
		TransactionID: req.TransactionID,
		PaymentStatus: req.PaymentStatus,
		Status:        entity.WebhookStatusProcessing,
		RetryCount:    0,
		RawPayload:    string(rawPayload),
		CreatedAt:     now,
	}

	if err := uc.webhookRepo.Create(ctx, webhookLog); err != nil {
		return fmt.Errorf("Failed to create webhook log: %w", err)
	}

	order.PaymentStatus = req.PaymentStatus

	if req.PaymentStatus == entity.Paid {
		order.Status = entity.Completed
	}

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		// In case something wrong happened, mark webhook as failed
		webhookLog.Status = entity.WebhookStatusFailed
		webhookLog.RetryCount++
		nextRetry := time.Now().Add(5 * time.Minute)
		webhookLog.NextRetryAt = &nextRetry
		uc.webhookRepo.Update(ctx, webhookLog)
		return fmt.Errorf("Failed to update order: %w", err)
	}

	// Mark webhook as completed
	webhookLog.Status = entity.WebhookStatusCompleted
	webhookLog.ProcessedAt = &now
	if err := uc.webhookRepo.Update(ctx, webhookLog); err != nil {
		fmt.Printf("Failed to update webhook log status: %v\n", err)
	}

	return nil
}

func (uc *PaymentUseCase) GetWebhookHistory(ctx context.Context, orderID string) ([]entity.WebhookLog, error) {
	return uc.webhookRepo.GetByOrderID(ctx, orderID)
}
