package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/types"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type Monitor struct {
	sc sync.Mutex
	log *zap.Logger

	serverAddress  string
	pollInterval   time.Duration
	reportInterval time.Duration

	pollCounter types.Counter
}

func New(ctf config.Config, logger *zap.Logger) (Monitor, error){
	return Monitor{
		log: logger,
		serverAddress: ctf.Address,
		pollInterval: ctf.PollInterval,
		reportInterval: ctf.ReportInterval,
	}, nil
}

// Start - запустить мониторинг
func (m *Monitor) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chanErr := make(chan error)

	go m.polling(ctx, chanErr)
	go m.reporting(ctx, chanErr)

	var errCount int
	var err error

	for {
		select {
		case err = <-chanErr:
			errCount++
			m.log.Info(err.Error())

			if errCount > 10 { // TODO: вынести в константу
				// TODO: реализовать отправку количества ошибко на сервер
				// TODO: реализовать новый тип метрики - стринг отправлять описание ошибки на сервер

				m.log.Info(fmt.Sprintf("превышено количество (%v) допустимых ошибок", 10))
				cancel()
			}
			case <-ctx.Done():
				m.log.Info("завершаем работу агента")
				m.log.Info("ожидаем перед окончательным завершением")
				time.Sleep(m.reportInterval)
				goto exit
		}
	}

	exit:
	return err
}

// pooling - инициирует забор данных с заданным интервалом monitor.pollInterval
func (m *Monitor) polling(ctx context.Context, chanError chan error) {
	ticker := time.NewTicker(m.pollInterval)

	for {
		select {
		case <-ticker.C:
			err := m.poll()
			if err != nil {
				chanError <- err
			}
		case <-ctx.Done():
			m.log.Info("сбор метрик приостановлен")
			return
		}
	}
}

func (m *Monitor) poll() error {
	m.sc.Lock()
	defer m.sc.Unlock()

	m.pollCounter++
	return nil
}

// reporting - инициирует отправку данных с заданным интервалом monitor.reportInterval
func (m *Monitor) reporting(ctx context.Context, chanError chan error) {
	ticker := time.NewTicker(m.reportInterval)

	for {
		select {
		case <-ticker.C:
			err := m.reportV2()
			if err != nil {
				chanError <- err
			}
		case <-ctx.Done():
			m.log.Info("отправка метрики приостановлена")
			return
		}
	}
}

// report - отправить данные
func (m *Monitor) report() error {
	m.sc.Lock()
	defer m.sc.Unlock()

	for key, value := range m.GetStats() {
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

func (m *Monitor) reportV2() error {
	for key, value := range m.GetStats() {
		endpoint := fmt.Sprintf("%s/update/", m.serverAddress)
		metric := models.Metrics{
			ID: key,
			MType: value.Type(),
		}
		switch v := value.(type) {
			case types.Counter:
				metric.Delta = &v
			case types.Gauge:
				metric.Value = &v
		}

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
func (m *Monitor) GetStats() map[string]types.Metricer {
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
