package bugsnagda

import (
	"time"
)

// rateLimitHeaderKey is the number of requests allowed per time window.
const rateLimitHeaderKey = "X-RateLimit-Limit"

// rateLimitRemainingHeaderKey is the number of requests remaining in the current time window.
const rateLimitRemainingHeaderKey = "X-RateLimit-Remaining"

// rateLimitTimeFrame is the duration for rate limits.
const rateLimitTimeFrame = time.Minute
