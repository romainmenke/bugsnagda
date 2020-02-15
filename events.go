package bugsnagda

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Event struct {
	ID            string                                `json:"id"`
	IsFullReport  bool                                  `json:"is_full_report"`
	URL           string                                `json:"url"`
	ProjectURL    string                                `json:"project_url"`
	ProjectID     string                                `json:"project_id"`
	ErrorReportID string                                `json:"error_id"`
	ReceivedAt    time.Time                             `json:"received_at"`
	Severity      errorReportSeverity                   `json:"severity"`
	Exceptions    []Exception                           `json:"exceptions"`
	Unhandled     bool                                  `json:"unhandled"`
	Context       string                                `json:"context"`
	App           App                                   `json:"app"`
	Threads       []Thread                              `json:"threads"`
	MetaData      json.RawMessage                       `json:"metaData"`
	Breadcrumbs   []Breadcrumb                          `json:"breadcrumbs"`
	FullReport    func(context.Context) (*Event, error) `json:"-"`
}

type EventsResponse struct {
	Next       func(context.Context) (*EventsResponse, error)
	Events     []Event
	TotalCount int
}

const eventEndpoint = address + "/projects/%s/events/%s"
const eventsOnProjectEndpoint = address + "/projects/%s/events"
const eventsOnErrorEndpoint = address + "/projects/%s/errors/%s/events"

type EventsOptions struct {
	// Direction to sort the results by
	Direction direction

	// ProjectID is the ID of the project
	ProjectID string

	// ErrorID is the ID of the error
	ErrorID string

	// PerPage is the pagination limit
	// Example: 10. Default: 30.
	PerPage int

	// Default: current time Only Event Events occurring before this time will be considered.
	Base time.Time

	// Sort the results by which field
	Sort sortBy

	FullReports bool
}

func (o EventsOptions) setQuery(u *url.URL) {
	if o.PerPage == 0 {
		o.PerPage = 30
	}

	if o.PerPage > 100 {
		o.PerPage = 100
	}

	q := u.Query()

	q.Set("direction", fmt.Sprint(o.Direction))
	q.Set("per_page", fmt.Sprint(o.PerPage))
	q.Set("sort", fmt.Sprint(o.Sort))

	if o.Base.IsZero() {
		o.Base = time.Now()
	}

	if o.FullReports {
		q.Set("full_reports", "true")
	}

	q.Set("base", o.Base.Format(timeFormat))

	u.RawQuery = q.Encode()
}

func (c *Client) EventsAll(ctx context.Context, opts EventsOptions) (*EventsResponse, error) {
	var (
		combinedEventsResponse = &EventsResponse{}
		eventsResponse         *EventsResponse
		err                    error
	)

	for true {
		if eventsResponse == nil {
			eventsResponse, err = c.Events(ctx, opts)
		} else {
			eventsResponse, err = eventsResponse.Next(ctx)
		}

		if err != nil {
			return nil, err
		}

		if eventsResponse == nil {
			break
		}

		combinedEventsResponse.TotalCount = eventsResponse.TotalCount
		combinedEventsResponse.Events = append(combinedEventsResponse.Events, eventsResponse.Events...)
	}

	return combinedEventsResponse, nil
}

func (c *Client) Events(ctx context.Context, opts EventsOptions) (*EventsResponse, error) {
	if opts.ErrorID != "" {
		return c.events(ctx, fmt.Sprintf(eventsOnErrorEndpoint, opts.ProjectID, opts.ErrorID), opts)
	}

	return c.events(ctx, fmt.Sprintf(eventsOnProjectEndpoint, opts.ProjectID), opts)
}

func (c *Client) events(ctx context.Context, u string, opts EventsOptions) (*EventsResponse, error) {
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

	events := []Event{}
	err = decoder.Decode(&events)
	if err != nil {
		return nil, err
	}

	for i := range events {
		events[i].ProjectID = opts.ProjectID
		events[i].FullReport = func(eventCtx context.Context) (*Event, error) {
			return c.event(eventCtx, opts.ProjectID, events[i].ID)
		}
	}

	out := &EventsResponse{
		Events:     events,
		TotalCount: totalCount(resp),
		Next: func(nextCtx context.Context) (*EventsResponse, error) {
			nextURL, err := next(resp)
			if err != nil {
				return nil, err
			}

			if nextURL == "" {
				return nil, nil
			}

			return c.events(nextCtx, nextURL, opts)
		},
	}

	return out, nil
}

func (c *Client) event(ctx context.Context, projectID string, eventID string) (*Event, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(eventEndpoint, projectID, eventID), nil)
	if err != nil {
		return nil, err
	}

	req = req.Clone(ctx)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, errorFromResponse(resp)
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
