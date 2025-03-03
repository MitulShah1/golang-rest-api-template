package product

import (
	"golang-rest-api-template/internal/handlers/product/model"
	"golang-rest-api-template/internal/response"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Product godoc
// @Summary Delete Product example
// @Schemes
// @Description Delete Product example
// @Tags Product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 	 200  {object}  model.StandardResponse
// @Failure      401  {object}  model.StandardResponse
// @Failure      400  {object}  model.StandardResponse
// @Failure      404  {string}  "404 page not found"
// @Failure      500  {object}  model.StandardResponse
// @Router /v1/product/{id} [DELETE]
func (p *ProductAPI) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := model.StandardResponse{}

	productID := mux.Vars(r)["id"]
	if productID == "" {
		p.sendErrorResponse(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	pid, err := strconv.Atoi(productID)
	if err != nil || pid <= 0 {
		p.sendErrorResponse(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := p.prdService.GetProductDetail(ctx, pid)
	if err != nil {
		p.logger.Error("error while fetching product details", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	if product == nil {
		res.Message = "Product not found"
		p.sendJSONResponse(w, res, http.StatusOK)
		return
	}

	if err := p.prdService.DeleteProduct(ctx, pid); err != nil {
		p.logger.Error("error while delete product", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	res.IsSuccess = true
	res.Message = "Product deleted successfully"
	p.sendJSONResponse(w, res, http.StatusOK)
}
