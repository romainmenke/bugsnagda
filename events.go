package bugsnagda

import (
	"context"

	"github.com/romainmenke/bugsnagda/event"
	"github.com/romainmenke/bugsnagda/options"
)

func (c *Client) EventsAll(ctx context.Context, opts options.Events) (*event.Response, error) {
	return event.All(ctx, c.http, opts)
}

func (c *Client) Events(ctx context.Context, opts options.Events) (*event.Response, error) {
	return event.Paginated(ctx, c.http, opts)
}
