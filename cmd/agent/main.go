package main

import (
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

	stop chan bool
}

func (monitor *Monitor) pool() {
	ticker := time.NewTicker(monitor.pollInterval)

	for {
		select {
			case <-ticker.C:
				fmt.Println("Метрика собрана")
			case <-monitor.stop:
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
			case <-monitor.stop:
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
	monitor.stop <- true
}

func NewMonitor(pollInterval, reportInterval time.Duration) Monitor {
	return Monitor{
		pollInterval: pollInterval,
		reportInterval: reportInterval,
		stop: make(chan bool, 0),
	}
}

//TODO: сделать остановшик
//TODO: сделать остановшик через определенное время
//TODO: сделать остановшик на опредленоое время

func main() {
	monitor := NewMonitor(time.Second * 2, time.Second * 10)
	monitor.Start()

	//monitor.Stop()

	timer := time.NewTimer(time.Second * 15)
	timer2 := time.NewTimer(time.Second * 20)

	select {
		case <-timer.C:
			fmt.Println("Сбор данных приостановлен!")
			monitor.Stop()
			monitor.Stop()
	}

	<-timer2.C
}
