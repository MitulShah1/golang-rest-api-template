package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang-rest-api-template/internal/repository/model"

	"github.com/Masterminds/squirrel"
)

const PRODUCT_TABLE_NAME = "products"

// ProductRepositoryInterface defines the methods for interacting with the product repository.
// The methods allow for retrieving product details, creating new products, updating existing products,
// and deleting products.
type ProductRepositoryInterface interface {
	GetProductDetail(ctx context.Context, id int) (product *model.Product, err error)
	CreateProduct(ctx context.Context, product *model.Product) (err error)
	UpdateProduct(ctx context.Context, pid int, product *model.Product) (err error)
	DeleteProduct(ctx context.Context, id int) (err error)
}

func (r *NewRepository) GetProductDetail(ctx context.Context, id int) (product *model.Product, err error) {
	// Implement the GetProductDetail method
	builder := squirrel.Select("*").From(PRODUCT_TABLE_NAME).Where("id = ?", id)
	query, args, err := builder.Limit(1).ToSql()

	if err != nil {
		return nil, err
	}

	var products model.Product
	err = r.db.DB.QueryRowxContext(ctx, query, args...).StructScan(&products)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &products, nil

}

// CreateProduct creates a new product in the database using the provided product data.
// It returns an error if the creation fails.
func (r *NewRepository) CreateProduct(ctx context.Context, product *model.Product) (err error) {
	builder := squirrel.Insert(PRODUCT_TABLE_NAME).
		Columns("name", "description", "price", "stock", "category_id").
		Values(product.Name, product.Description, product.Price, product.Stock, product.CategoryID)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build sql query: %s", err.Error())
	}

	_, err = r.db.DB.ExecContext(ctx, query, args...)

	return err
}

func (r *NewRepository) UpdateProduct(ctx context.Context, pid int, product *model.Product) (err error) {

	builder := squirrel.Update(PRODUCT_TABLE_NAME)

	if product.Name != "" {
		builder = builder.Set("name", product.Name)
	}
	if product.Description != "" {
		builder = builder.Set("description", product.Description)
	}
	if product.Price > 0 {
		builder = builder.Set("price", product.Price)
	}
	builder = builder.Where("id = ?", pid)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build sql query: %s", err.Error())
	}

	_, err = r.db.DB.ExecContext(ctx, query, args...)

	return err
}

// DeleteProduct deletes a product from the database by the given ID.
// It returns an error if the deletion fails.
func (r *NewRepository) DeleteProduct(ctx context.Context, id int) (err error) {
	builder := squirrel.Delete(PRODUCT_TABLE_NAME).Where("id = ?", id)
	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build sql query: %s", err.Error())
	}
	_, err = r.db.DB.ExecContext(ctx, query, args...)
	return err
}
