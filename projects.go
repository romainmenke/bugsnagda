package bugsnagda

import (
	"context"

	"github.com/romainmenke/bugsnagda/options"
	"github.com/romainmenke/bugsnagda/project"
)

func (c *Client) ProjectsAll(ctx context.Context, opts options.Projects) (*project.Response, error) {
	return project.All(ctx, c.http, opts)
}

func (c *Client) Projects(ctx context.Context, opts options.Projects) (*project.Response, error) {
	return project.Paginated(ctx, c.http, opts)
}
