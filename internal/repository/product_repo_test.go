package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"golang-rest-api-template/internal/repository/model"
	"golang-rest-api-template/package/database"

	"golang-rest-api-template/package/database/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetProductDetail(t *testing.T) {
	mockDB, mock, err := mocks.NewMockDBWithRegEx()
	assert.NoError(t, err)
	defer mockDB.Close()

	db := &database.Database{DB: mockDB}
	repo := &NewRepository{db: db}
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		expectedProduct := &model.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
		}

		rows := sqlmock.NewRows([]string{"id", "name", "description", "price"}).
			AddRow(expectedProduct.ID, expectedProduct.Name, expectedProduct.Description, expectedProduct.Price)

		mock.ExpectQuery("SELECT (.+) FROM products WHERE").
			WithArgs(1).
			WillReturnRows(rows)

		product, err := repo.GetProductDetail(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, product)
	})

	t.Run("No Rows Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM products WHERE").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		product, err := repo.GetProductDetail(ctx, 999)
		assert.NoError(t, err)
		assert.Nil(t, product)
	})

	t.Run("Database Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM products WHERE").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		product, err := repo.GetProductDetail(ctx, 1)
		assert.Error(t, err)
		assert.Nil(t, product)
	})
}
func TestRepository_CreateProduct(t *testing.T) {
	mockDB, mock, err := mocks.NewMockDBWithRegEx()
	assert.NoError(t, err)
	defer mockDB.Close()

	db := &database.Database{DB: mockDB}
	repo := &NewRepository{db: db}
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		product := &model.Product{
			Name:        "New Product",
			Description: "New Description",
			Price:       199.99,
			Stock:       100,
			CategoryID:  1,
		}

		mock.ExpectExec("INSERT INTO products").
			WithArgs(product.Name, product.Description, product.Price, product.Stock, product.CategoryID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateProduct(ctx, product)
		assert.NoError(t, err)
	})

	t.Run("Database Insert Error", func(t *testing.T) {
		product := &model.Product{
			Name:        "Failed Product",
			Description: "Failed Description",
			Price:       199.99,
			Stock:       100,
			CategoryID:  1,
		}

		mock.ExpectExec("INSERT INTO products").
			WithArgs(product.Name, product.Description, product.Price, product.Stock, product.CategoryID).
			WillReturnError(errors.New("database insert error"))

		err := repo.CreateProduct(ctx, product)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database insert error")
	})
}

func TestRepository_UpdateProduct(t *testing.T) {
	mockDB, mock, err := mocks.NewMockDBWithRegEx()
	assert.NoError(t, err)
	defer mockDB.Close()

	db := &database.Database{DB: mockDB}
	repo := &NewRepository{db: db}
	ctx := context.Background()

	t.Run("Success Update All Fields", func(t *testing.T) {
		product := &model.Product{
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       299.99,
		}

		mock.ExpectExec("UPDATE products").
			WithArgs(product.Name, product.Description, product.Price, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateProduct(ctx, 1, product)
		assert.NoError(t, err)
	})

	t.Run("Success Update Single Field", func(t *testing.T) {
		product := &model.Product{
			Name: "Only Name Update",
		}

		mock.ExpectExec("UPDATE products").
			WithArgs(product.Name, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateProduct(ctx, 1, product)
		assert.NoError(t, err)
	})

	t.Run("Update With Zero Price", func(t *testing.T) {
		product := &model.Product{
			Name:  "Test Product",
			Price: 0,
		}

		mock.ExpectExec("UPDATE products").
			WithArgs(product.Name, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateProduct(ctx, 1, product)
		assert.NoError(t, err)
	})

	t.Run("Update Non-Existent Product", func(t *testing.T) {
		product := &model.Product{
			Name: "Non-existent Product",
		}

		mock.ExpectExec("UPDATE products").
			WithArgs(product.Name, 999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateProduct(ctx, 999, product)
		assert.NoError(t, err)
	})

	t.Run("Database Error", func(t *testing.T) {
		product := &model.Product{
			Name: "Error Product",
		}

		mock.ExpectExec("UPDATE products").
			WithArgs(product.Name, 1).
			WillReturnError(errors.New("database error"))

		err := repo.UpdateProduct(ctx, 1, product)
		assert.Error(t, err)
	})
}

func TestRepository_DeleteProduct(t *testing.T) {
	mockDB, mock, err := mocks.NewMockDBWithRegEx()
	assert.NoError(t, err)
	defer mockDB.Close()

	db := &database.Database{DB: mockDB}
	repo := &NewRepository{db: db}
	ctx := context.Background()

	t.Run("Success Delete Existing Product", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM products WHERE").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteProduct(ctx, 1)
		assert.NoError(t, err)
	})

	t.Run("Delete Non-Existent Product", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM products WHERE").
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DeleteProduct(ctx, 999)
		assert.NoError(t, err)
	})

	t.Run("Database Connection Error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM products WHERE").
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		err := repo.DeleteProduct(ctx, 1)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
	})

	t.Run("Invalid ID Format", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM products WHERE").
			WithArgs(-1).
			WillReturnError(errors.New("invalid id format"))

		err := repo.DeleteProduct(ctx, -1)
		assert.Error(t, err)
	})

	t.Run("Database Constraint Violation", func(t *testing.T) {
		constraintErr := errors.New("foreign key constraint violation")
		mock.ExpectExec("DELETE FROM products WHERE").
			WithArgs(2).
			WillReturnError(constraintErr)

		err := repo.DeleteProduct(ctx, 2)
		assert.Error(t, err)
		assert.Equal(t, constraintErr, err)
	})
}
