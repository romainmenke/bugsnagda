package bugsnagda

import (
	"net/http"

	"golang.org/x/time/rate"
)

type transportOptions struct {
	token string
}

type roundTripper func(*http.Request) (*http.Response, error)

func (r roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}

func newTransport(opts transportOptions) http.RoundTripper {
	ratelimiter := rate.NewLimiter(rate.Limit(10)/60, 10)

	return roundTripper(func(req *http.Request) (*http.Response, error) {
		err := ratelimiter.Wait(req.Context())
		if err != nil {
			return nil, err
		}

		setAuthorizationHeader(req.Header, opts.token)
		setVersionHeader(req.Header)

		resp, err := http.DefaultTransport.RoundTrip(req)

		// TODO : adjust rate limiter here

		return resp, err
	})
}
