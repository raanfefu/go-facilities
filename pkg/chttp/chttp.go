package chttp

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/raanfefu/go-facilities/pkg/configure"
	"github.com/sony/gobreaker"
)

const (
	DEFAULT_RETRY_COUNT         = 1
	DEFAULT_RETRY_BACKOFF       = 500 * time.Millisecond
	DEFAULT_BREAKER_INTERVAL    = 10 * time.Second
	DEFAULT_BREAKER_TIMEOUT     = 10 * time.Second
	DEFAULT_BREAKER_PREFIX_NAME = "default"
	DEFAULT_BREAKER_MAX_REQUEST = 4
)

type (
	HttpClient interface {
		GetClient() *http.Client
		GetSettings() *Settings
	}

	RetrySettings struct {
		Retry       uint                             `json:"retry"`
		ShouldRetry func(error, *http.Response) bool `json:"-"`
		Backoff     func(retries uint) time.Duration `json:"-"`
	}

	// Settings configures CircuitBreaker:
	//
	// Name is the name of the CircuitBreaker.
	//
	// MaxRequests is the maximum number of requests allowed to pass through
	// when the CircuitBreaker is half-open.
	// If MaxRequests is 0, the CircuitBreaker allows only 1 request.
	//
	// Interval is the cyclic period of the closed state
	// for the CircuitBreaker to clear the internal Counts.
	// If Interval is less than or equal to 0, the CircuitBreaker doesn't clear internal Counts during the closed state.
	//
	// Timeout is the period of the open state,
	// after which the state of the CircuitBreaker becomes half-open.
	// If Timeout is less than or equal to 0, the timeout value of the CircuitBreaker is set to 60 seconds.
	//
	// ReadyToTrip is called with a copy of Counts whenever a request fails in the closed state.
	// If ReadyToTrip returns true, the CircuitBreaker will be placed into the open state.
	// If ReadyToTrip is nil, default ReadyToTrip is used.
	// Default ReadyToTrip returns true when the number of consecutive failures is more than 5.
	//
	// OnStateChange is called whenever the state of the CircuitBreaker changes.
	//
	// IsSuccessful is called with the error returned from a request.
	// If IsSuccessful returns true, the error is counted as a success.
	// Otherwise the error is counted as a failure.
	// If IsSuccessful is nil, default IsSuccessful is used, which returns false for all non-nil errors.
	BreakerSettings struct {
		Name          string                                                      `json:"name"`
		MaxRequests   uint32                                                      `json:"max-request"`
		Interval      time.Duration                                               `json:"interval"`
		Timeout       time.Duration                                               `json:"timeout"`
		ReadyToTrip   func(counts gobreaker.Counts) bool                          `json:"-"`
		OnStateChange func(name string, from gobreaker.State, to gobreaker.State) `json:"-"`
		IsSuccessful  func(err error) bool                                        `json:"-"`
	}

	Settings struct {
		RetrySettings   RetrySettings   `json:"retrySettings"`
		BreakerSettings BreakerSettings `json:"breakerSettings"`
	}

	impl struct {
		configure.DefaultConfiguraionService
		Settings *Settings
	}

	transport struct {
		transport   http.RoundTripper
		breaker     *gobreaker.CircuitBreaker
		retry       uint
		shouldRetry func(error, *http.Response) bool
		backoff     func(retries uint) time.Duration
	}
)

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {

	res, err := t.breaker.Execute(func() (interface{}, error) {
		var bodyBytes []byte
		if req.Body != nil {
			bodyBytes, _ = io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		res, err := t.transport.RoundTrip(req)

		retries := uint(0)
		for t.shouldRetry(err, res) && retries < t.retry {

			time.Sleep(t.backoff(retries))
			drainBody(res)

			if req.Body != nil {
				req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			res, err = t.transport.RoundTrip(req)
			retries++
			/*go func() {
				log.Println("retrying")
			}()*/

		}
		return res, err
	})

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}

	return res.(*http.Response), err
}

func New(settings *Settings) HttpClient {

	r := &impl{
		Settings: NewHttpSettings(settings),
	}
	return r
}

func (i *impl) GetSettings() *Settings {
	return i.Settings
}

func (i *impl) GetClient() *http.Client {
	return &http.Client{
		Transport: &transport{
			transport:   http.DefaultTransport,
			retry:       i.GetSettings().RetrySettings.Retry,
			shouldRetry: i.GetSettings().RetrySettings.ShouldRetry,
			backoff:     i.GetSettings().RetrySettings.Backoff,
			breaker:     gobreaker.NewCircuitBreaker(*mapBreakerSettings(i.GetSettings().BreakerSettings)),
		},
	}
}

func NewHttpSettings(settings *Settings) *Settings {

	if settings == nil {
		settings = &Settings{}
	}

	if settings.RetrySettings.Retry == 0 {
		settings.RetrySettings.Retry = DEFAULT_RETRY_COUNT
	}

	if settings.RetrySettings.Backoff == nil {
		settings.RetrySettings.Backoff = func(retries uint) time.Duration {
			return DEFAULT_RETRY_BACKOFF
		}
	}

	if settings.RetrySettings.ShouldRetry == nil {
		settings.RetrySettings.ShouldRetry = func(err error, rs *http.Response) bool {
			return err != nil
		}
	}

	if settings.BreakerSettings.Name == "" {
		settings.BreakerSettings.Name = fmt.Sprintf("%s-%d", DEFAULT_BREAKER_PREFIX_NAME, time.Now().UnixNano()/1e6)
	}
	if settings.BreakerSettings.Timeout == 0 {
		settings.BreakerSettings.Timeout = DEFAULT_BREAKER_TIMEOUT
		fmt.Println(settings.BreakerSettings.Timeout)
	}

	if settings.BreakerSettings.Interval == 0 {
		settings.BreakerSettings.Interval = DEFAULT_BREAKER_INTERVAL
	}

	if settings.BreakerSettings.MaxRequests == 0 {
		settings.BreakerSettings.MaxRequests = DEFAULT_BREAKER_MAX_REQUEST
	}

	if settings.BreakerSettings.ReadyToTrip == nil {
		settings.BreakerSettings.ReadyToTrip = func(counts gobreaker.Counts) bool {
			return false
		}
	}

	if settings.BreakerSettings.OnStateChange == nil {
		settings.BreakerSettings.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		}
	}

	if settings.BreakerSettings.IsSuccessful == nil {
		settings.BreakerSettings.IsSuccessful = func(err error) bool {
			return err == nil
		}
	}
	return settings
}

func mapBreakerSettings(bs BreakerSettings) *gobreaker.Settings {
	return &gobreaker.Settings{
		Name:        bs.Name,
		MaxRequests: bs.MaxRequests,
		Interval:    bs.Interval,
		Timeout:     bs.Timeout,
		ReadyToTrip: bs.ReadyToTrip,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("circuit breake (%s) state change: %s to  %s  ", name, from, to)
			if bs.OnStateChange != nil {
				bs.OnStateChange(name, from, to)
			}
		},
	}
}

func drainBody(resp *http.Response) {
	if resp != nil {
		if resp.Body != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
}
