package category

import "github.com/sriniously/tasker/internal/model"

type Category struct {
	model.Base
	UserID      string  `json:"userId" db:"user_id"`
	Name        string  `json:"name" db:"name"`
	Color       string  `json:"color" db:"color"`
	Description *string `json:"description" db:"description"`
}
