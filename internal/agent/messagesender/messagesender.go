package messagesender

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/agent/jobqueue"
	"github.com/jbakhtin/rtagent/internal/types"
)

type Sender interface {
	Send(key string, value types.Metricer) error
}

type Jober interface {
	Dequeue() *jobqueue.QNode
	IsEmpty() bool
}

type MessageSender struct {
	Sender Sender
	Jober Jober
}

func (jm *MessageSender) Do(ctx context.Context) error {
	if jm.Jober.IsEmpty() {
		return nil
	}

	node := jm.Jober.Dequeue()
	jm.Sender.Send(node.Key(), node.Value())

	return nil
}
