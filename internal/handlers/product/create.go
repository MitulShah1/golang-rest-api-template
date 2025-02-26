package product

import (
	"encoding/json"
	"golang-rest-api-template/internal/handlers/product/model"
	"golang-rest-api-template/internal/response"
	"golang-rest-api-template/package/validation"
	"io"
	"net/http"
)

func (p *ProductAPI) CreateProductDetail(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	res := model.StandardResponse{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		p.logger.Error("error while decoding request", err)
		response.SendResponseRaw(w, http.StatusBadRequest, nil)
		return
	}
	defer r.Body.Close()

	// Parse the request body
	var req model.CreateProductRequest
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

	// Create the product
	if err := p.prdService.CreateProduct(ctx, req); err != nil {
		p.logger.Error("error while creating product", err)
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
