package worker

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/agent/jobsqueue"
	"github.com/jbakhtin/rtagent/internal/types"
)

type IAPI interface {
	Send(key string, value types.Metricer) error
}

type IQueue interface {
	Dequeue() *jobqueue.Job
	IsEmpty() bool
}

type worker struct {
	queue IQueue
	api   IAPI
}

func New(queue IQueue, api IAPI) *worker {
	return &worker{
		queue,
		api,
	}
}

func (jm *worker) Do(ctx context.Context) error {
	if jm.queue.IsEmpty() {
		return nil
	}

	node := jm.queue.Dequeue()
	err := jm.api.Send(node.Key(), node.Value())
	if err != nil {
		return nil //ToDo: need to forward to main to circuit breaker
	}

	return nil
}
