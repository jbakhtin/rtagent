package agent

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/models"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jbakhtin/rtagent/internal/types"
)

// Metricer - интерфейс для сущности метрики
type Metricer interface {
	Type() string
}

type Monitor struct {
	serverAddress  string
	pollInterval   time.Duration
	reportInterval time.Duration

	memStats *runtime.MemStats

	pollCounter models.Counter
	randomValue models.Gauge
}

func New(serverAddress string, pollInterval, reportInterval time.Duration) (Monitor, error){
	return Monitor{
		serverAddress:serverAddress,
		pollInterval: pollInterval,
		reportInterval: reportInterval,
		memStats: &runtime.MemStats{},
	}, nil
}

// Start - запустить мониторинг
func (m Monitor) Start() error {
	ctx, cancel := context.WithCancel(context.Background())

	chanErr := make(chan error)

	go m.polling(ctx, chanErr)
	go m.reporting(ctx, chanErr)
	err := <-chanErr
	// TODO: реализовать счетчик  ошибок, если ошибок 10-20 - завершаем работу
	cancel()

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
			fmt.Println("Сбор метрик приостановлен!")
			return
		}
	}
}

func (m *Monitor) poll() error {
	runtime.ReadMemStats(m.memStats)
	m.randomValue.MValue = types.Gauge(rand.Intn(12)) //TODO: реализовать сеттер
	m.pollCounter.Increment()

	return nil
}

// reporting - инициирует отправку данных с заданным интервалом monitor.reportInterval
func (m *Monitor) reporting(ctx context.Context, chanError chan error) {
	ticker := time.NewTicker(m.reportInterval)

	for {
		select {
		case <-ticker.C:
			err := m.report()
			if err != nil {
				chanError <- err
			}
		case <-ctx.Done():
			fmt.Println("Отправка метрики приостановлена!")
			return
		}
	}
}

// report - отправить данные
func (m *Monitor) report() error {
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
			fmt.Println(err)
			return err
		}
	}

	m.pollCounter.Flush()

	return nil
}

// GetStats - Поулчить слайс содержщий последние акутальные данные
func (m Monitor) GetStats() map[string]Metricer {
	result := map[string]Metricer{}

	// memStats
	result["Alloc"] = types.Gauge(m.memStats.Alloc)
	result["Frees"] = types.Gauge(m.memStats.Frees)
	result["HeapAlloc"] = types.Gauge(m.memStats.HeapAlloc)
	result["BuckHashSys"] = types.Gauge(m.memStats.BuckHashSys)
	result["GCSys"] = types.Gauge(m.memStats.GCSys)
	result["HeapIdle"] = types.Gauge(m.memStats.HeapIdle)
	result["HeapInuse"] = types.Gauge(m.memStats.HeapInuse)
	result["HeapObjects"] = types.Gauge(m.memStats.HeapObjects)
	result["HeapReleased"] = types.Gauge(m.memStats.HeapReleased)
	result["HeapSys"] = types.Gauge(m.memStats.HeapSys)
	result["LastGC"] = types.Gauge(m.memStats.LastGC)
	result["Lookups"] = types.Gauge(m.memStats.Lookups)
	result["MCacheInuse"] = types.Gauge(m.memStats.MCacheInuse)
	result["MCacheSys"] = types.Gauge(m.memStats.MCacheSys)
	result["MSpanInuse"] = types.Gauge(m.memStats.MSpanInuse)
	result["MSpanSys"] = types.Gauge(m.memStats.MSpanSys)
	result["Mallocs"] = types.Gauge(m.memStats.Mallocs)
	result["NextGC"] = types.Gauge(m.memStats.NextGC)
	result["NumForcedGC"] = types.Gauge(m.memStats.NumForcedGC)
	result["NumGC"] = types.Gauge(m.memStats.NumGC)
	result["OtherSys"] = types.Gauge(m.memStats.OtherSys)
	result["PauseTotalNs"] = types.Gauge(m.memStats.PauseTotalNs)
	result["StackInuse"] = types.Gauge(m.memStats.StackInuse)
	result["StackSys"] = types.Gauge(m.memStats.StackSys)
	result["Sys"] = types.Gauge(m.memStats.Sys)
	result["TotalAlloc"] = types.Gauge(m.memStats.TotalAlloc)

	// Custom stats
	result["PollCount"] = m.pollCounter.MValue
	result["RandomValue"] = m.randomValue.MValue

	return result
}
