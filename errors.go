package client

import (
	"fmt"
)

type APIError struct {
	ErrorMessage string `json:"error_message,omitempty"`
	ErrorCode    string `json:"error_code,omitempty"`
	StatusCode   int
}

func (e *APIError) Error() string {
	return fmt.Sprintf(
		`"error_message": "%s", "error_code": "%s", "status_code": %d`,
		e.ErrorMessage, e.ErrorCode, e.StatusCode,
	)
}

func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)
	if !ok {
		return false
	}
	return e.ErrorMessage == t.ErrorMessage &&
		e.ErrorCode == t.ErrorCode &&
		e.StatusCode == t.StatusCode
}
