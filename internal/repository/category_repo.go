// Package repository provides data access layer for the application.
// It includes database operations for categories, products, and other entities.
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/MitulShah1/golang-rest-api-template/internal/repository/model"
)

var ErrCategoryNotFound = errors.New("category not found")

const CategoryTableName = "categories"

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *model.Category) (int64, error)
	GetCategoryByID(ctx context.Context, id int) (*model.Category, error)
	UpdateCategory(ctx context.Context, id int, category *model.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

// CreateCategory creates a new category in the database.
// It returns the ID of the created category or an error.
func (r *NewRepository) CreateCategory(ctx context.Context, category *model.Category) (int64, error) {
	query, args, err := squirrel.Insert(CategoryTableName).
		Columns("name", "parent_id", "description").
		Values(category.Name, category.ParentID, category.Description).
		ToSql()
	if err != nil {
		return 0, err
	}

	result, err := r.db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetCategoryByID retrieves a category by its ID from the database.
// It returns the category or an error if not found.
func (r *NewRepository) GetCategoryByID(ctx context.Context, id int) (*model.Category, error) {
	query, args, err := squirrel.Select("*").From(CategoryTableName).Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.DB.QueryRowxContext(ctx, query, args...)
	var category model.Category
	err = row.Scan(&category.ID, &category.Name, &category.ParentID, &category.Description, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	return &category, nil
}

// UpdateCategory updates an existing category in the database.
// It returns an error if the update fails.
func (r *NewRepository) UpdateCategory(ctx context.Context, id int, category *model.Category) error {
	query, args, err := squirrel.Update(CategoryTableName).
		Set("name", category.Name).
		Set("parent_id", category.ParentID).
		Set("description", category.Description).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.DB.ExecContext(ctx, query, args...)
	return err
}

// DeleteCategory removes a category from the database by its ID.
// It returns an error if the deletion fails.
func (r *NewRepository) DeleteCategory(ctx context.Context, id int) error {
	query, args, err := squirrel.Delete(CategoryTableName).Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.DB.ExecContext(ctx, query, args...)
	return err
}
