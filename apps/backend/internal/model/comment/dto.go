package comment

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ------------------------------------------------------------

type AddCommentPayload struct {
	TodoID  uuid.UUID `param:"id" validate:"required,uuid"`
	Content string    `json:"content" validate:"required,min=1,max=1000"`
}

func (p *AddCommentPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ------------------------------------------------------------

type GetCommentsByTodoIDPayload struct {
	TodoID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *GetCommentsByTodoIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ------------------------------------------------------------

type UpdateCommentPayload struct {
	ID      uuid.UUID `param:"id" validate:"required,uuid"`
	Content string    `json:"content" validate:"required,min=1,max=1000"`
}

func (p *UpdateCommentPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ------------------------------------------------------------

type DeleteCommentPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *DeleteCommentPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
