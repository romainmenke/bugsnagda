package bugsnagda

import (
	"net/http"
	"testing"
)

func TestSetAuthorizationHeader(t *testing.T) {
	h := http.Header{}

	SetAuthorizationHeader(h, "foo")

	if h.Get(AuthorizationHeaderKey) != "token foo" {
		t.Errorf("expected : \"token foo\", got \"%s\"", h.Get(AuthorizationHeaderKey))
	}
}
