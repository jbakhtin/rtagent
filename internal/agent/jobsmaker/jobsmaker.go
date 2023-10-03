package jobsmaker

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/agent/jobsqueue"
	"github.com/jbakhtin/rtagent/internal/types"
)

type Slicer interface {
	GetAll() map[string]types.Metricer
}

type Jober interface {
	Dequeue() *jobqueue.Job
	Enqueue(key string, metric types.Metricer)
}

type jobsMaker struct {
	slicer Slicer
	jober  Jober
}

func New(slicer Slicer, jober Jober) *jobsMaker {
	return &jobsMaker{
		slicer,
		jober,
	}
}

func (jm *jobsMaker) Do(ctx context.Context) error {
	stats := jm.slicer.GetAll()

	for key, metric := range stats {
		jm.jober.Enqueue(key, metric)
	}

	return nil
}
