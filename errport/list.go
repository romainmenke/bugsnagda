package errport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/romainmenke/bugsnagda/apiaddress"
	"github.com/romainmenke/bugsnagda/apierrors"
	"github.com/romainmenke/bugsnagda/event"
	"github.com/romainmenke/bugsnagda/options"
	"github.com/romainmenke/bugsnagda/pagination"
)

type Response struct {
	Next       func(context.Context) (*Response, error)
	Reports    []Report
	TotalCount int
}

const errorsEndpoint = apiaddress.Address + "/projects/%s/errors"

func All(ctx context.Context, client *http.Client, opts options.Reports) (*Response, error) {
	var (
		combinedResponse = &Response{}
		errorsResponse   *Response
		err              error
	)

	for true {
		if errorsResponse == nil {
			errorsResponse, err = Paginated(ctx, client, opts)
		} else {
			errorsResponse, err = errorsResponse.Next(ctx)
		}

		if err != nil {
			return nil, err
		}

		if errorsResponse == nil {
			break
		}

		combinedResponse.TotalCount = errorsResponse.TotalCount
		combinedResponse.Reports = append(combinedResponse.Reports, errorsResponse.Reports...)
	}

	return combinedResponse, nil
}

func Paginated(ctx context.Context, client *http.Client, opts options.Reports) (*Response, error) {
	return paginated(ctx, client, fmt.Sprintf(errorsEndpoint, opts.ProjectID), opts)
}

func paginated(ctx context.Context, client *http.Client, u string, opts options.Reports) (*Response, error) {
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

	reports := []Report{}
	err = decoder.Decode(&reports)
	if err != nil {
		return nil, err
	}

	for i := range reports {
		// Events
		reports[i].Events = func(eventsCtx context.Context, eventOpts options.Events) (*event.Response, error) {
			eventOpts.ProjectID = opts.ProjectID
			eventOpts.ErrorID = reports[i].ID

			return event.Paginated(eventsCtx, client, eventOpts)
		}

		// EventsAll
		reports[i].EventsAll = func(eventsCtx context.Context, eventOpts options.Events) (*event.Response, error) {
			eventOpts.ProjectID = opts.ProjectID
			eventOpts.ErrorID = reports[i].ID

			return event.All(eventsCtx, client, eventOpts)
		}
	}

	out := &Response{
		Reports:    reports,
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
