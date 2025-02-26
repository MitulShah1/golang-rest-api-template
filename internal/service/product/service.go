package product

import (
	"context"
	"golang-rest-api-template/internal/handlers/product/model"
	"golang-rest-api-template/internal/repository"
	sqlModel "golang-rest-api-template/internal/repository/model"
	"golang-rest-api-template/package/logger"
)

type ProductServiceInterface interface {
	GetProductDetail(ctx context.Context, id int) (product *model.ProductDetailResponse, err error)
	CreateProduct(ctx context.Context, product model.CreateProductRequest) (err error)
	UpdateProduct(ctx context.Context, pid int, product model.UpdateProductRequest) (err error)
	DeleteProduct(ctx context.Context, id int) (err error)
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

func (s *ProductService) GetProductDetail(ctx context.Context, id int) (product *model.ProductDetailResponse, err error) {
	prodDetail, err := s.repo.GetProductDetail(ctx, id)
	if err != nil {
		s.logger.Error("error while fetch product information", err)
		return nil, err
	}

	if prodDetail == nil {
		s.logger.Warn("product not found", "product id", id)
		return nil, nil
	}

	// Send the product details as the response
	product = &model.ProductDetailResponse{
		Id:          prodDetail.ID,
		Name:        prodDetail.Name,
		Description: prodDetail.Description,
		Price:       prodDetail.Price,
		Stock:       prodDetail.Stock,
		CategoryID:  prodDetail.CategoryID,
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
	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) (err error) {
	err = s.repo.DeleteProduct(ctx, id)
	if err != nil {
		s.logger.Error("error while delete product", err)
		return err
	}
	return nil
}
