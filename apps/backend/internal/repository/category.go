package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sriniously/tasker/internal/model"
	"github.com/sriniously/tasker/internal/model/category"
	"github.com/sriniously/tasker/internal/server"
)

type CategoryRepository struct {
	server *server.Server
}

func NewCategoryRepository(server *server.Server) *CategoryRepository {
	return &CategoryRepository{server: server}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, userID string,
	payload *category.CreateCategoryPayload,
) (*category.Category, error) {
	stmt := `
		INSERT INTO
			todo_categories (
				user_id,
				name,
				color,
				description
			)
		VALUES
			(
				@user_id,
				@name,
				@color,
				@description
			)
		RETURNING
		*
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":     userID,
		"name":        payload.Name,
		"color":       payload.Color,
		"description": payload.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create category query for user_id=%s name=%s: %w", userID, payload.Name, err)
	}

	categoryItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[category.Category])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:todo_categories for user_id=%s name=%s: %w", userID, payload.Name, err)
	}

	return &categoryItem, nil
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, userID string, categoryID uuid.UUID) (*category.Category, error) {
	stmt := `
		SELECT
			*
		FROM
			todo_categories
		WHERE
			id=@id
			AND user_id=@user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      categoryID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get category by id query for category_id=%s user_id=%s: %w", categoryID.String(), userID, err)
	}

	categoryItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[category.Category])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:todo_categories for category_id=%s user_id=%s: %w", categoryID.String(), userID, err)
	}

	return &categoryItem, nil
}

func (r *CategoryRepository) GetCategories(ctx context.Context, userID string,
	query *category.GetCategoriesQuery,
) (*model.PaginatedResponse[category.Category], error) {
	stmt := `
		SELECT
			*
		FROM
			todo_categories
		WHERE
			user_id=@user_id
	`

	args := pgx.NamedArgs{
		"user_id": userID,
	}

	// Add search filter if provided
	if query.Search != nil {
		stmt += ` AND name ILIKE '%' || @search || '%'`
		args["search"] = *query.Search
	}

	// Add sorting
	sortColumn := "name"
	if query.Sort != nil {
		sortColumn = *query.Sort
	}
	sortOrder := "asc"
	if query.Order != nil {
		sortOrder = *query.Order
	}
	stmt += fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortOrder)

	// Add pagination
	stmt += ` LIMIT @limit OFFSET @offset`
	args["limit"] = *query.Limit
	args["offset"] = (*query.Page - 1) * (*query.Limit)

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get categories query for user_id=%s: %w", userID, err)
	}

	categories, err := pgx.CollectRows(rows, pgx.RowToStructByName[category.Category])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[category.Category]{
				Data:       []category.Category{},
				Page:       *query.Page,
				Limit:      *query.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:todo_categories for user_id=%s: %w", userID, err)
	}

	// Get total count
	countStmt := `
		SELECT
			COUNT(*)
		FROM
			todo_categories
		WHERE
			user_id=@user_id
	`

	countArgs := pgx.NamedArgs{
		"user_id": userID,
	}

	if query.Search != nil {
		countStmt += ` AND name ILIKE '%' || @search || '%'`
		countArgs["search"] = *query.Search
	}

	var total int
	err = r.server.DB.Pool.QueryRow(ctx, countStmt, countArgs).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count of categories for user_id=%s: %w", userID, err)
	}

	return &model.PaginatedResponse[category.Category]{
		Data:       categories,
		Page:       *query.Page,
		Limit:      *query.Limit,
		Total:      total,
		TotalPages: (total + *query.Limit - 1) / *query.Limit,
	}, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, userID string,
	categoryID uuid.UUID, payload *category.UpdateCategoryPayload,
) (*category.Category, error) {
	stmt := `UPDATE todo_categories SET `
	args := pgx.NamedArgs{
		"id":      categoryID,
		"user_id": userID,
	}
	setClauses := []string{}

	if payload.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *payload.Name
	}
	if payload.Color != nil {
		setClauses = append(setClauses, "color = @color")
		args["color"] = *payload.Color
	}
	if payload.Description != nil {
		setClauses = append(setClauses, "description = @description")
		args["description"] = *payload.Description
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	stmt += strings.Join(setClauses, ", ")
	stmt += ` WHERE id = @id AND user_id = @user_id RETURNING *`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update category query for category_id=%s user_id=%s: %w", categoryID.String(), userID, err)
	}

	categoryItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[category.Category])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:todo_categories for category_id=%s user_id=%s: %w", categoryID.String(), userID, err)
	}

	return &categoryItem, nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, userID string, categoryID uuid.UUID) error {
	result, err := r.server.DB.Pool.Exec(ctx, `
		DELETE FROM todo_categories
		WHERE id = @id AND user_id = @user_id
	`, pgx.NamedArgs{
		"id":      categoryID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}
