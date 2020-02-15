package bugsnagda

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Error struct {
	Severity                   errorSeverity          `json:"severity"`
	AssignedCollaboratorID     string                 `json:"assigned_collaborator_id"`
	ID                         string                 `json:"id"`
	ProjectID                  string                 `json:"project_id"`
	URL                        string                 `json:"url"`
	ProjectURL                 string                 `json:"project_url"`
	ErrorClass                 string                 `json:"error_class"`
	Message                    string                 `json:"message"`
	Context                    string                 `json:"context"`
	OriginalSeverity           errorSeverity          `json:"original_severity"`
	OverriddenSeverity         errorSeverity          `json:"overridden_severity"`
	EventsCount                int                    `json:"events"`
	EventsURL                  string                 `json:"events_url"`
	UnthrottledOccurrenceCount int                    `json:"unthrottled_occurrence_count"`
	Users                      int                    `json:"users"`
	FirstSeen                  time.Time              `json:"first_seen"`
	LastSeen                   time.Time              `json:"last_seen"`
	FirstSeenUnfiltered        time.Time              `json:"first_seen_unfiltered"`
	LastSeenUnfiltered         time.Time              `json:"last_seen_unfiltered"`
	ReopenRules                ErrorReopenRules       `json:"reopen_rules"`
	Status                     errorStatus            `json:"status"`
	CommentCount               int                    `json:"comment_count"`
	CreatedIssue               ErrorIssue             `json:"created_issue"`
	MissingDSYMs               []string               `json:"missing_dsyms"`
	ReleaseStages              []string               `json:"release_stages"`
	GroupingReason             string                 `json:"grouping_reason"`
	GroupingFields             map[string]interface{} `json:"grouping_fields"`
}

type errorSeverity string

const SeverityError errorSeverity = "error"
const SeverityInfo errorSeverity = "info"
const SeverityWarning errorSeverity = "warning"

type errorStatus string

const StatusFixed errorStatus = "fixed"
const StatusForReview errorStatus = "for_review"
const StatusIgnored errorStatus = "ignored"
const StatusInProgress errorStatus = "in progress"
const StatusOpen errorStatus = "open"
const StatusSnoozed errorStatus = "snoozed"

type ErrorIssue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Number int    `json:"number"`
	Type   string `json:"type"`
	URL    string `json:"url"`
}

type ErrorReopenRules struct {
	ReopenIf              string `json:"reopen_if"`
	ReopenAfter           string `json:"reopen_after"`
	Seconds               int    `json:"seconds"`
	Occurences            int    `json:"occurences"`
	Hours                 int    `json:"hours"`
	OccurrenceThreshold   int    `json:"occurrence_threshold"`
	AdditionalOccurrences int    `json:"additional_occurrences"`
}

type ErrorsResponse struct {
	Next       func(context.Context) (*ErrorsResponse, error)
	Errors     []Error
	TotalCount int
}

const errorsEndpoint = address + "/projects/%s/errors"

type ErrorsOptions struct {
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

func (o ErrorsOptions) setQuery(u *url.URL) {
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

func (c *Client) ErrorsAll(ctx context.Context, opts ErrorsOptions) (*ErrorsResponse, error) {
	var (
		combinedErrorsResponse = &ErrorsResponse{}
		errorsResponse         *ErrorsResponse
		err                    error
	)

	for true {
		if errorsResponse == nil {
			errorsResponse, err = c.Errors(ctx, opts)
		} else {
			errorsResponse, err = errorsResponse.Next(ctx)
		}

		if err != nil {
			return nil, err
		}

		if errorsResponse == nil {
			break
		}

		combinedErrorsResponse.TotalCount = errorsResponse.TotalCount
		combinedErrorsResponse.Errors = append(combinedErrorsResponse.Errors, errorsResponse.Errors...)
	}

	return combinedErrorsResponse, nil
}

func (c *Client) Errors(ctx context.Context, opts ErrorsOptions) (*ErrorsResponse, error) {
	return c.errors(ctx, fmt.Sprintf(errorsEndpoint, opts.ProjectID), opts)
}

func (c *Client) errors(ctx context.Context, u string, opts ErrorsOptions) (*ErrorsResponse, error) {
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

	errors := []Error{}
	err = decoder.Decode(&errors)
	if err != nil {
		return nil, err
	}

	out := &ErrorsResponse{
		Errors:     errors,
		TotalCount: totalCount(resp),
		Next: func(nextCtx context.Context) (*ErrorsResponse, error) {
			nextURL, err := next(resp)
			if err != nil {
				return nil, err
			}

			if nextURL == "" {
				return nil, nil
			}

			return c.errors(nextCtx, nextURL, opts)
		},
	}

	return out, nil
}
