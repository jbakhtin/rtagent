package agent // Package agent TODO: rename to agent

import (
	"context"
	"fmt"
	"net/http"
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
	pollInterval time.Duration
	reportInterval time.Duration

	memStats runtime.MemStats

	pollCounter counter
	randomValue gauge

	ctx    context.Context // TODO: For what?
	cancel context.CancelFunc
}

// NewMonitor - констурктор для Monitor
func NewMonitor(ctx context.Context, pollInterval, reportInterval time.Duration) Monitor {
	ctx, cancel := context.WithCancel(ctx)

	return Monitor{
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		ctx:            ctx,
		cancel:         cancel,
	}
}

// pooling - инициирует забор данных с заданным интервалом monitor.pollInterval
func (monitor *Monitor) polling() {
	ticker := time.NewTicker(monitor.pollInterval)

	for {
		select {
		case <-ticker.C:
			monitor.poll()
		case <-monitor.ctx.Done():
			fmt.Println("Сбор метрик приостановлен!")
			return
		}
	}
}

func (monitor *Monitor) poll() {
	runtime.ReadMemStats(&monitor.memStats)
	monitor.randomValue = 12 // it is human randomizer
	monitor.pollCounter += 1 // TODO: исправить гонку

}

// reporting - инициирует отправку данных с заданным интервалом monitor.reportInterval
func (monitor *Monitor) reporting() () {
	ticker := time.NewTicker(monitor.reportInterval)

	for {
		select {
		case <-ticker.C:
			monitor.report()
		case <-monitor.ctx.Done():
			fmt.Println("Отправка метрики приостановлена!")
			return
		}
	}
}

// report - отправить данные
func (monitor *Monitor) report() () {
	for key, value := range monitor.GetStats() {
		endpoint := "http://127.0.0.1:8080/update/" + value.Type() + "/" + key + "/" + fmt.Sprint(value)
		req, err := http.NewRequest(http.MethodPost, endpoint, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Content-Type", "text/plain")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = res.Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	monitor.pollCounter = 0
}

// Start - запустить мониторинг
func (monitor *Monitor) Start () {
	go monitor.polling()
	go monitor.reporting()
}

// Stop - остановить мониторинг
func (monitor *Monitor) Stop () {
	monitor.cancel()
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