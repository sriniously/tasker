package todo

import (
	"time"

	"github.com/google/uuid"
	"github.com/sriniously/tasker/internal/model"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusActive    Status = "active"
	StatusCompleted Status = "completed"
	StatusArchived  Status = "archived"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Todo struct {
	model.Base
	UserID       string     `json:"userId" db:"user_id"`
	Title        string     `json:"title" db:"title"`
	Description  *string    `json:"description" db:"description"`
	Status       Status     `json:"status" db:"status"`
	Priority     Priority   `json:"priority" db:"priority"`
	DueDate      *time.Time `json:"dueDate" db:"due_date"`
	CompletedAt  *time.Time `json:"completedAt" db:"completed_at"`
	ParentTodoID *uuid.UUID `json:"parentTodoId" db:"parent_todo_id"`
	CategoryID   *uuid.UUID `json:"categoryId" db:"category_id"`
	Metadata     *Metadata  `json:"metadata" db:"metadata"`
	SortOrder    int        `json:"sortOrder" db:"sort_order"`
}

type Metadata struct {
	Tags       []string `json:"tags"`
	Reminder   *string  `json:"reminder"`
	Color      *string  `json:"color"`
	Difficulty *int     `json:"difficulty"`
}
