package category

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ------------------------------------------------------------
type CreateCategoryPayload struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Color       string  `json:"color" validate:"required,hexcolor"`
	Description *string `json:"description" validate:"omitempty,max=255"`
}

func (p *CreateCategoryPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ------------------------------------------------------------

type UpdateCategoryPayload struct {
	ID          uuid.UUID `param:"id" validate:"required,uuid"`
	Name        *string   `json:"name" validate:"omitempty,min=1,max=100"`
	Color       *string   `json:"color" validate:"omitempty,hexcolor"`
	Description *string   `json:"description" validate:"omitempty,max=255"`
}

func (p *UpdateCategoryPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ------------------------------------------------------------

type GetCategoriesQuery struct {
	Page   *int    `query:"page" validate:"omitempty,min=1"`
	Limit  *int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Sort   *string `query:"sort" validate:"omitempty,oneof=created_at updated_at name"`
	Order  *string `query:"order" validate:"omitempty,oneof=asc desc"`
	Search *string `query:"search" validate:"omitempty,min=1"`
}

func (q *GetCategoriesQuery) Validate() error {
	validate := validator.New()

	if err := validate.Struct(q); err != nil {
		return err
	}

	// Set defaults
	if q.Page == nil {
		defaultPage := 1
		q.Page = &defaultPage
	}
	if q.Limit == nil {
		defaultLimit := 50
		q.Limit = &defaultLimit
	}
	if q.Sort == nil {
		defaultSort := "name"
		q.Sort = &defaultSort
	}
	if q.Order == nil {
		defaultOrder := "asc"
		q.Order = &defaultOrder
	}

	return nil
}


type DeleteCategoryPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *DeleteCategoryPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
