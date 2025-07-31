package category

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/category/model"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryAPI_UpdateCategory(t *testing.T) {
	testLogger := logger.NewLogger(logger.DefaultOptions())
	api := &CategoryAPI{
		catSrvc: mockCategoryService,
		logger:  testLogger,
	}

	t.Run("Empty Category ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/categories/", http.NoBody)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{})
		api.UpdateCategory(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Category ID Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/categories/abc", http.NoBody)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		api.UpdateCategory(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Zero Category ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/categories/0", http.NoBody)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		api.UpdateCategory(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid JSON Body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		api.UpdateCategory(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Validation Error", func(t *testing.T) {
		invalidCategory := model.UpdateCategoryRequest{
			Name:        "",
			Description: "",
		}
		body, _ := json.Marshal(invalidCategory)
		req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewReader(body))
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		api.UpdateCategory(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response model.StandardResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Validation error", response.Message)
	})

	t.Run("Successful Update", func(t *testing.T) {
		validCategory := model.UpdateCategoryRequest{
			Name:        "Updated Category",
			Description: "Updated Description",
		}
		body, _ := json.Marshal(validCategory)
		req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewReader(body))
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		mockCategoryService.On("UpdateCategory", mock.Anything, 1, validCategory).Return(nil).Once()

		api.UpdateCategory(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.StandardResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.True(t, response.IsSuccess)
		assert.Equal(t, "Category updated successfully", response.Message)
		mockCategoryService.AssertExpectations(t)
	})
}
