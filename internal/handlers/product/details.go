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

	// Get the product ID from the request URL
	productID := mux.Vars(r)["id"]

	// Convert productID string to int
	pid, err := strconv.Atoi(productID)
	if err != nil {
		p.logger.Error("invalid product id", err)
		response.SendResponseRaw(w, http.StatusBadRequest, []byte("Invalid product ID"))
		return
	}

	// Call the service to get the product details
	product, err := p.prdService.GetProductDetail(ctx, pid)
	if err != nil {
		p.logger.Error("error while fetching product details", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	// Send the product details as the response
	res.IsSuccess = true
	res.Data = product

	resp, err := json.Marshal(res)
	if err != nil {
		p.logger.Error("error while marshalling response", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	response.SendResponseRaw(w, http.StatusOK, resp)
}
