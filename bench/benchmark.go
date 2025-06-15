package bench

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type BenchmarkDurations struct {
	durations []time.Duration
}

type BenchmarkAggregatedResult struct {
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

func (brr BenchmarkDurations) AggregateBenchmarkResult() (BenchmarkAggregatedResult, error) {
	benchmarkAggResult := BenchmarkAggregatedResult{}

	var avg float64
	for i, duration := range brr.durations {
		if i == 0 {
			benchmarkAggResult.Min = duration
			benchmarkAggResult.Avg = duration
			benchmarkAggResult.Max = duration
		} else {
			if duration < benchmarkAggResult.Min {
				benchmarkAggResult.Min = duration
			}
			if duration > benchmarkAggResult.Max {
				benchmarkAggResult.Max = duration
			}
			avg += float64(benchmarkAggResult.Avg) / float64(len(brr.durations))
		}
	}

	benchmarkAggResult.Avg = time.Duration(avg)
	sort.Slice(brr.durations, func(i, j int) bool {
		return brr.durations[i] < brr.durations[j]
	})
	benchmarkAggResult.P95 = brr.durations[int(0.95*float64(len(brr.durations)))]

	return benchmarkAggResult, nil
}

func RunBenchmarkParallel(fn func() (time.Duration, error), iterations int, parallelism int) (BenchmarkResult, error) {
	iterationsPerThread := int(float64(iterations) / float64(parallelism))
	if iterationsPerThread <= 0 {
		return BenchmarkResult{},
			fmt.Errorf("RunBenchmarkParallel failed - iterationsPerThread must be > 0, but is %d. [iterations=%d] [parallelism=%d]",
				iterationsPerThread, iterations, parallelism)
	}

	var wg sync.WaitGroup
	results := make(chan BenchmarkResult, parallelism)
	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go func() {
			benchmarkResult, err := RunBenchmark(fn, iterationsPerThread)
			if err != nil {
				fmt.Printf("Error ocurred in RunBenchmarkParallel: %v ", err)
			}
			results <- benchmarkResult
		}()
	}
	wg.Wait()
	close(results)

	for result := range results {
		fmt.Printf("Result: %v", result)
	}

	return BenchmarkResult{}, nil
}

func RunBenchmark(fn func() (time.Duration, error), iterations int) (BenchmarkDurations, error) {
	durations := make([]time.Duration, iterations)
	benchmarkDurations := BenchmarkDurations{}

	for i := 0; i < iterations; i++ {
		var err error
		durations[i], err = fn()

		if err != nil {
			return benchmarkDurations, fmt.Errorf("RunBenchmark() - Failed to run benchmark at iteration [%d]; error: %v", i, err)
		}
	}

	return benchmarkDurations, nil
}
