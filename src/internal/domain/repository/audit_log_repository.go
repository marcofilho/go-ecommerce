package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
)

type AuditLogRepository interface {
	// Create creates a new audit log entry
	Create(ctx context.Context, log *entity.AuditLog) error

	// List returns audit logs with optional filters
	List(ctx context.Context, filters AuditLogFilters, page, pageSize int) ([]*entity.AuditLog, int, error)

	// GetByResourceID returns all audit logs for a specific resource
	GetByResourceID(ctx context.Context, resourceType string, resourceID uuid.UUID) ([]*entity.AuditLog, error)
}

type AuditLogFilters struct {
	UserID       *uuid.UUID
	Action       *string
	ResourceType *string
	ResourceID   *uuid.UUID
	StartDate    *string
	EndDate      *string
}
