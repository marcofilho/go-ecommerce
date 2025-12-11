package testing

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/audit"
)

// MockServices implements the Services interface for testing
type MockServices struct {
	AuditService audit.AuditService
}

func (m *MockServices) GetAuditService() audit.AuditService {
	if m.AuditService != nil {
		return m.AuditService
	}
	return &MockAuditService{}
}

// MockAuditService is a mock implementation of audit.AuditService
type MockAuditService struct{}

func (m *MockAuditService) LogChange(ctx context.Context, userID *uuid.UUID, action, resourceType string, resourceID uuid.UUID, before, after interface{}) error {
	return nil
}
