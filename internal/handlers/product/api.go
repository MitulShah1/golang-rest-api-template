package product

import (
	"golang-rest-api-template/internal/service/product"
	"golang-rest-api-template/package/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	GET_PRODUCT_DETAIL_PATH = "/product/{id}"
	CREATE_PRODUCT_PATH     = "/product"
	UPDATE_PRODUCT_PATH     = "/product"
	DELETE_PRODUCT_PATH     = "/product/{id}"
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
	router.Handle(GET_PRODUCT_DETAIL_PATH, http.HandlerFunc(p.GetProductDetail)).Methods("POST")
}
