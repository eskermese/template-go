package core

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound             = errors.New("not found")
	ErrProductNameDuplicate = errors.New("product name already exists")
	ErrInvalidType          = errors.New("invalid type")
)

type Errors []FieldError

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e Errors) Error() string {
	var b strings.Builder

	for i, fieldError := range e {
		b.WriteString(
			fmt.Sprintf("%s: %s", fieldError.Field, fieldError.Message),
		)

		if i != len(e)-1 {
			b.WriteString(", ")
		}
	}

	return b.String()
}
