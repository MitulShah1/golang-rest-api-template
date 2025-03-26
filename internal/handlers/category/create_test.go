package category

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/category/model"
	"github.com/MitulShah1/golang-rest-api-template/internal/services/category/mocks"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"

	"github.com/stretchr/testify/assert"
)

var mockCategoryService = new(mocks.CategoryServiceInterface)

func TestCategoryAPI_CreateCategoryDetail(t *testing.T) {
	testLogger := logger.NewLogger(logger.DefaultOptions())
	api := &CategoryAPI{
		catSrvc: mockCategoryService,
		logger:  testLogger,
	}

	t.Run("Service Error", func(t *testing.T) {
		validCategory := model.CreateCategoryRequest{
			Name:        "Test Category",
			Description: "Test Description",
		}
		body, _ := json.Marshal(validCategory)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		w := httptest.NewRecorder()

		mockCategoryService.On("CreateCategory", context.Background(), validCategory).Return(int64(0), errors.New("service error")).Once()

		api.CreateCategoryDetail(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockCategoryService.AssertExpectations(t)
	})

	t.Run("Marshal Response Error", func(t *testing.T) {
		validCategory := model.CreateCategoryRequest{
			Name:        "Test Category",
			Description: "Test Description",
		}
		body, _ := json.Marshal(validCategory)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		w := httptest.NewRecorder()

		mockCategoryService.On("CreateCategory", context.Background(), validCategory).Return(int64(1), nil).Once()

		api.CreateCategoryDetail(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockCategoryService.AssertExpectations(t)

		var response model.StandardResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.True(t, response.IsSuccess)
		assert.NotNil(t, response.Data)
	})

	t.Run("Successful Category Creation", func(t *testing.T) {
		validCategory := model.CreateCategoryRequest{
			Name:        "Test Category",
			Description: "Test Description",
		}
		body, _ := json.Marshal(validCategory)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		w := httptest.NewRecorder()

		mockCategoryService.On("CreateCategory", context.Background(), validCategory).Return(int64(1), nil).Once()

		api.CreateCategoryDetail(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockCategoryService.AssertExpectations(t)

		var response model.StandardResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.True(t, response.IsSuccess)

		data, ok := response.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(1), data["category_id"])
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		api.CreateCategoryDetail(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Empty Request Body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/categories", nil)
		w := httptest.NewRecorder()

		api.CreateCategoryDetail(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Validation Error", func(t *testing.T) {
		invalidCategory := model.CreateCategoryRequest{
			Name:        "",
			Description: "",
		}
		body, _ := json.Marshal(invalidCategory)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		w := httptest.NewRecorder()

		api.CreateCategoryDetail(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response model.StandardResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Validation error", response.Message)
		assert.NotNil(t, response.Data)
	})
}
