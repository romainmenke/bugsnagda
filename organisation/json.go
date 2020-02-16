package organisation

import (
	"encoding/json"
	"time"
)

type organisation struct {
	AutoUpgrade      bool      `json:"auto_upgrade"`
	BillingEmails    []string  `json:"billing_emails"`
	CollaboratorsURL string    `json:"collaborators_url"`
	CreatedAt        time.Time `json:"created_at"`
	Creator          Creator   `json:"creator"`
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	ProjectsURL      string    `json:"projects_url"`
	Slug             string    `json:"slug"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpgradeURL       string    `json:"upgrade_url"`
}

func (x *Organisation) UnmarshalJSON(data []byte) error {
	t := organisation{}
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	x.AutoUpgrade = t.AutoUpgrade
	x.CollaboratorsURL = t.CollaboratorsURL
	x.CreatedAt = t.CreatedAt
	x.Creator = t.Creator
	x.ID = t.ID
	x.Name = t.Name
	x.ProjectsURL = t.ProjectsURL
	x.Slug = t.Slug
	x.UpdatedAt = t.UpdatedAt
	x.UpgradeURL = t.UpgradeURL

	for _, y := range t.BillingEmails {
		x.BillingEmails = append(x.BillingEmails, &BillingEmail{
			OrganisationID: t.ID,
			Value:          y,
		})
	}

	return nil
}

func (x Organisation) MarshalJSON() ([]byte, error) {
	t := organisation{}

	t.AutoUpgrade = x.AutoUpgrade
	t.CollaboratorsURL = x.CollaboratorsURL
	t.CreatedAt = x.CreatedAt
	t.Creator = x.Creator
	t.ID = x.ID
	t.Name = x.Name
	t.ProjectsURL = x.ProjectsURL
	t.Slug = x.Slug
	t.UpdatedAt = x.UpdatedAt
	t.UpgradeURL = x.UpgradeURL

	for _, y := range x.BillingEmails {
		t.BillingEmails = append(t.BillingEmails, y.Value)
	}

	return json.Marshal(t)
}
