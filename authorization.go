package bugsnagda

import (
	"fmt"
	"net/http"
)

// authorizationHeaderKey is used to authorize the client to access the Bugsnag Data Access API.
const authorizationHeaderKey = "Authorization"

// setAuthorizationHeader sets the Authorization token on a http header.
func setAuthorizationHeader(h http.Header, token string) {
	h.Set(authorizationHeaderKey, fmt.Sprintf("token %s", token))
}
