package agent

import "time"

type Configer interface {
	GetReportInterval() time.Duration
	GetPollInterval() time.Duration
	GetAcceptableCountAgentErrors() int
}
