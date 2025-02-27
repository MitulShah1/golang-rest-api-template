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

	// Get and validate product ID
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

	// Read and parse request body
	var req model.UpdateProductRequest
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		p.logger.Error("error while reading request body", err)
		response.SendResponseRaw(w, http.StatusBadRequest, nil)
		return
	}

	if err = json.Unmarshal(body, &req); err != nil {
		p.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if errors := validation.ValidateStruct(req); len(errors) > 0 {
		res.Message = "Validation error"
		res.Data = errors
		p.sendJSONResponse(w, res, http.StatusBadRequest)
		return
	}

	// Update product
	if err := p.prdService.UpdateProduct(ctx, pid, req); err != nil {
		p.logger.Error("error while updating product", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	// Send success response
	res.IsSuccess = true
	p.sendJSONResponse(w, res, http.StatusOK)
}
