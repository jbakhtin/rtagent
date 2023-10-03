package ratelimiter

import (
	"reflect"
	"testing"
	"time"
)

func TestLimiter_Close(t *testing.T) {
	type fields struct {
		waiter       chan struct{}
		resetCounter *time.Ticker
		maxCount     int
		counter      int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Close rete limiter",
			fields{
				make(chan struct{}),
				time.NewTicker(time.Second * 2),
				10,
				0,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Limiter{
				waiter:       tt.fields.waiter,
				resetCounter: tt.fields.resetCounter,
				maxCount:     tt.fields.maxCount,
				counter:      tt.fields.counter,
			}
			if err := l.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLimiter_Run(t *testing.T) {
	type fields struct {
		waiter       chan struct{}
		resetCounter *time.Ticker
		maxCount     int
		counter      int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"Close rete limiter",
			fields{
				make(chan struct{}),
				time.NewTicker(time.Millisecond * 100),
				10,
				0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Limiter{
				waiter:       tt.fields.waiter,
				resetCounter: tt.fields.resetCounter,
				maxCount:     tt.fields.maxCount,
				counter:      tt.fields.counter,
			}

			l.Run()

			var counter int

			timer := time.NewTimer(time.Second)

			for {
				select {
				case <-l.Wait():
					counter++
				case <-timer.C:
					goto Exit
				}
			}
		Exit:

			if counter < 90 && counter > 110 {
				t.Errorf("Error %v", counter)
			}
		})
	}
}

func TestLimiter_Wait(t *testing.T) {
	type fields struct {
		waiter       chan struct{}
		resetCounter *time.Ticker
		maxCount     int
		counter      int
	}
	tests := []struct {
		name   string
		fields fields
		want   chan struct{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Limiter{
				waiter:       tt.fields.waiter,
				resetCounter: tt.fields.resetCounter,
				maxCount:     tt.fields.maxCount,
				counter:      tt.fields.counter,
			}
			if got := l.Wait(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wait() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		timeInterval time.Duration
		count        int
	}
	tests := []struct {
		name string
		args args
		want *Limiter
	}{
		{
			"Build new ratelimiter",
			args{
				time.Second,
				1,
			},
			&Limiter{
				make(chan struct{}),
				time.NewTicker(time.Second),
				1,
				0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.timeInterval, tt.args.count)

			if !reflect.DeepEqual(got.counter, tt.want.counter) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got.maxCount, tt.want.maxCount) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
