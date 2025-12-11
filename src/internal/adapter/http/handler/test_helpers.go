package handler

import (
	"context"

	"github.com/google/uuid"
)

// mockAuditService is a mock implementation of audit.AuditService for testing
type mockAuditService struct{}

func (m *mockAuditService) LogChange(ctx context.Context, userID *uuid.UUID, action, resourceType string, resourceID uuid.UUID, before, after interface{}) error {
	return nil
}
