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
func (monitor *Monitor) polling(ctx context.Context, chanError chan error) {
	ticker := time.NewTicker(monitor.pollInterval)

	for {
		select {
		case <-ticker.C:
			err := monitor.poll()
			if err != nil {
				chanError <- err
			}
		case <-ctx.Done():
			fmt.Println("Сбор метрик приостановлен!")
			return
		}
	}
}

func (monitor *Monitor) poll() error {
	runtime.ReadMemStats(monitor.memStats)
	monitor.randomValue = 12 // TODO: реализовать рандомайзер
	monitor.pollCounter++    // TODO: исправить гонку

	return nil
}

// reporting - инициирует отправку данных с заданным интервалом monitor.reportInterval
func (monitor *Monitor) reporting(ctx context.Context, chanError chan error) {
	ticker := time.NewTicker(monitor.reportInterval)

	for {
		select {
		case <-ticker.C:
			err := monitor.report()
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
func (monitor *Monitor) report() error {
	for key, value := range monitor.GetStats() {
		client := resty.New()
		_, err := client.R().SetHeaders(map[string]string{
			"Content-Type": "text/plain",
		}).SetPathParams(map[string]string{
			"type":  value.Type(),
			"key":   key,
			"value": fmt.Sprint(value),
		}).Post(monitor.serverAddress + "/update/{type}/{key}/{value}")
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	monitor.pollCounter = 0

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
	// TODO: Выяснить правильный ли подход
	// принимаем решение, продолжить сбор и отправку метрик или останоить выполение программы
	// Решаю прикратить выполение программы
	cancel()

	return err
}

// GetStats - Поулчить слайс содержщий последние акутальные данные
func (monitor Monitor) GetStats() map[string]Metricer {
	result := map[string]Metricer{}

	// memStats
	result["Alloc"] = models.Gauge(monitor.memStats.Alloc)
	result["Frees"] = models.Gauge(monitor.memStats.Frees)
	result["HeapAlloc"] = models.Gauge(monitor.memStats.HeapAlloc)
	result["BuckHashSys"] = models.Gauge(monitor.memStats.BuckHashSys)
	result["GCSys"] = models.Gauge(monitor.memStats.GCSys)
	result["HeapIdle"] = models.Gauge(monitor.memStats.HeapIdle)
	result["HeapInuse"] = models.Gauge(monitor.memStats.HeapInuse)
	result["HeapObjects"] = models.Gauge(monitor.memStats.HeapObjects)
	result["HeapReleased"] = models.Gauge(monitor.memStats.HeapReleased)
	result["HeapSys"] = models.Gauge(monitor.memStats.HeapSys)
	result["LastGC"] = models.Gauge(monitor.memStats.LastGC)
	result["Lookups"] = models.Gauge(monitor.memStats.Lookups)
	result["MCacheInuse"] = models.Gauge(monitor.memStats.MCacheInuse)
	result["MCacheSys"] = models.Gauge(monitor.memStats.MCacheSys)
	result["MSpanInuse"] = models.Gauge(monitor.memStats.MSpanInuse)
	result["MSpanSys"] = models.Gauge(monitor.memStats.MSpanSys)
	result["Mallocs"] = models.Gauge(monitor.memStats.Mallocs)
	result["NextGC"] = models.Gauge(monitor.memStats.NextGC)
	result["NumForcedGC"] = models.Gauge(monitor.memStats.NumForcedGC)
	result["NumGC"] = models.Gauge(monitor.memStats.NumGC)
	result["OtherSys"] = models.Gauge(monitor.memStats.OtherSys)
	result["PauseTotalNs"] = models.Gauge(monitor.memStats.PauseTotalNs)
	result["StackInuse"] = models.Gauge(monitor.memStats.StackInuse)
	result["StackSys"] = models.Gauge(monitor.memStats.StackSys)
	result["Sys"] = models.Gauge(monitor.memStats.Sys)
	result["TotalAlloc"] = models.Gauge(monitor.memStats.TotalAlloc)

	// Custom stats
	result["PollCount"] = monitor.pollCounter
	result["RandomValue"] = monitor.randomValue

	return result
}
