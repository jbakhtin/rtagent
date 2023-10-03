package aggregator

import "github.com/jbakhtin/rtagent/internal/types"

type PoolCounter struct {
	Count int
}

func (pc *PoolCounter) PoolCount() (map[string]types.Metricer, error) {
	pc.Count++
	return map[string]types.Metricer{"PollCount": types.Counter(pc.Count)}, nil
}
