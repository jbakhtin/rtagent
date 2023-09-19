package jobqueue

import (
	"github.com/jbakhtin/rtagent/internal/types"
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
	head * QNode
	tail * QNode
	count int
}
func GetMyQueue() * MyQueue {
	// return new MyQueue
	return &MyQueue {
		nil,
		nil,
		0,
	}
}

func(this MyQueue) Size() int {
	return this.count
}
func(this MyQueue) IsEmpty() bool {
	return this.count == 0
}
// Add new node of queue
func(this *MyQueue) Enqueue(key string, metric types.Metricer) {
	var node * QNode = GetQNode(key, metric)
	if this.head == nil {
		// Add first element into queue
		this.head = node
	} else {
		// Add node at the end using tail
		this.tail.next = node
	}
	this.count++
	this.tail = node
}
// Delete a element into queue
func(this *MyQueue) Dequeue() *QNode {
	if this.head == nil {
		return nil
	}
	// Pointer variable which are storing
	// the address of deleted node
	var temp * QNode = this.head
	// Visit next node
	this.head = this.head.next
	this.count--
	if this.head == nil {
		// When deleting a last node of linked list
		this.tail = nil
	}
	return temp
}

