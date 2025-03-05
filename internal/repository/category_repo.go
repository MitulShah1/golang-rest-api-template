package repository

import (
	"context"
	"database/sql"
	"golang-rest-api-template/internal/repository/model"

	"github.com/Masterminds/squirrel"
)

const CATEGORY_TABLE_NAME = "categories"

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category model.Category) (int64, error)
	GetCategoryByID(ctx context.Context, id int) (*model.Category, error)
	UpdateCategory(ctx context.Context, id int, category model.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

// Create Category
func (r *NewRepository) CreateCategory(ctx context.Context, category model.Category) (int64, error) {
	query, args, err := squirrel.Insert(CATEGORY_TABLE_NAME).
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

// Get Category By ID
func (r *NewRepository) GetCategoryByID(ctx context.Context, id int) (*model.Category, error) {
	query, args, err := squirrel.Select("*").From(CATEGORY_TABLE_NAME).Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.DB.QueryRowxContext(ctx, query, args...)
	var category model.Category
	err = row.Scan(&category.ID, &category.Name, &category.ParentID, &category.Description, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

// Update Category
func (r *NewRepository) UpdateCategory(ctx context.Context, id int, category model.Category) error {
	query, args, err := squirrel.Update(CATEGORY_TABLE_NAME).
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

// Delete Category
func (r *NewRepository) DeleteCategory(ctx context.Context, id int) error {
	query, args, err := squirrel.Delete(CATEGORY_TABLE_NAME).Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.DB.ExecContext(ctx, query, args...)
	return err
}
