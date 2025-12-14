package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/cache"
)

// CachedProductRepository wraps a ProductRepository with Redis caching
type CachedProductRepository struct {
	repo  repository.ProductRepository
	cache *cache.RedisClient
}

func NewCachedProductRepository(repo repository.ProductRepository, cache *cache.RedisClient) repository.ProductRepository {
	return &CachedProductRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *CachedProductRepository) Create(ctx context.Context, product *entity.Product) error {
	if err := r.repo.Create(ctx, product); err != nil {
		return err
	}
	
	_ = r.cache.DeleteByPattern(ctx, "products:list:*")
	
	return nil
}

func (r *CachedProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	cacheKey := fmt.Sprintf("product:id:%s", id.String())
	
	cached, err := r.cache.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var product entity.Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			return &product, nil
		}
	}
	
	product, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if data, err := json.Marshal(product); err == nil {
		_ = r.cache.Set(ctx, cacheKey, data)
	}
	
	return product, nil
}

func (r *CachedProductRepository) GetAll(ctx context.Context, page, pageSize int, inStockOnly bool) ([]*entity.Product, int, error) {
	cacheKey := fmt.Sprintf("products:list:page:%d:size:%d:instock:%t", page, pageSize, inStockOnly)
	
	cached, err := r.cache.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var result struct {
			Products []*entity.Product
			Total    int
		}
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result.Products, result.Total, nil
		}
	}
	
	products, total, err := r.repo.GetAll(ctx, page, pageSize, inStockOnly)
	if err != nil {
		return nil, 0, err
	}
	
	result := struct {
		Products []*entity.Product
		Total    int
	}{
		Products: products,
		Total:    total,
	}
	if data, err := json.Marshal(result); err == nil {
		_ = r.cache.Set(ctx, cacheKey, data)
	}
	
	return products, total, nil
}

func (r *CachedProductRepository) Update(ctx context.Context, product *entity.Product) error {
	if err := r.repo.Update(ctx, product); err != nil {
		return err
	}
	
	cacheKey := fmt.Sprintf("product:id:%s", product.ID.String())
	_ = r.cache.Del(ctx, cacheKey)
	_ = r.cache.DeleteByPattern(ctx, "products:list:*")
	
	return nil
}

func (r *CachedProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.repo.Delete(ctx, id); err != nil {
		return err
	}
	
	cacheKey := fmt.Sprintf("product:id:%s", id.String())
	_ = r.cache.Del(ctx, cacheKey)
	_ = r.cache.DeleteByPattern(ctx, "products:list:*")
	
	return nil
}
