// Package category provides HTTP handlers for category-related operations.
// It includes endpoints for creating, reading, updating, and deleting categories.
package category

import (
	"encoding/json"
	"net/http"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/product/model"
	"github.com/MitulShah1/golang-rest-api-template/internal/response"
	"github.com/MitulShah1/golang-rest-api-template/internal/services/category"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"
	"github.com/gorilla/mux"
)

const (
	// CategoryByIDPath is the path for getting category by ID
	CategoryByIDPath = "/category/{id}"
	// CreateCategoryPath is the path for creating a category
	CreateCategoryPath = "/create-category"
	// UpdateCategoryPath is the path for updating a category
	UpdateCategoryPath = "/update-category/{id}"
	// DeleteCategoryPath is the path for deleting a category
	DeleteCategoryPath = "/category/{id}"
)

type CategoryAPI struct {
	logger  *logger.Logger
	catSrvc category.CategoryServiceInterface
}

func NewCategoryAPI(logger *logger.Logger, catSrvc category.CategoryServiceInterface) *CategoryAPI {
	return &CategoryAPI{
		logger:  logger,
		catSrvc: catSrvc,
	}
}

func (c *CategoryAPI) RegisterHandlers(router *mux.Router) {
	router.HandleFunc(CreateCategoryPath, c.CreateCategoryDetail).Methods(http.MethodPost)
	router.HandleFunc(CategoryByIDPath, c.GetCategoryByID).Methods(http.MethodGet)
	router.HandleFunc(UpdateCategoryPath, c.UpdateCategory).Methods(http.MethodPut)
	router.HandleFunc(DeleteCategoryPath, c.DeleteCategory).Methods(http.MethodDelete)
}

func (c *CategoryAPI) sendErrorResponse(w http.ResponseWriter, message string, status int) {
	res := model.StandardResponse{Message: message}
	resp, err := json.Marshal(res)
	if err != nil {
		c.logger.Error("error while marshalling error response", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}
	response.SendResponseRaw(w, status, resp)
}

func (c *CategoryAPI) sendJSONResponse(w http.ResponseWriter, data any, status int) {
	resp, err := json.Marshal(data)
	if err != nil {
		c.logger.Error("error while marshalling response", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}
	response.SendResponseRaw(w, status, resp)
}
