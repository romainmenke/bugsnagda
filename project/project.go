package project

import (
	"context"
	"time"

	"github.com/romainmenke/bugsnagda/errport"
	"github.com/romainmenke/bugsnagda/event"
	"github.com/romainmenke/bugsnagda/options"
)

type Project struct {
	APIKey                 string                                                            `json:"api_key"`
	CollaboratorsCount     int                                                               `json:"collaborators_count"`
	CreatedAt              time.Time                                                         `json:"created_at"`
	CustomEventFieldsUsed  int                                                               `json:"custom_event_fields_used"`
	DiscardedAppVersions   []*DiscardedAppVersion                                            `json:"discarded_app_versions"`
	DiscardedErrors        []*DiscardedError                                                 `json:"discarded_errors"`
	ErrorReports           func(context.Context, options.Reports) (*errport.Response, error) `json:"-" gorm:"-"`
	ErrorReportsAll        func(context.Context, options.Reports) (*errport.Response, error) `json:"-" gorm:"-"`
	ErrorsURL              string                                                            `json:"errors_url"`
	Events                 func(context.Context, options.Events) (*event.Response, error)    `json:"-" gorm:"-"`
	EventsAll              func(context.Context, options.Events) (*event.Response, error)    `json:"-" gorm:"-"`
	EventsURL              string                                                            `json:"events_url"`
	GlobalGrouping         []*GlobalGroup                                                    `json:"global_grouping"`
	HTMLURL                string                                                            `json:"html_url"`
	ID                     string                                                            `json:"id"`
	IgnoredBrowserVersions []*IgnoredBrowserVersion                                          `json:"ignored_browser_versions"`
	IgnoreOldBrowsers      bool                                                              `json:"ignore_old_browsers"`
	IsFullView             bool                                                              `json:"is_full_view"`
	Language               string                                                            `json:"language"`
	LocationGrouping       []*LocationGroup                                                  `json:"location_grouping"`
	Name                   string                                                            `json:"name"`
	OpenErrorCount         int                                                               `json:"open_error_count"`
	ReleaseStages          []*ProjectReleaseStage                                            `json:"release_stages"`
	ResolveOnDeploy        bool                                                              `json:"resolve_on_deploy"`
	Slug                   string                                                            `json:"slug"`
	UpdatedAt              time.Time                                                         `json:"updated_at"`
	URL                    string                                                            `json:"url"`
	URLWhitelist           []*WhitelistedURL                                                 `json:"url_whitelist"`
}
