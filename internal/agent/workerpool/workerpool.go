package workerpool

import "github.com/jbakhtin/rtagent/internal/types"

type Job struct {
	Key string
	Value types.Metricer
}

func NewJob(key string, value types.Metricer) Job {
	return Job{
		key,
		value,
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




