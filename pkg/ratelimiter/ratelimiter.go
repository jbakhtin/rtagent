// Package ratelimiter - это простая реализация ограничителя скорости, которая ограничивает
// количество разрешенных действий в течение определенного периода времени.
package ratelimiter

import (
	"time"
)

type Limiter struct {
	waiter       chan struct{} // Канал для увеличения счетчика после завершения действия.
	resetCounter *time.Ticker  // Тикер для сброса счетчика.
	maxCount     int           // Максимальное количество разрешенных действий.
	counter      int           // Текущее количество выполненных действий.
}

func New(timeInterval time.Duration, count int) *Limiter {
	l := &Limiter{
		maxCount:     count,
		counter:      0,
		resetCounter: time.NewTicker(timeInterval),
		waiter:       make(chan struct{}),
	}

	return l
}

// Wait ожидает окончания действия.
// Устанавливается в конце повторяющейся операции, количество выполнений которой нужно ограничить.
func (l *Limiter) Wait() chan struct{} {
	return l.waiter
}

// Run запускает внутренний цикл счетчика.
func (l *Limiter) Run() {
	go l.run()
}

// run обрабатывает состояния счетчика.
func (l *Limiter) run() {
	for {
		if l.counter >= l.maxCount {
			<-l.resetCounter.C
			l.counter = 0
		}

		l.waiter <- struct{}{}
		l.counter++
	}
}

func (l *Limiter) Close() error {
	l.resetCounter.Stop()
	close(l.waiter)
	return nil
}
