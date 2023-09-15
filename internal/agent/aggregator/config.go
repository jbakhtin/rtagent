package aggregator

import "time"

type Config interface{
	GetPollInterval() time.Duration
}