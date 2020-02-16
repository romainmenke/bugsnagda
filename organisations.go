package bugsnagda

import (
	"context"

	"github.com/romainmenke/bugsnagda/options"
	"github.com/romainmenke/bugsnagda/organisation"
)

func (c *Client) OrganisationsAll(ctx context.Context, opts options.Organisations) (*organisation.Response, error) {
	return organisation.All(ctx, c.http, opts)
}

func (c *Client) Organisations(ctx context.Context, opts options.Organisations) (*organisation.Response, error) {
	return organisation.Paginated(ctx, c.http, opts)
}
