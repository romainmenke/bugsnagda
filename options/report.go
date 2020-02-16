package options

import (
	"fmt"
	"net/url"
	"time"

	"github.com/romainmenke/bugsnagda/apitime"
)

type Reports struct {
	// Direction to sort the results by
	Direction direction

	// ProjectID is the ID of the organization
	ProjectID string

	// PerPage is the pagination limit
	// Example: 10. Default: 30.
	PerPage int

	// Default: current time Only Error Events occurring before this time will be considered.
	Base time.Time

	// Sort the results by which field
	Sort sortBy
}

func (o Reports) SetQuery(u *url.URL) {
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

	if !o.Base.IsZero() {
		q.Set("base", o.Base.Format(apitime.Format))
	}

	u.RawQuery = q.Encode()
}
