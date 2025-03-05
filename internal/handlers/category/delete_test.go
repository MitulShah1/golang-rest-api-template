package category

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-rest-api-template/internal/handlers/category/model"
	sqlModel "golang-rest-api-template/internal/repository/model"
	"golang-rest-api-template/package/logger"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryAPI_DeleteCategory(t *testing.T) {
	testLogger := logger.NewLogger(logger.DefaultOptions())
	api := &CategoryAPI{
		catSrvc: mockCategoryService,
		logger:  testLogger,
	}

	t.Run("Missing Category ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/categories/", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{})
		api.DeleteCategory(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Category ID Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/categories/abc", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		api.DeleteCategory(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Negative Category ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/categories/-1", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "-1"})
		api.DeleteCategory(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetCategoryByID Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockCategoryService.On("GetCategoryByID", mock.Anything, 1).Return(nil, errors.New("database error")).Once()

		api.DeleteCategory(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockCategoryService.AssertExpectations(t)
	})

	t.Run("Category Not Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockCategoryService.On("GetCategoryByID", mock.Anything, 1).Return(nil, nil).Once()

		api.DeleteCategory(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.StandardResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Category not found", response.Message)
		mockCategoryService.AssertExpectations(t)
	})

	t.Run("Delete Category Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockCategoryService.On("GetCategoryByID", mock.Anything, 1).Return(&sqlModel.Category{ID: 1}, nil).Once()
		mockCategoryService.On("DeleteCategory", mock.Anything, 1).Return(errors.New("delete error")).Once()

		api.DeleteCategory(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockCategoryService.AssertExpectations(t)
	})

	t.Run("Successful Delete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockCategoryService.On("GetCategoryByID", mock.Anything, 1).Return(&sqlModel.Category{ID: 1}, nil).Once()
		mockCategoryService.On("DeleteCategory", mock.Anything, 1).Return(nil).Once()

		api.DeleteCategory(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.StandardResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.True(t, response.IsSuccess)
		assert.Equal(t, "Category deleted successfully", response.Message)
		mockCategoryService.AssertExpectations(t)
	})
}
