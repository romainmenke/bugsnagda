package bugsnagda

import "net/http"

// versionHeaderKey is used to request a specific version of the Bugsnag Data Access API.
const versionHeaderKey = "X-Version"

// version is the current version used by the client.
const version = "2"

// setVersionHeader sets the Version Header on a http header.
func setVersionHeader(h http.Header) {
	h.Set(versionHeaderKey, version)
}
