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
	Email string `json:"email"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

type Organisation struct {
	AutoUpgrade      bool                                                              `json:"auto_upgrade"`
	BillingEmails    []string                                                          `json:"billing_emails"`
	CollaboratorsURL string                                                            `json:"collaborators_url"`
	CreatedAt        time.Time                                                         `json:"created_at"`
	Creator          Creator                                                           `json:"creator"`
	ID               string                                                            `json:"id"`
	Name             string                                                            `json:"name"`
	Projects         func(context.Context, ProjectsOptions) (*ProjectsResponse, error) `json:"-"`
	ProjectsAll      func(context.Context, ProjectsOptions) (*ProjectsResponse, error) `json:"-"`
	ProjectsURL      string                                                            `json:"projects_url"`
	Slug             string                                                            `json:"slug"`
	UpdatedAt        time.Time                                                         `json:"updated_at"`
	UpgradeURL       string                                                            `json:"upgrade_url"`
}

type OrganisationsResponse struct {
	Organisations []Organisation
	TotalCount    int
	Next          func(context.Context) (*OrganisationsResponse, error)
}

const organisationsEndpoint = address + "/user/organizations"

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

func (c *Client) OrganisationsAll(ctx context.Context, opts OrganisationsOptions) (*OrganisationsResponse, error) {
	var (
		combinedOrganisationsResponse = &OrganisationsResponse{}
		organisationsResponse         *OrganisationsResponse
		err                           error
	)

	for true {
		if organisationsResponse == nil {
			organisationsResponse, err = c.Organisations(ctx, opts)
		} else {
			organisationsResponse, err = organisationsResponse.Next(ctx)
		}

		if err != nil {
			return nil, err
		}

		if organisationsResponse == nil {
			break
		}

		combinedOrganisationsResponse.TotalCount = organisationsResponse.TotalCount
		combinedOrganisationsResponse.Organisations = append(combinedOrganisationsResponse.Organisations, organisationsResponse.Organisations...)
	}

	return combinedOrganisationsResponse, nil
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

	for i := range organisations {
		// Projects
		organisations[i].Projects = func(projectsCtx context.Context, projectOpts ProjectsOptions) (*ProjectsResponse, error) {
			projectOpts.OrganisationID = organisations[i].ID

			return c.Projects(projectsCtx, projectOpts)
		}

		// ProjectsAll
		organisations[i].ProjectsAll = func(projectsCtx context.Context, projectOpts ProjectsOptions) (*ProjectsResponse, error) {
			projectOpts.OrganisationID = organisations[i].ID

			return c.ProjectsAll(projectsCtx, projectOpts)
		}
	}

	out := &OrganisationsResponse{
		Organisations: organisations,
		TotalCount:    totalCount(resp),
		Next: func(nextCtx context.Context) (*OrganisationsResponse, error) {
			nextURL, err := next(resp)
			if err != nil {
				return nil, err
			}

			if nextURL == "" {
				return nil, nil
			}

			return c.organisations(nextCtx, nextURL, opts)
		},
	}

	return out, nil
}
