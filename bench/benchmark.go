package bench

import (
	"fmt"
	"sort"
	"time"
)

type BenchmarkResult struct {
	Min time.Duration
	Avg time.Duration
	Max time.Duration
	P95 time.Duration
}

type Benchmark interface {
	Open() error
	BenchmarkOnce() (time.Duration, error)
	Close()
}

func RunBenchmark(fn func() (time.Duration, error), iterations int) (BenchmarkResult, error) {
	durations := make([]time.Duration, iterations)
	benchmarkResult := BenchmarkResult{}
	var avg float64

	for i := 0; i < iterations; i++ {
		var err error
		durations[i], err = fn()

		if err != nil {
			return benchmarkResult, fmt.Errorf("RunBenchmark() - Failed to run benchmark at iteration [%d]; error: %v", i, err)
		}

		if i == 0 {
			benchmarkResult.Min = durations[i]
			benchmarkResult.Avg = durations[i]
			benchmarkResult.Max = durations[i]
		} else {
			if durations[i] < benchmarkResult.Min {
				benchmarkResult.Min = durations[i]
			}
			if durations[i] > benchmarkResult.Max {
				benchmarkResult.Max = durations[i]
			}
			avg += float64(benchmarkResult.Avg) / float64(iterations)
		}
	}

	benchmarkResult.Avg = time.Duration(avg)
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})
	benchmarkResult.P95 = durations[int(0.95*float64(len(durations)))]

	return benchmarkResult, nil
}
