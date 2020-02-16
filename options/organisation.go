package options

import (
	"fmt"
	"net/url"
)

type Organisations struct {
	// Admin set to true if only Organizations the Current User is an admin of should be returned
	Admin bool
	// PerPage is the pagination limit
	// Example: 10. Default: 30.
	PerPage int
}

func (o Organisations) SetQuery(u *url.URL) {
	if o.PerPage == 0 {
		o.PerPage = 30
	}

	if o.PerPage > 100 {
		o.PerPage = 100
	}

	q := u.Query()

	q.Set("per_page", fmt.Sprint(o.PerPage))

	if o.Admin {
		q.Set("admin", "true")
	}

	u.RawQuery = q.Encode()
}
