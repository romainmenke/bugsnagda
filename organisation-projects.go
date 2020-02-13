package bugsnagda

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Project struct {
	ID                     string      `json:"id"`
	Slug                   string      `json:"slug"`
	Name                   string      `json:"name"`
	CreatedAt              time.Time   `json:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at"`
	GlobalGrouping         []string    `json:"global_grouping"`
	LocationGrouping       []string    `json:"location_grouping"`
	DiscardedAppVersions   []string    `json:"discarded_app_versions"`
	DiscardedErrors        []string    `json:"discarded_errors"`
	URLWhitelist           []string    `json:"url_whitelist"`
	IgnoreOldBrowsers      bool        `json:"ignore_old_browsers"`
	IgnoredBrowserVersions interface{} `json:"ignored_browser_versions"` // TODO : type
	ResolveOnDeploy        bool        `json:"resolve_on_deploy"`
	APIKey                 string      `json:"api_key"`
	IsFullView             bool        `json:"is_full_view"`
	ReleaseStages          []string    `json:"release_stages"`
	Language               string      `json:"language"`
	URL                    string      `json:"url"`
	HTMLURL                string      `json:"html_url"`
	ErrorsURL              string      `json:"errors_url"`
	EventsURL              string      `json:"events_url"`
	OpenErrorCount         int         `json:"open_error_count"`
	CollaboratorsCount     int         `json:"collaborators_count"`
	CustomEventFieldsUsed  int         `json:"custom_event_fields_used"`
}

type OrganisationProjectsResponse struct {
	Projects   []Project
	TotalCount int
	Next       func(context.Context) (*OrganisationProjectsResponse, error)
}

const organisationProjectsEndpoint = "https://api.bugsnag.com/organizations/%s/projects"

type OrganisationProjectsOptions struct {
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

func (o OrganisationProjectsOptions) setQuery(u *url.URL) {
	if o.PerPage == 0 {
		o.PerPage = 30
	}

	q := u.Query()

	q.Set("direction", fmt.Sprint(o.Direction))
	q.Set("organization_id", fmt.Sprint(o.OrganisationID))
	q.Set("per_page", fmt.Sprint(o.PerPage))
	q.Set("q", fmt.Sprint(o.Q))
	q.Set("sort", fmt.Sprint(o.Sort))

	u.RawQuery = q.Encode()
}

func (c *Client) OrganisationProjects(ctx context.Context, opts OrganisationProjectsOptions) (*OrganisationProjectsResponse, error) {
	return c.organisationProjects(ctx, fmt.Sprintf(organisationProjectsEndpoint, opts.OrganisationID), opts)
}

func (c *Client) organisationProjects(ctx context.Context, u string, opts OrganisationProjectsOptions) (*OrganisationProjectsResponse, error) {
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

	projects := []Project{}
	err = decoder.Decode(&projects)
	if err != nil {
		return nil, err
	}

	out := &OrganisationProjectsResponse{
		Projects:   projects,
		TotalCount: totalCount(resp),
		Next: func(nextCtx context.Context) (*OrganisationProjectsResponse, error) {
			nextURL, err := next(resp)
			if err != nil {
				return nil, err
			}

			if nextURL == "" {
				return nil, nil
			}

			return c.organisationProjects(nextCtx, nextURL, opts)
		},
	}

	return out, nil
}
