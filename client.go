package bugsnagda

import "net/http"

type Client struct {
	http *http.Client
}

type Options struct {
	AuthorizationToken string
}

func New(opts Options) *Client {
	httpClient := http.DefaultClient
	httpClient.Transport = newTransport(transportOptions{
		token: opts.AuthorizationToken,
	})

	return &Client{
		http: httpClient,
	}
}
