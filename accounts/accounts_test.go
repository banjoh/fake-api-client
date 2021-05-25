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
)

func TestCreateAccountSuccess(t *testing.T) {
	id := uuid.New()
	oId := uuid.New()

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
	  }`, id, oId)
	resp := io.NopCloser(bytes.NewReader([]byte(json)))

	mock := client.MockClient{}
	mock.DoImpl = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       resp,
		}, nil
	}

	accCreate := AccountCreate{
		Type:           "accounts",
		ID:             &id,
		OrganisationID: &oId,
		Attributes: &Attributes{
			Country: "GB",
			Name:    []string{"John Doe"},
		},
	}

	accClient := NewWithClient(&mock)
	ctx := context.Background()

	acc, err := accClient.Create(ctx, &accCreate)

	assert.NotNil(t, acc)
	assert.NoError(t, err)
	assert.Equal(t, id, *acc.ID)
	assert.Equal(t, oId, *acc.OrganisationID)
	assert.Equal(t, "accounts", acc.Type)
	assert.Equal(t, "GB", acc.Attributes.Country)
	assert.Equal(t, []string{"John Doe"}, acc.Attributes.Name)
	assert.Equal(t, "2021-05-25T04:29:11.898Z", acc.CreatedOn)
	assert.Equal(t, "2021-05-25T04:29:11.898Z", acc.ModifiedOn)
}

func TestFetchAccountSuccess(t *testing.T) {
	id := uuid.New()
	oId := uuid.New()

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
	  }`, id, oId)
	resp := io.NopCloser(bytes.NewReader([]byte(json)))

	mock := client.MockClient{}
	mock.DoImpl = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       resp,
		}, nil
	}

	accClient := NewWithClient(&mock)
	ctx := context.Background()

	acc, err := accClient.Fetch(ctx, id)

	assert.NotNil(t, acc)
	assert.NoError(t, err)
	assert.Equal(t, id, *acc.ID)
	assert.Equal(t, oId, *acc.OrganisationID)
	assert.Equal(t, "accounts", acc.Type)
	assert.Equal(t, "GB", acc.Attributes.Country)
	assert.Equal(t, []string{"John Doe"}, acc.Attributes.Name)
	assert.Equal(t, "2021-05-25T04:29:11.898Z", acc.CreatedOn)
	assert.Equal(t, "2021-05-25T04:29:11.898Z", acc.ModifiedOn)
}

func TestDeleteAccountSuccess(t *testing.T) {
	mock := client.MockClient{}
	mock.DoImpl = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 204,
		}, nil
	}

	accClient := NewWithClient(&mock)
	ctx := context.Background()

	err := accClient.Delete(ctx, uuid.New(), 0)

	assert.NoError(t, err)
}

func TestDeleteAccountFailures(t *testing.T) {
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
