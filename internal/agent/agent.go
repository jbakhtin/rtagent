package agent

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/agent/sender"
	"sync"
	"time"

	"github.com/jbakhtin/rtagent/pkg/ratelimiter"

	"github.com/jbakhtin/rtagent/internal/agent/aggregator"
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

type Monitor struct {
	aggregator                 Aggregator
	sender					Sender
	workerPool                 workerpool.WorkerPool
	logger                     *zap.Logger
	serverAddress              string
	acceptableCountAgentErrors int
	pollCounter                types.Counter
	sc                         sync.Mutex
	reportInterval             time.Duration
	pollInterval               time.Duration
}

func New(cfg config.Config, logger *zap.Logger) (*Monitor, error) {
	workerPool, err := workerpool.NewWorkerPool()
	if err != nil {
		return nil, err
	}

	aggregator, err := aggregator.New().
		WithConfig(cfg).
		WithDefaultCollectors().
		Build()

	sender, err := sender.New().
		WithConfig(cfg).
		Build()

	if err != nil {
		return nil, err
	}

	return &Monitor{
		logger:                     logger,
		serverAddress:              fmt.Sprintf("http://%s", cfg.Address), //TODO: переделать зависимость от http/https
		pollInterval:               cfg.PollInterval,
		reportInterval:             cfg.ReportInterval,
		acceptableCountAgentErrors: cfg.AcceptableCountAgentErrors,
		workerPool:                 workerPool,
		aggregator:                 aggregator,
		sender:                 sender,
	}, nil
}

// Start - запустить мониторинг
func (m *Monitor) Start(cfg config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chanErr := make(chan error)

	go m.polling(ctx, cfg, chanErr)
	go m.reporting(ctx, cfg, chanErr)
	go m.Run(ctx, cfg, chanErr)

	var errCount int
	var err error

	err = func() error {
		for {
			select {
			case err = <-m.aggregator.Err():
				errCount++
				m.logger.Info(err.Error())

				if errCount > m.acceptableCountAgentErrors {

					m.logger.Info(fmt.Sprintf("превышено количество (%v) допустимых ошибок", m.acceptableCountAgentErrors))
					cancel()
				}
			case err = <-chanErr:
				errCount++
				m.logger.Info(err.Error())

				if errCount > m.acceptableCountAgentErrors {

					m.logger.Info(fmt.Sprintf("превышено количество (%v) допустимых ошибок", m.acceptableCountAgentErrors))
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
func (m *Monitor) polling(ctx context.Context, cfg config.Config, chanError chan error) {
	ticker := time.NewTicker(m.pollInterval)
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
func (m *Monitor) reporting(ctx context.Context, cfg config.Config, chanError chan error) {
	ticker := time.NewTicker(m.reportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := m.report()
			if err != nil {
				chanError <- err
			}

			m.pollCounter = 0
		case <-ctx.Done():
			m.logger.Info("отправка метрики приостановлена")
			return
		}
	}
}

func (m *Monitor) report() error {
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

func (m *Monitor) Run(ctx context.Context, cfg config.Config, chanError chan error) {
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
