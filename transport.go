package bugsnagda

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

type transportOptions struct {
	token string
}

type roundTripper func(*http.Request) (*http.Response, error)

func (r roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}

func newTransport(opts transportOptions) (http.RoundTripper, error) {
	var ratelimiter *rate.Limiter
	var waitForNextWindow uint64
	var limit rate.Limit
	var burst int

	{
		pingReq, err := http.NewRequest(http.MethodHead, organisationsEndpoint, nil)
		if err != nil {
			return nil, err
		}

		setAuthorizationHeader(pingReq.Header, opts.token)
		setVersionHeader(pingReq.Header)

		pingResp, err := http.DefaultClient.Do(pingReq)
		if err != nil {
			return nil, err
		}

		limit = rate.Limit(rateLimit(pingResp)) / 60
		burst = rateLimitRemaining(pingResp)

		ratelimiter = rate.NewLimiter(limit, burst)
	}

	return roundTripper(func(req *http.Request) (*http.Response, error) {
		if atomic.LoadUint64(&waitForNextWindow) > 0 {
			r := ratelimiter.ReserveN(time.Now(), burst)
			time.Sleep(r.Delay() + (time.Second * 5))

			atomic.StoreUint64(&waitForNextWindow, 0)
		}

		err := ratelimiter.Wait(req.Context())
		if err != nil {
			return nil, err
		}

		setAuthorizationHeader(req.Header, opts.token)
		setVersionHeader(req.Header)

		resp, err := http.DefaultTransport.RoundTrip(req)

		if rateLimitRemaining(resp) <= 1 {
			atomic.AddUint64(&waitForNextWindow, 1)
		}

		log.Println(req.URL)
		log.Println(rateLimitRemaining(resp))

		return resp, err
	}), nil
}
