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

func (t *taskmanager) DoIt(ctx context.Context) (err error) {
	defer func() {
		err = errors.Wrap(err, "task manager")
	}()
	var msgs = make([]string, 0, len(t.funcs))
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
	wg.Add(len(t.funcs))

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

	for i := range t.funcs {
		i := i
		eg.Go(func() error {
			defer wg.Done()
			tempErr := t.funcs[i](newCtx)
			if tempErr != nil {
				msgs = append(msgs, tempErr.Error())
			}
			return tempErr
		})
	}

	<-ctx.Done()

	return err
}
