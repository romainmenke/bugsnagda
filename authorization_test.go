package bugsnagda

import (
	"net/http"
	"testing"
)

func TestSetAuthorizationHeader(t *testing.T) {
	h := http.Header{}

	setAuthorizationHeader(h, "foo")

	if h.Get(authorizationHeaderKey) != "token foo" {
		t.Errorf("expected : \"token foo\", got \"%s\"", h.Get(authorizationHeaderKey))
	}
}
