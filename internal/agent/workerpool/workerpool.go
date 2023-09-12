package workerpool

import "github.com/jbakhtin/rtagent/internal/types"

type Job struct {
	Value types.Metricer
	Key   string
}

func NewJob(key string, value types.Metricer) Job {
	return Job{
		value,
		key,
	}
}

type WorkerPool struct {
	Jobs chan Job
}

func NewWorkerPool() (WorkerPool, error) {
	return WorkerPool{
		make(chan Job),
	}, nil
}
