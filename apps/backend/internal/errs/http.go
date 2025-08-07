package errs

import (
	"strings"
)

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ActionType string

const (
	ActionTypeRedirect ActionType = "redirect"
)

type Action struct {
	Type    ActionType `json:"type"`
	Message string     `json:"message"`
	Value   string     `json:"value"`
}

type HTTPError struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Status   int    `json:"status"`
	Override bool   `json:"override"`
	// field level errors
	Errors []FieldError `json:"errors"`
	// action to be taken
	Action *Action `json:"action"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) Is(target error) bool {
	_, ok := target.(*HTTPError)

	return ok
}

func (e *HTTPError) WithMessage(message string) *HTTPError {
	return &HTTPError{
		Code:     e.Code,
		Message:  message,
		Status:   e.Status,
		Override: e.Override,
		Errors:   e.Errors,
		Action:   e.Action,
	}
}

func MakeUpperCaseWithUnderscores(str string) string {
	return strings.ToUpper(strings.ReplaceAll(str, " ", "_"))
}
