package product

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-rest-api-template/internal/handlers/product/model"
	"golang-rest-api-template/package/logger"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductAPI_DeleteProduct(t *testing.T) {
	testLogger := logger.NewLogger(logger.DefaultOptions())

	api := &ProductAPI{
		prdService: mockService,
		logger:     testLogger,
	}

	t.Run("Missing Product ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/products/", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{})
		api.DeleteProduct(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Product ID Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/products/abc", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		api.DeleteProduct(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Zero Product ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/products/0", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		api.DeleteProduct(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Get Product Detail Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockService.On("GetProductDetail", mock.Anything, 1).Return(nil, errors.New("database error")).Once()

		api.DeleteProduct(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Product Not Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockService.On("GetProductDetail", mock.Anything, 1).Return(nil, nil).Once()

		api.DeleteProduct(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.StandardResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Product not found", response.Message)
		mockService.AssertExpectations(t)
	})

	t.Run("Delete Product Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockService.On("GetProductDetail", mock.Anything, 1).Return(&model.ProductDetailResponse{Id: 1}, nil).Once()
		mockService.On("DeleteProduct", mock.Anything, 1).Return(errors.New("delete error")).Once()

		api.DeleteProduct(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Successful Delete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockService.On("GetProductDetail", mock.Anything, 1).Return(&model.ProductDetailResponse{Id: 1}, nil).Once()
		mockService.On("DeleteProduct", mock.Anything, 1).Return(nil).Once()

		api.DeleteProduct(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.StandardResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.True(t, response.IsSuccess)
		mockService.AssertExpectations(t)
	})
}
