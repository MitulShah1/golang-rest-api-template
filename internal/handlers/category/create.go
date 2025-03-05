package category

import (
	"encoding/json"
	"golang-rest-api-template/internal/handlers/category/model"
	"golang-rest-api-template/internal/response"
	"golang-rest-api-template/package/validation"
	"io"
	"net/http"
)

// Category godoc
// @Summary Create Category example
// @Schemes
// @Description Create Category example
// @Tags Category
// @Accept json
// @Produce json
// @Param category body model.CreateCategoryRequest true "Category"
// @Success 	 200  {object}  model.CreateCategoryResponse
// @Failure      401  {object}  model.StandardResponse
// @Failure      400  {object}  model.StandardResponse
// @Failure      404  {string}  "404 page not found"
// @Failure      500  {object}  model.StandardResponse
// @Router /v1/create-category [post]
func (c *CategoryAPI) CreateCategoryDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := model.StandardResponse{}

	// Read and parse request body
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		c.logger.Error("error while reading request body", err)
		response.SendResponseRaw(w, http.StatusBadRequest, nil)
		return
	}

	var req model.CreateCategoryRequest
	if err = json.Unmarshal(body, &req); err != nil {
		c.logger.Error("error while parsing request body", err)
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

	// Create Category
	cateId, err := c.catSrvc.CreateCategory(ctx, req)
	if err != nil {
		c.logger.Error("error while creating category", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	// Send success response
	res.IsSuccess = true
	res.Data = struct {
		CategoryID int64 `json:"category_id"`
	}{
		CategoryID: cateId,
	}
	resp, err := json.Marshal(res)
	if err != nil {
		c.logger.Error("error while marshalling response", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	response.SendResponseRaw(w, http.StatusOK, resp)
}
