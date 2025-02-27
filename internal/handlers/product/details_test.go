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

func TestProductAPI_GetProductDetail(t *testing.T) {
	testLogger := logger.NewLogger(logger.DefaultOptions())
	api := &ProductAPI{
		prdService: mockService,
		logger:     testLogger,
	}

	t.Run("Successful Get With Data", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		w := httptest.NewRecorder()

		expectedProduct := &model.ProductDetailResponse{
			Id:          1,
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
			Stock:       10,
			CategoryID:  1,
		}

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockService.On("GetProductDetail", mock.Anything, 1).Return(expectedProduct, nil).Once()

		api.GetProductDetail(w, req)

		var data = map[string]interface{}{
			"id":          expectedProduct.Id,
			"name":        expectedProduct.Name,
			"description": expectedProduct.Description,
			"price":       expectedProduct.Price,
			"category_id": expectedProduct.CategoryID,
			"stock":       expectedProduct.Stock,
		}

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.StandardResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.True(t, response.IsSuccess)
		assert.IsType(t, data, response.Data)
		mockService.AssertExpectations(t)
	})

	t.Run("Negative Product ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/-1", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "-1"})
		api.GetProductDetail(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Returns Error With Custom Message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockService.On("GetProductDetail", mock.Anything, 1).Return(nil, errors.New("custom service error")).Once()

		api.GetProductDetail(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
