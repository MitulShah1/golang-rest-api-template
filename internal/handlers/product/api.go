package product

import (
	"golang-rest-api-template/internal/service/product"
	"golang-rest-api-template/package/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	PRODUCT_DETAIL_PATH = "/product/{id}"
	CREATE_PRODUCT_PATH = "/create-product"
	UPDATE_PRODUCT_PATH = "/update-product/{id}"
	DELETE_PRODUCT_PATH = "/product/{id}"
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
	router.Handle(PRODUCT_DETAIL_PATH, http.HandlerFunc(p.GetProductDetail)).Methods(http.MethodGet)
	router.Handle(CREATE_PRODUCT_PATH, http.HandlerFunc(p.CreateProductDetail)).Methods(http.MethodPost)
	router.Handle(UPDATE_PRODUCT_PATH, http.HandlerFunc(p.UpdateProductDetail)).Methods(http.MethodPut)
	router.Handle(DELETE_PRODUCT_PATH, http.HandlerFunc(p.DeleteProduct)).Methods(http.MethodDelete)
}
