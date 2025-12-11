package audit

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
	"gorm.io/datatypes"
)

// AuditService handles audit logging for entity changes
type AuditService interface {
	LogChange(ctx context.Context, userID *uuid.UUID, action, resourceType string, resourceID uuid.UUID, before, after interface{}) error
}

type auditService struct {
	repo repository.AuditLogRepository
}

func NewAuditService(repo repository.AuditLogRepository) AuditService {
	return &auditService{repo: repo}
}

func (s *auditService) LogChange(ctx context.Context, userID *uuid.UUID, action, resourceType string, resourceID uuid.UUID, before, after interface{}) error {
	var payloadBefore, payloadAfter datatypes.JSON

	// Convert before payload to JSON
	if before != nil {
		beforeBytes, err := json.Marshal(before)
		if err != nil {
			return err
		}
		payloadBefore = datatypes.JSON(beforeBytes)
	}

	// Convert after payload to JSON
	if after != nil {
		afterBytes, err := json.Marshal(after)
		if err != nil {
			return err
		}
		payloadAfter = datatypes.JSON(afterBytes)
	}

	// Create audit log entry
	log := &entity.AuditLog{
		UserID:        userID,
		Action:        action,
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		PayloadBefore: payloadBefore,
		PayloadAfter:  payloadAfter,
	}

	return s.repo.Create(ctx, log)
}
