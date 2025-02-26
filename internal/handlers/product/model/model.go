package model

type StandardResponse struct {
	IsSuccess bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	CategoryID  int     `json:"category_id" validate:"required"`
	Stock       int     `json:"stock" validate:"required"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"omitempty,required"`
	Description string  `json:"description" validate:"omitempty,required"`
	Price       float64 `json:"price" validate:"omitempty,required"`
	CategoryID  int     `json:"category_id" validate:"omitempty,required"`
	Stock       int     `json:"stock" validate:"omitempty,required"`
}

type ProductDetailResponse struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  int     `json:"category_id"`
	Stock       int     `json:"stock"`
}
