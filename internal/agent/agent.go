package agent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jbakhtin/rtagent/pkg/ratelimiter"

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
	sender					Sender
	workerPool                 workerpool.WorkerPool
	serverAddress              string
	sc                         sync.Mutex
	cfg Configer
	errorChan chan error
	softShuttingDown bool
	isShuttingDown bool
}

// Run - запустить мониторинг
func (m *agent) Run(ctx context.Context, cfg config.Config) {
	go m.polling(ctx)
	go m.reporting(ctx)
	go m.run(ctx, cfg)
}

func (m *agent) Err() chan error {
	return m.errorChan
}

// pooling - инициирует забор данных с заданным интервалом monitor.pollInterval
func (m *agent) polling(ctx context.Context) {
	ticker := time.NewTicker(m.cfg.GetPollInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go m.aggregator.Pool(ctx)
		case err := <-m.aggregator.Err():
			m.errorChan<- err
		case <-ctx.Done():
			return
		}
	}
}

// reporting - инициирует отправку данных с заданным интервалом monitor.reportInterval
func (m *agent) reporting(ctx context.Context) {
	ticker := time.NewTicker(m.cfg.GetReportInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			timer := time.NewTimer(time.Second * 3)
			<-timer.C
			err := m.report()
			if err != nil {
				m.errorChan<- err
			}

		case <-ctx.Done():
			return
		}
	}
}

func (m *agent) report() error {
	stats, err := m.aggregator.GetAll()
	if err != nil {
		return err
	}

	for key, value := range stats {
		if m.isShuttingDown && !m.softShuttingDown {
			close(m.workerPool.Jobs)
			return errors.New("sending the report was interrupted")
		}

		job := workerpool.NewJob(key, value)
		m.workerPool.Jobs <- job
	}

	return nil
}

func (m *agent) run(ctx context.Context, cfg config.Config) {
	limiter := ratelimiter.New(1*time.Second, cfg.RateLimit)
	err := limiter.Run(ctx)
	if err != nil {
		m.errorChan<- err
	}

	for {
		if m.isShuttingDown && !m.softShuttingDown {
			_ = limiter.Close(ctx)
			goto Exit
		}

		select {
		case job, ok := <-m.workerPool.Jobs:
			if !ok {
				_ = limiter.Close(ctx)
				goto Exit
			}
			limiter.Wait()

			go func() {
				err = m.sender.Send(job.Key, job.Value)
				if err != nil {
					m.errorChan<- err
				}
			}()
		}
	}
	Exit:
}

func (m *agent) Close(ctx context.Context) error {
	fmt.Println("agent closing")
	m.isShuttingDown = true
	return nil
}