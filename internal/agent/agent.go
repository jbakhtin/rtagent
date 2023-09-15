package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jbakhtin/rtagent/pkg/ratelimiter"

	"github.com/jbakhtin/rtagent/internal/agent/aggregator"
	"github.com/jbakhtin/rtagent/internal/agent/workerpool"

	handlerModels "github.com/jbakhtin/rtagent/internal/server/models"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/types"
	"go.uber.org/zap"
)

type Aggregator interface {
	Run(ctx context.Context)
	GetAll() (map[string]types.Metricer, error)
	Err() chan error
}

type Monitor struct {
	aggregator                 Aggregator
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
	}, nil
}

// Start - запустить мониторинг
func (m *Monitor) Start(cfg config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chanErr := make(chan error)

	go m.aggregator.Run(ctx)
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
				err = m.sendJSON(ctx, cfg, job)
				if err != nil {
					m.logger.Info(err.Error())
					chanError <- err
				}
			}()
		}
	}
}

func (m *Monitor) sendJSON(ctx context.Context, cfg config.Config, job workerpool.Job) error {
	endpoint := fmt.Sprintf("%s/update/", fmt.Sprintf("http://%s", cfg.Address))
	metric, err := handlerModels.ToJSON(cfg, job.Key, job.Value)
	if err != nil {
		return err
	}

	metric.Hash, err = metric.CalcHash([]byte(cfg.KeyApp))
	if err != nil {
		return err
	}

	hash, err := metric.CalcHash([]byte(cfg.KeyApp))
	if err != nil {
		return err
	}
	metric.Hash = hash

	buf, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if err = response.Body.Close(); err != nil {
		return err
	}

	return nil
}
