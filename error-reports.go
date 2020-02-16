package bugsnagda

import (
	"context"

	"github.com/romainmenke/bugsnagda/errport"
	"github.com/romainmenke/bugsnagda/options"
)

func (c *Client) ErrorReportsAll(ctx context.Context, opts options.Reports) (*errport.Response, error) {
	return errport.All(ctx, c.http, opts)
}

func (c *Client) ErrorReports(ctx context.Context, opts options.Reports) (*errport.Response, error) {
	return errport.Paginated(ctx, c.http, opts)
}
