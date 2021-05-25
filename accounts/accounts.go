package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"time"

	client "github.com/banjoh/fake-api-client"
	"github.com/google/uuid"
)

const (
	defultBaseURL         = "http://localhost:8080"
	accountsPath          = "v1/organisation/accounts"
	defaultContentType    = "application/vnd.api+json"
	initialBackoffSeconds = float64(0)
	maxTimeoutSeconds     = float64(15)
	defaultRetryCount     = 5
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// New creates a new instance of the accounts resource API
// This client utilizes a default http client
func New() (*Resource, error) {
	return NewWithClient(client.DefaultClient)
}

// NewWithClient creates a new instance of the accounts resource API
// This client utilizes a dependency injected http client
func NewWithClient(c client.HTTPClient) (*Resource, error) {
	if c == nil {
		return nil, fmt.Errorf("accounts.NewWithClient: nil client.HTTPClient")
	}
	return &Resource{
		BaseURL: defultBaseURL,
		client:  c,
	}, nil
}

func (r *Resource) Create(ctx context.Context, acc *AccountCreate) (*Account, error) {
	if ctx == nil {
		return nil, fmt.Errorf("accounts.Create: nil Context")
	}

	if acc == nil {
		return nil, fmt.Errorf("nil AccountCreate")
	}

	dto := AccountCreateDTO{Data: *acc}

	data, err := json.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("marshalling error: %w", err)
	}

	url := fmt.Sprintf("%s/%s", r.BaseURL, accountsPath)

	// The io stream will be closed by the client
	body := bytes.NewReader(data)
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	setPostDefaultHeaders(req)

	// We only retry idempotent requests
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		return unmarshalAccount(resp)
	}

	return nil, unmarshalErrorResponse(resp)
}

func (r *Resource) Fetch(ctx context.Context, accID uuid.UUID) (*Account, error) {
	if ctx == nil {
		return nil, fmt.Errorf("accounts.Fetch: nil Context")
	}

	url := fmt.Sprintf("%s/%s/%s", r.BaseURL, accountsPath, accID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	setDefaultHeaders(req)

	resp, err := retriedDo(req, r.client)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return unmarshalAccount(resp)
	}

	return nil, unmarshalErrorResponse(resp)
}

func (r *Resource) Delete(ctx context.Context, accID uuid.UUID, version int) error {
	if ctx == nil {
		return fmt.Errorf("accounts.Delete: nil Context")
	}

	url := fmt.Sprintf("%s/%s/%s?version=%d", r.BaseURL,
		accountsPath, accID, version,
	)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	setDefaultHeaders(req)

	resp, err := retriedDo(req, r.client)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
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
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var got AccountDTO
	err = json.Unmarshal(b, &got)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling err: %w", err)
	}

	return &got.Data, nil
}

func unmarshalErrorResponse(resp *http.Response) error {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if string(b) == "" {
		return &client.APIError{
			StatusCode: resp.StatusCode,
		}
	}

	var apiErr client.APIError
	err = json.Unmarshal(b, &apiErr)
	if err != nil {
		return &client.APIError{
			StatusCode:   resp.StatusCode,
			ErrorMessage: string(b),
		}
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

func containsStatus(statuses []int, status int) bool {
	for _, s := range statuses {
		if s == status {
			return true
		}
	}

	return false
}

func isTemporaryOrTimeout(err error) bool {
	if ne, ok := err.(net.Error); ok && (ne.Temporary() || ne.Timeout()) { // nolint: errorlint
		return true
	}

	return false
}

func retriedDo(req *http.Request, c client.HTTPClient) (*http.Response, error) {
	tryableStatuses := []int{500, 502, 503, 504}

	var resp *http.Response
	var err error
	durationSecs := nextBackoff(initialBackoffSeconds)
	for i := 0; i < defaultRetryCount; i++ {
		resp, err = c.Do(req)
		if err != nil {
			// Retry network errors deemed retryable
			if !isTemporaryOrTimeout(err) {
				duration := time.Duration(durationSecs * float64(time.Millisecond))
				time.Sleep(duration)

				// Get the next back off duration
				durationSecs = nextBackoff(durationSecs)
				continue
			} else {
				return nil, err
			}
		}

		if !containsStatus(tryableStatuses, resp.StatusCode) {
			return resp, nil
		}

		duration := time.Duration(durationSecs * float64(time.Millisecond))
		time.Sleep(duration)

		// Get the next back off duration
		durationSecs = nextBackoff(durationSecs)
	}

	return resp, nil
}

func nextBackoff(curr float64) float64 {
	if curr <= 0 {
		return 1 + rand.Float64() // nolint: gosec
	}

	next := math.Floor(curr)

	if next < maxTimeoutSeconds {
		next *= 2
	}

	// Next back-off + some jitter. The jitter is necessary so as to
	// to avoid many clients retrying at the exact same time. The many
	// concurrent requests can exhaust the servers TCP connections
	return next + rand.Float64() // nolint: gosec
}
