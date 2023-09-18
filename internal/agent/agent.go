package agent

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/agentv2/workerPool"
	"github.com/pkg/errors"
	"sync"
	"time"

	"github.com/jbakhtin/rtagent/internal/agent/workerpool"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/types"
)

type Aggregator interface {
	Pool(ctx context.Context)
	GetAll() (map[string]types.Metricer, error)
	Err() chan error
}

type Sender interface {
	Send(key string, value types.Metricer) error
}

type agent struct {
	aggregator                 Aggregator
	sender						Sender
	workerPool                 workerpool.WorkerPool
	serverAddress              string
	sc                         sync.Mutex
	cfg Configer
	errorChan chan error
}

// Run - запустить мониторинг
func (m *agent) Run(ctx context.Context, cfg config.Config) {
	wp, _ := workerPool.New()

	// Task 1
	wp.AddJob(func() error {
		ticker := time.NewTicker(m.cfg.GetPollInterval())
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				go m.aggregator.Pool(ctx)
			case err := <-m.aggregator.Err():
				m.errorChan<- err
			case <-ctx.Done():
				return nil
			}
		}
	})

	// Task 1
	wp.AddJob(func() error {
		ticker := time.NewTicker(m.cfg.GetReportInterval())
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				stats, err := m.aggregator.GetAll()
				err = errors.New("test error")
				if err != nil {
					m.errorChan<- err
				}

				for key, metric := range stats {
					go func(key string, metric types.Metricer) {
						err = m.sender.Send(key, metric)
						if err != nil {
							m.errorChan<- err
						}
					}(key, metric)
				}

			case <-ctx.Done():
				return nil
			}
		}
	})
}

func (m *agent) Err() chan error {
	return m.errorChan
}