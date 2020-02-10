package bugsnagda

import (
	"time"
)

// RateLimitHeaderKey is the number of requests allowed per time window.
const RateLimitHeaderKey = "X-RateLimit-Limit"

// RateLimitRemainingHeaderKey is the number of requests remaining in the current time window.
const RateLimitRemainingHeaderKey = "X-RateLimit-Remaining"

// RateLimitTimeFrame is the duration for rate limits.
const RateLimitTimeFrame = time.Minute
