package bugsnagda

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// parseLinkHeader turns Link: <https://api.bugsnag.com/example?offset=590bce131f7314d98eac23ba&per_page=5>; rel="next" into https://api.bugsnag.com/example?offset=590bce131f7314d98eac23ba&per_page=5
// It returns an empty string when:
// - the length limit is exceeded
// - the 'rel' directive does not have a value of 'next'
func parseLinkHeader(h string) string {

	var linkStart int
	var linkEnd int

	// Don't parse anything longer than the limit
	if len(h) > 2048 {
		return ""
	}

RUNELOOP:
	for index, runeValue := range h {
		switch runeValue {

		// This indicates the start of a url.
		case '<':
			linkStart = index + 1

			// This indicates the end of a url.
		case '>':
			linkEnd = index
			break RUNELOOP

			// This indicates the start of a directive.
			// Reset the indices
		case ';':
			linkStart = 0
			linkEnd = 0
		}

	}

	// If either one of the indices is 0, no url was found.
	if linkStart == 0 || linkEnd == 0 {
		return ""
	}

	// If the remainder does not contain the next directive return an empty string.
	if !strings.Contains(h[linkEnd+1:], "rel=next") && !strings.Contains(h[linkEnd+1:], `rel="next"`) {
		return ""
	}

	// Trim spaced from the url.
	url := strings.TrimSpace(h[linkStart:linkEnd])

	return url
}

// next returns a URL to be used when paginating
func next(resp *http.Response) (*url.URL, error) {
	if resp == nil {
		return nil, errors.New("Response was nil, this is unexpected. Please open an issue at https://github.com/romainmenke/bugsnagda")
	}

	linkHeader := resp.Header.Get("Link")
	nextStr := parseLinkHeader(linkHeader)

	return url.Parse(nextStr)
}

// perPage sets the pagination limit on a request
func perPage(req *http.Request, perPage int) {
	q := req.URL.Query()
	q.Set("per_page", fmt.Sprint(perPage))

	req.URL.RawQuery = q.Encode()
}

// totalCountHeaderKey is the total count for the requested resource
const totalCountHeaderKey = "X-Total-Count"

// totalCount returns the total count for the requested resource
func totalCount(resp *http.Response) (int, error) {
	if resp == nil {
		return 0, errors.New("Response was nil, this is unexpected. Please open an issue at https://github.com/romainmenke/bugsnagda")
	}

	totalCountHeader := resp.Header.Get(totalCountHeaderKey)
	totalCount, err := strconv.ParseInt(totalCountHeader, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(totalCount), nil
}
