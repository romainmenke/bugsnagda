package options

import (
	"fmt"
	"net/url"
	"time"

	"github.com/romainmenke/bugsnagda/apitime"
)

type Events struct {
	// Direction to sort the results by
	Direction direction

	// ProjectID is the ID of the project
	ProjectID string

	// ReportID is the ID of the error
	ReportID string

	// PerPage is the pagination limit
	// Example: 10. Default: 30.
	PerPage int

	// Default: current time Only Event Events occurring before this time will be considered.
	Base time.Time

	// Sort the results by which field
	Sort sortBy

	FullReports bool
}

func (o Events) SetQuery(u *url.URL) {
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

	if o.FullReports {
		q.Set("full_reports", "true")
	}

	u.RawQuery = q.Encode()
}
