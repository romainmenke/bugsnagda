package bugsnagda

import (
	"net/http"
	"strconv"
	"time"
)

// rateLimitHeaderKey is the number of requests allowed per time window.
const rateLimitHeaderKey = "X-RateLimit-Limit"

// rateLimitRemainingHeaderKey is the number of requests remaining in the current time window.
const rateLimitRemainingHeaderKey = "X-RateLimit-Remaining"

// rateLimitTimeFrame is the duration for rate limits.
const rateLimitTimeFrame = time.Minute

func rateLimit(resp *http.Response) int {
	if resp == nil {
		return 1
	}

	rateLimitHeader := resp.Header.Get(rateLimitHeaderKey)
	rateLimit, err := strconv.ParseInt(rateLimitHeader, 10, 64)
	if err != nil {
		return 1
	}

	return int(rateLimit)
}

func rateLimitRemaining(resp *http.Response) int {
	if resp == nil {
		return 1
	}

	rateLimitRemainingHeader := resp.Header.Get(rateLimitRemainingHeaderKey)
	rateLimitRemaining, err := strconv.ParseInt(rateLimitRemainingHeader, 10, 64)
	if err != nil {
		return 1
	}

	return int(rateLimitRemaining)
}
