// Package model provides data structures for category-related operations.
// It includes request and response models for category API endpoints.
package model

type StandardResponse struct {
	IsSuccess bool   `json:"success"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name"        validate:"required"`
	ParentID    *int   `json:"parentId"    validate:"omitempty,required"`
	Description string `json:"description" validate:"required"`
}

type CreateCategoryResponse struct {
	IsSuccess bool   `json:"success"`
	Message   string `json:"message"`
	Data      struct {
		ID int64 `json:"categoryId"`
	} `json:"data"`
}

type CategoryDetail struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ParentID    *int   `json:"parentId"`
	Description string `json:"description"`
}

type CategoryByIDResponse struct {
	IsSuccess bool           `json:"success"`
	Message   string         `json:"message"`
	Data      CategoryDetail `json:"data"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name"        validate:"required"`
	ParentID    *int   `json:"parentId"    validate:"omitempty,required"`
	Description string `json:"description" validate:"required"`
}
