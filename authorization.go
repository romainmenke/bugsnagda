package bugsnagda

import (
	"fmt"
	"net/http"
)

// AuthorizationHeaderKey is used to authorize the client to access the Bugsnag Data Access API.
const AuthorizationHeaderKey = "Authorization"

// SetAuthorizationHeader sets the Authorization token on a http header.
func SetAuthorizationHeader(h http.Header, token string) {
	h.Set(AuthorizationHeaderKey, fmt.Sprintf("token %s", token))
}
