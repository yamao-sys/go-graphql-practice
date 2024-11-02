package view

import (
	"fmt"
	"net/http"
)

type ViewError struct {
	Code    int64
	Message error
}

func (e ViewError) Error() string {
	return fmt.Sprintf("[%d]%s", e.Code, e.Message)
}

func NewBadRequestUserView(err error) ViewError {
	return ViewError{
		Code:    http.StatusBadRequest,
		Message: err,
	}
}

func NewInternalServerErrorUserView(err error) ViewError {
	return ViewError{
		Code:    http.StatusInternalServerError,
		Message: err,
	}
}
