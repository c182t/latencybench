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
	Min time.Duration `yaml:"min"`
	Avg time.Duration `yaml:"avg"`
	Max time.Duration `yaml:"max"`
	P95 time.Duration `yaml:"p96"`
}

type Benchmark interface {
	Setup() error
	RunOnce() (time.Duration, error)
	Teardown()
	Clone() Benchmark
	GetOptions() *BenchmarkOptions
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

func RunBenchmarkSerial(b Benchmark) (BenchmarkAggregatedResult, error) {
	benchmarkDurations, err := RunBenchmark(b, b.GetOptions().Iterations)
	benchmarkAggRes, err := AggregateBenchmarkDurations(benchmarkDurations)
	if err != nil {
		fmt.Printf("Error ocurred in RunBenchmarkSerial: %v ", err)
	}

	return benchmarkAggRes, nil
}

func RunBenchmarkParallel(b Benchmark) (BenchmarkAggregatedResult, error) {
	iterationsPerThread := int(float64(b.GetOptions().Iterations) / float64(b.GetOptions().Parallelism))
	if iterationsPerThread <= 0 {
		return BenchmarkAggregatedResult{},
			fmt.Errorf("RunBenchmarkParallel failed - iterationsPerThread must be > 0, but is %d. [iterations=%d] [parallelism=%d]",
				iterationsPerThread, b.GetOptions().Iterations, b.GetOptions().Parallelism)
	}

	var wg sync.WaitGroup
	becnhmarkDurationsChan := make(chan BenchmarkDurations, b.GetOptions().Parallelism)
	for i := 0; i < b.GetOptions().Parallelism; i++ {
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

	benchmarkDurationsArray := make([]BenchmarkDurations, b.GetOptions().Parallelism)
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
