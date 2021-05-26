package accounts

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	client "github.com/banjoh/fake-api-client"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type timeoutErr struct{}

func (e *timeoutErr) Error() string { return "timeout" }
func (e *timeoutErr) Timeout() bool { return true }

type temporaryErr struct{}

func (e *temporaryErr) Error() string   { return "timeout" }
func (e *temporaryErr) Temporary() bool { return true }

func TestReturningNonAPIError(t *testing.T) {
	body := io.NopCloser(bytes.NewReader([]byte("")))
	mock := client.MockClient{}
	genericErr := errors.New("Generic error")

	mock.DoImpl = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 0,
			Body:       body,
		}, genericErr
	}

	accClient, err := NewWithClient(&mock, &client.MockRetrySleeper{})
	require.NoError(t, err)

	ctx := context.Background()
	acc, err := accClient.Fetch(ctx, uuid.New())

	assert.ErrorIs(t, err, genericErr)
	assert.Nil(t, acc)

	err = accClient.Delete(ctx, uuid.New(), 0)
	assert.ErrorIs(t, err, genericErr)

	accCreate := AccountCreate{}
	acc, err = accClient.Create(ctx, &accCreate)

	assert.ErrorIs(t, err, genericErr)
	assert.Nil(t, acc)
}

func TestRetryingCalls(t *testing.T) {
	tests := map[string]struct {
		code        int
		expectedErr error
		err         error
	}{
		"500": {
			code: 500, expectedErr: &client.APIError{StatusCode: 500}, err: nil,
		},
		"502": {
			code: 502, expectedErr: &client.APIError{StatusCode: 502}, err: nil,
		},
		"503": {
			code: 503, expectedErr: &client.APIError{StatusCode: 503}, err: nil,
		},
		"504": {
			code: 504, expectedErr: &client.APIError{StatusCode: 504}, err: nil,
		},
		"network timeout":         {expectedErr: &timeoutErr{}, err: &timeoutErr{}},
		"temporary network error": {expectedErr: &temporaryErr{}, err: &temporaryErr{}},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			body := io.NopCloser(bytes.NewReader([]byte("")))
			mock := client.MockClient{}

			calls := 0
			mock.DoImpl = func(*http.Request) (*http.Response, error) {
				calls++
				return &http.Response{
					StatusCode: tc.code,
					Body:       body,
				}, tc.err
			}

			accClient, err := NewWithClient(&mock, &client.MockRetrySleeper{})
			require.NoError(t, err)

			ctx := context.Background()
			acc, err := accClient.Fetch(ctx, uuid.New())

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Nil(t, acc)

			assert.Equal(t, RetryCount, calls)
		})
	}
}
