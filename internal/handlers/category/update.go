// Package category provides HTTP handlers for category-related operations.
// It includes endpoints for creating, reading, updating, and deleting categories.
package category

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/category/model"
	"github.com/MitulShah1/golang-rest-api-template/internal/response"
	"github.com/MitulShah1/golang-rest-api-template/package/validation"
	"github.com/gorilla/mux"
)

// UpdateCategory godoc
// @Summary Update Category example
// @Schemes
// @Description Update Category example
// @Tags Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body model.UpdateCategoryRequest true "Category"
// @Success 	 200  {object}  model.StandardResponse
// @Failure      401  {object}  model.StandardResponse
// @Failure      400  {object}  model.StandardResponse
// @Failure      404  {string} string "404 page not found"
// @Failure      500  {object}  model.StandardResponse
// @Router /v1/category/{id} [put]
// UpdateCategory handles HTTP requests for updating existing categories.
// It validates the request and updates the category in the database.
func (c *CategoryAPI) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := model.StandardResponse{}

	// Get and validate category ID
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

	// Read and parse request body
	var req model.UpdateCategoryRequest
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		c.logger.Error("error while reading request body", err)
		response.SendResponseRaw(w, http.StatusBadRequest, nil)
		return
	}

	if err = json.Unmarshal(body, &req); err != nil {
		c.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if errors := validation.ValidateStruct(req); len(errors) > 0 {
		res.Message = "Validation error"
		res.Data = errors
		c.sendJSONResponse(w, res, http.StatusBadRequest)
		return
	}

	// Update category
	if err := c.catSrvc.UpdateCategory(ctx, cid, req); err != nil {
		c.logger.Error("error while updating category", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	// Send success response
	res.IsSuccess = true
	res.Message = "Category updated successfully"
	c.sendJSONResponse(w, res, http.StatusOK)
}
