package rtagent

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

const (
	pollInterval = time.Second * 2
	reportInterval = time.Second * 10
)

type Gauge float64

type Counter struct {
	value int64
}

func (c *Counter) Increment() int64 {
	c.value++
	return c.value
}

func (c *Counter) SetZero() int64 {
	c.value = 0
	return c.value
}

// TODO: реализовать кастомный тип метрики PollCount - counter
// TODO: реализовать кастомный тип метрики RandomValue - gauge

type Monitor struct {
	pollInterval time.Duration
	reportInterval time.Duration

	memStats runtime.MemStats
	poolCounter Counter
	randomValue Gauge

	ctx    context.Context
	cancel context.CancelFunc
}

func (monitor *Monitor) pool() {
	ticker := time.NewTicker(monitor.pollInterval)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Метрика собрана")
			runtime.ReadMemStats(&monitor.memStats) // TODO: вызов промежуточной функции
			monitor.poolCounter.Increment()
			monitor.randomValue = 12
		case <-monitor.ctx.Done():
			fmt.Println("Сбор метрик приостановлен!")
			return
		}
	}
}

func (monitor *Monitor) report() () {
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
			for key, value := range monitor.GetMemStats() {
				endpoint = "http://127.0.0.1:8080/update/" + fmt.Sprintf("%T", value) + "/" + key + "/" + fmt.Sprint(value)
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

			// Send randomValue metrics
			endpoint = "http://127.0.0.1:8080/update/gauge/RandomValue/" + fmt.Sprint(monitor.randomValue)

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

			// Send pullCounter metrics
			endpoint = "http://127.0.0.1:8080/update/counter/PulCounter/" + fmt.Sprint(monitor.poolCounter.value)
			monitor.poolCounter.SetZero()

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

		case <-monitor.ctx.Done():
			fmt.Println("Отправка метрики приостановлена!")
			return
		}
	}

	// TODO: вызов функции пост обработки метрик
}

func (monitor *Monitor) Start () {
	go monitor.pool()
	go monitor.report()
}

func (monitor *Monitor) Stop () {
	monitor.cancel()
}

func(monitor Monitor) GetMemStats() map[string]Gauge {
	result := make(map[string]Gauge, 20)

	result["Alloc"] = Gauge(monitor.memStats.Alloc)
	result["Frees"] = Gauge(monitor.memStats.Frees)
	result["HeapAlloc"] = Gauge(monitor.memStats.HeapAlloc)
	result["BuckHashSys"] = Gauge(monitor.memStats.BuckHashSys)
	result["GCSys"] = Gauge(monitor.memStats.GCSys)
	result["HeapIdle"] =Gauge( monitor.memStats.HeapIdle)

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
