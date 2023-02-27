package config

import "time"

type Config struct {
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	Address        string        `env:"ADDRESS"`
}
