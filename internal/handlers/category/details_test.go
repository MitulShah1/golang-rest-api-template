package category

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/category/model"
	sqlModel "github.com/MitulShah1/golang-rest-api-template/internal/repository/model"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryAPI_GetCategoryById(t *testing.T) {
	testLogger := logger.NewLogger(logger.DefaultOptions())
	api := &CategoryAPI{
		catSrvc: mockCategoryService,
		logger:  testLogger,
	}

	t.Run("Missing Category ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/categories/", nil)
		w := httptest.NewRecorder()

		api.GetCategoryById(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Category ID Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/categories/abc", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		api.GetCategoryById(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Zero Category ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/categories/0", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		api.GetCategoryById(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Category Not Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockCategoryService.On("GetCategoryByID", mock.Anything, 1).Return(nil, nil).Once()

		api.GetCategoryById(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.StandardResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.True(t, response.IsSuccess)
		assert.Equal(t, "Category not found", response.Message)
		mockCategoryService.AssertExpectations(t)
	})

	t.Run("Successful Category Retrieval", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
		w := httptest.NewRecorder()

		expectedCategory := &sqlModel.Category{
			ID:          1,
			Name:        "Test Category",
			Description: "Test Description",
		}

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockCategoryService.On("GetCategoryByID", mock.Anything, 1).Return(expectedCategory, nil).Once()

		api.GetCategoryById(w, req)

		var data = map[string]interface{}{
			"id":          expectedCategory.ID,
			"name":        expectedCategory.Name,
			"description": expectedCategory.Description,
		}

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.StandardResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.True(t, response.IsSuccess)
		assert.IsType(t, data, response.Data)
		mockCategoryService.AssertExpectations(t)
	})
}
