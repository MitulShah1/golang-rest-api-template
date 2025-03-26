package category

import (
	"context"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/category/model"
	"github.com/MitulShah1/golang-rest-api-template/internal/repository"
	sqlModel "github.com/MitulShah1/golang-rest-api-template/internal/repository/model"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"
)

type CategoryServiceInterface interface {
	CreateCategory(ctx context.Context, category model.CreateCategoryRequest) (int64, error)
	GetCategoryByID(ctx context.Context, id int) (*sqlModel.Category, error)
	UpdateCategory(ctx context.Context, id int, category model.UpdateCategoryRequest) error
	DeleteCategory(ctx context.Context, id int) error
}

type CategoryService struct {
	repo   repository.DBRepository
	logger *logger.Logger
}

func NewCategoryService(repo repository.DBRepository, logger *logger.Logger) CategoryServiceInterface {
	return &CategoryService{
		repo:   repo,
		logger: logger,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category model.CreateCategoryRequest) (int64, error) {

	s.logger.Info("Creating category", "category", category)

	cat := sqlModel.Category{
		Name:        category.Name,
		ParentID:    category.ParentID,
		Description: category.Description,
	}

	return s.repo.CreateCategory(ctx, cat)
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (*sqlModel.Category, error) {

	s.logger.Info("Getting category by ID", "id", id)

	return s.repo.GetCategoryByID(ctx, id)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id int, category model.UpdateCategoryRequest) error {

	s.logger.Info("Updating category", "category", category)

	updCat := sqlModel.Category{
		Name:        category.Name,
		ParentID:    category.ParentID,
		Description: category.Description,
	}

	return s.repo.UpdateCategory(ctx, id, updCat)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {

	s.logger.Info("Deleting category", "id", id)

	return s.repo.DeleteCategory(ctx, id)
}
