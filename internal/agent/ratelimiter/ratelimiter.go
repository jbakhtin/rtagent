package ratelimiter

import "time"

type Limiter struct {
	maxCount     int
	counter      int
	resetCounter *time.Ticker
	waiter       chan struct{}
}

func (l *Limiter) run() {
	for {
		if l.counter > l.maxCount { // Если количество операций во временном интервале израсходовано, то
			<-l.resetCounter.C // дожидаемся завершения временного интервала
			l.counter = 0      // обновляем счетчик
		}

		select {
		case l.waiter <- struct{}{}: // На каждое выполненное действие
			// уменьшаем счетчик допустимых действий в заданном интервале
			l.counter++

		case <-l.resetCounter.C: // если не успеваем выполнить заданное количество действий в заданном интервале, то
			// обнуляем счетчик
			l.counter = 0
		}
	}
}

func (l *Limiter) Wait() {
	<-l.waiter
}

func NewLimiter(timeInterval time.Duration, count int) *Limiter {
	l := &Limiter{
		maxCount:     count,
		counter:      0,
		resetCounter: time.NewTicker(timeInterval),
		waiter:       make(chan struct{}),
	}
	go l.run()

	return l
}
