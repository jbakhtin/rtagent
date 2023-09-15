package aggregator

import (
	"github.com/jbakhtin/rtagent/internal/types"
	gopsutil "github.com/shirou/gopsutil/v3/mem"
	"runtime"
)

// Runtime - Поулчить слайс содержщий последние акутальные данные
func Runtime() (map[string]types.Metricer, error) {
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	result := map[string]types.Metricer{}

	// memStats
	result["Alloc"] = types.Gauge(memStats.Alloc)
	result["Frees"] = types.Gauge(memStats.Frees)
	result["HeapAlloc"] = types.Gauge(memStats.HeapAlloc)
	result["BuckHashSys"] = types.Gauge(memStats.BuckHashSys)
	result["GCSys"] = types.Gauge(memStats.GCSys)
	result["GCCPUFraction"] = types.Gauge(memStats.GCCPUFraction)
	result["HeapIdle"] = types.Gauge(memStats.HeapIdle)
	result["HeapInuse"] = types.Gauge(memStats.HeapInuse)
	result["HeapObjects"] = types.Gauge(memStats.HeapObjects)
	result["HeapReleased"] = types.Gauge(memStats.HeapReleased)
	result["HeapSys"] = types.Gauge(memStats.HeapSys)
	result["LastGC"] = types.Gauge(memStats.LastGC)
	result["Lookups"] = types.Gauge(memStats.Lookups)
	result["MCacheInuse"] = types.Gauge(memStats.MCacheInuse)
	result["MCacheSys"] = types.Gauge(memStats.MCacheSys)
	result["MSpanInuse"] = types.Gauge(memStats.MSpanInuse)
	result["MSpanSys"] = types.Gauge(memStats.MSpanSys)
	result["Mallocs"] = types.Gauge(memStats.Mallocs)
	result["NextGC"] = types.Gauge(memStats.NextGC)
	result["NumForcedGC"] = types.Gauge(memStats.NumForcedGC)
	result["NumGC"] = types.Gauge(memStats.NumGC)
	result["OtherSys"] = types.Gauge(memStats.OtherSys)
	result["PauseTotalNs"] = types.Gauge(memStats.PauseTotalNs)
	result["StackInuse"] = types.Gauge(memStats.StackInuse)
	result["StackSys"] = types.Gauge(memStats.StackSys)
	result["Sys"] = types.Gauge(memStats.Sys)
	result["TotalAlloc"] = types.Gauge(memStats.TotalAlloc)

	// Custom stats
	//result["PollCount"] = m.pollCounter
	//result["RandomValue"] = types.Gauge(rand.Int())

	return result, nil
}

func Gopsutil() (map[string]types.Metricer, error) {
	memStats, err := gopsutil.VirtualMemory()
	if err != nil {
		return nil, err
	}

	result := map[string]types.Metricer{}

	// memStats
	result["TotalMemory"] = types.Gauge(memStats.Total)
	result["FreeMemory"] = types.Gauge(memStats.Free)
	result["CPUutilization1"] = types.Gauge(memStats.Used)

	return result, nil
}
