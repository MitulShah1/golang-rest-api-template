package product

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/product/model"
	"github.com/MitulShah1/golang-rest-api-template/internal/response"
	"github.com/MitulShah1/golang-rest-api-template/package/validation"
)

// Product godoc
// @Summary Create Product example
// @Schemes
// @Description Create Product example
// @Tags Product
// @Accept json
// @Produce json
// @Param product body model.CreateProductRequest true "Product"
// @Success 	 200  {object}  model.ProductDetailResponse
// @Failure      401  {object}  model.StandardResponse
// @Failure      400  {object}  model.StandardResponse
// @Failure      404  {string}  "404 page not found"
// @Failure      500  {object}  model.StandardResponse
// @Router /v1/create-product [post]
func (p *ProductAPI) CreateProductDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := model.StandardResponse{}

	// Read and parse request body
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		p.logger.Error("error while reading request body", err)
		response.SendResponseRaw(w, http.StatusBadRequest, nil)
		return
	}

	var req model.CreateProductRequest
	if err = json.Unmarshal(body, &req); err != nil {
		p.logger.Error("error while parsing request body", err)
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

	// Create product
	if err = p.prdService.CreateProduct(ctx, req); err != nil {
		p.logger.Error("error while creating product", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	// Send success response
	res.IsSuccess = true
	resp, err := json.Marshal(res)
	if err != nil {
		p.logger.Error("error while marshalling response", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	response.SendResponseRaw(w, http.StatusOK, resp)
}
