package periodic

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/jbakhtin/rtagent/internal/agent/tasker"
	"time"
)

type task struct {
	name string
	duration time.Duration
	f tasker.Func
}

func New(name string, duration time.Duration, f tasker.Func) *task {
	return &task{
		name,
		duration,
		f,
	}
}

func (t *task) Do(ctx context.Context) error {
	ticker := time.NewTicker(t.duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := t.f(ctx)
			if err != nil {
				return errors.Wrap(ctx.Err(), t.name)
			}

		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), t.name)
		}
	}
}


