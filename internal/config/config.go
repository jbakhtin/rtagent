package config

import "time"

type Config struct {
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreFile      string        `env:"STORE_FILE" envDefault:"devops-metrics-db.json"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL" envDefault:"40s"`
	Restore        bool          `env:"RESTORE" envDefault:"true"`
}
