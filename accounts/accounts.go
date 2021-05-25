package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	client "github.com/banjoh/fake-api-client"
	"github.com/google/uuid"
)

const (
	defultBaseURL      = "http://localhost:8080"
	accountsPath       = "v1/organisation/accounts"
	defaultContentType = "application/vnd.api+json"
)

func New() *AccountsResource {
	return NewWithClient(client.DefaultClient)
}

func NewWithClient(c client.HttpClient) *AccountsResource {
	return &AccountsResource{
		BaseURL: defultBaseURL,
		client:  c,
	}
}

func (r *AccountsResource) Create(ctx context.Context, acc *AccountCreate) (*Account, error) {
	dto := AccountCreateDTO{Data: *acc}

	data, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s", r.BaseURL, accountsPath)

	// The io stream will be closed by the client
	body := bytes.NewReader(data)
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	setPostDefaultHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		return unmarshalAccount(resp)
	}

	return nil, unmarshalErrorResponse(resp)
}

func (r *AccountsResource) Fetch(ctx context.Context, accID uuid.UUID) (*Account, error) {
	url := fmt.Sprintf("%s/%s/%s", r.BaseURL, accountsPath, accID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	setDefaultHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return unmarshalAccount(resp)
	}

	return nil, unmarshalErrorResponse(resp)
}

func (r *AccountsResource) Delete(ctx context.Context, accID uuid.UUID, version int) error {
	url := fmt.Sprintf("%s/%s/%s?version=%d", r.BaseURL,
		accountsPath, accID, version,
	)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	setDefaultHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		// Deletion succeded
		return nil
	}

	return unmarshalErrorResponse(resp)
}

func unmarshalAccount(resp *http.Response) (*Account, error) {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var got AccountDTO
	err = json.Unmarshal(b, &got)
	if err != nil {
		return nil, err
	}

	return &got.Data, nil
}

func unmarshalErrorResponse(resp *http.Response) error {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(b) == "" {
		return &client.APIError{
			StatusCode: resp.StatusCode,
		}
	}

	var apiErr client.APIError
	err = json.Unmarshal(b, &apiErr)
	if err != nil {
		return err
	}

	apiErr.StatusCode = resp.StatusCode
	return &apiErr
}

func setDefaultHeaders(req *http.Request) {
	req.Header.Set("Accept", defaultContentType)
	ts := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	req.Header.Set("Date", ts)
}

func setPostDefaultHeaders(req *http.Request) {
	req.Header.Set("Content-Type", defaultContentType)

	setDefaultHeaders(req)
}
