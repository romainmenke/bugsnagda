package errport

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type report struct {
	Severity                   severity       `json:"severity"`
	AssignedCollaboratorID     string         `json:"assigned_collaborator_id"`
	ID                         string         `json:"id"`
	ProjectID                  string         `json:"project_id"`
	URL                        string         `json:"url"`
	ProjectURL                 string         `json:"project_url"`
	ErrorClass                 string         `json:"error_class"`
	Message                    string         `json:"message"`
	Context                    string         `json:"context"`
	OriginalSeverity           severity       `json:"original_severity"`
	OverriddenSeverity         severity       `json:"overridden_severity"`
	EventsCount                int            `json:"events"`
	EventsURL                  string         `json:"events_url"`
	UnthrottledOccurrenceCount int            `json:"unthrottled_occurrence_count"`
	Users                      int            `json:"users"`
	FirstSeen                  time.Time      `json:"first_seen"`
	LastSeen                   time.Time      `json:"last_seen"`
	FirstSeenUnfiltered        time.Time      `json:"first_seen_unfiltered"`
	LastSeenUnfiltered         time.Time      `json:"last_seen_unfiltered"`
	ReopenRules                *ReopenRules   `json:"reopen_rules"`
	Status                     status         `json:"status"`
	CommentCount               int            `json:"comment_count"`
	CreatedIssue               *Issue         `json:"created_issue"`
	MissingDSYMs               []string       `json:"missing_dsyms"`
	ReleaseStages              []string       `json:"release_stages"`
	GroupingReason             string         `json:"grouping_reason"`
	GroupingFields             jsonRawMessage `json:"grouping_fields"`
}

func (x *Report) UnmarshalJSON(data []byte) error {
	t := report{}
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	x.Severity = t.Severity
	x.AssignedCollaboratorID = t.AssignedCollaboratorID
	x.ID = t.ID
	x.ProjectID = t.ProjectID
	x.URL = t.URL
	x.ProjectURL = t.ProjectURL
	x.ErrorClass = t.ErrorClass
	x.Message = t.Message
	x.Context = t.Context
	x.OriginalSeverity = t.OriginalSeverity
	x.OverriddenSeverity = t.OverriddenSeverity
	x.EventsCount = t.EventsCount
	x.EventsURL = t.EventsURL
	x.UnthrottledOccurrenceCount = t.UnthrottledOccurrenceCount
	x.Users = t.Users
	x.FirstSeen = t.FirstSeen
	x.LastSeen = t.LastSeen
	x.FirstSeenUnfiltered = t.FirstSeenUnfiltered
	x.LastSeenUnfiltered = t.LastSeenUnfiltered
	x.ReopenRules = t.ReopenRules
	x.Status = t.Status
	x.CommentCount = t.CommentCount
	x.CreatedIssue = t.CreatedIssue
	x.GroupingReason = t.GroupingReason
	x.GroupingFields = t.GroupingFields

	if x.ReopenRules != nil {
		x.ReopenRules.ReportID = x.ID
	}

	if x.CreatedIssue != nil {
		x.CreatedIssue.ReportID = x.ID
	}

	for _, y := range t.MissingDSYMs {
		x.MissingDSYMs = append(x.MissingDSYMs, &MissingDSYM{
			ReportID: t.ID,
			Value:    y,
		})
	}

	for _, y := range t.ReleaseStages {
		x.ReleaseStages = append(x.ReleaseStages, &ReportReleaseStage{
			ReportID: t.ID,
			Value:    y,
		})
	}

	return nil
}

func (x Report) MarshalJSON() ([]byte, error) {
	t := report{}

	t.Severity = x.Severity
	t.AssignedCollaboratorID = x.AssignedCollaboratorID
	t.ID = x.ID
	t.ProjectID = x.ProjectID
	t.URL = x.URL
	t.ProjectURL = x.ProjectURL
	t.ErrorClass = x.ErrorClass
	t.Message = x.Message
	t.Context = x.Context
	t.OriginalSeverity = x.OriginalSeverity
	t.OverriddenSeverity = x.OverriddenSeverity
	t.EventsCount = x.EventsCount
	t.EventsURL = x.EventsURL
	t.UnthrottledOccurrenceCount = x.UnthrottledOccurrenceCount
	t.Users = x.Users
	t.FirstSeen = x.FirstSeen
	t.LastSeen = x.LastSeen
	t.FirstSeenUnfiltered = x.FirstSeenUnfiltered
	t.LastSeenUnfiltered = x.LastSeenUnfiltered
	t.ReopenRules = x.ReopenRules
	t.Status = x.Status
	t.CommentCount = x.CommentCount
	t.CreatedIssue = x.CreatedIssue
	t.GroupingReason = x.GroupingReason
	t.GroupingFields = x.GroupingFields

	for _, y := range x.MissingDSYMs {
		t.MissingDSYMs = append(t.MissingDSYMs, y.Value)
	}

	for _, y := range x.ReleaseStages {
		t.ReleaseStages = append(t.ReleaseStages, y.Value)
	}

	return json.Marshal(t)
}

type jsonRawMessage json.RawMessage

func (x *jsonRawMessage) UnmarshalJSON(data []byte) error {
	*x = jsonRawMessage(data)

	return nil
}

func (x jsonRawMessage) MarshalJSON() ([]byte, error) {
	return []byte(x), nil
}

func (x jsonRawMessage) Value() (driver.Value, error) {
	return []byte(x), nil
}

func (x *jsonRawMessage) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		*x = []byte(v)
		return nil
	case []byte:
		*x = v
		return nil
	}
	// otherwise, return an error
	return errors.New("failed to scan []byte")
}
