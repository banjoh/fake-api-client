package client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructAPIError(t *testing.T) {
	msg := `
	{
		"error_message": "Error message from API",
		"error_code": "b5930880-1001-453b-86bd-2c5e29bd98d7"
	}`

	var apiErr APIError
	err := json.Unmarshal([]byte(msg), &apiErr)

	assert.NoError(t, err)

	apiErr.StatusCode = 404
	expect := APIError{
		ErrorMessage: "Error message from API",
		ErrorCode:    "b5930880-1001-453b-86bd-2c5e29bd98d7",
		StatusCode:   404,
	}
	assert.ErrorIs(t, &apiErr, &expect)
}

func TestConstructEmptyAPIError(t *testing.T) {
	var apiErr APIError
	err := json.Unmarshal([]byte(`{}`), &apiErr)

	assert.NoError(t, err)

	apiErr.StatusCode = 404
	expect := APIError{
		StatusCode: 404,
	}
	assert.ErrorIs(t, &apiErr, &expect)
}
