package storages

import (
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/storages/dbstorage"
	"github.com/jbakhtin/rtagent/internal/storages/memstorage"
	"github.com/jbakhtin/rtagent/internal/types"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func genRandString(length int) (string) {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func Benchmark(b *testing.B) {
	cfg, _ := config.NewConfigBuilder().WithAllFromEnv().Build()
	memStorage, _ := memstorage.NewMemStorage(cfg)
	dbStorage, _ := dbstorage.New(cfg)

	storageMapSize := 1000
	testData := make(map[string]models.Metricer, storageMapSize)

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
			for _ ,model := range testData {
				memStorage.Set(model)
			}
		}
	})

	b.Run("mem_get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for key := range testData {
				memStorage.Get(key)
			}
		}
	})

	b.Run("db_set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _ ,model := range testData {
				dbStorage.Set(model)
			}
		}
	})

	b.Run("db_get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for key := range testData {
				dbStorage.Get(key)
			}
		}
	})
}
