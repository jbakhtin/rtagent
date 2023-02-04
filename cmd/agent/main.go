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

	ctx    context.Context
	cancel context.CancelFunc
}

func (monitor *Monitor) pool() {
	ticker := time.NewTicker(monitor.pollInterval)

	for {
		select {
			case <-ticker.C:
				fmt.Println("Метрика собрана")
			case <-monitor.ctx.Done():
				fmt.Println("Сбор метрик приостановлен!")
				return
		}
	}
}

func (monitor *Monitor) report() {
	ticker := time.NewTicker(monitor.reportInterval)

	for {
		select {
			case <-ticker.C:
				fmt.Println("Метрика отрпавлена!")
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

func (monitor *Monitor) Done () {
	monitor.cancel()
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

//TODO: сделать остановшик
//TODO: сделать остановшик через определенное время
//TODO: сделать остановшик на опредленоое время

func main() {
	ctx := context.Background()

	monitor := NewMonitor(ctx, time.Second * 2, time.Second * 10)

	monitor.Start()
	time.AfterFunc(time.Second * 20, monitor.Done)

	timer := time.NewTimer(time.Second * 60)
	<-timer.C
}
