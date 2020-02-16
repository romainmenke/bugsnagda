package event

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type event struct {
	ID           string         `json:"id"`
	IsFullReport bool           `json:"is_full_report"`
	URL          string         `json:"url"`
	ProjectURL   string         `json:"project_url"`
	ProjectID    string         `json:"project_id"`
	ReportID     string         `json:"error_id"`
	ReceivedAt   time.Time      `json:"received_at"`
	Severity     severity       `json:"severity"`
	Exceptions   []*Exception   `json:"exceptions"`
	Unhandled    bool           `json:"unhandled"`
	Context      string         `json:"context"`
	App          *App           `json:"app"`
	Threads      []*Thread      `json:"threads"`
	MetaData     jsonRawMessage `json:"metaData"`
	Breadcrumbs  []*Breadcrumb  `json:"breadcrumbs"`
}

func (x *Event) UnmarshalJSON(data []byte) error {
	t := event{}
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	x.ID = t.ID
	x.IsFullReport = t.IsFullReport
	x.URL = t.URL
	x.ProjectURL = t.ProjectURL
	x.ProjectID = t.ProjectID
	x.ReportID = t.ReportID
	x.ReceivedAt = t.ReceivedAt
	x.Severity = t.Severity
	x.Exceptions = t.Exceptions
	x.Unhandled = t.Unhandled
	x.Context = t.Context
	x.App = t.App
	x.Threads = t.Threads
	x.MetaData = t.MetaData
	x.Breadcrumbs = t.Breadcrumbs

	if x.App != nil {
		x.App.EventID = x.ID
	}

	for i := range x.Exceptions {
		if x.Exceptions[i].ID == "" {
			x.Exceptions[i].ID = uuid.New().String()
		}

		x.Exceptions[i].EventID = x.ID

		for j := range x.Exceptions[i].Stacktrace {
			x.Exceptions[i].Stacktrace[j].Chain = j
			x.Exceptions[i].Stacktrace[j].EventID = x.ID
			x.Exceptions[i].Stacktrace[j].ExceptionID = x.Exceptions[i].ID
		}
	}

	for i := range x.Threads {
		x.Threads[i].EventID = x.ID

		for j := range x.Threads[i].Stacktrace {
			x.Threads[i].Stacktrace[j].Chain = j
			x.Threads[i].Stacktrace[j].EventID = x.ID
			x.Threads[i].Stacktrace[j].ThreadID = x.Threads[i].NumericID
		}
	}

	for i := range x.Breadcrumbs {
		x.Breadcrumbs[i].EventID = x.ID
	}

	return nil
}

func (x Event) MarshalJSON() ([]byte, error) {
	t := event{}

	t.ID = x.ID
	t.IsFullReport = x.IsFullReport
	t.URL = x.URL
	t.ProjectURL = x.ProjectURL
	t.ProjectID = x.ProjectID
	t.ReportID = x.ReportID
	t.ReceivedAt = x.ReceivedAt
	t.Severity = x.Severity
	t.Exceptions = x.Exceptions
	t.Unhandled = x.Unhandled
	t.Context = x.Context
	t.App = x.App
	t.Threads = x.Threads
	t.MetaData = x.MetaData
	t.Breadcrumbs = x.Breadcrumbs

	return json.Marshal(t)
}

type stacktrace struct {
	ID                string
	EventID           string
	ThreadID          int
	ExceptionID       string
	LineNumber        int            `json:"lineNumber"`
	ColumnNumber      int            `json:"columnNumber"`
	File              string         `json:"file"`
	InProject         bool           `json:"inProject"`
	Method            string         `json:"method"`
	Code              map[int]string `json:"code"`
	MachoUUID         string         `json:"machoUUID"`
	SourceControlLink string         `json:"sourceControlLink"`
	SourceControlName string         `json:"sourceControlName"`
}

func (x *Stacktrace) UnmarshalJSON(data []byte) error {
	t := stacktrace{}
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	x.ID = t.ID
	x.EventID = t.EventID
	x.ThreadID = t.ThreadID
	x.ExceptionID = t.ExceptionID
	x.LineNumber = t.LineNumber
	x.ColumnNumber = t.ColumnNumber
	x.File = t.File
	x.InProject = t.InProject
	x.Method = t.Method
	x.MachoUUID = t.MachoUUID
	x.SourceControlLink = t.SourceControlLink
	x.SourceControlName = t.SourceControlName

	if x.ID == "" {
		x.ID = uuid.New().String()
	}

	for k, v := range t.Code {
		x.Code = append(x.Code, &StacktraceCode{
			Key:          k,
			Value:        v,
			StacktraceID: x.ID,
		})
	}

	return nil
}

func (x Stacktrace) MarshalJSON() ([]byte, error) {
	t := stacktrace{}

	t.ID = x.ID
	t.EventID = x.EventID
	t.ThreadID = x.ThreadID
	t.ExceptionID = x.ExceptionID
	t.LineNumber = x.LineNumber
	t.ColumnNumber = x.ColumnNumber
	t.File = x.File
	t.InProject = x.InProject
	t.Method = x.Method
	t.MachoUUID = x.MachoUUID
	t.SourceControlLink = x.SourceControlLink
	t.SourceControlName = x.SourceControlName

	for _, v := range x.Code {
		t.Code[v.Key] = v.Value
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
