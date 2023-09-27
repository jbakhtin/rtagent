package limited

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/agent/tasker"
	"reflect"
	"sync"
	"testing"
	"time"
)

type Counter struct {
	Counter int
	sync.RWMutex
}

func (c *Counter) Increment(ctx context.Context) error {
	c.Lock()
	defer c.Unlock()
	c.Counter++
	return nil
}

func (c *Counter) Get() int {
	c.Lock()
	defer c.Unlock()
	return c.Counter
}

func TestNew(t *testing.T) {
	counter := Counter{}

	type args struct {
		name     string
		limit    int
		duration time.Duration
		f        tasker.Func
	}
	tests := []struct {
		name string
		args args
		want task
	}{
		{
			"New limited task",
			args{
				"increment counter not more then 10 for 10 seconds",
				10,
				time.Second * 10,
				counter.Increment,
			},
			task{
				counter.Increment,
				"increment counter not more then 10 for 10 seconds",
				time.Second * 10,
				10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.name, tt.args.limit, tt.args.duration, tt.args.f)

			if !reflect.DeepEqual(got.name, tt.want.name) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got.duration, tt.want.duration) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got.limit, tt.want.limit) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_task_Do(t1 *testing.T) {
	counter := Counter{}

	type fields struct {
		f        tasker.Func
		name     string
		duration time.Duration
		limit    int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Test limited task",
			fields{
				counter.Increment,
				"increment counter not more then 10 for 10 seconds",
				time.Second * 10,
				10,
			},
			true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &task{
				f:        tt.fields.f,
				name:     tt.fields.name,
				duration: tt.fields.duration,
				limit:    tt.fields.limit,
			}
			ctx, cancel := context.WithCancel(context.TODO())
			timer := time.NewTimer(time.Second * 1)
			go func() {
				if err := t.Do(ctx); (err != nil) != tt.wantErr {
					t1.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			<-timer.C
			cancel()

			if counter.Get() != t.limit {
				t1.Errorf("A different result was expected: %v, want: %v", counter.Counter, t.limit)
			}
		})
	}
}
