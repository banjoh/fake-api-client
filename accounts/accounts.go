package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"

	client "github.com/banjoh/fake-api-client"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	defultBaseURL         = "http://localhost:8080"
	accountsPath          = "v1/organisation/accounts"
	defaultContentType    = "application/vnd.api+json"
	defaultRetrySleepSecs = 2
	defaultRetryCount     = 5
)

var retryableStatusCodes = []int{500, 502, 503, 504}

// RetryCount denotes the number of times to retry requests
// When RetryCount == 0, requests are not retried
var RetryCount = defaultRetryCount

// RetryDurationSecs is the number of seconds to sleep.
// A random jitter is added to each sleep interval
var RetryDurationSecs float64 = defaultRetrySleepSecs

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// New creates a new instance of the accounts resource API
// This client utilizes a default http client
func New() (*Resource, error) {
	return NewWithClient(client.DefaultClient, &client.DefaultRetrySleeper{})
}

// NewWithClient creates a new instance of the accounts resource API
// This client requires a dependency injected http client and retry sleeper
func NewWithClient(c client.HTTPClient, s client.RetrySleeper) (*Resource, error) {
	if c == nil {
		return nil, fmt.Errorf("accounts.NewWithClient: nil client.HTTPClient")
	}
	if s == nil {
		return nil, fmt.Errorf("accounts.NewWithClient: nil client.RetrySleeper")
	}

	return &Resource{
		BaseURL:      defultBaseURL,
		client:       c,
		retrySleeper: s,
	}, nil
}

// Create an account resource
// This API is not idempotent and will therefore not be retried when errors occur.
// * On success, an *Account is returns an the error will be nil
// * On failure, the returned *Account will be nil. The error variable will contain
//		* client.APIError if the response contained API specific errors
//		* any other error that occured. This includes json marshaling errors,
//		  network specific errors etc
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

	// We only retry idempotent requests i.e GET, DELETE
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

// Fetch an account resource
// This API is idempotent and will therefore be retried when some specific errors occur.
// * On success, the queried account is returned in *Account and the error will be nil
// * On failure, the returned *Account will be nil. The error variable will contain
//		* client.APIError if the response contained API specific errors
//		* any other error that occured. This includes json marshaling errors,
//		  network specific errors etc
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

	resp, err := retriedDo(req, r.client, r.retrySleeper)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return unmarshalAccount(resp)
	}

	return nil, unmarshalErrorResponse(resp)
}

// Delete an account resource
// This API is idempotent and will therefore be retried when some specific errors occur.
// * On success, the account resource will be deleted and the error will be nil
// * On failure, the returned error variable will contain
//		* client.APIError if the response contained API specific errors
//		* any other error that occured. This includes json marshaling errors,
//		  network specific errors etc
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

	resp, err := retriedDo(req, r.client, r.retrySleeper)
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if string(body) == "" {
		return &client.APIError{
			StatusCode: resp.StatusCode,
		}
	}

	var apiErr client.APIError
	err = json.Unmarshal(body, &apiErr)
	if err != nil {
		return &client.APIError{
			StatusCode:   resp.StatusCode,
			ErrorMessage: string(body),
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

// retriedDo implements a simple retry logic for temporary error situations
func retriedDo(req *http.Request, c client.HTTPClient, s client.RetrySleeper) (*http.Response, error) {
	if RetryCount < 1 || RetryDurationSecs <= 0 {
		return c.Do(req)
	}

	var resp *http.Response
	var err error

	for i := 0; i < RetryCount; i++ {

		// sleep + jitter. An additional jitter is necessary so as to
		// avoid many clients retrying at the exact same time. The many
		// concurrent requests can exhaust server TCP connection resources
		duration := (RetryDurationSecs + rand.Float64()) * 1000 // nolint: gosec

		resp, err = c.Do(req)
		if err != nil {
			// Retry network errors deemed retryable
			if !isTemporaryOrTimeout(err) {
				logrus.Debugf("Network error caught. Retry request after %.0fms: err=%s", duration, err)
				s.Sleep(time.Duration(duration) * time.Millisecond)

				continue
			} else {
				return nil, err
			}
		}

		// Retry API errors safe for retrying
		if !containsStatus(retryableStatusCodes, resp.StatusCode) {
			break
		}

		logrus.Debugf("Server responded with error. Retry request after %.0fms: code=%d, status=%s",
			duration, resp.StatusCode, resp.Status,
		)
		s.Sleep(time.Duration(duration) * time.Millisecond)
	}

	return resp, err
}
