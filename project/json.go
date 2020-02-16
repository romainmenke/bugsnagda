package project

import (
	"encoding/json"
	"fmt"
	"time"
)

type project struct {
	APIKey                 string                 `json:"api_key"`
	CollaboratorsCount     int                    `json:"collaborators_count"`
	CreatedAt              time.Time              `json:"created_at"`
	CustomEventFieldsUsed  int                    `json:"custom_event_fields_used"`
	DiscardedAppVersions   []string               `json:"discarded_app_versions"`
	DiscardedErrors        []string               `json:"discarded_errors"`
	ErrorsURL              string                 `json:"errors_url"`
	EventsURL              string                 `json:"events_url"`
	GlobalGrouping         []string               `json:"global_grouping"`
	HTMLURL                string                 `json:"html_url"`
	ID                     string                 `json:"id"`
	IgnoredBrowserVersions map[string]interface{} `json:"ignored_browser_versions"`
	IgnoreOldBrowsers      bool                   `json:"ignore_old_browsers"`
	IsFullView             bool                   `json:"is_full_view"`
	Language               string                 `json:"language"`
	LocationGrouping       []string               `json:"location_grouping"`
	Name                   string                 `json:"name"`
	OpenErrorCount         int                    `json:"open_error_count"`
	ReleaseStages          []string               `json:"release_stages"`
	ResolveOnDeploy        bool                   `json:"resolve_on_deploy"`
	Slug                   string                 `json:"slug"`
	UpdatedAt              time.Time              `json:"updated_at"`
	URL                    string                 `json:"url"`
	URLWhitelist           []string               `json:"url_whitelist"`
}

func (x *Project) UnmarshalJSON(data []byte) error {
	t := project{}
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	x.APIKey = t.APIKey
	x.CollaboratorsCount = t.CollaboratorsCount
	x.CreatedAt = t.CreatedAt
	x.CustomEventFieldsUsed = t.CustomEventFieldsUsed
	x.ErrorsURL = t.ErrorsURL
	x.EventsURL = t.EventsURL
	x.HTMLURL = t.HTMLURL
	x.ID = t.ID
	x.IgnoreOldBrowsers = t.IgnoreOldBrowsers
	x.IsFullView = t.IsFullView
	x.Language = t.Language
	x.Name = t.Name
	x.OpenErrorCount = t.OpenErrorCount
	x.ResolveOnDeploy = t.ResolveOnDeploy
	x.Slug = t.Slug
	x.UpdatedAt = t.UpdatedAt
	x.URL = t.URL

	for _, y := range t.DiscardedAppVersions {
		x.DiscardedAppVersions = append(x.DiscardedAppVersions, &DiscardedAppVersion{
			ProjectID: t.ID,
			Value:     y,
		})
	}

	for _, y := range t.DiscardedErrors {
		x.DiscardedErrors = append(x.DiscardedErrors, &DiscardedError{
			ProjectID: t.ID,
			Value:     y,
		})
	}

	for _, y := range t.GlobalGrouping {
		x.GlobalGrouping = append(x.GlobalGrouping, &GlobalGroup{
			ProjectID: t.ID,
			Value:     y,
		})
	}

	for _, y := range t.LocationGrouping {
		x.LocationGrouping = append(x.LocationGrouping, &LocationGroup{
			ProjectID: t.ID,
			Value:     y,
		})
	}

	for _, y := range t.LocationGrouping {
		x.LocationGrouping = append(x.LocationGrouping, &LocationGroup{
			ProjectID: t.ID,
			Value:     y,
		})
	}

	for _, y := range t.ReleaseStages {
		x.ReleaseStages = append(x.ReleaseStages, &ProjectReleaseStage{
			ProjectID: t.ID,
			Value:     y,
		})
	}

	for k, v := range t.IgnoredBrowserVersions {
		x.IgnoredBrowserVersions = append(x.IgnoredBrowserVersions, &IgnoredBrowserVersion{
			ProjectID: t.ID,
			Key:       k,
			Value:     fmt.Sprint(v),
		})
	}

	return nil
}

func (x Project) MarshalJSON() ([]byte, error) {
	t := project{}

	t.APIKey = x.APIKey
	t.CollaboratorsCount = x.CollaboratorsCount
	t.CreatedAt = x.CreatedAt
	t.CustomEventFieldsUsed = x.CustomEventFieldsUsed
	t.ErrorsURL = x.ErrorsURL
	t.EventsURL = x.EventsURL
	t.HTMLURL = x.HTMLURL
	t.ID = x.ID
	t.IgnoreOldBrowsers = x.IgnoreOldBrowsers
	t.IsFullView = x.IsFullView
	t.Language = x.Language
	t.Name = x.Name
	t.OpenErrorCount = x.OpenErrorCount
	t.ResolveOnDeploy = x.ResolveOnDeploy
	t.Slug = x.Slug
	t.UpdatedAt = x.UpdatedAt
	t.URL = x.URL

	for _, y := range x.DiscardedAppVersions {
		t.DiscardedAppVersions = append(t.DiscardedAppVersions, y.Value)
	}

	for _, y := range x.DiscardedErrors {
		t.DiscardedErrors = append(t.DiscardedErrors, y.Value)
	}

	for _, y := range x.GlobalGrouping {
		t.GlobalGrouping = append(t.GlobalGrouping, y.Value)
	}

	for _, y := range x.LocationGrouping {
		t.LocationGrouping = append(t.LocationGrouping, y.Value)
	}

	for _, y := range x.ReleaseStages {
		t.ReleaseStages = append(t.ReleaseStages, y.Value)
	}

	for _, y := range x.URLWhitelist {
		t.URLWhitelist = append(t.URLWhitelist, y.Value)
	}

	for _, y := range x.IgnoredBrowserVersions {
		t.IgnoredBrowserVersions[y.Key] = y.Value
	}

	return json.Marshal(t)
}
