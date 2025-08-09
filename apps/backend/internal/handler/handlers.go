package handler

import (
	"github.com/sriniously/tasker/internal/server"
	"github.com/sriniously/tasker/internal/service"
)

type Handlers struct {
	Health   *HealthHandler
	OpenAPI  *OpenAPIHandler
	Todo     *TodoHandler
	Comment  *CommentHandler
	Category *CategoryHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:   NewHealthHandler(s),
		OpenAPI:  NewOpenAPIHandler(s),
		Todo:     NewTodoHandler(s, services.Todo),
		Category: NewCategoryHandler(s, services.Category),
		Comment:  NewCommentHandler(s, services.Comment),
	}
}
