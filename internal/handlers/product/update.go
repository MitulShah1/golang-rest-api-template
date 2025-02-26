package product

import (
	"encoding/json"
	"golang-rest-api-template/internal/handlers/product/model"
	"golang-rest-api-template/internal/response"
	"golang-rest-api-template/package/validation"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (p *ProductAPI) UpdateProductDetail(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	res := model.StandardResponse{}

	// Get the product ID from the request URL
	productID := mux.Vars(r)["id"]

	if productID == "" {
		p.logger.Error("product id is empty")
		res.Message = "Product ID is required"
		resp, _ := json.Marshal(res)
		response.SendResponseRaw(w, http.StatusBadRequest, resp)
		return
	}

	// Convert productID string to int
	pid, err := strconv.Atoi(productID)
	if err != nil || pid == 0 {
		p.logger.Error("invalid product id", err)
		res.Message = "Invalid product ID"
		resp, _ := json.Marshal(res)
		response.SendResponseRaw(w, http.StatusBadRequest, resp)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		p.logger.Error("error while decoding request", err)
		response.SendResponseRaw(w, http.StatusBadRequest, nil)
		return
	}
	defer r.Body.Close()

	// Parse the request body
	var req model.UpdateProductRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		res.Message = "Invalid request body"
		p.logger.Error("error while decoding request", err)
		resp, _ := json.Marshal(res)
		response.SendResponseRaw(w, http.StatusBadRequest, resp)
		return
	}

	// validate the request
	errors := validation.ValidateStruct(req)
	if len(errors) > 0 {
		res.Message = "Validation error"
		res.Data = errors
		resp, _ := json.Marshal(res)
		response.SendResponseRaw(w, http.StatusBadRequest, resp)
		return
	}

	// Update the product
	if err := p.prdService.UpdateProduct(ctx, pid, req); err != nil {
		p.logger.Error("error while updating product", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	// Send the product details as the response
	res.IsSuccess = true
	resp, err := json.Marshal(res)
	if err != nil {
		p.logger.Error("error while marshalling response", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	response.SendResponseRaw(w, http.StatusOK, resp)
}
