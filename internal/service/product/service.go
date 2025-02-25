package product

import (
	"context"
	"golang-rest-api-template/internal/repository"
	"golang-rest-api-template/internal/repository/model"
	"golang-rest-api-template/package/logger"
)

type ProductServiceInterface interface {
	GetProductDetail(ctx context.Context, id int) (product *model.Product, err error)
	CreateProduct(ctx context.Context, product *model.Product) (err error)
	UpdateProduct(ctx context.Context, product *model.Product) (err error)
	DeleteProduct(ctx context.Context, id string) (err error)
}

type ProductService struct {
	repo   repository.DBRepository
	logger *logger.Logger
}

func NewProductService(repo repository.DBRepository, logger *logger.Logger) ProductServiceInterface {
	return &ProductService{
		repo:   repo,
		logger: logger,
	}
}

func (s *ProductService) GetProductDetail(ctx context.Context, id int) (product *model.Product, err error) {
	product, err = s.repo.GetProductDetail(ctx, id)
	if err != nil {
		s.logger.Error("error while fetch product information", err)
		return nil, err
	}
	return product, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, product *model.Product) (err error) {
	err = s.repo.CreateProduct(ctx, product)
	if err != nil {
		s.logger.Error("error while create product", err)
		return err
	}
	return nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *model.Product) (err error) {
	err = s.repo.UpdateProduct(ctx, product)
	if err != nil {
		s.logger.Error("error while update product", err)
		return err
	}
	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) (err error) {
	err = s.repo.DeleteProduct(ctx, id)
	if err != nil {
		s.logger.Error("error while delete product", err)
		return err
	}
	return nil
}
