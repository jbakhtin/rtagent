package agent // Package agent TODO: rename to agent

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"reflect"
	"runtime"
	"time"
)

// Metricer - интерфейс для сущности метрики
type Metricer interface {
	Type() string
}

// gauge - тип сущности метрики
type gauge float64

// Type - возвращает тип метрики
func (g gauge) Type() string {
	return reflect.TypeOf(g).Name()
}
// counter - тип сущности метрики
type counter int64

// Type - возвращает тип мерики
func (c counter) Type() string {
	return reflect.TypeOf(c).Name()
}

// increment - имплементирует нестандартную логику изменения значения
func (c *counter) increment(count counter) {
	 *c += count
}

// flush - имплементирует логику изменения значения
func (c *counter) flush() {
	*c = 0
}

type Monitor struct {
	serverAddress string
	pollInterval time.Duration
	reportInterval time.Duration

	memStats runtime.MemStats

	pollCounter counter
	randomValue gauge
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
		case <- ctx.Done():
			fmt.Println("Сбор метрик приостановлен!")
			return
		}
	}
}

func (monitor *Monitor) poll() error {
	runtime.ReadMemStats(&monitor.memStats)
	monitor.randomValue = 12 // TODO: реализовать рандомайзер
	monitor.pollCounter += 1 // TODO: исправить гонку

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
		case <- ctx.Done():
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
			"type": value.Type(),
			"key": key,
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
func Start (serverAddress string, pollInterval, reportInterval time.Duration) error {
	ctx, cancel := context.WithCancel(context.Background())

	monitor := Monitor{
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		serverAddress: serverAddress,
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

//GetStats - Поулчить слайс содержщий последние акутальные данные
func(monitor Monitor) GetStats() map[string]Metricer {
	result := map[string]Metricer{}

	// memStats
	result["Alloc"] = gauge(monitor.memStats.Alloc)
	result["Frees"] = gauge(monitor.memStats.Frees)
	result["HeapAlloc"] = gauge(monitor.memStats.HeapAlloc)
	result["BuckHashSys"] = gauge(monitor.memStats.BuckHashSys)
	result["GCSys"] = gauge(monitor.memStats.GCSys)
	result["HeapIdle"] = gauge( monitor.memStats.HeapIdle)
	result["HeapInuse"] = gauge( monitor.memStats.HeapInuse)
	result["HeapObjects"] = gauge( monitor.memStats.HeapObjects)
	result["HeapReleased"] = gauge( monitor.memStats.HeapReleased)
	result["HeapSys"] = gauge( monitor.memStats.HeapSys)
	result["LastGC"] = gauge( monitor.memStats.LastGC)
	result["Lookups"] = gauge( monitor.memStats.Lookups)
	result["MCacheInuse"] = gauge( monitor.memStats.MCacheInuse)
	result["MCacheSys"] = gauge( monitor.memStats.MCacheSys)
	result["MSpanInuse"] = gauge( monitor.memStats.MSpanInuse)
	result["MSpanSys"] = gauge( monitor.memStats.MSpanSys)
	result["Mallocs"] = gauge( monitor.memStats.Mallocs)
	result["NextGC"] = gauge( monitor.memStats.NextGC)
	result["NumForcedGC"] = gauge( monitor.memStats.NumForcedGC)
	result["NumGC"] = gauge( monitor.memStats.NumGC)
	result["OtherSys"] = gauge( monitor.memStats.OtherSys)
	result["PauseTotalNs"] = gauge( monitor.memStats.PauseTotalNs)
	result["StackInuse"] = gauge( monitor.memStats.StackInuse)
	result["StackSys"] = gauge( monitor.memStats.StackSys)
	result["Sys"] = gauge( monitor.memStats.Sys)
	result["TotalAlloc"] = gauge( monitor.memStats.TotalAlloc)

	// Custom stats
	result["PollCount"] = monitor.pollCounter
	result["RandomValue"] = monitor.randomValue

	return result
}