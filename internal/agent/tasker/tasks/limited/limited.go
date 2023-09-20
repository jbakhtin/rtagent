package limited

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/jbakhtin/rtagent/internal/agent/tasker"
	"github.com/jbakhtin/rtagent/pkg/ratelimiter"
	"time"
)

type task struct {
	name string
	limit int
	duration time.Duration
	f tasker.Func
}

func New(name string, limit int, duration time.Duration, f tasker.Func) *task {
	return &task{
		name,
		limit,
		duration,
		f,
	}
}

func (t *task) Do(ctx context.Context) error {
	rl := ratelimiter.New(t.duration, t.limit)
	rl.Run()

	for {
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), t.name)
		case <-rl.Wait():
			err := t.f(ctx)
			if err != nil {
				return errors.Wrap(ctx.Err(), t.name)
			}
		}
	}
}
