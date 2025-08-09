package repository

import "github.com/sriniously/tasker/internal/server"

type Repositories struct {
	Todo     *TodoRepository
	Comment  *CommentRepository
	Category *CategoryRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Todo:     NewTodoRepository(s),
		Comment:  NewCommentRepository(s),
		Category: NewCategoryRepository(s),
	}
}
