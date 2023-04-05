package ratelimiter

import "time"

type Limiter struct {
	maxCount int
	count    int
	ticker   *time.Ticker
	ch       chan struct{}
}

func (l *Limiter) run() {
	for {
		// if counter has reached 0: block until next tick
		if l.count <= 0 {
			<-l.ticker.C
			l.count = l.maxCount
		}

		// otherwise:
		// decrement 'count' each time a message is sent on channel,
		// reset 'count' to 'maxCount' when ticker says so
		select {
		case l.ch <- struct{}{}:
			l.count--

		case <-l.ticker.C:
			l.count = l.maxCount
		}
	}
}

func (l *Limiter) Wait() {
	<-l.ch
}

func NewLimiter(d time.Duration, count int) *Limiter {
	l := &Limiter{
		maxCount: count,
		count:    count,
		ticker:   time.NewTicker(d),
		ch:       make(chan struct{}),
	}
	go l.run()

	return l
}
