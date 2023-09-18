package agent

import "github.com/jbakhtin/rtagent/internal/agent/workerpool"

type builder struct {
	agent agent
	err error
}

func New() *builder {
	return &builder{}
}

func (b *builder) WithConfig(cfg Configer) *builder {
	b.agent.cfg = cfg
	return b
}

func (b *builder) WithAggregator(aggregator Aggregator) *builder {
	b.agent.aggregator = aggregator
	return b
}

func (b *builder) WithSender(sender Sender) *builder {
	b.agent.sender = sender
	return b
}

func (b *builder) WithSoftShuttingDown() *builder {
	b.agent.softShuttingDown = true
	return b
}

func (b *builder) Build() (*agent, error) {
	b.agent.workerPool, b.err = workerpool.NewWorkerPool()

	if b.err != nil {
		return nil, b.err
	}
	return &b.agent, nil
}
