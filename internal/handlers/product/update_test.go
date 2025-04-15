package product

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/product/model"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductAPI_UpdateProductDetail(t *testing.T) {
	testLogger := logger.NewLogger(logger.DefaultOptions())
	api := &ProductAPI{
		prdService: mockService,
		logger:     testLogger,
	}

	t.Run("Missing Product ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/products/", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{})
		api.UpdateProductDetail(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Product ID Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/products/abc", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		api.UpdateProductDetail(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Body Read Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/products/1", &ErrorReader{Err: errors.New("read error")})
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		api.UpdateProductDetail(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid JSON Body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		api.UpdateProductDetail(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Validation Error", func(t *testing.T) {
		invalidProduct := model.UpdateProductRequest{
			Name:  "",
			Price: -1,
		}
		body, _ := json.Marshal(invalidProduct)
		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewReader(body))
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		api.UpdateProductDetail(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Update Error", func(t *testing.T) {
		validProduct := model.UpdateProductRequest{
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       199.99,
			Stock:       20,
			CategoryID:  2,
		}
		body, _ := json.Marshal(validProduct)
		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewReader(body))
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockService.On("UpdateProduct", mock.Anything, 1, validProduct).Return(errors.New("update error")).Once()

		api.UpdateProductDetail(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Successful Update", func(t *testing.T) {
		validProduct := model.UpdateProductRequest{
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       199.99,
			Stock:       20,
			CategoryID:  2,
		}
		body, _ := json.Marshal(validProduct)
		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewReader(body))
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockService.On("UpdateProduct", mock.Anything, 1, validProduct).Return(nil).Once()

		api.UpdateProductDetail(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.StandardResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.True(t, response.IsSuccess)
		mockService.AssertExpectations(t)
	})
}
