// Package product provides HTTP handlers for product-related operations.
// It includes endpoints for creating, reading, updating, and deleting products.
package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/product/model"
	"github.com/MitulShah1/golang-rest-api-template/internal/response"
	"github.com/gorilla/mux"
)

// GetProductDetail godoc
// @Summary Get Product details example
// @Description Get Product details by ID
// @Tags Product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} model.StandardResponse
// @Failure 400 {object} model.StandardResponse
// @Failure 401 {object} model.StandardResponse
// @Failure 404 {string} string "404 page not found"
// @Failure 500 {object} model.StandardResponse
// @Router /v1/product/{id} [get]
// GetProductDetail handles HTTP requests for retrieving product details by ID.
// It validates the ID and returns the product information.
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
