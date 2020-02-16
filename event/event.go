package event

import (
	"context"
	"time"
)

type Event struct {
	ID           string                                `json:"id"`
	IsFullReport bool                                  `json:"is_full_report"`
	URL          string                                `json:"url"`
	ProjectURL   string                                `json:"project_url"`
	ProjectID    string                                `json:"project_id"`
	ReportID     string                                `json:"error_id"`
	ReceivedAt   time.Time                             `json:"received_at"`
	Severity     severity                              `json:"severity"`
	Exceptions   []*Exception                          `json:"exceptions"`
	Unhandled    bool                                  `json:"unhandled"`
	Context      string                                `json:"context"`
	App          *App                                  `json:"app"`
	Threads      []*Thread                             `json:"threads"`
	MetaData     jsonRawMessage                        `json:"metaData"`
	Breadcrumbs  []*Breadcrumb                         `json:"breadcrumbs"`
	FullReport   func(context.Context) (*Event, error) `json:"-" gorm:"-"`
}
