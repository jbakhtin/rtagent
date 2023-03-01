package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"time"
)

type Config struct {
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreFile      string        `env:"STORE_FILE" envDefault:"tmp/devops-metrics-db.json"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL" envDefault:"40s"`
	Restore        bool          `env:"RESTORE" envDefault:"true"`
}

func (c *Config) InitFromENV() error {
	err := env.Parse(c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) InitFromFlag() error {
	c.PollInterval = *flag.Duration("pollInterval", 2, "a Duration")
	c.ReportInterval = *flag.Duration("reportInterval", 10, "a Duration")
	c.StoreInterval = *flag.Duration("storeInterval", 300, "a Duration")
	c.StoreFile = *flag.String("storeFile", "tmp/devops-metrics-db.json", "a String")
	c.Address = *flag.String("dddress", "127.0.0.1:8080", "a String")
	c.Restore = *flag.Bool("restore", true, "a Bool")

	return nil
}