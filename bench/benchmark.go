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
	Setup() error
	RunOnce() (time.Duration, error)
	Teardown()
	Clone() Benchmark
}

func AggregateBenchmarkDurationsMultiple(brrs []BenchmarkDurations) (BenchmarkAggregatedResult, error) {
	var flatBencharkDurations []time.Duration
	for _, benchmarkDurationsChunk := range brrs {
		flatBencharkDurations = append(flatBencharkDurations, benchmarkDurationsChunk.durations...)
	}

	return AggregateBenchmarkDurations(BenchmarkDurations{flatBencharkDurations})
}

func AggregateBenchmarkDurations(brr BenchmarkDurations) (BenchmarkAggregatedResult, error) {
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

func RunBenchmarkSerial(b Benchmark, iterations int) (BenchmarkAggregatedResult, error) {
	benchmarkDurations, err := RunBenchmark(b, iterations)
	benchmarkAggRes, err := AggregateBenchmarkDurations(benchmarkDurations)
	if err != nil {
		fmt.Printf("Error ocurred in RunBenchmarkSerial: %v ", err)
	}

	return benchmarkAggRes, nil
}

func RunBenchmarkParallel(b Benchmark, iterations int, parallelism int) (BenchmarkAggregatedResult, error) {
	iterationsPerThread := int(float64(iterations) / float64(parallelism))
	if iterationsPerThread <= 0 {
		return BenchmarkAggregatedResult{},
			fmt.Errorf("RunBenchmarkParallel failed - iterationsPerThread must be > 0, but is %d. [iterations=%d] [parallelism=%d]",
				iterationsPerThread, iterations, parallelism)
	}

	var wg sync.WaitGroup
	becnhmarkDurationsChan := make(chan BenchmarkDurations, parallelism)
	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			benchmarkDurations, err := RunBenchmark(b.Clone(), iterationsPerThread)
			if err != nil {
				fmt.Printf("Error ocurred in RunBenchmarkParallel: %v\n", err)
			}
			becnhmarkDurationsChan <- benchmarkDurations
		}()
	}
	wg.Wait()
	close(becnhmarkDurationsChan)

	benchmarkDurationsArray := make([]BenchmarkDurations, parallelism)
	for benchmarkDurations := range becnhmarkDurationsChan {
		benchmarkDurationsArray = append(benchmarkDurationsArray, benchmarkDurations)
	}

	benchmarkAggRes, err := AggregateBenchmarkDurationsMultiple(benchmarkDurationsArray)
	if err != nil {
		fmt.Printf("Error ocurred in RunBenchmarkParallel: %v\n", err)
	}

	return benchmarkAggRes, nil
}

func RunBenchmark(b Benchmark, iterations int) (BenchmarkDurations, error) {
	b.Setup()
	defer b.Teardown()

	durations := make([]time.Duration, iterations)
	benchmarkDurations := BenchmarkDurations{}

	for i := 0; i < iterations; i++ {
		var err error
		durations[i], err = b.RunOnce()

		if err != nil {
			return benchmarkDurations, fmt.Errorf("RunBenchmark() - Failed to run benchmark at iteration [%d]; error: %v", i, err)
		}
	}
	benchmarkDurations.durations = durations

	return benchmarkDurations, nil
}
