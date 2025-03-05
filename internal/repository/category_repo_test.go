package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"golang-rest-api-template/internal/repository/model"
	"golang-rest-api-template/package/database"
	"golang-rest-api-template/package/database/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_CreateCategory(t *testing.T) {
	mockDB, mock, err := mocks.NewMockDBWithRegEx()
	assert.NoError(t, err)
	defer mockDB.Close()

	db := &database.Database{DB: mockDB}
	repo := &NewRepository{db: db}
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		parentID := 1
		category := model.Category{
			Name:        "Test Category",
			ParentID:    &parentID,
			Description: "Test Description",
		}

		mock.ExpectExec("INSERT INTO categories").
			WithArgs(category.Name, category.ParentID, category.Description).
			WillReturnResult(sqlmock.NewResult(1, 1))

		id, err := repo.CreateCategory(ctx, category)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})

	t.Run("Success Without Parent", func(t *testing.T) {
		category := model.Category{
			Name:        "Root Category",
			ParentID:    nil,
			Description: "Root Description",
		}

		mock.ExpectExec("INSERT INTO categories").
			WithArgs(category.Name, category.ParentID, category.Description).
			WillReturnResult(sqlmock.NewResult(2, 1))

		id, err := repo.CreateCategory(ctx, category)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), id)
	})

	t.Run("Database Error", func(t *testing.T) {
		category := model.Category{
			Name:        "Error Category",
			Description: "Error Description",
		}

		mock.ExpectExec("INSERT INTO categories").
			WithArgs(category.Name, category.ParentID, category.Description).
			WillReturnError(errors.New("database error"))

		id, err := repo.CreateCategory(ctx, category)
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
	})

	t.Run("Duplicate Entry", func(t *testing.T) {
		category := model.Category{
			Name:        "Duplicate Category",
			Description: "Duplicate Description",
		}

		mock.ExpectExec("INSERT INTO categories").
			WithArgs(category.Name, category.ParentID, category.Description).
			WillReturnError(errors.New("Error 1062: Duplicate entry"))

		id, err := repo.CreateCategory(ctx, category)
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
	})

	t.Run("Invalid Parent ID", func(t *testing.T) {
		parentID := 999
		category := model.Category{
			Name:        "Invalid Parent",
			ParentID:    &parentID,
			Description: "Invalid Parent Description",
		}

		mock.ExpectExec("INSERT INTO categories").
			WithArgs(category.Name, category.ParentID, category.Description).
			WillReturnError(errors.New("foreign key constraint fails"))

		id, err := repo.CreateCategory(ctx, category)
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
	})
}
func TestRepository_GetCategoryByID(t *testing.T) {
	mockDB, mock, err := mocks.NewMockDBWithRegEx()
	assert.NoError(t, err)
	defer mockDB.Close()

	db := &database.Database{DB: mockDB}
	repo := &NewRepository{db: db}
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		expectedCategory := &model.Category{
			ID:          1,
			Name:        "Test Category",
			Description: "Test Description",
			ParentID:    nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		rows := sqlmock.NewRows([]string{"id", "name", "parent_id", "description", "created_at", "updated_at"}).
			AddRow(expectedCategory.ID, expectedCategory.Name, expectedCategory.ParentID, expectedCategory.Description, expectedCategory.CreatedAt, expectedCategory.UpdatedAt)

		mock.ExpectQuery("SELECT \\* FROM categories WHERE").
			WithArgs(expectedCategory.ID).
			WillReturnRows(rows)

		category, err := repo.GetCategoryByID(ctx, expectedCategory.ID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCategory.ID, category.ID)
		assert.Equal(t, expectedCategory.Name, category.Name)
		assert.Equal(t, expectedCategory.Description, category.Description)
	})

	t.Run("Category Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM categories WHERE").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		category, err := repo.GetCategoryByID(ctx, 999)
		assert.NoError(t, err)
		assert.Nil(t, category)
	})

	t.Run("Database Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM categories WHERE").
			WithArgs(1).
			WillReturnError(errors.New("database connection error"))

		category, err := repo.GetCategoryByID(ctx, 1)
		assert.Error(t, err)
		assert.Nil(t, category)
	})

	t.Run("SQL Query Building Error", func(t *testing.T) {
		// Simulate a case where query building fails
		// This is a rare case but good for coverage
		invalidID := "invalid"
		mock.ExpectQuery("SELECT \\* FROM categories WHERE").
			WithArgs(invalidID).
			WillReturnError(errors.New("invalid query"))

		category, err := repo.GetCategoryByID(ctx, 1)
		assert.Error(t, err)
		assert.Nil(t, category)
	})

	t.Run("Scan Error", func(t *testing.T) {
		// Return invalid data type to cause scan error
		rows := sqlmock.NewRows([]string{"id", "name", "parent_id", "description", "created_at", "updated_at"}).
			AddRow("invalid", "Test Category", nil, "Test Description", time.Now(), time.Now())

		mock.ExpectQuery("SELECT \\* FROM categories WHERE").
			WithArgs(1).
			WillReturnRows(rows)

		category, err := repo.GetCategoryByID(ctx, 1)
		assert.Error(t, err)
		assert.Nil(t, category)
	})
}
func TestRepository_UpdateCategory(t *testing.T) {
	mockDB, mock, err := mocks.NewMockDBWithRegEx()
	assert.NoError(t, err)
	defer mockDB.Close()

	db := &database.Database{DB: mockDB}
	repo := &NewRepository{db: db}
	ctx := context.Background()

	t.Run("Success Update All Fields", func(t *testing.T) {
		parentID := 2
		category := model.Category{
			Name:        "Updated Category",
			ParentID:    &parentID,
			Description: "Updated Description",
		}

		mock.ExpectExec("UPDATE categories").
			WithArgs(category.Name, category.ParentID, category.Description, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateCategory(ctx, 1, category)
		assert.NoError(t, err)
	})

	t.Run("Update With Null Parent ID", func(t *testing.T) {
		category := model.Category{
			Name:        "Updated Root Category",
			ParentID:    nil,
			Description: "Updated Root Description",
		}

		mock.ExpectExec("UPDATE categories").
			WithArgs(category.Name, category.ParentID, category.Description, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateCategory(ctx, 1, category)
		assert.NoError(t, err)
	})

	t.Run("No Rows Affected", func(t *testing.T) {
		category := model.Category{
			Name:        "Non-existent Category",
			Description: "Description",
		}

		mock.ExpectExec("UPDATE categories").
			WithArgs(category.Name, category.ParentID, category.Description, 999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateCategory(ctx, 999, category)
		assert.NoError(t, err)
	})

	t.Run("Foreign Key Violation", func(t *testing.T) {
		parentID := 999
		category := model.Category{
			Name:        "Invalid Parent Category",
			ParentID:    &parentID,
			Description: "Description",
		}

		mock.ExpectExec("UPDATE categories").
			WithArgs(category.Name, category.ParentID, category.Description, 1).
			WillReturnError(errors.New("foreign key constraint fails"))

		err := repo.UpdateCategory(ctx, 1, category)
		assert.Error(t, err)
	})

	t.Run("Database Connection Error", func(t *testing.T) {
		category := model.Category{
			Name:        "Error Category",
			Description: "Description",
		}

		mock.ExpectExec("UPDATE categories").
			WithArgs(category.Name, category.ParentID, category.Description, 1).
			WillReturnError(errors.New("connection refused"))

		err := repo.UpdateCategory(ctx, 1, category)
		assert.Error(t, err)
	})
}
func TestRepository_DeleteCategory(t *testing.T) {
	mockDB, mock, err := mocks.NewMockDBWithRegEx()
	assert.NoError(t, err)
	defer mockDB.Close()

	db := &database.Database{DB: mockDB}
	repo := &NewRepository{db: db}
	ctx := context.Background()

	t.Run("Success Delete Category", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM categories").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteCategory(ctx, 1)
		assert.NoError(t, err)
	})

	t.Run("Delete Non-Existent Category", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM categories").
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DeleteCategory(ctx, 999)
		assert.NoError(t, err)
	})

	t.Run("Delete With Referenced Foreign Key", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM categories").
			WithArgs(1).
			WillReturnError(errors.New("foreign key constraint fails"))

		err := repo.DeleteCategory(ctx, 1)
		assert.Error(t, err)
	})

	t.Run("SQL Query Building Error", func(t *testing.T) {
		// Simulate a case where query building might fail
		// This is an edge case where the squirrel library might fail
		mock.ExpectExec("DELETE FROM categories").
			WithArgs(-1).
			WillReturnError(errors.New("invalid query"))

		err := repo.DeleteCategory(ctx, -1)
		assert.Error(t, err)
	})

	t.Run("Database Connection Error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM categories").
			WithArgs(1).
			WillReturnError(errors.New("connection refused"))

		err := repo.DeleteCategory(ctx, 1)
		assert.Error(t, err)
	})
}
