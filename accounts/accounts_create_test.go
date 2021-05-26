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

func TestCreateAccountSuccess(t *testing.T) {
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
			StatusCode: http.StatusCreated,
			Body:       resp,
		}, nil
	}

	accCreate := AccountCreate{
		Type:           "accounts",
		ID:             &id,
		OrganisationID: &oID,
		Attributes: &Attributes{
			Country: "GB",
			Name:    []string{"John Doe"},
		},
	}

	accClient, err := NewWithClient(&mock, &client.MockRetrySleeper{})
	require.NoError(t, err)

	ctx := context.Background()

	acc, err := accClient.Create(ctx, &accCreate)

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

func TestCreateAccountErrors(t *testing.T) {
	tests := map[string]struct {
		code int
		body string
		err  error
	}{
		"bad request": {
			code: 400,
			body: `{"error_message": "validation error"}`,
			err: &client.APIError{
				ErrorMessage: "validation error",
				StatusCode:   http.StatusBadRequest,
			},
		},
		"unauthorized": {code: 401, body: "unauthorized",
			err: &client.APIError{
				StatusCode:   401,
				ErrorMessage: "unauthorized",
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

			accCreate := AccountCreate{}
			ctx := context.Background()
			acc, err := accClient.Create(ctx, &accCreate)

			assert.ErrorIs(t, err, tc.err)
			assert.Nil(t, acc)
		})
	}
}

func TestCreateAccountNilAccountCreate(t *testing.T) {
	accClient, err := NewWithClient(&client.MockClient{}, &client.MockRetrySleeper{})
	require.NoError(t, err)

	ctx := context.Background()
	acc, err := accClient.Create(ctx, nil)

	assert.Error(t, err)
	assert.Nil(t, acc)
}
