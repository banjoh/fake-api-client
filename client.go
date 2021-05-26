package client

import (
	"net/http"
	"time"
)

var (
	// DefaultClient is the client all resource API instances
	// default to when instanciated using New()
	DefaultClient HTTPClient
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

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MockClient struct {
	DoImpl func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoImpl(req)
}

type RetrySleeper interface {
	Sleep(d time.Duration)
}

type DefaultRetrySleeper struct{}

func (s *DefaultRetrySleeper) Sleep(d time.Duration) { time.Sleep(d) }

type MockRetrySleeper struct{}

func (s *MockRetrySleeper) Sleep(d time.Duration) {}
