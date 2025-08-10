package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sriniously/tasker/internal/errs"
	"github.com/sriniously/tasker/internal/middleware"
	"github.com/sriniously/tasker/internal/model"
	"github.com/sriniously/tasker/internal/model/todo"
	"github.com/sriniously/tasker/internal/server"
	"github.com/sriniously/tasker/internal/service"
)

type TodoHandler struct {
	Handler
	todoService *service.TodoService
}

func NewTodoHandler(s *server.Server, todoService *service.TodoService) *TodoHandler {
	return &TodoHandler{
		Handler:     NewHandler(s),
		todoService: todoService,
	}
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *todo.CreateTodoPayload) (*todo.Todo, error) {
			userID := middleware.GetUserID(c)
			return h.todoService.CreateTodo(c, userID, payload)
		},
		http.StatusCreated,
		&todo.CreateTodoPayload{},
	)(c)
}

func (h *TodoHandler) GetTodoByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *todo.GetTodoByIDPayload) (*todo.PopulatedTodo, error) {
			userID := middleware.GetUserID(c)
			return h.todoService.GetTodoByID(c, userID, payload.ID)
		},
		http.StatusOK,
		&todo.GetTodoByIDPayload{},
	)(c)
}

func (h *TodoHandler) GetTodos(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, query *todo.GetTodosQuery) (*model.PaginatedResponse[todo.PopulatedTodo], error) {
			userID := middleware.GetUserID(c)
			return h.todoService.GetTodos(c, userID, query)
		},
		http.StatusOK,
		&todo.GetTodosQuery{},
	)(c)
}

func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *todo.UpdateTodoPayload) (*todo.Todo, error) {
			userID := middleware.GetUserID(c)
			return h.todoService.UpdateTodo(c, userID, payload)
		},
		http.StatusOK,
		&todo.UpdateTodoPayload{},
	)(c)
}

func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *todo.DeleteTodoPayload) error {
			userID := middleware.GetUserID(c)
			return h.todoService.DeleteTodo(c, userID, payload.ID)
		},
		http.StatusNoContent,
		&todo.DeleteTodoPayload{},
	)(c)
}

func (h *TodoHandler) GetTodoStats(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *todo.GetTodoStatsPayload) (*todo.TodoStats, error) {
			userID := middleware.GetUserID(c)
			return h.todoService.GetTodoStats(c, userID)
		},
		http.StatusOK,
		&todo.GetTodoStatsPayload{},
	)(c)
}

func (h *TodoHandler) UploadTodoAttachment(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *todo.UploadTodoAttachmentPayload) (*todo.TodoAttachment, error) {
			userID := middleware.GetUserID(c)

			form, err := c.MultipartForm()
			if err != nil {
				return nil, errs.NewBadRequestError("multipart form not found", false, nil, nil, nil)
			}

			files := form.File["file"]
			if len(files) == 0 {
				return nil, errs.NewBadRequestError("no file found", false, nil, nil, nil)
			}

			if len(files) > 1 {
				return nil, errs.NewBadRequestError("only one file allowed per upload", false, nil, nil, nil)
			}

			return h.todoService.UploadTodoAttachment(c, userID, payload.TodoID, files[0])
		},
		http.StatusCreated,
		&todo.UploadTodoAttachmentPayload{},
	)(c)
}

func (h *TodoHandler) DeleteTodoAttachment(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *todo.DeleteTodoAttachmentPayload) error {
			userID := middleware.GetUserID(c)
			return h.todoService.DeleteTodoAttachment(c, userID, payload.TodoID, payload.AttachmentID)
		},
		http.StatusNoContent,
		&todo.DeleteTodoAttachmentPayload{},
	)(c)
}

func (h *TodoHandler) GetAttachmentPresignedURL(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *todo.GetAttachmentPresignedURLPayload) (*struct {
			URL string `json:"url"`
		}, error,
		) {
			userID := middleware.GetUserID(c)
			url, err := h.todoService.GetAttachmentPresignedURL(c, userID, payload.TodoID, payload.AttachmentID)
			if err != nil {
				return nil, err
			}
			return &struct {
				URL string `json:"url"`
			}{URL: url}, nil
		},
		http.StatusOK,
		&todo.GetAttachmentPresignedURLPayload{},
	)(c)
}
