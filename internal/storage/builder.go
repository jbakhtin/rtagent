package storage

import (
	"database/sql"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/storage/dbstorage"
	"github.com/jbakhtin/rtagent/internal/storage/filestorage"
	"github.com/jbakhtin/rtagent/internal/storage/memstorage"
	"github.com/pressly/goose/v3"
	"sync"
)

// MetricRepository интерфейс реализации хранилища.
type MetricRepository interface {
	GetAll() (map[string]models.Metricer, error)
	Get(key string) (models.Metricer, error)
	Set(models.Metricer) (models.Metricer, error)
	SetBatch([]models.Metricer) ([]models.Metricer, error)
	TestPing() error
}

type builder struct {
	repository MetricRepository
	err error
}

func New() *builder {
	return &builder{}
}

type postgresConfig interface {
	GetDatabaseDriver() string
	GetDatabaseDSN() string
}

func (b *builder) Postgres(cfg postgresConfig) *builder {
	db, err := sql.Open(cfg.GetDatabaseDriver(), cfg.GetDatabaseDSN())
	if err != nil {
		b.err = err
		return b
	}

	goose.SetBaseFS(dbstorage.EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		b.err = err
		return b
	}

	if err := goose.Up(db, "migrations"); err != nil {
		b.err = err
		return b
	}

	b.repository = &dbstorage.DBStorage{
		DB:     db,
	}

	return b
}

type inMemoryConfig interface {}

func (b *builder) InMemory(cfg inMemoryConfig) *builder {
	b.repository = &memstorage.MemStorage{
		Items:  make(map[string]models.Metricer, 0),
		Mx:     &sync.RWMutex{},
	}

	return b
}

type fileConfig interface {
	inMemoryConfig
}

func (b *builder) File(cfg fileConfig) *builder {
	b.repository = &filestorage.FileStorage{
		MemStorage: memstorage.MemStorage{
			Items: make(map[string]models.Metricer, 0),
			Mx:    &sync.RWMutex{},
		},
	}
	return b
}

func (b *builder) Build() (MetricRepository, error) {
	if b.err != nil {
		return nil, b.err
	}

	return b.repository, nil
}