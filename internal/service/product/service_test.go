package product

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang-rest-api-template/internal/handlers/product/model"
	"golang-rest-api-template/internal/service/product/mocks"
)

var mockService = new(mocks.ProductServiceInterface)

func TestProductService_GetProductDetail(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		expectedProduct := &model.ProductDetailResponse{
			Id:          1,
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
		}

		mockService.On("GetProductDetail", ctx, 1).Return(expectedProduct, nil)

		product, err := mockService.GetProductDetail(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, product)
		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService.On("GetProductDetail", ctx, 999).Return(nil, errors.New("product not found"))

		product, err := mockService.GetProductDetail(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, product)
		mockService.AssertExpectations(t)
	})
}

func TestProductService_CreateProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		newProduct := model.CreateProductRequest{
			Name:        "New Product",
			Description: "New Description",
			Price:       199.99,
		}

		mockService.On("CreateProduct", ctx, newProduct).Return(nil)

		err := mockService.CreateProduct(ctx, newProduct)

		assert.NoError(t, err)
		mockService.AssertExpectations(t)
	})

	t.Run("DB Error", func(t *testing.T) {
		invalidProduct := model.CreateProductRequest{
			Name:  "",
			Price: -1,
		}

		mockService.On("CreateProduct", ctx, invalidProduct).Return(errors.New("DB Error"))

		err := mockService.CreateProduct(ctx, invalidProduct)

		assert.Error(t, err)
		mockService.AssertExpectations(t)
	})
}

func TestProductService_UpdateProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		updateProduct := model.UpdateProductRequest{
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       299.99,
		}

		mockService.On("UpdateProduct", ctx, 1, updateProduct).Return(nil)

		err := mockService.UpdateProduct(ctx, 1, updateProduct)

		assert.NoError(t, err)
		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		updateProduct := model.UpdateProductRequest{
			Name: "Updated Product",
		}

		mockService.On("UpdateProduct", ctx, 999, updateProduct).Return(errors.New("product not found"))

		err := mockService.UpdateProduct(ctx, 999, updateProduct)

		assert.Error(t, err)
		mockService.AssertExpectations(t)
	})
}

func TestProductService_DeleteProduct(t *testing.T) {

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockService.On("DeleteProduct", ctx, 1).Return(nil)

		err := mockService.DeleteProduct(ctx, 1)

		assert.NoError(t, err)
		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService.On("DeleteProduct", ctx, 999).Return(errors.New("product not found"))

		err := mockService.DeleteProduct(ctx, 999)

		assert.Error(t, err)
		mockService.AssertExpectations(t)
	})
}
