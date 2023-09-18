package limited

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/jbakhtin/rtagent/internal/agent/tasker"
	"github.com/jbakhtin/rtagent/pkg/ratelimiter"
	"time"
)

type Task struct {
	Name string
	Limit int
	Duration time.Duration
	F tasker.Func
}

func (t *Task) Do(ctx context.Context) error {
	rl := ratelimiter.New(t.Duration, t.Limit)
	rl.Run()

	for {
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), t.Name)
		case <-rl.Wait():
			err := t.F(ctx)
			if err != nil {
				return errors.Wrap(ctx.Err(), t.Name)
			}
		}
	}
}
