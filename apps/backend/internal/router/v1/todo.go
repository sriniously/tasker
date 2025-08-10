package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/sriniously/tasker/internal/handler"
	"github.com/sriniously/tasker/internal/middleware"
)

func registerTodoRoutes(r *echo.Group, h *handler.TodoHandler, ch *handler.CommentHandler, auth *middleware.AuthMiddleware) {
	// Todo operations
	todos := r.Group("/todos")
	todos.Use(auth.RequireAuth)

	// Collection operations
	todos.POST("", h.CreateTodo)
	todos.GET("", h.GetTodos)
	todos.GET("/stats", h.GetTodoStats)

	// Individual todo operations
	dynamicTodo := todos.Group("/:id")
	dynamicTodo.GET("", h.GetTodoByID)
	dynamicTodo.PATCH("", h.UpdateTodo)
	dynamicTodo.DELETE("", h.DeleteTodo)

	// Todo comments
	todoComments := dynamicTodo.Group("/comments")
	todoComments.POST("", ch.AddComment)
	todoComments.GET("", ch.GetCommentsByTodoID)

	// Todo attachments
	todoAttachments := dynamicTodo.Group("/attachments")
	todoAttachments.POST("", h.UploadTodoAttachment)
	todoAttachments.DELETE("/:attachmentId", h.DeleteTodoAttachment)
	todoAttachments.GET("/:attachmentId/download", h.GetAttachmentPresignedURL)
}
