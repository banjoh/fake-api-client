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

	accClient := NewWithClient(&mock)
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
