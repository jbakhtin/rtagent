package jobqueue

import (
	"github.com/jbakhtin/rtagent/internal/types"
	"sync"
)

type Job struct {
	value types.Metricer
	next  *Job
	key   string
}

func NewJob(key string, value types.Metricer) *Job {
	return &Job{
		key:   key,
		value: value,
		next:  nil,
	}
}

func (n Job) Key() string {
	return n.key
}

func (n Job) Value() types.Metricer {
	return n.value
}

type queue struct {
	head *Job
	tail *Job
	sync.RWMutex
	count int
}

func NewQueue() *queue {
	return &queue{}
}

func (q *queue) Size() int {
	q.Lock()
	defer q.Unlock()
	return q.count
}
func (q *queue) IsEmpty() bool {
	q.Lock()
	defer q.Unlock()
	return q.count == 0
}

func (q *queue) Enqueue(key string, metric types.Metricer) {
	q.Lock()
	defer q.Unlock()

	node := NewJob(key, metric)
	if q.head == nil {
		q.head = node
	} else {
		q.tail.next = node
	}

	q.count++
	q.tail = node
}

func (q *queue) Dequeue() *Job {
	q.Lock()
	defer q.Unlock()

	if q.head == nil {
		return nil
	}
	temp := q.head

	q.head = q.head.next
	q.count--
	if q.head == nil {
		q.tail = nil
	}
	return temp
}
