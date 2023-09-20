package storages

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/storages/dbstorage"
	"github.com/jbakhtin/rtagent/internal/storages/memstorage"
	"github.com/jbakhtin/rtagent/internal/types"
)

func genRandString(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

// Benchmark - сравнивает производительность хранилищ.
func Benchmark(b *testing.B) {
	// Инициализация хранилищ
	cfg, _ := config.NewConfigBuilder().WithAllFromEnv().Build()
	memStorage, err := memstorage.NewMemStorage(cfg)
	if err != nil {
		b.Fatal(err)
	}

	dbStorage, err := dbstorage.New(cfg)
	if err != nil {
		b.Fatal(err)
	}

	storageMapSize := 1000
	testData := make(map[string]models.Metricer, storageMapSize)

	// Создаем тестовые данные.
	for i := 0; i < storageMapSize; i++ {
		gauge := types.Gauge(12)
		model := models.Gauge{
			Description: models.Description{
				MKey:  genRandString(8),
				MType: gauge.Type(),
			},
			MValue: gauge,
		}

		testData[model.MKey] = model
	}

	b.ResetTimer()

	b.Run("mem_set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, model := range testData {
				_, err := memStorage.Set(model)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("mem_get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for key := range testData {
				_, err := memStorage.Get(key)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("db_set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, model := range testData {
				_, err := dbStorage.Set(model)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("db_get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for key := range testData {
				_, err := dbStorage.Get(key)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}
