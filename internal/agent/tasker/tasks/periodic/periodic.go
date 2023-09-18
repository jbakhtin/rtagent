package periodic

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/jbakhtin/rtagent/internal/agent/tasker"
	"time"
)

type Task struct {
	Name string
	Duration time.Duration
	F tasker.Func
}

func (t *Task) Do(ctx context.Context) error {
	ticker := time.NewTicker(t.Duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := t.F(ctx)
			if err != nil {
				return errors.Wrap(ctx.Err(), t.Name)
			}

		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), t.Name)
		}
	}
}


