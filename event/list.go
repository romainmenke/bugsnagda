package event

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/romainmenke/bugsnagda/apiaddress"
	"github.com/romainmenke/bugsnagda/apierrors"
	"github.com/romainmenke/bugsnagda/options"
	"github.com/romainmenke/bugsnagda/pagination"
)

type Response struct {
	Next       func(context.Context) (*Response, error)
	Events     []Event
	TotalCount int
}

const eventEndpoint = apiaddress.Address + "/projects/%s/events/%s"
const eventsOnProjectEndpoint = apiaddress.Address + "/projects/%s/events"
const eventsOnErrorEndpoint = apiaddress.Address + "/projects/%s/errors/%s/events"

func All(ctx context.Context, client *http.Client, opts options.Events) (*Response, error) {
	var (
		combinedResponse = &Response{}
		eventsResponse   *Response
		err              error
	)

	for true {
		if eventsResponse == nil {
			eventsResponse, err = Paginated(ctx, client, opts)
		} else {
			eventsResponse, err = eventsResponse.Next(ctx)
		}

		if err != nil {
			return nil, err
		}

		if eventsResponse == nil {
			break
		}

		combinedResponse.TotalCount = eventsResponse.TotalCount
		combinedResponse.Events = append(combinedResponse.Events, eventsResponse.Events...)
	}

	return combinedResponse, nil
}

func Paginated(ctx context.Context, client *http.Client, opts options.Events) (*Response, error) {
	if opts.ErrorID != "" {
		return paginated(ctx, client, fmt.Sprintf(eventsOnErrorEndpoint, opts.ProjectID, opts.ErrorID), opts)
	}

	return paginated(ctx, client, fmt.Sprintf(eventsOnProjectEndpoint, opts.ProjectID), opts)
}

func paginated(ctx context.Context, client *http.Client, u string, opts options.Events) (*Response, error) {
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

	events := []Event{}
	err = decoder.Decode(&events)
	if err != nil {
		return nil, err
	}

	for i := range events {
		events[i].ProjectID = opts.ProjectID
		events[i].FullReport = func(eventCtx context.Context) (*Event, error) {
			return event(eventCtx, client, opts.ProjectID, events[i].ID)
		}
	}

	out := &Response{
		Events:     events,
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

func event(ctx context.Context, client *http.Client, projectID string, eventID string) (*Event, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(eventEndpoint, projectID, eventID), nil)
	if err != nil {
		return nil, err
	}

	req = req.Clone(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, apierrors.FromResponse(resp)
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	event := &Event{}
	err = decoder.Decode(event)
	if err != nil {
		return nil, err
	}

	event.FullReport = func(context.Context) (*Event, error) {
		return event, nil
	}

	return event, nil
}
