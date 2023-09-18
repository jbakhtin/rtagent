package workerPool

import (
	"sync"
)

type Job func() error

type WorkerPool struct {
	errChan chan error
	wg  sync.WaitGroup
}

func New() (*WorkerPool, error){
	return &WorkerPool{
		make(chan error),
		sync.WaitGroup{},
	}, nil
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.wg.Add(1)
	go func(job Job) {
		err := job()
		if err != nil {
			wp.errChan <- err
		}

		wp.wg.Done()
	}(job)
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

func (wp *WorkerPool) Err() chan error {
	return wp.errChan
}