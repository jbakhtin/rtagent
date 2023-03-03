package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"time"
)

type Config struct {
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	Address        string        `env:"ADDRESS"`
	StoreFile      string        `env:"STORE_FILE"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL"`
	Restore        bool          `env:"RESTORE"`
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
	cb.config.PollInterval = *flag.Duration("p", time.Second * 2, "test")
	flag.Parse()

	return cb
}


func (cb *Builder) WithReportInterval(reportInterval time.Duration) *Builder {
	cb.config.ReportInterval = reportInterval

	return cb
}

func (cb *Builder) WithReportIntervalFromFlag() *Builder {
	cb.config.ReportInterval = *flag.Duration("r", time.Second * 10, "test")
	flag.Parse()

	return cb
}


func (cb *Builder) WithAddress(address string) *Builder {
	cb.config.Address = address

	return cb
}

func (cb *Builder) WithAddressFromFlag() *Builder {
	cb.config.Address = *flag.String("a", "localhost:8080", "test")
	flag.Parse()

	return cb
}


func (cb *Builder) WithStoreFile(storeFile string) *Builder {
	cb.config.StoreFile = storeFile

	return cb
}

func (cb *Builder) WithStoreFileFromFlag() *Builder {
	cb.config.StoreFile = *flag.String("f", "tmp/devops-metrics-db.json", "test")
	flag.Parse()

	return cb
}


func (cb *Builder) WithStoreInterval(storeInterval time.Duration) *Builder {
	cb.config.StoreInterval = storeInterval

	return cb
}

func (cb *Builder) WithStoreIntervalFromFlag() *Builder {
	cb.config.StoreInterval = *flag.Duration("i", time.Second * 300, "test")
	flag.Parse()

	return cb
}


func (cb *Builder) WithRestore(restore bool) *Builder {
	cb.config.Restore = restore

	return cb
}

func (cb *Builder) WithRestoreFromFlag() *Builder {
	cb.config.Restore = *flag.Bool("r", true, "test")
	flag.Parse()

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
