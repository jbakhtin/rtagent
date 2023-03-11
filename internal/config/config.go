package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	_pollInterval               = time.Second * 2
	_reportInterval             = time.Second * 10
	_address                    = "127.0.0.1:8080"
	_storeFile                  = "tmp/devops-metrics-db.json"
	_storeInterval              = time.Second * 300
	_restore                    = true
	_acceptableCountAgentErrors = 10
)

const (
	_pollIntervalLabel               = "Период чтения состояния метрик"
	_reportIntervalLabel             = "Период отправки состояния метрик на сервер"
	_addressLabel                    = "Адрес сервера"
	_storeFileLabel                  = "Файл хранения слепков состояния MemStorage"
	_storeIntervalLabel              = "Период создания слепков MemStorage в секундах"
	_restoreLabel                    = "Загрузить последний слепок MemStorage перед стартом сервиса"
	_acceptableCountAgentErrorsLabel = "Допустимое количество ошибок от агента"
)

type Config struct {
	PollInterval               time.Duration `env:"POLL_INTERVAL"`
	ReportInterval             time.Duration `env:"REPORT_INTERVAL"`
	Address                    string        `env:"ADDRESS"`
	StoreFile                  string        `env:"STORE_FILE"`
	StoreInterval              time.Duration `env:"STORE_INTERVAL"`
	Restore                    bool          `env:"RESTORE"`
	AcceptableCountAgentErrors int           `env:"ACCEPTABLE_COUNT_AGENT_ERRORS"`
}

type Builder struct {
	config Config
	err    error
}

func NewConfigBuilder() *Builder {
	return &Builder{
		Config{
			_pollInterval,
			_reportInterval,
			_address,
			_storeFile,
			_storeInterval,
			_restore,
			_acceptableCountAgentErrors,
		},
		nil,
	}
}

func (cb *Builder) WithAllFromFlagsS() *Builder {
	address := flag.String("a", _address, _addressLabel)
	storeFile := flag.String("f", _storeFile, _storeFileLabel)
	storeInterval := flag.Duration("i", _storeInterval, _storeIntervalLabel)
	restore := flag.Bool("r", _restore, _restoreLabel)
	flag.Parse()

	cb.config.Address = *address
	cb.config.StoreFile = *storeFile
	cb.config.StoreInterval = *storeInterval
	cb.config.Restore = *restore

	return cb
}

// WithAllFromFlagsA
// TODO: назвал методы для инициализации конфигов аггента и сервера по разному,
// но есть желание вынести конфиги в разные файлы. Как лучше?
// Лучше иметь два метода, один для чтения из env, второй из флагов, или вынести чтение окружения в main,
// а конфиг инициализировать примитивнымими методами, например метод WithReportInterval()
func (cb *Builder) WithAllFromFlagsA() *Builder {
	pollInterval := flag.Duration("p", _pollInterval, _pollIntervalLabel)
	reportInterval := flag.Duration("r", _reportInterval, _reportIntervalLabel)
	address := flag.String("a", _address, _addressLabel)
	acceptableCountAgentErrors := flag.Int("e", _acceptableCountAgentErrors, _acceptableCountAgentErrorsLabel)
	flag.Parse()

	cb.config.PollInterval = *pollInterval
	cb.config.ReportInterval = *reportInterval
	cb.config.Address = *address
	cb.config.AcceptableCountAgentErrors = *acceptableCountAgentErrors

	return cb
}

func (cb *Builder) WithAllFromEnv() *Builder {
	err := env.Parse(&cb.config)
	if err != nil {
		cb.err = err
	}

	return cb
}

func (cb *Builder) Build() (Config, error) {
	return cb.config, cb.err
}
