package accounts

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	client "github.com/banjoh/fake-api-client"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchAccountSuccess(t *testing.T) {
	id := uuid.New()
	oID := uuid.New()

	json := fmt.Sprintf(`{
		"data": {
		  "type": "accounts",
		  "id": "%s",
		  "version": 0,
		  "organisation_id": "%s",
		  "attributes": {
			"country": "GB",
			"name": [
			  "John Doe"
			]
		  },
		  "created_on": "2021-05-25T04:29:11.898Z",
		  "modified_on": "2021-05-25T04:29:11.898Z"
		}
	  }`, id, oID)
	resp := io.NopCloser(bytes.NewReader([]byte(json)))

	mock := client.MockClient{}
	mock.DoImpl = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       resp,
		}, nil
	}

	accClient, err := NewWithClient(&mock, &client.MockRetrySleeper{})
	require.NoError(t, err)

	ctx := context.Background()
	acc, err := accClient.Fetch(ctx, id)

	assert.NotNil(t, acc)
	assert.NoError(t, err)
	assert.Equal(t, id, *acc.ID)
	assert.Equal(t, oID, *acc.OrganisationID)
	assert.Equal(t, "accounts", acc.Type)
	assert.Equal(t, "GB", acc.Attributes.Country)
	assert.Equal(t, []string{"John Doe"}, acc.Attributes.Name)
	assert.Equal(t, "2021-05-25T04:29:11.898Z", acc.CreatedOn)
	assert.Equal(t, "2021-05-25T04:29:11.898Z", acc.ModifiedOn)
}

func TestFetchAccountErrors(t *testing.T) {
	tests := map[string]struct {
		code int
		body string
		err  error
	}{
		"not found": {code: 404, body: "", err: &client.APIError{StatusCode: 404}},
		"not implemented": {
			code: 501,
			body: "Not Implemented",
			err: &client.APIError{
				ErrorMessage: "Not Implemented",
				ErrorCode:    "",
				StatusCode:   http.StatusNotImplemented,
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

			accClient, err := NewWithClient(&mock, &client.MockRetrySleeper{})
			require.NoError(t, err)

			ctx := context.Background()
			acc, err := accClient.Fetch(ctx, uuid.New())

			assert.ErrorIs(t, err, tc.err)
			assert.Nil(t, acc)
		})
	}
}
