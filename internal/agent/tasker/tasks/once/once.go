package once

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/jbakhtin/rtagent/internal/agent/tasker"
)

type Task struct {
	Name string
	F tasker.Func
}

func (t *Task) Do(ctx context.Context) error {
	if err := t.F(ctx); err != nil {
		return errors.Wrap(ctx.Err(), t.Name)
	}

	return nil
}
