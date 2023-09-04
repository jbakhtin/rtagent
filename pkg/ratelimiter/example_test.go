package ratelimiter

import (
	"context"
	"fmt"
	"time"
)

func Example() {
	// Определим счетчик с временным интервалом в 10 секунд и количеством операций, за временной интервал, равное двум.
	limiter := New(time.Second*10, 2)
	_ = limiter.Run(context.TODO())

	// Цикл будет бесконечно выполняться, выводя в консоль "Draw" два раза каждые 10 секунд.
	for {
		fmt.Println("Draw")
		limiter.Wait()
	}
}