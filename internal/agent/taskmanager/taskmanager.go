package taskmanager

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
)

type Func func(ctx context.Context) error

type taskmanager struct {
	funcs []Func
}

func (c *taskmanager) DoIt(ctx context.Context) (err error) {
	defer func() {
		err = errors.Wrap(err, "task manager")
	}()
	var msgs = make([]string, 0, len(c.funcs))
	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
		if len(msgs) > 0 {
			err = fmt.Errorf(
				"tasks finished with error(s): \n\t-%s",
				strings.Join(msgs, "\n\t-"),
			)
		}
	}()
	wg.Add(len(c.funcs))

	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg := errgroup.Group{}
	defer func() {
		tempEer := eg.Wait()
		if tempEer != nil && err != nil {
			err = tempEer
		}

		cancel()
	}()

	for i := range c.funcs {
		i := i
		eg.Go(func() error {
			defer wg.Done()
			tempErr := c.funcs[i](newCtx)
			if tempErr != nil {
				msgs = append(msgs, tempErr.Error())
			}
			return tempErr
		})
	}

	<-ctx.Done()

	return err
}
