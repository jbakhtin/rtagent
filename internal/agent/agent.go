package agent // Package agent TODO: rename to agent

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jbakhtin/rtagent/internal/models"
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
	m.randomValue = 12 // TODO: реализовать рандомайзер
	m.pollCounter++

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

	m.pollCounter = 0

	return nil
}

// Start - запустить мониторинг
func Start(serverAddress string, pollInterval, reportInterval time.Duration) error {
	ctx, cancel := context.WithCancel(context.Background())

	monitor := Monitor{
		serverAddress:  serverAddress,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		memStats: &runtime.MemStats{},
	}

	chanErr := make(chan error)

	go monitor.polling(ctx, chanErr)
	go monitor.reporting(ctx, chanErr)
	err := <-chanErr
	// TODO: реализовать счетчик  ошибок, если ошибок 10-20 - завершаем работу
	cancel()

	return err
}

// GetStats - Поулчить слайс содержщий последние акутальные данные
func (m Monitor) GetStats() map[string]Metricer {
	result := map[string]Metricer{}

	// memStats
	result["Alloc"] = models.Gauge(m.memStats.Alloc)
	result["Frees"] = models.Gauge(m.memStats.Frees)
	result["HeapAlloc"] = models.Gauge(m.memStats.HeapAlloc)
	result["BuckHashSys"] = models.Gauge(m.memStats.BuckHashSys)
	result["GCSys"] = models.Gauge(m.memStats.GCSys)
	result["HeapIdle"] = models.Gauge(m.memStats.HeapIdle)
	result["HeapInuse"] = models.Gauge(m.memStats.HeapInuse)
	result["HeapObjects"] = models.Gauge(m.memStats.HeapObjects)
	result["HeapReleased"] = models.Gauge(m.memStats.HeapReleased)
	result["HeapSys"] = models.Gauge(m.memStats.HeapSys)
	result["LastGC"] = models.Gauge(m.memStats.LastGC)
	result["Lookups"] = models.Gauge(m.memStats.Lookups)
	result["MCacheInuse"] = models.Gauge(m.memStats.MCacheInuse)
	result["MCacheSys"] = models.Gauge(m.memStats.MCacheSys)
	result["MSpanInuse"] = models.Gauge(m.memStats.MSpanInuse)
	result["MSpanSys"] = models.Gauge(m.memStats.MSpanSys)
	result["Mallocs"] = models.Gauge(m.memStats.Mallocs)
	result["NextGC"] = models.Gauge(m.memStats.NextGC)
	result["NumForcedGC"] = models.Gauge(m.memStats.NumForcedGC)
	result["NumGC"] = models.Gauge(m.memStats.NumGC)
	result["OtherSys"] = models.Gauge(m.memStats.OtherSys)
	result["PauseTotalNs"] = models.Gauge(m.memStats.PauseTotalNs)
	result["StackInuse"] = models.Gauge(m.memStats.StackInuse)
	result["StackSys"] = models.Gauge(m.memStats.StackSys)
	result["Sys"] = models.Gauge(m.memStats.Sys)
	result["TotalAlloc"] = models.Gauge(m.memStats.TotalAlloc)

	// Custom stats
	result["PollCount"] = m.pollCounter
	result["RandomValue"] = m.randomValue

	return result
}
