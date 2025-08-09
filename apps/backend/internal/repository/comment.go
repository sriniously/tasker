package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sriniously/tasker/internal/model/comment"
	"github.com/sriniously/tasker/internal/server"
)

type CommentRepository struct {
	server *server.Server
}

func NewCommentRepository(server *server.Server) *CommentRepository {
	return &CommentRepository{server: server}
}

func (r *CommentRepository) AddComment(ctx context.Context, userID string, todoID uuid.UUID,
	payload *comment.AddCommentPayload,
) (*comment.Comment, error) {
	stmt := `
		INSERT INTO
			todo_comments (
				todo_id,
				user_id,
				content
			)
		VALUES
			(
				@todo_id,
				@user_id,
				@content
			)
		RETURNING
		*
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"todo_id": todoID,
		"user_id": userID,
		"content": payload.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute add comment query for todo_id=%s user_id=%s: %w", todoID.String(), userID, err)
	}

	commentItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[comment.Comment])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:todo_comments for todo_id=%s user_id=%s: %w", todoID.String(), userID, err)
	}

	return &commentItem, nil
}

func (r *CommentRepository) GetCommentsByTodoID(ctx context.Context, userID string, todoID uuid.UUID) ([]comment.Comment, error) {
	stmt := `
		SELECT
			*
		FROM
			todo_comments
		WHERE
			todo_id=@todo_id
			AND user_id=@user_id
		ORDER BY
			created_at ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"todo_id": todoID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get comments by todo id query for todo_id=%s user_id=%s: %w", todoID.String(), userID, err)
	}

	comments, err := pgx.CollectRows(rows, pgx.RowToStructByName[comment.Comment])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows from table:todo_comments for todo_id=%s user_id=%s: %w", todoID.String(), userID, err)
	}

	return comments, nil
}

func (r *CommentRepository) GetCommentByID(ctx context.Context, userID string, commentID uuid.UUID) (*comment.Comment, error) {
	stmt := `
		SELECT
			*
		FROM
			todo_comments
		WHERE
			id=@id
			AND user_id=@user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      commentID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get comment by id query for comment_id=%s user_id=%s: %w", commentID.String(), userID, err)
	}

	commentItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[comment.Comment])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:todo_comments for comment_id=%s user_id=%s: %w", commentID.String(), userID, err)
	}

	return &commentItem, nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, userID string, commentID uuid.UUID, content string) (*comment.Comment, error) {
	stmt := `
		UPDATE
			todo_comments
		SET
			content=@content
		WHERE
			id=@id
			AND user_id=@user_id
		RETURNING
		*
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      commentID,
		"user_id": userID,
		"content": content,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute update comment query for comment_id=%s user_id=%s: %w", commentID.String(), userID, err)
	}

	commentItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[comment.Comment])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:todo_comments for comment_id=%s user_id=%s: %w", commentID.String(), userID, err)
	}

	return &commentItem, nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, userID string, commentID uuid.UUID) error {
	result, err := r.server.DB.Pool.Exec(ctx, `
		DELETE FROM todo_comments
		WHERE id = @id AND user_id = @user_id
	`, pgx.NamedArgs{
		"id":      commentID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}
