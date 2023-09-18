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
	_acceptableCountAgentErrors = 100
	_keyApp                     = ""
	_databaseDSN                = ""
	_databaseDriver             = "pgx"
	_rateLimit                  = 100
	_shutdownTimeout                  = 10
)

const (
	_pollIntervalLabel               = "Период чтения состояния метрик"
	_reportIntervalLabel             = "Период отправки состояния метрик на сервер"
	_addressLabel                    = "Адрес сервера"
	_storeFileLabel                  = "Файл хранения слепков состояния MemStorage"
	_storeIntervalLabel              = "Период создания слепков MemStorage в секундах"
	_restoreLabel                    = "Загрузить последний слепок MemStorage перед стартом сервиса"
	_acceptableCountAgentErrorsLabel = "Допустимое количество ошибок от агента"
	_keyAppLabel                     = "Ключ приложения"
	_databaseDSNLabel                = "DSN БД"
	_databaseDriverLabel             = "Драйвер подключения к БД"
	_rateLimitLabel                  = "Количество исходящих запросов в секунду"
	_shutdownTimeoutLabel                  = "Время на заерщение всех процессов перед отключением"
)

type Config struct {
	Address                    string        `env:"ADDRESS"`
	StoreFile                  string        `env:"STORE_FILE"`
	KeyApp                     string        `env:"KEY"`
	DatabaseDSN                string        `env:"DATABASE_DSN"`
	DatabaseDriver             string        `env:"DATABASE_DRIVER" envDefault:"pgx"`
	AcceptableCountAgentErrors int           `env:"ACCEPTABLE_COUNT_AGENT_ERRORS"`
	RateLimit                  int           `env:"RATE_LIMIT" envDefault:"10"`
	PollInterval               time.Duration `env:"POLL_INTERVAL"`
	ReportInterval             time.Duration `env:"REPORT_INTERVAL"`
	StoreInterval              time.Duration `env:"STORE_INTERVAL"`
	Restore                    bool          `env:"RESTORE"`
	ShutdownTimeout                    time.Duration          `env:"SHUTDOWN_TIMEOUT"`
}

type Builder struct {
	err    error
	config Config
}

func NewConfigBuilder() *Builder {
	return &Builder{
		nil,
		Config{
			_address,
			_storeFile,
			_keyApp,
			_databaseDSN,
			_databaseDriver,
			_acceptableCountAgentErrors,
			_rateLimit,
			_pollInterval,
			_reportInterval,
			_storeInterval,
			_restore,
			_shutdownTimeout,
		},
	}
}

func (cb *Builder) WithAllFromFlagsS() *Builder {
	address := flag.String("a", _address, _addressLabel)
	storeFile := flag.String("f", _storeFile, _storeFileLabel)
	storeInterval := flag.Duration("i", _storeInterval, _storeIntervalLabel)
	restore := flag.Bool("r", _restore, _restoreLabel)
	keyApp := flag.String("k", _keyApp, _keyAppLabel)
	databaseDSN := flag.String("d", _databaseDSN, _databaseDSNLabel)
	databaseDriver := flag.String("dbDriver", _databaseDriver, _databaseDriverLabel)
	shutdownTimeout := flag.Duration("shutdownTimeout", _shutdownTimeout, _shutdownTimeoutLabel)
	flag.Parse()

	cb.config.Address = *address
	cb.config.StoreFile = *storeFile
	cb.config.StoreInterval = *storeInterval
	cb.config.Restore = *restore
	cb.config.KeyApp = *keyApp
	cb.config.DatabaseDSN = *databaseDSN
	cb.config.DatabaseDriver = *databaseDriver
	cb.config.ShutdownTimeout = *shutdownTimeout

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
	keyApp := flag.String("k", _keyApp, _keyAppLabel)
	rateLimit := flag.Int("l", _rateLimit, _rateLimitLabel)
	shutdownTimeout := flag.Duration("shutdownTimeout", _shutdownTimeout, _shutdownTimeoutLabel)
	flag.Parse()

	cb.config.PollInterval = *pollInterval
	cb.config.ReportInterval = *reportInterval
	cb.config.Address = *address
	cb.config.AcceptableCountAgentErrors = *acceptableCountAgentErrors
	cb.config.KeyApp = *keyApp
	cb.config.RateLimit = *rateLimit
	cb.config.ShutdownTimeout = *shutdownTimeout

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

func (c Config) GetPollInterval() time.Duration  {
	return c.PollInterval
}

func (c Config) GetServerAddress() string  {
	return c.Address
}

func (c Config) GetReportInterval() time.Duration  {
	return c.ReportInterval
}

func (c Config) GetKeyApp() string  {
	return c.KeyApp
}

func (c Config) GetAcceptableCountAgentErrors() int  {
	return c.AcceptableCountAgentErrors
}

func (c Config) GetShutdownTimeout() time.Duration  {
	return c.ShutdownTimeout
}