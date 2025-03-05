package category

import (
	"encoding/json"
	"golang-rest-api-template/internal/handlers/product/model"
	"golang-rest-api-template/internal/response"
	"golang-rest-api-template/internal/services/category"
	"golang-rest-api-template/package/logger"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	//GET_ALL_CATEGORIES   = "/category"
	CATEGORY_BY_ID_PATH  = "/category/{id}"
	CREATE_CATEGORY_PATH = "/create-category"
	UPDATE_CATEGORY_PATH = "/update-category/{id}"
	DELETE_CATEGORY_PATH = "/category/{id}"
)

type CategoryAPI struct {
	logger  *logger.Logger
	catSrvc category.CategoryServiceInterface
}

func NewCategoryAPI(logger *logger.Logger, catSrvc category.CategoryServiceInterface) *CategoryAPI {
	return &CategoryAPI{
		logger:  logger,
		catSrvc: catSrvc,
	}
}

func (c *CategoryAPI) RegisterHandlers(router *mux.Router) {
	router.HandleFunc(CREATE_CATEGORY_PATH, c.CreateCategoryDetail).Methods(http.MethodPost)
	router.HandleFunc(CATEGORY_BY_ID_PATH, c.GetCategoryById).Methods(http.MethodGet)
	router.HandleFunc(UPDATE_CATEGORY_PATH, c.UpdateCategory).Methods(http.MethodPut)
	router.HandleFunc(DELETE_CATEGORY_PATH, c.DeleteCategory).Methods(http.MethodDelete)

}

func (c *CategoryAPI) sendErrorResponse(w http.ResponseWriter, message string, status int) {
	res := model.StandardResponse{Message: message}
	resp, _ := json.Marshal(res)
	response.SendResponseRaw(w, status, resp)
}

func (c *CategoryAPI) sendJSONResponse(w http.ResponseWriter, data interface{}, status int) {
	resp, err := json.Marshal(data)
	if err != nil {
		c.logger.Error("error while marshalling response", err)
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}
	response.SendResponseRaw(w, status, resp)
}
