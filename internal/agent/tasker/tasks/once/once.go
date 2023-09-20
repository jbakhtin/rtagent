package once

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/jbakhtin/rtagent/internal/agent/tasker"
)

type task struct {
	name string
	f tasker.Func
}
func New(name string, f tasker.Func) *task {
	return &task{
		name,
		f,
	}
}

func (t *task) Do(ctx context.Context) error {
	if err := t.f(ctx); err != nil {
		return errors.Wrap(ctx.Err(), t.name)
	}

	return nil
}
