// Package product provides business logic for product operations.
// It includes service layer functionality for product management with Redis caching.
package product

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/product/model"
	"github.com/MitulShah1/golang-rest-api-template/internal/repository"
	sqlModel "github.com/MitulShah1/golang-rest-api-template/internal/repository/model"
	"github.com/MitulShah1/golang-rest-api-template/package/cache"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"
)

var ErrProductNotFound = errors.New("product not found")

type ProductServiceInterface interface {
	GetProductDetail(ctx context.Context, id int) (product *model.ProductDetailResponse, err error)
	CreateProduct(ctx context.Context, product model.CreateProductRequest) (err error)
	UpdateProduct(ctx context.Context, pid int, product model.UpdateProductRequest) (err error)
	DeleteProduct(ctx context.Context, id int) (err error)
}

type ProductService struct {
	repo   repository.DBRepository
	logger *logger.Logger
	cache  *cache.Cache
}

func NewProductService(repo repository.DBRepository, logger *logger.Logger, cache *cache.Cache) ProductServiceInterface {
	return &ProductService{
		repo:   repo,
		logger: logger,
		cache:  cache,
	}
}

func (s *ProductService) GetProductDetail(ctx context.Context, id int) (product *model.ProductDetailResponse, err error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("product:%d", id)
	var cachedProduct model.ProductDetailResponse

	if err := s.cache.Get(ctx, cacheKey, &cachedProduct); err == nil {
		s.logger.Debug("product retrieved from cache", "product_id", id)
		return &cachedProduct, nil
	}

	// Cache miss, get from database
	prodDetail, err := s.repo.GetProductDetail(ctx, id)
	if err != nil {
		s.logger.Error("error while fetch product information", err)
		return nil, err
	}

	if prodDetail == nil {
		s.logger.Warn("product not found", "product id", id)
		return nil, ErrProductNotFound
	}

	// Send the product details as the response
	product = &model.ProductDetailResponse{
		ID:          prodDetail.ID,
		Name:        prodDetail.Name,
		Description: prodDetail.Description,
		Price:       prodDetail.Price,
		Stock:       prodDetail.Stock,
		CategoryID:  prodDetail.CategoryID,
	}

	// Cache the result for future requests
	if err := s.cache.Set(ctx, cacheKey, product, 30*time.Minute); err != nil {
		s.logger.Warn("failed to cache product", "product_id", id, "error", err)
	}

	return product, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, product model.CreateProductRequest) (err error) {
	productd := &sqlModel.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
	}

	err = s.repo.CreateProduct(ctx, productd)
	if err != nil {
		s.logger.Error("error while create product", err)
		return err
	}

	// Invalidate product cache patterns
	s.invalidateProductCache(ctx)

	return nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, pid int, product model.UpdateProductRequest) (err error) {
	productd := &sqlModel.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
	}

	err = s.repo.UpdateProduct(ctx, pid, productd)
	if err != nil {
		s.logger.Error("error while update product", err)
		return err
	}

	// Invalidate specific product cache and patterns
	s.invalidateProductCache(ctx)
	if err := s.cache.Delete(ctx, fmt.Sprintf("product:%d", pid)); err != nil {
		s.logger.Warn("failed to delete product cache", "product_id", pid, "error", err)
	}

	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) (err error) {
	err = s.repo.DeleteProduct(ctx, id)
	if err != nil {
		s.logger.Error("error while delete product", err)
		return err
	}

	// Invalidate product cache patterns
	s.invalidateProductCache(ctx)
	if err := s.cache.Delete(ctx, fmt.Sprintf("product:%d", id)); err != nil {
		s.logger.Warn("failed to delete product cache", "product_id", id, "error", err)
	}

	return nil
}

// invalidateProductCache removes all product-related cache entries
func (s *ProductService) invalidateProductCache(ctx context.Context) {
	// Delete all product cache patterns
	if err := s.cache.DeletePattern(ctx, "product:*"); err != nil {
		s.logger.Warn("failed to invalidate product cache", "error", err)
	}
}
