package accounts

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	client "github.com/banjoh/fake-api-client"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAccountSuccess(t *testing.T) {
	mock := client.MockClient{}
	mock.DoImpl = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Body:       io.NopCloser(bytes.NewReader([]byte(""))),
		}, nil
	}

	accClient := NewWithClient(&mock)

	ctx := context.Background()
	err := accClient.Delete(ctx, uuid.New(), 0)

	assert.Nil(t, err)
}

func TestDeleteAccountErrors(t *testing.T) {
	tests := map[string]struct {
		code int
		body string
		err  error
	}{
		"not found": {code: 404, body: "", err: &client.APIError{StatusCode: 404}},
		"conflict": {
			code: 409,
			body: `{"error_message": "invalid version"}`,
			err: &client.APIError{
				ErrorMessage: "invalid version",
				ErrorCode:    "",
				StatusCode:   http.StatusConflict,
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			body := io.NopCloser(bytes.NewReader([]byte(tc.body)))
			mock := client.MockClient{}
			mock.DoImpl = func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: tc.code,
					Body:       body,
				}, nil
			}

			accClient := NewWithClient(&mock)

			ctx := context.Background()
			err := accClient.Delete(ctx, uuid.New(), 0)

			assert.ErrorIs(t, err, tc.err)
		})
	}
}
