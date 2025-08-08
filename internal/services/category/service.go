// Package category provides business logic for category operations.
// It includes service layer functionality for category management with Redis caching.
package category

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/category/model"
	"github.com/MitulShah1/golang-rest-api-template/internal/repository"
	sqlModel "github.com/MitulShah1/golang-rest-api-template/internal/repository/model"
	"github.com/MitulShah1/golang-rest-api-template/package/cache"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"
)

type CategoryServiceInterface interface {
	CreateCategory(ctx context.Context, category model.CreateCategoryRequest) (int64, error)
	GetCategoryByID(ctx context.Context, id int) (*sqlModel.Category, error)
	UpdateCategory(ctx context.Context, id int, category model.UpdateCategoryRequest) error
	DeleteCategory(ctx context.Context, id int) error
}

type CategoryService struct {
	repo   repository.DBRepository
	logger *logger.Logger
	cache  *cache.Cache
}

func NewCategoryService(repo repository.DBRepository, logger *logger.Logger, cache *cache.Cache) CategoryServiceInterface {
	return &CategoryService{
		repo:   repo,
		logger: logger,
		cache:  cache,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category model.CreateCategoryRequest) (int64, error) {
	s.logger.Info("Creating category", "category", category)

	cat := sqlModel.Category{
		Name:        category.Name,
		ParentID:    category.ParentID,
		Description: category.Description,
	}

	id, err := s.repo.CreateCategory(ctx, &cat)
	if err != nil {
		return 0, err
	}

	// Invalidate category cache patterns
	s.invalidateCategoryCache(ctx)

	return id, nil
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (*sqlModel.Category, error) {
	s.logger.Info("Getting category by ID", "id", id)

	// Try to get from cache first
	cacheKey := fmt.Sprintf("category:%d", id)
	var cachedCategory sqlModel.Category

	if err := s.cache.Get(ctx, cacheKey, &cachedCategory); err == nil {
		s.logger.Debug("category retrieved from cache", "category_id", id)
		return &cachedCategory, nil
	}

	// Cache miss, get from database
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		s.logger.Error("error while fetch category information", err)
		return nil, err
	}

	if category == nil {
		s.logger.Warn("category not found", "category id", id)
		return nil, errors.New("category not found")
	}

	// Cache the result for future requests
	if err := s.cache.Set(ctx, cacheKey, category, 30*time.Minute); err != nil {
		s.logger.Warn("failed to cache category", "category_id", id, "error", err)
	}

	return category, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id int, category model.UpdateCategoryRequest) error {
	s.logger.Info("Updating category", "category", category)

	updCat := sqlModel.Category{
		Name:        category.Name,
		ParentID:    category.ParentID,
		Description: category.Description,
	}

	err := s.repo.UpdateCategory(ctx, id, &updCat)
	if err != nil {
		return err
	}

	// Invalidate category cache patterns
	s.invalidateCategoryCache(ctx)
	if err := s.cache.Delete(ctx, fmt.Sprintf("category:%d", id)); err != nil {
		s.logger.Warn("failed to delete category cache", "category_id", id, "error", err)
	}

	return nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	s.logger.Info("Deleting category", "id", id)

	err := s.repo.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate category cache patterns
	s.invalidateCategoryCache(ctx)
	if err := s.cache.Delete(ctx, fmt.Sprintf("category:%d", id)); err != nil {
		s.logger.Warn("failed to delete category cache", "category_id", id, "error", err)
	}

	return nil
}

// invalidateCategoryCache removes all category-related cache entries
func (s *CategoryService) invalidateCategoryCache(ctx context.Context) {
	// Delete all category cache patterns
	if err := s.cache.DeletePattern(ctx, "category:*"); err != nil {
		s.logger.Warn("failed to invalidate category cache", "error", err)
	}
}
