package model

type StandardResponse struct {
	IsSuccess bool        `json:"issuccess"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

type ProductDetailResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
