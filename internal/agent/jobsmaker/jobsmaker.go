package jobsmaker

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/agent/jobqueue"
	"github.com/jbakhtin/rtagent/internal/types"
)

type Slicer interface {
	GetAll() (map[string]types.Metricer)
}

type Jober interface {
	Dequeue() *jobqueue.QNode
	Enqueue(key string, metric types.Metricer)
}

type JobsMaker struct {
	Slicer Slicer
	Jober Jober
}

func (jm *JobsMaker) Do(ctx context.Context) error {
	stats := jm.Slicer.GetAll()

	for key, metric := range stats {
		jm.Jober.Enqueue(key, metric)
	}

	return nil
}
