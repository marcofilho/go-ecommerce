package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
	"gorm.io/gorm"
)

type AuditLogRepositoryPostgres struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) repository.AuditLogRepository {
	return &AuditLogRepositoryPostgres{db: db}
}

func (r *AuditLogRepositoryPostgres) Create(ctx context.Context, log *entity.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AuditLogRepositoryPostgres) List(ctx context.Context, filters repository.AuditLogFilters, page, pageSize int) ([]*entity.AuditLog, int, error) {
	var logs []*entity.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.AuditLog{})

	// Apply filters
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.Action != nil {
		query = query.Where("action = ?", *filters.Action)
	}
	if filters.ResourceType != nil {
		query = query.Where("resource_type = ?", *filters.ResourceType)
	}
	if filters.ResourceID != nil {
		query = query.Where("resource_id = ?", *filters.ResourceID)
	}
	if filters.StartDate != nil {
		startTime, err := time.Parse(time.RFC3339, *filters.StartDate)
		if err == nil {
			query = query.Where("timestamp >= ?", startTime)
		}
	}
	if filters.EndDate != nil {
		endTime, err := time.Parse(time.RFC3339, *filters.EndDate)
		if err == nil {
			query = query.Where("timestamp <= ?", endTime)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Order("timestamp DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, int(total), nil
}

func (r *AuditLogRepositoryPostgres) GetByResourceID(ctx context.Context, resourceType string, resourceID uuid.UUID) ([]*entity.AuditLog, error) {
	var logs []*entity.AuditLog
	err := r.db.WithContext(ctx).
		Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).
		Order("timestamp DESC").
		Find(&logs).Error
	return logs, err
}
