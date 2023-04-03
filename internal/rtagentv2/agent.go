package rtagentv2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	gopsutil "github.com/shirou/gopsutil/v3/mem"

	handlerModels "github.com/jbakhtin/rtagent/internal/server/models"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/types"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type Monitor struct {
	sc    sync.Mutex
	loger *zap.Logger

	serverAddress              string
	pollInterval               time.Duration
	reportInterval             time.Duration
	acceptableCountAgentErrors int

	pollCounter types.Counter

	worker Worker
}

func NewMonitor(cfg config.Config, logger *zap.Logger) (Monitor, error) {
	worker, _ := NewWorker(cfg, 5)
	return Monitor{
		loger:                      logger,
		serverAddress:              fmt.Sprintf("http://%s", cfg.Address), //TODO: переделать зависимость от http/https
		pollInterval:               cfg.PollInterval,
		reportInterval:             cfg.ReportInterval,
		acceptableCountAgentErrors: cfg.AcceptableCountAgentErrors,
		worker:                     worker,
	}, nil
}

// Start - запустить мониторинг
func (m *Monitor) Start(cfg config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chanErr := make(chan error)

	go m.polling(ctx, cfg, chanErr)
	go m.reporting(ctx, cfg, chanErr)
	go m.worker.Run(ctx)

	var errCount int
	var err error

	err = func() error {
		for {
			select {
			case err = <-chanErr:
				errCount++
				m.loger.Info(err.Error())

				if errCount > m.acceptableCountAgentErrors {

					m.loger.Info(fmt.Sprintf("превышено количество (%v) допустимых ошибок", m.acceptableCountAgentErrors))
					cancel()
				}
			case <-ctx.Done():
				m.loger.Info("завершаем работу агента")
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
			err := m.poll(cfg)
			if err != nil {
				chanError <- err
			}
		case <-ctx.Done():
			m.loger.Info("сбор метрик приостановлен")
			return
		}
	}
}

func (m *Monitor) poll(cfg config.Config) error {
	m.sc.Lock()
	defer m.sc.Unlock()

	m.pollCounter++
	return nil
}

// reporting - инициирует отправку данных с заданным интервалом monitor.reportInterval
func (m *Monitor) reporting(ctx context.Context, cfg config.Config, chanError chan error) {
	ticker := time.NewTicker(m.reportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("test")
			err := m.reportRuntime()
			if err != nil {
				chanError <- err
			}

			err = m.reportGopsutil()
			if err != nil {
				chanError <- err
			}

			m.pollCounter = 0
		case <-ctx.Done():
			m.loger.Info("отправка метрики приостановлена")
			return
		}
	}
}

func (m *Monitor) reportRuntime() error {
	for key, value := range m.getStatsRuntime() {
		fmt.Println("test 1")
		job := Job{
			key,
			value,
		}

		m.worker.workerPool.Jobs <- job
	}
	return nil
}

func (m *Monitor) reportGopsutil() error {
	for key, value := range m.getStatsGopsutil() {
		fmt.Println("test 2")
		job := Job{
			key,
			value,
		}

		m.worker.workerPool.Jobs <- job
	}
	return nil
}

// GetStats - Поулчить слайс содержщий последние акутальные данные
func (m *Monitor) getStatsRuntime() map[string]types.Metricer {
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	result := map[string]types.Metricer{}

	// memStats
	result["Alloc"] = types.Gauge(memStats.Alloc)
	result["Frees"] = types.Gauge(memStats.Frees)
	result["HeapAlloc"] = types.Gauge(memStats.HeapAlloc)
	result["BuckHashSys"] = types.Gauge(memStats.BuckHashSys)
	result["GCSys"] = types.Gauge(memStats.GCSys)
	result["GCCPUFraction"] = types.Gauge(memStats.GCCPUFraction)
	result["HeapIdle"] = types.Gauge(memStats.HeapIdle)
	result["HeapInuse"] = types.Gauge(memStats.HeapInuse)
	result["HeapObjects"] = types.Gauge(memStats.HeapObjects)
	result["HeapReleased"] = types.Gauge(memStats.HeapReleased)
	result["HeapSys"] = types.Gauge(memStats.HeapSys)
	result["LastGC"] = types.Gauge(memStats.LastGC)
	result["Lookups"] = types.Gauge(memStats.Lookups)
	result["MCacheInuse"] = types.Gauge(memStats.MCacheInuse)
	result["MCacheSys"] = types.Gauge(memStats.MCacheSys)
	result["MSpanInuse"] = types.Gauge(memStats.MSpanInuse)
	result["MSpanSys"] = types.Gauge(memStats.MSpanSys)
	result["Mallocs"] = types.Gauge(memStats.Mallocs)
	result["NextGC"] = types.Gauge(memStats.NextGC)
	result["NumForcedGC"] = types.Gauge(memStats.NumForcedGC)
	result["NumGC"] = types.Gauge(memStats.NumGC)
	result["OtherSys"] = types.Gauge(memStats.OtherSys)
	result["PauseTotalNs"] = types.Gauge(memStats.PauseTotalNs)
	result["StackInuse"] = types.Gauge(memStats.StackInuse)
	result["StackSys"] = types.Gauge(memStats.StackSys)
	result["Sys"] = types.Gauge(memStats.Sys)
	result["TotalAlloc"] = types.Gauge(memStats.TotalAlloc)

	// Custom stats
	result["PollCount"] = m.pollCounter
	result["RandomValue"] = types.Gauge(rand.Int())

	return result
}

func (m *Monitor) getStatsGopsutil() map[string]types.Metricer {
	memStats, _ := gopsutil.VirtualMemory()

	result := map[string]types.Metricer{}

	// memStats
	result["TotalMemory"] = types.Gauge(memStats.Total)
	result["FreeMemory"] = types.Gauge(memStats.Free)
	result["CPUutilization1"] = types.Gauge(memStats.Used)

	return result
}

// --- worker pool

type Job struct {
	key   string
	value types.Metricer
}

type WorkerPool struct {
	Jobs chan Job
}

type Worker struct {
	config     config.Config
	workerPool WorkerPool
	rateLimit  int
}

func NewWorker(cfg config.Config, rateLimit int) (Worker, error) {
	workerPool := WorkerPool{
		make(chan Job),
	}

	return Worker{
		cfg,
		workerPool,
		rateLimit,
	}, nil
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		ticker := time.NewTicker(time.Second)

		for i := 0; i < w.rateLimit; i++ {

			select {
			case <-ctx.Done():
				break
			case <-ticker.C:
				break
			case job := <-w.workerPool.Jobs:
				//g, _ := errgroup.WithContext(ctx)

				go func() error {
					endpoint := fmt.Sprintf("%s/update/", "http://localhost:8080")
					metric, err := handlerModels.ToJSON(w.config, job.key, job.value)
					if err != nil {
						return err
					}

					metric.Hash, err = metric.CalcHash([]byte(w.config.KeyApp))
					if err != nil {
						return err
					}

					hash, err := metric.CalcHash([]byte(w.config.KeyApp))
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
				}()

				//err := g.Wait()
				//if err != nil {
				//	return err
				//}
			}
		}
	}
}
