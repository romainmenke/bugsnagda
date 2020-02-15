package bugsnagda

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ErrorReport struct {
	Severity                   errorReportSeverity                                           `json:"severity"`
	AssignedCollaboratorID     string                                                        `json:"assigned_collaborator_id"`
	ID                         string                                                        `json:"id"`
	ProjectID                  string                                                        `json:"project_id"`
	URL                        string                                                        `json:"url"`
	ProjectURL                 string                                                        `json:"project_url"`
	ErrorClass                 string                                                        `json:"error_class"`
	Message                    string                                                        `json:"message"`
	Context                    string                                                        `json:"context"`
	OriginalSeverity           errorReportSeverity                                           `json:"original_severity"`
	OverriddenSeverity         errorReportSeverity                                           `json:"overridden_severity"`
	EventsCount                int                                                           `json:"events"`
	EventsURL                  string                                                        `json:"events_url"`
	Events                     func(context.Context, EventsOptions) (*EventsResponse, error) `json:"-"`
	EventsAll                  func(context.Context, EventsOptions) (*EventsResponse, error) `json:"-"`
	UnthrottledOccurrenceCount int                                                           `json:"unthrottled_occurrence_count"`
	Users                      int                                                           `json:"users"`
	FirstSeen                  time.Time                                                     `json:"first_seen"`
	LastSeen                   time.Time                                                     `json:"last_seen"`
	FirstSeenUnfiltered        time.Time                                                     `json:"first_seen_unfiltered"`
	LastSeenUnfiltered         time.Time                                                     `json:"last_seen_unfiltered"`
	ReopenRules                ErrorReportReopenRules                                        `json:"reopen_rules"`
	Status                     errorReportStatus                                             `json:"status"`
	CommentCount               int                                                           `json:"comment_count"`
	CreatedIssue               ErrorReportIssue                                              `json:"created_issue"`
	MissingDSYMs               []string                                                      `json:"missing_dsyms"`
	ReleaseStages              []string                                                      `json:"release_stages"`
	GroupingReason             string                                                        `json:"grouping_reason"`
	GroupingFields             map[string]interface{}                                        `json:"grouping_fields"`
}

type errorReportSeverity string

const SeverityError errorReportSeverity = "error"
const SeverityInfo errorReportSeverity = "info"
const SeverityWarning errorReportSeverity = "warning"

type errorReportStatus string

const StatusFixed errorReportStatus = "fixed"
const StatusForReview errorReportStatus = "for_review"
const StatusIgnored errorReportStatus = "ignored"
const StatusInProgress errorReportStatus = "in progress"
const StatusOpen errorReportStatus = "open"
const StatusSnoozed errorReportStatus = "snoozed"

type ErrorReportIssue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Number int    `json:"number"`
	Type   string `json:"type"`
	URL    string `json:"url"`
}

type ErrorReportReopenRules struct {
	ReopenIf              string `json:"reopen_if"`
	ReopenAfter           string `json:"reopen_after"`
	Seconds               int    `json:"seconds"`
	Occurences            int    `json:"occurences"`
	Hours                 int    `json:"hours"`
	OccurrenceThreshold   int    `json:"occurrence_threshold"`
	AdditionalOccurrences int    `json:"additional_occurrences"`
}

type ErrorReportsResponse struct {
	Next       func(context.Context) (*ErrorReportsResponse, error)
	Errors     []ErrorReport
	TotalCount int
}

const errorsEndpoint = address + "/projects/%s/errors"

type ErrorReportsOptions struct {
	// Direction to sort the results by
	Direction direction

	// ProjectID is the ID of the organization
	ProjectID string

	// PerPage is the pagination limit
	// Example: 10. Default: 30.
	PerPage int

	// Default: current time Only Error Events occurring before this time will be considered.
	Base time.Time

	// Sort the results by which field
	Sort sortBy
}

func (o ErrorReportsOptions) setQuery(u *url.URL) {
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

	q.Set("base", o.Base.Format(timeFormat))

	u.RawQuery = q.Encode()
}

func (c *Client) ErrorReportsAll(ctx context.Context, opts ErrorReportsOptions) (*ErrorReportsResponse, error) {
	var (
		combinedErrorReportsResponse = &ErrorReportsResponse{}
		errorsResponse               *ErrorReportsResponse
		err                          error
	)

	for true {
		if errorsResponse == nil {
			errorsResponse, err = c.ErrorReports(ctx, opts)
		} else {
			errorsResponse, err = errorsResponse.Next(ctx)
		}

		if err != nil {
			return nil, err
		}

		if errorsResponse == nil {
			break
		}

		combinedErrorReportsResponse.TotalCount = errorsResponse.TotalCount
		combinedErrorReportsResponse.Errors = append(combinedErrorReportsResponse.Errors, errorsResponse.Errors...)
	}

	return combinedErrorReportsResponse, nil
}

func (c *Client) ErrorReports(ctx context.Context, opts ErrorReportsOptions) (*ErrorReportsResponse, error) {
	return c.errorReports(ctx, fmt.Sprintf(errorsEndpoint, opts.ProjectID), opts)
}

func (c *Client) errorReports(ctx context.Context, u string, opts ErrorReportsOptions) (*ErrorReportsResponse, error) {
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

	errors := []ErrorReport{}
	err = decoder.Decode(&errors)
	if err != nil {
		return nil, err
	}

	for i := range errors {
		// Events
		errors[i].Events = func(eventsCtx context.Context, eventOpts EventsOptions) (*EventsResponse, error) {
			eventOpts.ProjectID = opts.ProjectID
			eventOpts.ErrorID = errors[i].ID

			return c.Events(eventsCtx, eventOpts)
		}

		// EventsAll
		errors[i].EventsAll = func(eventsCtx context.Context, eventOpts EventsOptions) (*EventsResponse, error) {
			eventOpts.ProjectID = opts.ProjectID
			eventOpts.ErrorID = errors[i].ID

			return c.EventsAll(eventsCtx, eventOpts)
		}
	}

	out := &ErrorReportsResponse{
		Errors:     errors,
		TotalCount: totalCount(resp),
		Next: func(nextCtx context.Context) (*ErrorReportsResponse, error) {
			nextURL, err := next(resp)
			if err != nil {
				return nil, err
			}

			if nextURL == "" {
				return nil, nil
			}

			return c.errorReports(nextCtx, nextURL, opts)
		},
	}

	return out, nil
}
