package service

import (
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sriniously/tasker/internal/errs"
	"github.com/sriniously/tasker/internal/lib/aws"
	"github.com/sriniously/tasker/internal/middleware"
	"github.com/sriniously/tasker/internal/model"
	"github.com/sriniously/tasker/internal/model/todo"
	"github.com/sriniously/tasker/internal/repository"
	"github.com/sriniously/tasker/internal/server"
)

type TodoService struct {
	server       *server.Server
	todoRepo     *repository.TodoRepository
	categoryRepo *repository.CategoryRepository
	awsClient    *aws.AWS
}

func NewTodoService(server *server.Server, todoRepo *repository.TodoRepository,
	categoryRepo *repository.CategoryRepository, awsClient *aws.AWS,
) *TodoService {
	return &TodoService{
		server:       server,
		todoRepo:     todoRepo,
		categoryRepo: categoryRepo,
		awsClient:    awsClient,
	}
}

func (s *TodoService) CreateTodo(ctx echo.Context, userID string, payload *todo.CreateTodoPayload) (*todo.Todo, error) {
	logger := middleware.GetLogger(ctx)

	// Validate parent todo exists and belongs to user (if provided)
	if payload.ParentTodoID != nil {
		parentTodo, err := s.todoRepo.CheckTodoExists(ctx.Request().Context(), userID, *payload.ParentTodoID)
		if err != nil {
			logger.Error().Err(err).Msg("parent todo validation failed")
			return nil, err
		}

		if !parentTodo.CanHaveChildren() {
			err := errs.NewBadRequestError("Parent todo cannot have children (subtasks can't have subtasks)", false, nil, nil, nil)
			logger.Warn().Msg("parent todo cannot have children")
			return nil, err
		}
	}

	// Validate category exists and belongs to user (if provided)
	if payload.CategoryID != nil {
		_, err := s.categoryRepo.GetCategoryByID(ctx.Request().Context(), userID, *payload.CategoryID)
		if err != nil {
			logger.Error().Err(err).Msg("category validation failed")
			return nil, err
		}
	}

	todoItem, err := s.todoRepo.CreateTodo(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create todo")
		return nil, err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "todo_created").
		Str("todo_id", todoItem.ID.String()).
		Str("title", todoItem.Title).
		Str("category_id", func() string {
			if todoItem.CategoryID != nil {
				return todoItem.CategoryID.String()
			}
			return ""
		}()).
		Str("priority", string(todoItem.Priority)).
		Msg("Todo created successfully")

	return todoItem, nil
}

func (s *TodoService) GetTodoByID(ctx echo.Context, userID string, todoID uuid.UUID) (*todo.PopulatedTodo, error) {
	logger := middleware.GetLogger(ctx)

	todoItem, err := s.todoRepo.GetTodoByID(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch todo by ID")
		return nil, err
	}

	return todoItem, nil
}

func (s *TodoService) GetTodos(ctx echo.Context, userID string, query *todo.GetTodosQuery) (*model.PaginatedResponse[todo.PopulatedTodo], error) {
	logger := middleware.GetLogger(ctx)

	result, err := s.todoRepo.GetTodos(ctx.Request().Context(), userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch todos")
		return nil, err
	}

	return result, nil
}

func (s *TodoService) UpdateTodo(ctx echo.Context, userID string, payload *todo.UpdateTodoPayload) (*todo.Todo, error) {
	logger := middleware.GetLogger(ctx)

	// Validate parent todo exists and belongs to user (if provided)
	if payload.ParentTodoID != nil {
		parentTodo, err := s.todoRepo.CheckTodoExists(ctx.Request().Context(), userID, *payload.ParentTodoID)
		if err != nil {
			logger.Error().Err(err).Msg("parent todo validation failed")
			return nil, err
		}

		if parentTodo.ID == payload.ID {
			err := errs.NewBadRequestError("Todo cannot be its own parent", false, nil, nil, nil)
			logger.Warn().Msg("todo cannot be its own parent")
			return nil, err
		}

		if !parentTodo.CanHaveChildren() {
			err := errs.NewBadRequestError("Parent todo cannot have children (subtasks can't have subtasks)", false, nil, nil, nil)
			logger.Warn().Msg("parent todo cannot have children")
			return nil, err
		}

		logger.Debug().Msg("parent todo validation passed")
	}

	// Validate category exists and belongs to user (if provided)
	if payload.CategoryID != nil {
		_, err := s.categoryRepo.GetCategoryByID(ctx.Request().Context(), userID, *payload.CategoryID)
		if err != nil {
			logger.Error().Err(err).Msg("category validation failed")
			return nil, err
		}

		logger.Debug().Msg("category validation passed")
	}

	updatedTodo, err := s.todoRepo.UpdateTodo(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update todo")
		return nil, err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "todo_updated").
		Str("todo_id", updatedTodo.ID.String()).
		Str("title", updatedTodo.Title).
		Str("category_id", func() string {
			if updatedTodo.CategoryID != nil {
				return updatedTodo.CategoryID.String()
			}
			return ""
		}()).
		Str("priority", string(updatedTodo.Priority)).
		Str("status", string(updatedTodo.Status)).
		Msg("Todo updated successfully")

	return updatedTodo, nil
}

func (s *TodoService) DeleteTodo(ctx echo.Context, userID string, todoID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	err := s.todoRepo.DeleteTodo(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete todo")
		return err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "todo_deleted").
		Str("todo_id", todoID.String()).
		Msg("Todo deleted successfully")

	return nil
}

func (s *TodoService) GetTodoStats(ctx echo.Context, userID string) (*todo.TodoStats, error) {
	logger := middleware.GetLogger(ctx)

	stats, err := s.todoRepo.GetTodoStats(ctx.Request().Context(), userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch todo statistics")
		return nil, err
	}

	return stats, nil
}

func (s *TodoService) UploadTodoAttachment(
	ctx echo.Context,
	userID string,
	todoID uuid.UUID,
	file *multipart.FileHeader,
) (*todo.TodoAttachment, error) {
	logger := middleware.GetLogger(ctx)

	// Verify todo exists and belongs to user
	_, err := s.todoRepo.CheckTodoExists(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("todo validation failed")
		return nil, err
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		logger.Error().Err(err).Msg("failed to open uploaded file")
		return nil, errs.NewBadRequestError("failed to open uploaded file", false, nil, nil, nil)
	}
	defer src.Close()

	// Upload to S3
	s3Key, err := s.awsClient.S3.UploadFile(
		ctx.Request().Context(),
		s.server.Config.AWS.UploadBucket,
		"todos/attachments/"+file.Filename,
		src,
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to upload file to S3")
		return nil, errors.Wrap(err, "failed to upload file")
	}

	// Detect MIME type
	src, err = file.Open()
	if err != nil {
		logger.Error().Err(err).Msg("failed to reopen file for MIME detection")
		return nil, errs.NewBadRequestError("failed to process file", false, nil, nil, nil)
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		logger.Error().Err(err).Msg("failed to read file for MIME detection")
		return nil, errs.NewBadRequestError("failed to process file", false, nil, nil, nil)
	}
	mimeType := http.DetectContentType(buffer)

	// Create attachment record
	attachment, err := s.todoRepo.UploadTodoAttachment(
		ctx.Request().Context(),
		todoID,
		userID,
		s3Key,
		file.Filename,
		file.Size,
		mimeType,
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create attachment record")
		return nil, err
	}

	logger.Info().
		Str("attachment_id", attachment.ID.String()).
		Str("s3_key", s3Key).
		Msg("uploaded todo attachment")

	return attachment, nil
}

func (s *TodoService) DeleteTodoAttachment(
	ctx echo.Context,
	userID string,
	todoID uuid.UUID,
	attachmentID uuid.UUID,
) error {
	logger := middleware.GetLogger(ctx)

	// Verify todo exists and belongs to user
	_, err := s.todoRepo.CheckTodoExists(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("todo validation failed")
		return err
	}

	// Get attachment details for S3 deletion
	attachment, err := s.todoRepo.GetTodoAttachment(
		ctx.Request().Context(),
		todoID,
		attachmentID,
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get attachment details")
		return err
	}

	// Delete attachment record
	err = s.todoRepo.DeleteTodoAttachment(
		ctx.Request().Context(),
		todoID,
		attachmentID,
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete attachment record")
		return err
	}

	// Delete from S3 asynchronously
	go func() {
		err := s.awsClient.S3.DeleteObject(
			ctx.Request().Context(),
			s.server.Config.AWS.UploadBucket,
			attachment.DownloadKey,
		)
		if err != nil {
			s.server.Logger.Error().
				Err(err).
				Str("s3_key", attachment.DownloadKey).
				Msg("failed to delete attachment from S3")
		}
	}()

	logger.Info().Msg("deleted todo attachment")

	return nil
}

func (s *TodoService) GetAttachmentPresignedURL(
	ctx echo.Context,
	userID string,
	todoID uuid.UUID,
	attachmentID uuid.UUID,
) (string, error) {
	logger := middleware.GetLogger(ctx)

	// Verify todo exists and belongs to user
	_, err := s.todoRepo.CheckTodoExists(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("todo validation failed")
		return "", err
	}

	// Get attachment details
	attachment, err := s.todoRepo.GetTodoAttachment(
		ctx.Request().Context(),
		todoID,
		attachmentID,
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get attachment details")
		return "", err
	}

	// Generate presigned URL
	url, err := s.awsClient.S3.CreatePresignedUrl(
		ctx.Request().Context(),
		s.server.Config.AWS.UploadBucket,
		attachment.DownloadKey,
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate presigned URL")
		return "", err
	}

	return url, nil
}
