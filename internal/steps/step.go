package steps

import "context"

type Step interface {
	Execute(ctx context.Context, data *Context) error
}
