package errport

import (
	"context"
	"time"

	"github.com/romainmenke/bugsnagda/event"
	"github.com/romainmenke/bugsnagda/options"
)

type Report struct {
	Severity                   severity                                                       `json:"severity"`
	AssignedCollaboratorID     string                                                         `json:"assigned_collaborator_id"`
	ID                         string                                                         `json:"id"`
	ProjectID                  string                                                         `json:"project_id"`
	URL                        string                                                         `json:"url"`
	ProjectURL                 string                                                         `json:"project_url"`
	ErrorClass                 string                                                         `json:"error_class"`
	Message                    string                                                         `json:"message"`
	Context                    string                                                         `json:"context"`
	OriginalSeverity           severity                                                       `json:"original_severity"`
	OverriddenSeverity         severity                                                       `json:"overridden_severity"`
	EventsCount                int                                                            `json:"events"`
	EventsURL                  string                                                         `json:"events_url"`
	Events                     func(context.Context, options.Events) (*event.Response, error) `json:"-" gorm:"-"`
	EventsAll                  func(context.Context, options.Events) (*event.Response, error) `json:"-" gorm:"-"`
	UnthrottledOccurrenceCount int                                                            `json:"unthrottled_occurrence_count"`
	Users                      int                                                            `json:"users"`
	FirstSeen                  time.Time                                                      `json:"first_seen"`
	LastSeen                   time.Time                                                      `json:"last_seen"`
	FirstSeenUnfiltered        time.Time                                                      `json:"first_seen_unfiltered"`
	LastSeenUnfiltered         time.Time                                                      `json:"last_seen_unfiltered"`
	ReopenRules                *ReopenRules                                                   `json:"reopen_rules"`
	Status                     status                                                         `json:"status"`
	CommentCount               int                                                            `json:"comment_count"`
	CreatedIssue               *Issue                                                         `json:"created_issue"`
	MissingDSYMs               []*MissingDSYM                                                 `json:"missing_dsyms"`
	ReleaseStages              []*ReleaseStage                                                `json:"release_stages"`
	GroupingReason             string                                                         `json:"grouping_reason"`
	GroupingFields             map[string]interface{}                                         `json:"grouping_fields" gorm:"-"` // TODO
}
