package product

import (
	"encoding/json"
	"golang-rest-api-template/internal/handlers/product/model"
	"golang-rest-api-template/internal/response"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (p *ProductAPI) GetProductDetail(w http.ResponseWriter, r *http.Request) {
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

	res.IsSuccess = true
	if product == nil {
		res.Message = "Product not found"
	}
	res.Data = product

	resp, err := json.Marshal(res)
	if err != nil {
		p.logger.Error("error while marshalling response", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	response.SendResponseRaw(w, http.StatusOK, resp)
}
