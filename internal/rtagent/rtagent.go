package rtagent

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

const (
	pollInterval = time.Second * 2
	reportInterval = time.Second * 10
)

type Metricer interface {
	Type() string
}

type gauge float64
type counter int64

func (g gauge) Type() string {
	return reflect.TypeOf(g).Name()
}

func (c counter) Type() string {
	return reflect.TypeOf(c).Name()
}

func (c *counter) increment(count counter) {
	 *c += count
}

func (c *counter) flush() {
	*c = 0
}

// TODO: реализовать кастомный тип метрики PollCount - counter
// TODO: реализовать кастомный тип метрики RandomValue - gauge

type Monitor struct {
	pollInterval time.Duration
	reportInterval time.Duration

	memStats runtime.MemStats

	poolCounter counter
	randomValue gauge

	ctx    context.Context
	cancel context.CancelFunc
}

func (monitor *Monitor) pooling() {
	ticker := time.NewTicker(monitor.pollInterval)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Метрика собрана")
			runtime.ReadMemStats(&monitor.memStats) // TODO: вызов промежуточной функции
			monitor.randomValue = 12
			monitor.poolCounter.increment(1)
		case <-monitor.ctx.Done():
			fmt.Println("Сбор метрик приостановлен!")
			return
		}
	}
}

func (monitor *Monitor) reporting() () {
	ticker := time.NewTicker(monitor.reportInterval)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Метрика отрпавлена!")
			var endpoint string
			var client *http.Client
			var req *http.Request
			var err error
			var res *http.Response

			// Send base metrics
			for key, value := range monitor.GetStats() {
				endpoint = "http://127.0.0.1:8080/update/" + value.Type() + "/" + key + "/" + fmt.Sprint(value)
				//fmt.Println(endpoint)

				req, err = http.NewRequest(http.MethodPost, endpoint, nil)
				if err != nil {
					fmt.Println(err)
					return
				}

				req.Header.Add("Content-Type", "text/plain")

				// конструируем клиент
				client = &http.Client{}
				// отправляем запрос
				res, err = client.Do(req)
				if err != nil {
					fmt.Println(err)
					return
				}

				// печатаем код ответа
				fmt.Println("Статус-код ", res.Status)
				res.Body.Close()
			}
			monitor.poolCounter.flush()
		case <-monitor.ctx.Done():
			fmt.Println("Отправка метрики приостановлена!")
			return
		}
	}
}

func (monitor *Monitor) Start () {
	go monitor.pooling()
	go monitor.reporting()
}

func (monitor *Monitor) Stop () {
	monitor.cancel()
}

func(monitor Monitor) GetStats() map[string]Metricer {
	result := map[string]Metricer{}

	result["Alloc"] = gauge(monitor.memStats.Alloc)
	result["Frees"] = gauge(monitor.memStats.Frees)
	result["HeapAlloc"] = gauge(monitor.memStats.HeapAlloc)
	result["BuckHashSys"] = gauge(monitor.memStats.BuckHashSys)
	result["GCSys"] = gauge(monitor.memStats.GCSys)
	result["HeapIdle"] = gauge( monitor.memStats.HeapIdle)

	result["PoolCount"] = monitor.poolCounter
	result["RandomValue"] = monitor.randomValue

	return result
}

func NewMonitor(ctx context.Context, pollInterval, reportInterval time.Duration) Monitor {
	ctx, cancel := context.WithCancel(ctx)

	return Monitor{
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		ctx:            ctx,
		cancel:         cancel,
	}
}
