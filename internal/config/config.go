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

type Builder struct {
	config   Config
	err error
}

func NewConfigBuilder() *Builder {
	return &Builder{
		Config{},
		nil,
	}
}

func (cb *Builder) WithPollInterval(pollInterval time.Duration) *Builder{
	cb.config.PollInterval = pollInterval

	return cb
}

func (cb *Builder) WithPollIntervalFromFlag() *Builder{
	cb.config.PollInterval = *flag.Duration("p", 2, "test")

	return cb
}


func (cb *Builder) WithReportInterval(reportInterval time.Duration) *Builder {
	cb.config.ReportInterval = reportInterval

	return cb
}

func (cb *Builder) WithReportIntervalFromFlag() *Builder {
	cb.config.ReportInterval = *flag.Duration("r", 10, "test")

	return cb
}


func (cb *Builder) WithAddress(address string) *Builder {
	cb.config.Address = address

	return cb
}

func (cb *Builder) WithAddressFromFlag() *Builder {
	cb.config.Address = *flag.String("a", "localhost:8080", "test")

	return cb
}


func (cb *Builder) WithStoreFile(storeFile string) *Builder {
	cb.config.StoreFile = storeFile

	return cb
}

func (cb *Builder) WithStoreFileFromFlag() *Builder {
	cb.config.StoreFile = *flag.String("f", "tmp/devops-metrics-db.json", "test")

	return cb
}


func (cb *Builder) WithStoreInterval(storeInterval time.Duration) *Builder {
	cb.config.StoreInterval = storeInterval

	return cb
}

func (cb *Builder) WithStoreIntervalFromFlag() *Builder {
	cb.config.StoreInterval = *flag.Duration("i", 300, "test")

	return cb
}


func (cb *Builder) WithRestore(restore bool) *Builder {
	cb.config.Restore = restore

	return cb
}

func (cb *Builder) WithRestoreFromFlag() *Builder {
	cb.config.Restore = *flag.Bool("r", true, "test")

	return cb
}

// ---

func (cb *Builder) WithAllFromEnv() *Builder {
	err := env.Parse(&cb.config)
	if err != nil {
		cb.err = err
	}

	return cb
}

func (cb *Builder) Build() (Config, error) {
	return  cb.config, cb.err
}
