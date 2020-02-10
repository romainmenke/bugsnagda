package bugsnagda

import "net/http"

type transportOptions struct {
	token string
}

type roundTripper func(*http.Request) (*http.Response, error)

func (r roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}

func newTransport(opts transportOptions) http.RoundTripper {
	return roundTripper(func(req *http.Request) (*http.Response, error) {
		setAuthorizationHeader(req.Header, opts.token)
		setVersionHeader(req.Header)

		return http.DefaultTransport.RoundTrip(req)
	})
}
