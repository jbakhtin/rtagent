package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

const (
	pollInterval = time.Second * 2
	reportInterval = time.Second * 10
)

// TODO: реализовать кастомный тип метрики PollCount - counter
// TODO: реализовать кастомный тип метрики RandomValue - gauge

type Monitor struct {
	pollInterval time.Duration
	reportInterval time.Duration

	memStats runtime.MemStats

	ctx context.Context
}

func (monitor *Monitor) pool(ctx context.Context) {
	ticker := time.NewTicker(monitor.pollInterval)

	for {
		select {
			case <-ticker.C:
				fmt.Println("Метрика собрана")
			case <-ctx.Done():
				fmt.Println("Сбор метрик приостановлен!")
				return
		}
	}
}

func (monitor *Monitor) report(ctx context.Context) {
	ticker := time.NewTicker(monitor.reportInterval)

	for {
		select {
			case <-ticker.C:
				fmt.Println("Метрика отрпавлена!")
			case <-ctx.Done():
				fmt.Println("Отправка метрики приостановлена!")
				return
		}
	}
}

func (monitor *Monitor) Start (ctx context.Context) {
	go monitor.pool(ctx)
	go monitor.report(ctx)
}

func NewMonitor(pollInterval, reportInterval time.Duration) Monitor {
	return Monitor{
		pollInterval: pollInterval,
		reportInterval: reportInterval,
	}
}

//TODO: сделать остановшик
//TODO: сделать остановшик через определенное время
//TODO: сделать остановшик на опредленоое время

func main() {
	ctx := context.Background()
	ctx, cencel := context.WithCancel(ctx)

	monitor := NewMonitor(time.Second * 2, time.Second * 10)
	monitor.Start(ctx)

	time.AfterFunc(time.Second * 20, cencel)

	timer := time.NewTimer(time.Second * 60)
	<-timer.C
}
