// Package category provides HTTP handlers for category-related operations.
// It includes endpoints for creating, reading, updating, and deleting categories.
package category

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/category/model"
	"github.com/MitulShah1/golang-rest-api-template/internal/response"
	"github.com/gorilla/mux"
)

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get category details by ID
// @Tags Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} model.CategoryByIDResponse
// @Failure 400 {object} model.StandardResponse
// @Failure 404 {object} model.StandardResponse
// @Router /category/{id} [get]
func (c *CategoryAPI) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := model.StandardResponse{}

	categoryID := mux.Vars(r)["id"]
	if categoryID == "" {
		c.sendErrorResponse(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	cid, err := strconv.Atoi(categoryID)
	if err != nil || cid <= 0 {
		c.sendErrorResponse(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := c.catSrvc.GetCategoryByID(ctx, cid)
	if err != nil {
		c.logger.Error("error while fetching category details", err, "category_id", cid)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	res.IsSuccess = true
	if category == nil {
		res.Message = "Category not found"
	}
	res.Data = category

	resp, err := json.Marshal(res)
	if err != nil {
		c.logger.Error("error while marshalling response", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	response.SendResponseRaw(w, http.StatusOK, resp)
}
