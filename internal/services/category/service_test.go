package category

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/category/model"
	sqlModel "github.com/MitulShah1/golang-rest-api-template/internal/repository/model"
	"github.com/MitulShah1/golang-rest-api-template/internal/services/category/mocks"
)

var mockRepo = new(mocks.CategoryServiceInterface)

func TestCategoryService_CreateCategory(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		createReq := model.CreateCategoryRequest{
			Name:        "Test Category",
			ParentID:    nil,
			Description: "Test Description",
		}

		mockRepo.On("CreateCategory", ctx, createReq).Return(int64(1), nil)

		id, err := mockRepo.CreateCategory(ctx, createReq)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("With Parent ID", func(t *testing.T) {
		parentID := 5
		createReq := model.CreateCategoryRequest{
			Name:        "Child Category",
			ParentID:    &parentID,
			Description: "Child Description",
		}

		mockRepo.On("CreateCategory", ctx, createReq).Return(int64(2), nil)

		id, err := mockRepo.CreateCategory(ctx, createReq)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		createReq := model.CreateCategoryRequest{
			Name:        "Error Category",
			Description: "Error Description",
		}

		mockRepo.On("CreateCategory", ctx, createReq).Return(int64(0), errors.New("database error"))

		id, err := mockRepo.CreateCategory(ctx, createReq)

		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		mockRepo.AssertExpectations(t)
	})
}
func TestCategoryService_GetCategoryByID(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		expectedCategory := &sqlModel.Category{
			ID:          1,
			Name:        "Test Category",
			Description: "Test Description",
		}

		mockRepo.On("GetCategoryByID", ctx, 1).Return(expectedCategory, nil)

		category, err := mockRepo.GetCategoryByID(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, expectedCategory.ID, category.ID)
		assert.Equal(t, expectedCategory.Name, category.Name)
		assert.Equal(t, expectedCategory.Description, category.Description)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.On("GetCategoryByID", ctx, 999).Return(nil, nil)

		category, err := mockRepo.GetCategoryByID(ctx, 999)

		assert.NoError(t, err)
		assert.Nil(t, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		mockRepo.On("GetCategoryByID", ctx, -1).Return(nil, errors.New("invalid category ID"))

		category, err := mockRepo.GetCategoryByID(ctx, -1)

		assert.Error(t, err)
		assert.Nil(t, category)
		assert.Contains(t, err.Error(), "invalid category ID")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockRepo.On("GetCategoryByID", ctx, 2).Return(nil, errors.New("database connection error"))

		category, err := mockRepo.GetCategoryByID(ctx, 2)

		assert.Error(t, err)
		assert.Nil(t, category)
		assert.Contains(t, err.Error(), "database connection error")
		mockRepo.AssertExpectations(t)
	})
}
func TestCategoryService_DeleteCategory(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("DeleteCategory", ctx, 1).Return(nil)

		err := mockRepo.DeleteCategory(ctx, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Category In Use", func(t *testing.T) {
		mockRepo.On("DeleteCategory", ctx, 2).Return(errors.New("category is referenced by other entities"))

		err := mockRepo.DeleteCategory(ctx, 2)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "category is referenced by other entities")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Non-existent Category", func(t *testing.T) {
		mockRepo.On("DeleteCategory", ctx, 999).Return(errors.New("category not found"))

		err := mockRepo.DeleteCategory(ctx, 999)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "category not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Database Connection Error", func(t *testing.T) {
		mockRepo.On("DeleteCategory", ctx, 3).Return(errors.New("database connection failed"))

		err := mockRepo.DeleteCategory(ctx, 3)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection failed")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Category ID", func(t *testing.T) {
		mockRepo.On("DeleteCategory", ctx, -1).Return(errors.New("invalid category ID"))

		err := mockRepo.DeleteCategory(ctx, -1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid category ID")
		mockRepo.AssertExpectations(t)
	})
}
func TestCategoryService_UpdateCategory(t *testing.T) {
	ctx := context.Background()

	t.Run("Success with All Fields", func(t *testing.T) {
		parentID := 3
		updateReq := model.UpdateCategoryRequest{
			Name:        "Updated Category",
			ParentID:    &parentID,
			Description: "Updated Description",
		}

		mockRepo.On("UpdateCategory", ctx, 1, updateReq).Return(nil)

		err := mockRepo.UpdateCategory(ctx, 1, updateReq)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success with Partial Update", func(t *testing.T) {
		updateReq := model.UpdateCategoryRequest{
			Name:        "Only Name Update",
			Description: "",
		}

		mockRepo.On("UpdateCategory", ctx, 2, updateReq).Return(nil)

		err := mockRepo.UpdateCategory(ctx, 2, updateReq)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update with Empty Fields", func(t *testing.T) {
		updateReq := model.UpdateCategoryRequest{}

		mockRepo.On("UpdateCategory", ctx, 4, updateReq).Return(errors.New("invalid update data"))

		err := mockRepo.UpdateCategory(ctx, 4, updateReq)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid update data")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update Non-existent Category", func(t *testing.T) {
		updateReq := model.UpdateCategoryRequest{
			Name: "Non-existent Category",
		}

		mockRepo.On("UpdateCategory", ctx, 999, updateReq).Return(errors.New("category not found"))

		err := mockRepo.UpdateCategory(ctx, 999, updateReq)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "category not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update with Invalid Parent ID", func(t *testing.T) {
		parentID := -1
		updateReq := model.UpdateCategoryRequest{
			Name:     "Invalid Parent",
			ParentID: &parentID,
		}

		mockRepo.On("UpdateCategory", ctx, 5, updateReq).Return(errors.New("invalid parent ID"))

		err := mockRepo.UpdateCategory(ctx, 5, updateReq)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid parent ID")
		mockRepo.AssertExpectations(t)
	})
}
