package product

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-rest-api-template/internal/handlers/product/model"
	"golang-rest-api-template/internal/services/product/mocks"
	"golang-rest-api-template/package/logger"

	"github.com/stretchr/testify/assert"
)

var mockService = new(mocks.ProductServiceInterface)

func TestProductAPI_CreateProductDetail(t *testing.T) {
	testLogger := logger.NewLogger(logger.DefaultOptions())
	api := &ProductAPI{
		prdService: mockService,
		logger:     testLogger,
	}

	t.Run("Invalid JSON Body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		api.CreateProductDetail(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Empty Request Body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/products", nil)
		w := httptest.NewRecorder()

		api.CreateProductDetail(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Validation Error", func(t *testing.T) {
		invalidProduct := model.CreateProductRequest{
			Name:  "",
			Price: -1,
		}
		body, _ := json.Marshal(invalidProduct)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		w := httptest.NewRecorder()

		api.CreateProductDetail(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Success", func(t *testing.T) {
		validProduct := model.CreateProductRequest{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
			Stock:       10,
			CategoryID:  1,
		}
		body, _ := json.Marshal(validProduct)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		w := httptest.NewRecorder()

		mockService.On("CreateProduct", context.Background(), validProduct).Return(nil)

		api.CreateProductDetail(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var response model.StandardResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.True(t, response.IsSuccess)
	})

	t.Run("Body Read Error", func(t *testing.T) {
		errReader := &ErrorReader{Err: errors.New("read error")}
		req := httptest.NewRequest(http.MethodPost, "/products", errReader)
		w := httptest.NewRecorder()

		api.CreateProductDetail(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

type ErrorReader struct {
	Err error
}

func (e *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, e.Err
}
