package messagesender

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

type messageSender struct {
	queue IQueue
	api   IAPI
}

func New(queue IQueue, api IAPI) *messageSender {
	return &messageSender{
		queue,
		api,
	}
}

func (ms *messageSender) Do(ctx context.Context) error {
	if ms.queue.IsEmpty() {
		return nil
	}

	node := ms.queue.Dequeue()
	err := ms.api.Send(node.Key(), node.Value())
	if err != nil {
		return nil //ToDo: need to forward to main to circuit breaker
	}

	return nil
}
