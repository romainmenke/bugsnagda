package organisation

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/romainmenke/bugsnagda/apiaddress"
	"github.com/romainmenke/bugsnagda/apierrors"
	"github.com/romainmenke/bugsnagda/options"
	"github.com/romainmenke/bugsnagda/pagination"
	"github.com/romainmenke/bugsnagda/project"
)

type Response struct {
	Organisations []Organisation
	TotalCount    int
	Next          func(context.Context) (*Response, error)
}

const organisationsEndpoint = apiaddress.Address + "/user/organizations"

func All(ctx context.Context, client *http.Client, opts options.Organisations) (*Response, error) {
	var (
		combinedResponse      = &Response{}
		organisationsResponse *Response
		err                   error
	)

	for true {
		if organisationsResponse == nil {
			organisationsResponse, err = Paginated(ctx, client, opts)
		} else {
			organisationsResponse, err = organisationsResponse.Next(ctx)
		}

		if err != nil {
			return nil, err
		}

		if organisationsResponse == nil {
			break
		}

		combinedResponse.TotalCount = organisationsResponse.TotalCount
		combinedResponse.Organisations = append(combinedResponse.Organisations, organisationsResponse.Organisations...)
	}

	return combinedResponse, nil
}

func Paginated(ctx context.Context, client *http.Client, opts options.Organisations) (*Response, error) {
	return paginated(ctx, client, organisationsEndpoint, opts)
}

func paginated(ctx context.Context, client *http.Client, u string, opts options.Organisations) (*Response, error) {
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

	organisations := []Organisation{}
	err = decoder.Decode(&organisations)
	if err != nil {
		return nil, err
	}

	for i := range organisations {
		// ID
		organisationID := organisations[i].ID

		// Projects
		organisations[i].Projects = func(projectsCtx context.Context, projectOpts options.Projects) (*project.Response, error) {
			projectOpts.OrganisationID = organisationID

			return project.Paginated(projectsCtx, client, projectOpts)
		}

		// ProjectsAll
		organisations[i].ProjectsAll = func(projectsCtx context.Context, projectOpts options.Projects) (*project.Response, error) {
			projectOpts.OrganisationID = organisationID

			return project.All(projectsCtx, client, projectOpts)
		}
	}

	out := &Response{
		Organisations: organisations,
		TotalCount:    pagination.TotalCount(resp),
		Next: func(nextCtx context.Context) (*Response, error) {
			nextURL, err := pagination.Next(resp)
			if err != nil {
				return nil, err
			}

			if nextURL == "" {
				return nil, nil
			}

			return paginated(nextCtx, client, nextURL, options.Organisations{})
		},
	}

	return out, nil
}
