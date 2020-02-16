package organisation

import (
	"context"
	"time"

	"github.com/romainmenke/bugsnagda/options"
	"github.com/romainmenke/bugsnagda/project"
)

type Organisation struct {
	ID               string                                                             `json:"id"`
	AutoUpgrade      bool                                                               `json:"auto_upgrade"`
	BillingEmails    []*BillingEmail                                                    `json:"billing_emails"`
	CollaboratorsURL string                                                             `json:"collaborators_url"`
	CreatedAt        time.Time                                                          `json:"created_at"`
	Creator          Creator                                                            `json:"creator"`
	Name             string                                                             `json:"name"`
	Projects         func(context.Context, options.Projects) (*project.Response, error) `json:"-" gorm:"-"`
	ProjectsAll      func(context.Context, options.Projects) (*project.Response, error) `json:"-" gorm:"-"`
	ProjectsURL      string                                                             `json:"projects_url"`
	Slug             string                                                             `json:"slug"`
	UpdatedAt        time.Time                                                          `json:"updated_at"`
	UpgradeURL       string                                                             `json:"upgrade_url"`
}
