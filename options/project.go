package options

import (
	"fmt"
	"net/url"
)

type Projects struct {
	// Direction to sort the results by
	Direction direction

	// OrganisationID is the ID of the organization
	OrganisationID string

	// PerPage is the pagination limit
	// Example: 10. Default: 30.
	PerPage int

	// Q searches projects with names matching parameter
	Q string

	// Sort the results by which field
	Sort sortBy
}

func (o Projects) SetQuery(u *url.URL) {
	if o.PerPage > 100 {
		o.PerPage = 100
	}

	q := u.Query()

	if o.Direction != "" {
		q.Set("direction", fmt.Sprint(o.Direction))
	}

	if o.PerPage > 0 {
		q.Set("per_page", fmt.Sprint(o.PerPage))
	}

	if o.Sort != "" {
		q.Set("sort", fmt.Sprint(o.Sort))
	}

	if o.Q != "" {
		q.Set("q", fmt.Sprint(o.Q))
	}

	u.RawQuery = q.Encode()
}
