package service

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sriniously/tasker/internal/middleware"
	"github.com/sriniously/tasker/internal/model"
	"github.com/sriniously/tasker/internal/model/category"
	"github.com/sriniously/tasker/internal/repository"
	"github.com/sriniously/tasker/internal/server"
)

type CategoryService struct {
	server       *server.Server
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(server *server.Server, categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		server:       server,
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) CreateCategory(ctx echo.Context, userID string,
	payload *category.CreateCategoryPayload,
) (*category.Category, error) {
	logger := middleware.GetLogger(ctx)

	categoryItem, err := s.categoryRepo.CreateCategory(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create category")
		return nil, err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "category_created").
		Str("category_id", categoryItem.ID.String()).
		Str("name", categoryItem.Name).
		Str("color", categoryItem.Color).
		Msg("Category created successfully")

	return categoryItem, nil
}

func (s *CategoryService) GetCategories(ctx echo.Context, userID string,
	query *category.GetCategoriesQuery,
) (*model.PaginatedResponse[category.Category], error) {
	logger := middleware.GetLogger(ctx)

	categories, err := s.categoryRepo.GetCategories(ctx.Request().Context(), userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch categories")
		return nil, err
	}

	return categories, nil
}

func (s *CategoryService) GetCategoryByID(ctx echo.Context, userID string, categoryID uuid.UUID) (*category.Category, error) {
	logger := middleware.GetLogger(ctx)

	categoryItem, err := s.categoryRepo.GetCategoryByID(ctx.Request().Context(), userID, categoryID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch category by ID")
		return nil, err
	}

	return categoryItem, nil
}

func (s *CategoryService) UpdateCategory(ctx echo.Context, userID string, categoryID uuid.UUID,
	payload *category.UpdateCategoryPayload,
) (*category.Category, error) {
	logger := middleware.GetLogger(ctx)

	categoryItem, err := s.categoryRepo.UpdateCategory(ctx.Request().Context(), userID, categoryID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update category")
		return nil, err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "category_updated").
		Str("category_id", categoryItem.ID.String()).
		Str("name", categoryItem.Name).
		Msg("Category updated successfully")

	return categoryItem, nil
}

func (s *CategoryService) DeleteCategory(ctx echo.Context, userID string, categoryID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	err := s.categoryRepo.DeleteCategory(ctx.Request().Context(), userID, categoryID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete category")
		return err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "category_deleted").
		Str("category_id", categoryID.String()).
		Msg("Category deleted successfully")

	return nil
}
