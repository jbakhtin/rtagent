package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/server/models"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/types"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"

	"github.com/go-resty/resty/v2"
)

type Monitor struct {
	sc    sync.Mutex
	loger *zap.Logger

	serverAddress  string
	pollInterval   time.Duration
	reportInterval time.Duration
	acceptableCountAgentErrors int

	pollCounter types.Counter
}

func New(cfg config.Config, logger *zap.Logger) (Monitor, error) {
	return Monitor{
		loger:          logger,
		serverAddress:  fmt.Sprintf("http://%s", cfg.Address), //TODO: переделать зависимость от http/https
		pollInterval:   cfg.PollInterval,
		reportInterval: cfg.ReportInterval,
		acceptableCountAgentErrors: cfg.AcceptableCountAgentErrors,
	}, nil
}

// Start - запустить мониторинг
func (m *Monitor) Start(cfg config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chanErr := make(chan error)

	go m.polling(ctx, cfg, chanErr)
	go m.reporting(ctx, cfg, chanErr)

	var errCount int
	var err error

	err = func () error{
		for {
			select {
			case err = <-chanErr:
				errCount++
				m.loger.Info(err.Error())

				if errCount > m.acceptableCountAgentErrors {
					// TODO: реализовать отправку количества ошибко на сервер
					// TODO: реализовать новый тип метрики - стринг отправлять описание ошибки на сервер

					m.loger.Info(fmt.Sprintf("превышено количество (%v) допустимых ошибок", m.acceptableCountAgentErrors))
					cancel()
				}
			case <-ctx.Done():
				m.loger.Info("завершаем работу агента")
				err := m.report(cfg)
				if err != nil {
					return err
				}
				return nil
			}
		}
	} ()

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
			err := m.reportJSON(cfg)
			if err != nil {
				chanError <- err
			}
		case <-ctx.Done():
			m.loger.Info("отправка метрики приостановлена")
			return
		}
	}
}

// report - отправить данные
func (m *Monitor) report(cfg config.Config) error {
	m.sc.Lock()
	defer m.sc.Unlock()

	for key, value := range m.getStats() {
		endpoint := fmt.Sprintf("%s/update/{type}/{key}/{value}", m.serverAddress)
		client := resty.New()
		_, err := client.R().SetHeaders(map[string]string{
			"Content-Type": "text/plain",
		}).SetPathParams(map[string]string{
			"type":  value.Type(),
			"key":   key,
			"value": fmt.Sprint(value),
		}).Post(endpoint)
		if err != nil {
			return err
		}
	}

	m.pollCounter = 0

	return nil
}

func (m *Monitor) reportJSON(cfg config.Config) error {
	for key, value := range m.getStats() {
		endpoint := fmt.Sprintf("%s/update/", m.serverAddress)
		metric, err := models.ToJSON(cfg, key, value)
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
	}

	m.pollCounter = 0
	return nil
}

// GetStats - Поулчить слайс содержщий последние акутальные данные
func (m *Monitor) getStats() map[string]types.Metricer {
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
