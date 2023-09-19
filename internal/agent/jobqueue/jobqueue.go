package jobqueue

import (
	"github.com/jbakhtin/rtagent/internal/types"
	"sync"
)

type QNode struct {
	key string
	value types.Metricer
	next * QNode
}
func GetQNode(key string, value types.Metricer) * QNode {
	// return new QNode
	return &QNode {
		key,
		value,
		nil,
	}
}

func (n QNode) Key() string {
	return n.key
}

func (n QNode) Value() types.Metricer {
	return n.value
}

type MyQueue struct {
	sync.RWMutex
	head * QNode
	tail * QNode
	count int
}
func GetMyQueue() * MyQueue {
	// return new MyQueue
	return &MyQueue {}
}

func(q *MyQueue) Size() int {
	q.Lock()
	defer q.Unlock()
	return q.count
}
func(q *MyQueue) IsEmpty() bool {
	q.Lock()
	defer q.Unlock()
	return q.count == 0
}
// Add new node of queue
func(q *MyQueue) Enqueue(key string, metric types.Metricer) {
	q.Lock()
	defer q.Unlock()
	var node * QNode = GetQNode(key, metric)
	if q.head == nil {
		// Add first element into queue
		q.head = node
	} else {
		// Add node at the end using tail
		q.tail.next = node
	}
	q.count++
	q.tail = node
}
// Delete a element into queue
func(q *MyQueue) Dequeue() *QNode {
	q.Lock()
	defer q.Unlock()
	if q.head == nil {
		return nil
	}
	// Pointer variable which are storing
	// the address of deleted node
	var temp * QNode = q.head
	// Visit next node
	q.head = q.head.next
	q.count--
	if q.head == nil {
		// When deleting a last node of linked list
		q.tail = nil
	}
	return temp
}

