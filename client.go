package client

import (
	"net/http"
	"time"
)

var (
	DefaultClient HttpClient
)

func init() {
	// Some sane http client configuration parameters
	t := http.DefaultTransport.(*http.Transport).Clone()

	// Persistent connection pooling configurations
	t.MaxIdleConns = 50
	t.MaxConnsPerHost = 20
	t.MaxIdleConnsPerHost = 20

	// Close idle connections after a duration
	t.IdleConnTimeout = time.Second * 30

	DefaultClient = &http.Client{
		// TODO: What's the sanest timeout? Consider proxy timeouts
		Timeout:   time.Second * 59, // Connect + Read timeout
		Transport: t,
	}
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MockClient struct {
	DoImpl func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoImpl(req)
}
