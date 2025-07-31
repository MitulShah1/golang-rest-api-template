// Package model provides data structures for product-related operations.
// It includes request and response models for product API endpoints.
package model

type StandardResponse struct {
	IsSuccess bool   `json:"success"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"        validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price"       validate:"required"`
	CategoryID  int     `json:"categoryId"  validate:"required"`
	Stock       int     `json:"stock"       validate:"required"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"        validate:"omitempty,required"`
	Description string  `json:"description" validate:"omitempty,required"`
	Price       float64 `json:"price"       validate:"omitempty,required,min=1"`
	CategoryID  int     `json:"categoryId"  validate:"omitempty,required"`
	Stock       int     `json:"stock"       validate:"omitempty,required"`
}

type ProductDetailResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  int     `json:"categoryId"`
	Stock       int     `json:"stock"`
}
