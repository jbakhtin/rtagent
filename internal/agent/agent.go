package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/types"
	"golang.org/x/exp/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
)

type Monitor struct {
	serverAddress  string
	pollInterval   time.Duration
	reportInterval time.Duration

	pollCounter types.Counter
}

func New(serverAddress string, pollInterval, reportInterval time.Duration) (Monitor, error){
	return Monitor{
		serverAddress:serverAddress,
		pollInterval: pollInterval,
		reportInterval: reportInterval,
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
		request.Header.Set("Content-Type", "application/json")

		client := http.Client{}
		_, err = client.Do(request)
		//defer response.Body.Close()

	}

	m.pollCounter = 0

	return nil
}

// GetStats - Поулчить слайс содержщий последние акутальные данные
func (m Monitor) GetStats() map[string]types.Metricer {
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	result := map[string]types.Metricer{}

	// memStats
	result["Alloc"] = types.Gauge(memStats.Alloc)
	result["Frees"] = types.Gauge(memStats.Frees)
	result["HeapAlloc"] = types.Gauge(memStats.HeapAlloc)
	result["BuckHashSys"] = types.Gauge(memStats.BuckHashSys)
	result["GCSys"] = types.Gauge(memStats.GCSys)
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
