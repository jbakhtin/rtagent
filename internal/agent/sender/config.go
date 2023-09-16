package sender

import "time"

type Configer interface {
	GetReportInterval() time.Duration
	GetServerAddress() string
	GetKeyApp() string
}
