package model

type StandardResponse struct {
	IsSuccess bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	ParentID    *int   `json:"parent_id" validate:"omitempty,required"`
	Description string `json:"description" validate:"required"`
}

type CreateCategoryResponse struct {
	IsSuccess bool   `json:"success"`
	Message   string `json:"message"`
	Data      struct {
		ID int64 `json:"category_id"`
	} `json:"data"`
}

type CategoryDetail struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ParentID    *int   `json:"parent_id"`
	Description string `json:"description"`
}

type CategoryByIDResponse struct {
	IsSuccess bool           `json:"success"`
	Message   string         `json:"message"`
	Data      CategoryDetail `json:"data"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	ParentID    *int   `json:"parent_id" validate:"omitempty,required"`
	Description string `json:"description" validate:"required"`
}
