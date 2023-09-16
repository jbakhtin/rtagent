package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jbakhtin/rtagent/pkg/ratelimiter"

	"github.com/jbakhtin/rtagent/internal/agent/workerpool"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/types"
	"go.uber.org/zap"
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
	logger                     *zap.Logger
	serverAddress              string
	sc                         sync.Mutex
	cfg Configer
}

// Start - запустить мониторинг
func (m *agent) Start(cfg config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chanErr := make(chan error)

	go m.polling(ctx, cfg, chanErr)
	go m.reporting(ctx, cfg, chanErr)
	go m.run(ctx, cfg, chanErr)

	var errCount int
	var err error

	err = func() error {
		for {
			select {
			case err = <-m.aggregator.Err():
				errCount++
				m.logger.Info(err.Error())

				if errCount > m.cfg.GetAcceptableCountAgentErrors() {

					m.logger.Info(fmt.Sprintf("превышено количество (%v) допустимых ошибок", m.cfg.GetAcceptableCountAgentErrors()))
					cancel()
				}
			case err = <-chanErr:
				errCount++
				m.logger.Info(err.Error())

				if errCount > m.cfg.GetAcceptableCountAgentErrors() {

					m.logger.Info(fmt.Sprintf("превышено количество (%v) допустимых ошибок", m.cfg.GetAcceptableCountAgentErrors()))
					cancel()
				}
			case <-ctx.Done():
				m.logger.Info("завершаем работу агента")
				if err != nil {
					return err
				}
				return nil
			}
		}
	}()

	return err
}

// pooling - инициирует забор данных с заданным интервалом monitor.pollInterval
func (m *agent) polling(ctx context.Context, cfg config.Config, chanError chan error) {
	ticker := time.NewTicker(m.cfg.GetPollInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
				m.aggregator.Pool(ctx)
		case err := <-m.aggregator.Err():
			chanError<- err
		case <-ctx.Done():
			return
		}
	}
}

// reporting - инициирует отправку данных с заданным интервалом monitor.reportInterval
func (m *agent) reporting(ctx context.Context, cfg config.Config, chanError chan error) {
	ticker := time.NewTicker(m.cfg.GetReportInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := m.report()
			if err != nil {
				chanError <- err
			}

		case <-ctx.Done():
			m.logger.Info("отправка метрики приостановлена")
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
		job := workerpool.NewJob(key, value)
		m.workerPool.Jobs <- job
	}
	return nil
}

func (m *agent) run(ctx context.Context, cfg config.Config, chanError chan error) {
	limiter := ratelimiter.New(1*time.Second, cfg.RateLimit)
	err := limiter.Run(ctx)
	if err != nil {
		chanError <- err
	}

	for {
		select {
		case <-ctx.Done():
		case job := <-m.workerPool.Jobs:
			limiter.Wait()

			go func() {
				err = m.sender.Send(job.Key, job.Value)
				if err != nil {
					m.logger.Info(err.Error())
					chanError <- err
				}
			}()
		}
	}
}
