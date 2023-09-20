package closer

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"strings"
)

type Func func(ctx context.Context) error

type closer struct {
	funcs []Func
}

func (c *closer) Close(ctx context.Context) error {
	var (
		msgs     = make([]string, 0, len(c.funcs))
		complete = make(chan struct{}, 1)
	)

	go func() {
		for _, f := range c.funcs {
			if err := f(ctx); err != nil {
				msgs = append(msgs, fmt.Sprintf("[!] %v", err))
			}
		}

		complete <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		msgs = append(msgs, fmt.Sprintf("[!!] %v", errors.Wrap(ctx.Err(), "shutdown cancelled:")))
	case <-complete:
		break
	}

	if len(msgs) > 0 {
		return fmt.Errorf(
			"shutdown finished with error(s): \n%s",
			strings.Join(msgs, "\n"),
		)
	}

	return nil
}
