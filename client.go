package bugsnagda

import "net/http"

type Client struct {
	http *http.Client
}

type Options struct {
	AuthorizationToken string
}

func New(opts Options) (*Client, error) {
	transport, err := newTransport(transportOptions{
		token: opts.AuthorizationToken,
	})
	if err != nil {
		return nil, err
	}

	httpClient := http.DefaultClient
	httpClient.Transport = transport

	return &Client{
		http: httpClient,
	}, nil
}
