package project

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/romainmenke/bugsnagda/apiaddress"
	"github.com/romainmenke/bugsnagda/apierrors"
	"github.com/romainmenke/bugsnagda/errport"
	"github.com/romainmenke/bugsnagda/event"
	"github.com/romainmenke/bugsnagda/options"
	"github.com/romainmenke/bugsnagda/pagination"
)

type Response struct {
	Next       func(context.Context) (*Response, error)
	Projects   []Project
	TotalCount int
}

const endpoint = apiaddress.Address + "/organizations/%s/projects"

func All(ctx context.Context, client *http.Client, opts options.Projects) (*Response, error) {
	var (
		combinedResponse = &Response{}
		projectsResponse *Response
		err              error
	)

	for true {
		if projectsResponse == nil {
			projectsResponse, err = Paginated(ctx, client, opts)
		} else {
			projectsResponse, err = projectsResponse.Next(ctx)
		}

		if err != nil {
			return nil, err
		}

		if projectsResponse == nil {
			break
		}

		combinedResponse.TotalCount = projectsResponse.TotalCount
		combinedResponse.Projects = append(combinedResponse.Projects, projectsResponse.Projects...)
	}

	return combinedResponse, nil
}

func Paginated(ctx context.Context, client *http.Client, opts options.Projects) (*Response, error) {
	return paginated(ctx, client, fmt.Sprintf(endpoint, opts.OrganisationID), opts)
}

func paginated(ctx context.Context, client *http.Client, u string, opts options.Projects) (*Response, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req = req.Clone(ctx)
	opts.SetQuery(req.URL)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, apierrors.FromResponse(resp)
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
		projects[i].ErrorReports = func(errorsCtx context.Context, errorOpts options.Reports) (*errport.Response, error) {
			errorOpts.ProjectID = projects[i].ID

			return errport.Paginated(errorsCtx, client, errorOpts)
		}

		// ErrorReportsAll
		projects[i].ErrorReportsAll = func(errorsCtx context.Context, errorOpts options.Reports) (*errport.Response, error) {
			errorOpts.ProjectID = projects[i].ID

			return errport.All(errorsCtx, client, errorOpts)
		}

		// Events
		projects[i].Events = func(eventsCtx context.Context, eventOpts options.Events) (*event.Response, error) {
			eventOpts.ProjectID = projects[i].ID

			return event.Paginated(eventsCtx, client, eventOpts)
		}

		// EventsAll
		projects[i].EventsAll = func(eventsCtx context.Context, eventOpts options.Events) (*event.Response, error) {
			eventOpts.ProjectID = projects[i].ID

			return event.All(eventsCtx, client, eventOpts)
		}
	}

	out := &Response{
		Projects:   projects,
		TotalCount: pagination.TotalCount(resp),
		Next: func(nextCtx context.Context) (*Response, error) {
			nextURL, err := pagination.Next(resp)
			if err != nil {
				return nil, err
			}

			if nextURL == "" {
				return nil, nil
			}

			return paginated(nextCtx, client, nextURL, opts)
		},
	}

	return out, nil
}
