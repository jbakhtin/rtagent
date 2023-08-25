// Package ratelimiter - это простая реализация ограничителя скорости, которая ограничивает
// количество разрешенных действий в течение определенного периода времени.
package ratelimiter

import (
	"context"
	"time"
)

type Limiter struct {
	maxCount     int           // Максимальное количество разрешенных действий.
	counter      int           // Текущее количество выполненных действий.
	resetCounter *time.Ticker  // Тикер для сброса счетчика.
	waiter       chan struct{} // Канал для увеличения счетчика после завершения действия.
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
func (l *Limiter) Wait() {
	<-l.waiter
}

// Run запускает внутренний цикл счетчика.
func (l *Limiter) Run(ctx context.Context) error {
	go l.run(ctx)

	return nil
}

// run обрабатывает состояния счетчика.
func (l *Limiter) run(ctx context.Context) {
	for {
		if l.counter > l.maxCount {
			<-l.resetCounter.C
			l.counter = 0
		}

		select {
		case <-ctx.Done():
			l.resetCounter.Stop()
			return
		case l.waiter <- struct{}{}:
			l.counter++

		case <-l.resetCounter.C:
			l.counter = 0
		}
	}
}
