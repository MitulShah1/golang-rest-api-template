package category

import (
	"encoding/json"
	"golang-rest-api-template/internal/handlers/category/model"
	"golang-rest-api-template/internal/response"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Category godoc
// @Summary Get Category details example
// @Schemes
// @Description Get Category details example
// @Tags Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 	 200  {object}  model.CategoryByIDResponse
// @Failure      401  {object}  model.StandardResponse
// @Failure      400  {object}  model.StandardResponse
// @Failure      404  {string}  "404 page not found"
// @Failure      500  {object}  model.StandardResponse
// @Router /v1/category/{id} [get]
func (p *CategoryAPI) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := model.StandardResponse{}

	categoryID := mux.Vars(r)["id"]
	if categoryID == "" {
		p.sendErrorResponse(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	cid, err := strconv.Atoi(categoryID)
	if err != nil || cid <= 0 {
		p.sendErrorResponse(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := p.catSrvc.GetCategoryByID(ctx, cid)
	if err != nil {
		p.logger.Error("error while fetching category details", err, "category_id", cid)
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
		p.logger.Error("error while marshalling response", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	response.SendResponseRaw(w, http.StatusOK, resp)
}
