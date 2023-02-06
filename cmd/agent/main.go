package main

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

type gauge float64
type counter int64

type AnotherMetric struct {
	PollCount counter
	RandomValue gauge
}

// TODO: реализовать кастомный тип метрики PollCount - counter
// TODO: реализовать кастомный тип метрики RandomValue - gauge

type Monitor struct {
	pollInterval time.Duration
	reportInterval time.Duration

	memStats runtime.MemStats

	ctx    context.Context
	cancel context.CancelFunc
}

func (monitor *Monitor) pool() {
	ticker := time.NewTicker(monitor.pollInterval)



	for {
		select {
			case <-ticker.C:
				fmt.Println("Метрика собрана")
				runtime.ReadMemStats(&monitor.memStats)
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

				for key, value := range monitor.GetMemStats() {
					endpoint := "http://127.0.0.1:8080/update/" + fmt.Sprintf("%T", value) + "/" + key + "/" + fmt.Sprint(value)
					//fmt.Println(endpoint)

					_, err := http.NewRequest(http.MethodPost, endpoint, nil)
					if err != nil {
						fmt.Println(err)
					}
				}

			case <-monitor.ctx.Done():
				fmt.Println("Отправка метрики приостановлена!")
				return
		}
	}
}

func (monitor *Monitor) Start () {
	go monitor.pool()
	go monitor.report()
}

func (monitor *Monitor) Stop () {
	monitor.cancel()
}

func(monitor Monitor) GetMemStats() map[string]gauge {
	result := make(map[string]gauge, 20)

	result["Alloc"] = gauge(monitor.memStats.Alloc)
	result["Frees"] = gauge(monitor.memStats.Frees)
	result["HeapAlloc"] = gauge(monitor.memStats.HeapAlloc)
	result["BuckHashSys"] = gauge(monitor.memStats.BuckHashSys)
	result["GCSys"] = gauge(monitor.memStats.GCSys)
	result["HeapIdle"] =gauge( monitor.memStats.HeapIdle)

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

func main() {
	ctx := context.Background()

	monitor := NewMonitor(ctx, time.Second * 2, time.Second * 10)

	monitor.Start()
	time.AfterFunc(time.Second * 20, monitor.Stop)

	timer := time.NewTimer(time.Second * 60)
	<-timer.C
}
