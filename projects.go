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
	APIKey                 string                                                                    `json:"api_key"`
	CollaboratorsCount     int                                                                       `json:"collaborators_count"`
	CreatedAt              time.Time                                                                 `json:"created_at"`
	CustomEventFieldsUsed  int                                                                       `json:"custom_event_fields_used"`
	DiscardedAppVersions   []string                                                                  `json:"discarded_app_versions"`
	DiscardedErrors        []string                                                                  `json:"discarded_errors"`
	ErrorReports           func(context.Context, ErrorReportsOptions) (*ErrorReportsResponse, error) `json:"-"`
	ErrorReportsAll        func(context.Context, ErrorReportsOptions) (*ErrorReportsResponse, error) `json:"-"`
	ErrorsURL              string                                                                    `json:"errors_url"`
	Events                 func(context.Context, EventsOptions) (*EventsResponse, error)             `json:"-"`
	EventsAll              func(context.Context, EventsOptions) (*EventsResponse, error)             `json:"-"`
	EventsURL              string                                                                    `json:"events_url"`
	GlobalGrouping         []string                                                                  `json:"global_grouping"`
	HTMLURL                string                                                                    `json:"html_url"`
	ID                     string                                                                    `json:"id"`
	IgnoredBrowserVersions map[string]interface{}                                                    `json:"ignored_browser_versions"`
	IgnoreOldBrowsers      bool                                                                      `json:"ignore_old_browsers"`
	IsFullView             bool                                                                      `json:"is_full_view"`
	Language               string                                                                    `json:"language"`
	LocationGrouping       []string                                                                  `json:"location_grouping"`
	Name                   string                                                                    `json:"name"`
	OpenErrorCount         int                                                                       `json:"open_error_count"`
	ReleaseStages          []string                                                                  `json:"release_stages"`
	ResolveOnDeploy        bool                                                                      `json:"resolve_on_deploy"`
	Slug                   string                                                                    `json:"slug"`
	UpdatedAt              time.Time                                                                 `json:"updated_at"`
	URL                    string                                                                    `json:"url"`
	URLWhitelist           []string                                                                  `json:"url_whitelist"`
}

type ProjectsResponse struct {
	Next       func(context.Context) (*ProjectsResponse, error)
	Projects   []Project
	TotalCount int
}

const projectsEndpoint = address + "/organizations/%s/projects"

type ProjectsOptions struct {
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

func (o ProjectsOptions) setQuery(u *url.URL) {
	if o.PerPage == 0 {
		o.PerPage = 30
	}

	if o.PerPage > 100 {
		o.PerPage = 100
	}

	q := u.Query()

	q.Set("direction", fmt.Sprint(o.Direction))
	q.Set("per_page", fmt.Sprint(o.PerPage))
	q.Set("q", fmt.Sprint(o.Q))
	q.Set("sort", fmt.Sprint(o.Sort))

	u.RawQuery = q.Encode()
}

func (c *Client) ProjectsAll(ctx context.Context, opts ProjectsOptions) (*ProjectsResponse, error) {
	var (
		combinedProjectsResponse = &ProjectsResponse{}
		projectsResponse         *ProjectsResponse
		err                      error
	)

	for true {
		if projectsResponse == nil {
			projectsResponse, err = c.Projects(ctx, opts)
		} else {
			projectsResponse, err = projectsResponse.Next(ctx)
		}

		if err != nil {
			return nil, err
		}

		if projectsResponse == nil {
			break
		}

		combinedProjectsResponse.TotalCount = projectsResponse.TotalCount
		combinedProjectsResponse.Projects = append(combinedProjectsResponse.Projects, projectsResponse.Projects...)
	}

	return combinedProjectsResponse, nil
}

func (c *Client) Projects(ctx context.Context, opts ProjectsOptions) (*ProjectsResponse, error) {
	return c.projects(ctx, fmt.Sprintf(projectsEndpoint, opts.OrganisationID), opts)
}

func (c *Client) projects(ctx context.Context, u string, opts ProjectsOptions) (*ProjectsResponse, error) {
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

	for i := range projects {
		// ErrorReports
		projects[i].ErrorReports = func(errorsCtx context.Context, errorOpts ErrorReportsOptions) (*ErrorReportsResponse, error) {
			errorOpts.ProjectID = projects[i].ID

			return c.ErrorReports(errorsCtx, errorOpts)
		}

		// ErrorReportsAll
		projects[i].ErrorReportsAll = func(errorsCtx context.Context, errorOpts ErrorReportsOptions) (*ErrorReportsResponse, error) {
			errorOpts.ProjectID = projects[i].ID

			return c.ErrorReportsAll(errorsCtx, errorOpts)
		}

		// Events
		projects[i].Events = func(eventsCtx context.Context, eventOpts EventsOptions) (*EventsResponse, error) {
			eventOpts.ProjectID = projects[i].ID

			return c.Events(eventsCtx, eventOpts)
		}

		// EventsAll
		projects[i].EventsAll = func(eventsCtx context.Context, eventOpts EventsOptions) (*EventsResponse, error) {
			eventOpts.ProjectID = projects[i].ID

			return c.EventsAll(eventsCtx, eventOpts)
		}
	}

	out := &ProjectsResponse{
		Projects:   projects,
		TotalCount: totalCount(resp),
		Next: func(nextCtx context.Context) (*ProjectsResponse, error) {
			nextURL, err := next(resp)
			if err != nil {
				return nil, err
			}

			if nextURL == "" {
				return nil, nil
			}

			return c.projects(nextCtx, nextURL, opts)
		},
	}

	return out, nil
}
