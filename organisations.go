package bugsnagda

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Creator struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Organisation struct {
	ID               string    `json:"id"`
	Slug             string    `json:"slug"`
	Name             string    `json:"name"`
	BillingEmails    []string  `json:"billing_emails"`
	AutoUpgrade      bool      `json:"auto_upgrade"`
	Creator          Creator   `json:"creator"`
	CollaboratorsURL string    `json:"collaborators_url"`
	ProjectsURL      string    `json:"projects_url"`
	UpgradeURL       string    `json:"upgrade_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type OrganisationsResponse struct {
	Organisations []Organisation
	TotalCount    int
	Next          func(context.Context) (*OrganisationsResponse, error)
}

const organisationsEndpoint = "https://api.bugsnag.com/user/organizations"

type OrganisationsOptions struct {
	// Admin set to true if only Organizations the Current User is an admin of should be returned
	Admin bool
	// PerPage is the pagination limit
	// Example: 10. Default: 30.
	PerPage int
}

func (o OrganisationsOptions) setQuery(u *url.URL) {
	if o.PerPage == 0 {
		o.PerPage = 30
	}

	q := u.Query()

	q.Set("per_page", fmt.Sprint(o.PerPage))

	if o.Admin {
		q.Set("admin", "true")
	}

	u.RawQuery = q.Encode()
}

func (c *Client) Organisations(ctx context.Context, opts OrganisationsOptions) (*OrganisationsResponse, error) {
	return c.organisations(ctx, organisationsEndpoint, opts)
}

func (c *Client) organisations(ctx context.Context, u string, opts OrganisationsOptions) (*OrganisationsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req = req.Clone(ctx)
	opts.setQuery(req.URL)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, errorFromResponse(resp)
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	organisations := []Organisation{}
	err = decoder.Decode(&organisations)
	if err != nil {
		return nil, err
	}

	out := &OrganisationsResponse{
		Organisations: organisations,
		Next: func(nextCtx context.Context) (*OrganisationsResponse, error) {
			nextURL, err := next(resp)
			if err != nil {
				return nil, err
			}

			return c.organisations(nextCtx, nextURL, opts)
		},
	}

	return out, nil
}
