package client

import (
	"math/rand"
	"net/http"
	"time"
)

var (
	// DefaultClient is the client all resource API instances
	// default to when instantiated using New(). The variable
	// can also be used indipendently
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
		Timeout:   time.Second * 60, // Connect + Read timeout
		Transport: t,
	}

	// Random number generator for generating retry jitters
	rand.Seed(time.Now().UTC().UnixNano())
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
