package category

import (
	"net/http"
	"strconv"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/category/model"
	"github.com/MitulShah1/golang-rest-api-template/internal/response"

	"github.com/gorilla/mux"
)

// Category godoc
// @Summary Delete Category example
// @Schemes
// @Description Delete Category example
// @Tags Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 	 200  {object}  model.StandardResponse
// @Failure      401  {object}  model.StandardResponse
// @Failure      400  {object}  model.StandardResponse
// @Failure      404  {string}  "404 page not found"
// @Failure      500  {object}  model.StandardResponse
// @Router /v1/category/{id} [DELETE]
func (c *CategoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
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
		c.logger.Error("error while fetching category details", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	if category == nil {
		res.Message = "Category not found"
		c.sendJSONResponse(w, res, http.StatusOK)
		return
	}

	if err := c.catSrvc.DeleteCategory(ctx, cid); err != nil {
		c.logger.Error("error while delete category", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	res.IsSuccess = true
	res.Message = "Category deleted successfully"
	c.sendJSONResponse(w, res, http.StatusOK)
}
