// Package product provides HTTP handlers for product-related operations.
// It includes endpoints for creating, reading, updating, and deleting products.
package product

import (
	"encoding/json"
	"net/http"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/product/model"
	"github.com/MitulShah1/golang-rest-api-template/internal/response"
	"github.com/MitulShah1/golang-rest-api-template/internal/services/product"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"
	"github.com/gorilla/mux"
)

const (
	// ProductDetailPath is the path for getting product details
	ProductDetailPath = "/product/{id}"
	// CreateProductPath is the path for creating a product
	CreateProductPath = "/create-product"
	// UpdateProductPath is the path for updating a product
	UpdateProductPath = "/update-product/{id}"
	// DeleteProductPath is the path for deleting a product
	DeleteProductPath = "/product/{id}"
)

type ProductAPI struct {
	logger     *logger.Logger
	prdService product.ProductServiceInterface
}

func NewProductAPI(logger *logger.Logger, prdService product.ProductServiceInterface) *ProductAPI {
	return &ProductAPI{
		logger:     logger,
		prdService: prdService,
	}
}

func (p *ProductAPI) RegisterHandlers(router *mux.Router) {
	router.Handle(ProductDetailPath, http.HandlerFunc(p.GetProductDetail)).Methods(http.MethodGet)
	router.Handle(CreateProductPath, http.HandlerFunc(p.CreateProductDetail)).Methods(http.MethodPost)
	router.Handle(UpdateProductPath, http.HandlerFunc(p.UpdateProductDetail)).Methods(http.MethodPut)
	router.Handle(DeleteProductPath, http.HandlerFunc(p.DeleteProduct)).Methods(http.MethodDelete)
}

func (p *ProductAPI) sendErrorResponse(w http.ResponseWriter, message string, status int) {
	res := model.StandardResponse{Message: message}
	resp, err := json.Marshal(res)
	if err != nil {
		p.logger.Error("error while marshalling error response", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}
	response.SendResponseRaw(w, status, resp)
}

func (p *ProductAPI) sendJSONResponse(w http.ResponseWriter, data any, status int) {
	resp, err := json.Marshal(data)
	if err != nil {
		p.logger.Error("error while marshalling response", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}
	response.SendResponseRaw(w, status, resp)
}
